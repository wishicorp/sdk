package pool

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"testing"
	"time"
)

func TestNewWorkerPool(t *testing.T) {
	logger := hclog.New(&hclog.LoggerOptions{Name: "main-pool", Level: hclog.Trace})
	ctx, cancel := context.WithCancel(context.Background())
	pool := NewWorkerPool("main-pool", ctx, logger)

	pool.AddWorker(new("worker-1", ctx, logger))
	pool.AddWorker(new("worker-2", ctx, logger))
	pool.AddWorker(new("worker-3", ctx, logger))
	pool.AddWorker(new("worker-4", ctx, logger))
	pool.AddWorker(new("worker-5", ctx, logger))

	pool.StartWorkers()
	go pool.Start()

	go func() {
		for i := 0; i < 1000; i++ {
			sub := NewSubject(fmt.Sprintf("subject %d", i))
			sub.Observer(NewReader(fmt.Sprintf("reader %d", i)))
			pool.inputChan <- sub
		}
	}()

	time.AfterFunc(time.Second*5, func() {
		cancel()
	})

	time.Sleep(time.Second * 10)
}

func new(name string, ctx context.Context, logger hclog.Logger) *Worker {
	worker := NewWorker(name, ctx, factory, logger)

	return worker
}
