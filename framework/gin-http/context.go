package gin_http

import (
	"github.com/gin-gonic/gin"
	"github.com/wishicorp/sdk/helper/httputil"
)

type Context struct {
	*gin.Context
}

func (c *Context) SuccessJSON(data interface{}) {
	c.Context.SecureJSON(200, httputil.SuccessResponse(data))
}
func (c *Context) SuccessPaginationJSON(total int64, data interface{}) {
	pageable := httputil.ParsePageable(c.Context)
	pagination := httputil.Pagination(int(total), pageable)
	c.Context.SecureJSON(200, httputil.PageResponse(pagination, data))
}

func (c *Context) FailureJSON(code ResponseCode, message string) {
	c.Context.SecureJSON(200, httputil.ErrorResponse(int(code), message))
}
func (c *Context) DefaultFailureJSON(message string) {
	c.Context.SecureJSON(200, httputil.ErrorResponse(1, message))
}

type HandlerFunc func(*Context)

func ContextHandler(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := new(Context)
		context.Context = c
		handler(context)
	}
}
