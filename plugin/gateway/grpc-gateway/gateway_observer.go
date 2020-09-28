package grpc_gateway

import (
	"github.com/wishicorp/sdk/pool"
)

type poolResponse struct {
	err    error
	result interface{}
}
type reader struct {
	c chan<- *poolResponse
}

func NewReader(c chan<- *poolResponse) *reader {
	return &reader{c}
}

func (r *reader) Update(result interface{}, err error) {
	r.c <- &poolResponse{err: err, result: result}
}

func (m *GRPCGateway) NewObserver(c chan<- *poolResponse) pool.Observer {
	return NewReader(c)
}
