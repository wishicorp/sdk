package gateway

import (
	"github.com/wishicorp/sdk/plugin/gateway/consts"
)

// http返回数据结构
type Reply struct {
	Code    consts.ReplyCode `json:"code"`
	Result  interface{}      `json:"result,omitempty"`
	Message string           `json:"message,omitempty"`
}

func Error(code consts.ReplyCode, message string) *Reply {
	return &Reply{Code: code, Message: message}
}

func Success(result interface{}) *Reply {
	return &Reply{Code: 0, Result: result}
}
