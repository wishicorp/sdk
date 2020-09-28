package grpc_gateway

import (
	"context"
	"github.com/wishicorp/sdk/pool"
)

func (m *GRPCGateway) startWorkerPool(workerSize int) {
	poolLogger := m.logger.Named("pool-0")
	m.workerPool = pool.NewWorkerPool("pool-0", context.Background(), poolLogger)

	for i := 0; i < workerSize; i++ {
		m.workerPool.NewWorker(m.backend)
	}

	m.workerPool.StartWorkers()
	go m.workerPool.Start()

}
