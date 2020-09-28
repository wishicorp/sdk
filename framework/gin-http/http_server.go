package gin_http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

type RouterFactory interface {
	RegisterEndpoint(s *Server)
}

type Server struct {
	Router  *gin.Engine
	Handler http.Handler
	server  *http.Server
	addr    string
	port    int
	ls      net.Listener
}

func NewServer() *Server {
	s := gin.New()
	return &Server{Router: s}
}

func (s *Server) AddMiddleware(f func(c *gin.Context)) {
	s.Router.Use(f)
}

func (s *Server) Listen(address string, port uint) (err error) {
	addr := fmt.Sprintf("%s:%d", address, port)
	s.ls, err = net.Listen("tcp", addr)
	return err
}

func (s *Server) Serve() (err error) {
	addr := fmt.Sprintf("%s:%d", s.addr, s.port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.Router,
	}
	s.server = srv
	s.Handler = srv.Handler
	go func() {
		err = srv.Serve(s.ls)
	}()
	return err
}

func (s *Server) Shutdown(ctx context.Context) {
	if s == nil || s.ls == nil {
		return
	}
	_ = s.server.Shutdown(ctx)
}
