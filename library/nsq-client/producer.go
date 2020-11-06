package nsq_client

import (
	"github.com/nsqio/go-nsq"
	"time"
)

type Producer struct {
	*nsq.Producer
}
// stopCh 停止变量出管道
func (producer *Producer) PingWithStopCh(stopCh chan<- bool, duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		stopped := false
		for !stopped {
			select {
			case <-ticker.C:
				if err := producer.Producer.Ping(); err != nil {
					stopCh <- true
					stopped = true
					ticker.Stop()
				}
			}
		}
	}()
}