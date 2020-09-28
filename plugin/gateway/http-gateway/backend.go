package http_gateway

import (
	"context"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/pluginregister"
	"time"
)

type workData struct {
	backend string
	request *logical.Request
}

//网关处理客户端http请求
func (m *HttpGateway) backend() func(i interface{}) (interface{}, error) {
	return func(i interface{}) (result interface{}, err error) {
		data := i.(*workData)
		backendName := data.backend

		defer func(then time.Time) {
			if nil != err {
				m.logger.Error("backend", "id", data.request.ID,
					"name", backendName, "namespace", data.request.Namespace, "status", "finished",
					"resp", jsonutil.EncodeToString(result), "err", err, "took", time.Since(then))
			} else {
				if m.logger.IsTrace() {
					m.logger.Trace("backend", "id", data.request.ID,
						"name", backendName, "namespace", data.request.Namespace, "status", "finished",
						"resp", jsonutil.EncodeToString(result), "took", time.Since(then))
				}
			}
		}(time.Now())

		if m.logger.IsTrace() {
			m.logger.Trace("backend", "id", data.request.ID, "name", backendName,
				"namespace", data.request.Namespace,
				"status", "started", "request", jsonutil.EncodeToString(data.request))
		}

		backend, has := m.pm.GetBackend(backendName)
		if !has {
			return nil, pluginregister.PluginNotExists
		}

		backend.Incr()
		defer backend.DeIncr()

		_, err = m.authorization(backend, data.request)
		if err != nil {
			return nil, err
		}

		result, err = backend.HandleRequest(context.Background(), data.request)

		return result, err
	}
}

func (m *HttpGateway) authorization(backend logical.Backend, request *logical.Request) (logical.Authorized, error) {
	if m.authenticator == nil {
		return logical.Authorized{}, logical.ErrAuthenticatorNotSet
	}

	auth, err := m.authenticator.Authorization(backend, request)
	return auth, err
}
