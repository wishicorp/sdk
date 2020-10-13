package http_gateway

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
	gin_http "github.com/wishicorp/sdk/framework/gin-http"
	"github.com/wishicorp/sdk/plugin/gateway"
	"github.com/wishicorp/sdk/plugin/pluginregister"
	"github.com/wishicorp/sdk/pool"
	"strings"
	"time"
)

var _ gateway.Gateway = (*HttpGateway)(nil)

type HttpGateway struct {
	pm          *pluginregister.PluginManager
	logger      hclog.Logger
	workerPool  *pool.WorkerPool
	ctx         context.Context
	cancel      context.CancelFunc
	ginServer   *gin_http.Server
	running     chan bool
	workerSize  int
	authMethod  *gateway.Method
	authEnabled bool
	security    gateway.Security
}

func NewGateway(m *pluginregister.PluginManager, workerSize int, logger hclog.Logger) *HttpGateway {
	ctx, cancel := context.WithCancel(context.Background())
	gw := &HttpGateway{
		pm:          m,
		logger:      logger.Named("http-gateway"),
		ctx:         ctx,
		cancel:      cancel,
		running:     make(chan bool, 1),
		workerSize:  workerSize,
		ginServer:   gin_http.NewServer(),
		authEnabled: true,
	}
	return gw
}

func (m *HttpGateway) SetAuthEnabled() {
	m.authEnabled = true
}

func (m *HttpGateway) SetAuthDisabled() {
	m.authEnabled = false
}

func (m *HttpGateway) SetSecurity(security gateway.Security) {
	m.security = security
}

func (m *HttpGateway) SetAuthMethod(method string) error {
	if method == "" {
		return nil
	}
	methods := strings.Split(method, ".")[:]
	if len(methods) != 3 {
		return errors.New("auth method error")
	}
	m.authMethod = &gateway.Method{
		Backend:   methods[0],
		Namespace: methods[1],
		Operation: methods[2],
	}
	return nil
}

//关闭网关
func (m *HttpGateway) Shutdown() {
	defer func() {
		if m.logger.IsTrace() {
			m.logger.Trace("exited")
		}
	}()

	m.ginServer.Shutdown(context.Background())

	m.workerPool.Shutdown()

	select {
	case <-m.workerPool.Running():
	case <-time.After(time.Second * 1):
	}
	close(m.running)
}

//网关是否在运行(阻塞等待)
func (m *HttpGateway) Running() <-chan bool {
	return m.running
}

func (m *HttpGateway) AddRouter(method, router string, handleFunc func(*gin.Context) error) {
	m.ginServer.Router.Handle(method, router, func(c *gin.Context) {
		err := handleFunc(c)
		if nil != err {
			m.logger.Error("handle", "method", method, "router", router, err)
		}
	})
}

func (m *HttpGateway) Listen(addr string, port uint) error {
	if err := m.ginServer.Listen(addr, port); err != nil {
		return err
	}
	return nil
}
func (m *HttpGateway) Serve() error {
	m.startWorkerPool(m.workerSize)
	m.api()
	m.open()
	//m.schemas()
	return m.ginServer.Serve()
}
