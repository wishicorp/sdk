package http_gateway

import (
	"context"
	"fmt"
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
					"err", err, "took", time.Since(then))
			} else {
				if m.logger.IsTrace() {
					m.logger.Trace("backend", "id", data.request.ID,
						"name", backendName, "namespace", data.request.Namespace, "status", "finished",
						"took", time.Since(then))
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
		if m.authEnabled && m.authMethod != nil {
			authReply, err := m.authorization(backend, data.request)
			if err != nil {
				return nil, fmt.Errorf("auth: %s", err.Error())
			}
			if authReply.ResultCode != 0 {
				return authReply, nil
			}
			if err := jsonutil.Swap(authReply.Content.Data, &data.request.Authorized); err != nil {
				return nil, err
			}
		}

		result, err = backend.HandleRequest(context.Background(), data.request)
		return result, err
	}
}
