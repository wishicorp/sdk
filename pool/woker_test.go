package pool

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"strings"
	"testing"
	"time"
)

type Reader struct {
	name string
}

func NewReader(name string) *Reader {
	return &Reader{
		name: name,
	}
}

func (r *Reader) Update(result interface{}, err error) {
	fmt.Println("observer", r.name, result, err)
}

func factory() func(i interface{}) (interface{}, error) {
	return func(i interface{}) (interface{}, error) {
		if strings.Contains(i.(string), "subject 5") {
			fmt.Println("sleeping.......")
			//time.Sleep(time.Second * 15)
		}
		return i.(string) + " factory executed!", nil
	}
}
func TestNewWorker(t *testing.T) {
	logger := hclog.Default()
	logger.SetLevel(hclog.Trace)
	ctx, _ := context.WithCancel(context.Background())
	worker := NewWorker("worker-1", ctx, factory, logger)

	sub := NewSubject("worker test")
	sub.Observer(NewReader("reader"))

	go worker.Start()

	for i := 0; i < 1000; i++ {
		if err := worker.Input(sub); err != nil {
			return
		}
	}

	time.AfterFunc(time.Second*5, func() {
		worker.Stop()
	})

	time.Sleep(time.Second * 7)

}
