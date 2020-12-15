package pool

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"sort"
	"sync"
	"time"
)

//fanout模型的多任务分发池
//任务数据是一个发布订阅模式的实现
type WorkerPool struct {
	sync.Mutex
	name        string
	logger      hclog.Logger
	running     chan bool
	stopCtx     context.Context //用于从外部context直接关闭pool
	stopChan    chan bool       //用于pool自身的shutdown方法关闭
	workers     Workers
	inputChan   chan *subject //任务输入chan
	workerCount int
	stopped     bool
}

func NewWorkerPool(name string, ctx context.Context, logger hclog.Logger) *WorkerPool {
	return &WorkerPool{
		name:      name,
		stopCtx:   ctx,
		logger:    logger,
		inputChan: make(chan *subject, 163840),
		stopChan:  make(chan bool, 1),
		running:   make(chan bool, 1),
		workers:   make(Workers, 0),
		Mutex:     sync.Mutex{},
	}
}
func (p *WorkerPool) NewWorker(factory Factory) {
	p.Lock()
	p.Unlock()
	name := fmt.Sprintf("worker-%d", p.workerCount)
	p.AddWorker(NewWorker(name, context.Background(), factory, p.logger))
}

//输入数据
func (p *WorkerPool) Input(sub *subject) {
	if p.stopped {
		return
	}
	p.inputChan <- sub
}

//添加worker到pool
func (p *WorkerPool) AddWorker(worker *Worker) {
	p.Lock()
	p.Unlock()
	p.workerCount++
	p.workers = append(p.workers, worker)
}

func (p *WorkerPool) Shutdown() {
	p.stopped = true
	p.stopChan <- true
}

func (p *WorkerPool) Running() <-chan bool {
	return p.running
}

func (p *WorkerPool) Start() {
	if p.logger.IsTrace() {
		p.logger.Trace("started")
	}

	isRunning := true
	for isRunning || len(p.inputChan) > 0 {
		select {
		case sub := <-p.inputChan:
			sort.Sort(p.workers) //取work chan内任务最少的那个
			worker := p.workers[0]
			if err := worker.Input(sub); err != nil {
				p.logger.Error("dispatch worker",
					"worker", worker.name,
					"data", jsonutil.EncodeToString(sub.data),
					"err", err.Error())
			}
		case <-p.stopChan:
			isRunning = false
			p.stopped = true
		case <-p.stopCtx.Done():
			p.stopChan <- true
		}
	}

	p.cleanup()

}

func (p *WorkerPool) StartWorkers() {
	for _, worker := range p.workers {
		go worker.Start()
	}
}

func (p *WorkerPool) cleanup() {
	defer func() {
		if p.logger.IsTrace() {
			p.logger.Trace("exited")
		}
	}()

	if p.logger.IsTrace() {
		p.logger.Trace("cleanup...")
	}
	for i := 0; i < len(p.inputChan); {
		time.Sleep(time.Second)
	}

	for _, worker := range p.workers {
		worker.Stop()
		<-worker.Running()
	}
	close(p.inputChan)
	close(p.running)
	close(p.stopChan)
}
