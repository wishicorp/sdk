package http_gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"github.com/wishicorp/sdk/plugin"
	"github.com/wishicorp/sdk/plugin/gateway"
	"github.com/wishicorp/sdk/plugin/gateway/consts"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/pluginregister"
)

func (m *HttpGateway) writerReply(c *gin.Context, resp *workerReply) {

	if m.logger.IsTrace() {
		m.traceReply(c, resp)
	}

	if resp.err == pluginregister.PluginNotExists {
		reply := gateway.Error(
			consts.ReplyCodeBackendNotExists,
			resp.err.Error())

		c.SecureJSON(200, reply)
		return
	}

	if resp.err == plugin.ErrPluginShutdown {
		reply := gateway.Error(
			consts.ReplyCodeBackendShutdown,
			resp.err.Error())

		c.SecureJSON(200, reply)
		return
	}
	if resp.err == logical.ErrAuthorizationTokenRequired || resp.err == logical.ErrAuthorizationTokenInvalid {
		reply := gateway.Error(consts.ReplyCodeAuthorizedRequired, resp.err.Error())
		c.SecureJSON(200, reply)
		return
	}
	if nil != resp.err {
		reply := gateway.Error(consts.ReplyCodeFailure, resp.err.Error())
		c.SecureJSON(200, reply)
		return
	}
	c.SecureJSON(200, gateway.Success(resp.result))
}

func (m *HttpGateway) traceReply(c *gin.Context, resp *workerReply) {
	if m.logger.IsTrace() {
		m.logger.Trace(
			"http-gateway trace",
			"path", c.Request.RequestURI,
			"method", c.Request.Method,
			"client", getRemoteAddr(c),
			"response", jsonutil.EncodeToString(resp.result),
			"err", resp.err,
		)
	}
}
