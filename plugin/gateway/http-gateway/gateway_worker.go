package http_gateway

import (
	"context"
	"github.com/wishicorp/sdk/pool"
)

type workerReply struct {
	err    error
	result interface{}
}
type reader struct {
	c chan<- *workerReply
}

func NewReader(c chan<- *workerReply) *reader {
	return &reader{c}
}

func (r *reader) Update(result interface{}, err error) {
	r.c <- &workerReply{err: err, result: result}
}

func (m *HttpGateway) NewObserver(c chan<- *workerReply) pool.Observer {
	return NewReader(c)
}

func (m *HttpGateway) startWorkerPool(workerSize int) {
	poolLogger := m.logger.Named("pool-0")
	m.workerPool = pool.NewWorkerPool("pool-0", context.Background(), poolLogger)

	for i := 0; i < workerSize; i++ {
		m.workerPool.NewWorker(m.backend)
	}

	m.workerPool.StartWorkers()
	go m.workerPool.Start()

}
