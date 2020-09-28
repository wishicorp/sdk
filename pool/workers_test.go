package pool

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"sort"
	"testing"
)

func TestWorkers_Swap(t *testing.T) {
	workers := make(Workers, 0)
	logger := hclog.Default()
	logger.SetLevel(hclog.Trace)
	w1 := NewWorker("w1", context.Background(), factory, logger)
	w2 := NewWorker("w2", context.Background(), factory, logger)
	w3 := NewWorker("w3", context.Background(), factory, logger)
	w4 := NewWorker("w4", context.Background(), factory, logger)

	w1.Input(nil)
	w1.Input(nil)
	w1.Input(nil)
	w2.Input(nil)
	w3.Input(nil)
	w4.Input(nil)
	w2.Input(nil)
	w1.Input(nil)

	workers = append(workers, w1)
	workers = append(workers, w2)
	workers = append(workers, w3)
	workers = append(workers, w4)

	sort.Sort(workers)
	for _, worker := range workers {
		t.Log(worker.ChanSize())
	}
}
