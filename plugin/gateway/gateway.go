package gateway

import "github.com/wishicorp/sdk/plugin/logical"

type Client struct {
	RemoteAddr string
	Referer    string
	UserAgent  string
}

type Method struct {
	Backend   string
	Namespace string
	Operation string
}
type RequestArgs struct {
	Method    string `json:"method" binding:"required"` //${backend}.${namespace}.${operation}
	Version   string `json:"version" binding:"required"`
	Timestamp int    `json:"timestamp" binding:"required"`
	SignType  string `json:"sign_type" binding:"required"`
	Sign      string `json:"sign" binding:"required"`
	Data      string `json:"data"`
}

type Gateway interface {
	SetSecurity(security Security)
	SetPluginAuthorized(authenticator logical.PluginAuthenticator)
	Shutdown()
	Running() <-chan bool
	Listen(addr string, port uint) error
	Serve() error
}
