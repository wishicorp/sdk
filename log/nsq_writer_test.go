package log

import (
	"context"
	"testing"
	"time"
)

func TestWithNetWriters(t *testing.T) {
}

func p1(ch chan int) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ch:
			case <-time.After(time.Second * 3):
				cancel()
			}
		}
	}()
	return ctx
}
