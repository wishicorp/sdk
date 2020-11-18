package http_gateway

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"github.com/wishicorp/sdk/plugin/gateway"
	"github.com/wishicorp/sdk/plugin/gateway/consts"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/pool"
	"strings"
)

func (m *HttpGateway) open(basePath string) {
	m.ginServer.Router.POST(basePath+"/open", func(c *gin.Context) {
		request := new(gateway.RequestArgs)
		if err := c.ShouldBindJSON(request); err != nil {
			c.SecureJSON(200, gateway.Error(consts.ReplyCodeFailure, err.Error()))
			return
		}
		methods := strings.Split(request.Method, ".")[:]
		if len(methods) != 3 {
			c.SecureJSON(200,
				gateway.Error(consts.ReplyCodeMethodInvalid, "method error"))
			return
		}
		method := gateway.Method{
			Backend:   methods[0],
			Namespace: methods[1],
			Operation: methods[2],
		}

		if m.security != nil {
			client := &gateway.Client{
				RemoteAddr: GetRemoteAddr(c),
				Referer:    c.Request.Referer(),
				UserAgent:  c.Request.UserAgent(),
			}
			if err := m.security.Blocker(&method, client); err != nil {
				c.SecureJSON(200,
					gateway.Error(consts.ReplyCodeReqBlocked, err.Error()))
				return
			}

			if err := m.security.RateLimiter(&method, client); err != nil {
				c.SecureJSON(200, gateway.Error(consts.ReplyCodeRateLimited, err.Error()))
				return
			}
		}
		bs, _ := jsonutil.EncodeJSON(request)
		m.invokeRequest(c, method, string(bs))
	})
}
func (m *HttpGateway) api(basePath string) {
	m.ginServer.Router.POST(basePath+"/api", func(c *gin.Context) {
		request := new(gateway.RequestArgs)
		if err := c.ShouldBindJSON(request); err != nil {
			c.SecureJSON(200, gateway.Error(consts.ReplyCodeFailure, err.Error()))
			return
		}
		methods := strings.Split(request.Method, ".")[:]
		if len(methods) != 3 {
			c.SecureJSON(200,
				gateway.Error(consts.ReplyCodeMethodInvalid, "method error"))
			return
		}
		method := gateway.Method{
			Backend:   methods[0],
			Namespace: methods[1],
			Operation: methods[2],
		}

		if m.security != nil {
			if !m.security.SignVerify(request) {
				c.SecureJSON(200,
					gateway.Error(consts.ReplyCodeSignInvalid, "invalid sign"))
				return
			}
			client := &gateway.Client{
				RemoteAddr: GetRemoteAddr(c),
				Referer:    c.Request.Referer(),
				UserAgent:  c.Request.UserAgent(),
			}
			if err := m.security.Blocker(&method, client); err != nil {
				c.SecureJSON(200,
					gateway.Error(consts.ReplyCodeReqBlocked, err.Error()))
				return
			}

			if err := m.security.RateLimiter(&method, client); err != nil {
				c.SecureJSON(200, gateway.Error(consts.ReplyCodeRateLimited, err.Error()))
				return
			}
		}
		m.invokeRequest(c, method, request.Data)
	})

}

func (m *HttpGateway) schemas(basePath string) {
	m.ginServer.Router.GET(basePath+"/schemas", func(c *gin.Context) {
		backends := make([]map[string]string, 0)
		backendName := c.Query("backend")
		namespace := c.Query("namespace")
		schemas := map[string][]*logical.NamespaceSchema{}
		if backendName != "" {
			backends = append(backends, map[string]string{"name": backendName})
		} else {
			backends = m.pm.List()
		}

		for _, bMap := range backends {
			backend, has := m.pm.GetBackend(bMap["name"])
			if !has {
				continue
			}
			resp, err := backend.SchemaRequest(context.Background())
			if nil != err {
				c.SecureJSON(200, gateway.Error(consts.ReplyCodeFailure, err.Error()))
				return
			}
			if namespace != "" {
				for _, schema := range resp.NamespaceSchemas {
					if schema.Namespace == namespace {
						schemas[namespace] = []*logical.NamespaceSchema{schema}
					}
				}
				break
			} else {
				schemas[bMap["name"]] = resp.NamespaceSchemas
			}
		}

		c.SecureJSON(200, gateway.Success(schemas))

	})
}

func (m *HttpGateway) invokeRequest(c *gin.Context, method gateway.Method, data string) {
	request := &logical.Request{
		ID:        uuid.New().String(),
		Namespace: method.Namespace,
		Operation: logical.Operation(method.Operation),
		Data: map[string][]byte{
			"data": []byte(data),
		},
		Headers: c.Request.Header,
		Token:   c.Request.Header.Get(logical.AuthTokenName),
		Connection: &logical.Connection{
			RemoteAddr: GetRemoteAddr(c),
			ConnState:  c.Request.TLS,
		},
	}
	workerData := &workData{
		backend: method.Backend,
		request: request,
	}
	output := make(chan *workerReply, 1)
	subject := pool.NewSubject(workerData)

	subject.Observer(m.NewObserver(output))

	m.workerPool.Input(subject)

	select {
	case d := <-output:
		m.writerReply(c, d)
	}
}

func GetRemoteAddr(c *gin.Context) string {
	remoteAddr := c.Request.RemoteAddr
	if remoteAddr == "127.0.0.1" {
		remoteAddr = c.GetHeader("X-Forwarded-For")
	}
	return remoteAddr
}
