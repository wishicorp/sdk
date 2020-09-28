package nsq_client

import (
	"context"
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/nsqio/go-nsq"
	"time"
)

type publisher struct {
	ProducerAddr string
	Producer     *Producer
	topic        string
	inCh         chan []byte
	logger       hclog.Logger
}

func NewPublisher(topic string, p *Producer, logger hclog.Logger) *publisher {
	return &publisher{Producer: p, topic: topic, inCh: make(chan []byte, 163840), logger: logger}
}

func (s *publisher) Publish(in []byte) {
	s.inCh <- in
}

//启动消息发送
//inCh 数据输入chan
//routineCtx 子routine的context，用于通知父routine自身已停止
func (s *publisher) StarPublisher(useAsync bool, batchSize int,
	parentContext context.Context) (routineCtx context.Context, err error) {

	//用于ping，失败后通知publish routine，主动推出循环并且执行执行cancel方法通知父routine该routine已成退出
	stopCh := make(chan bool)

	s.Producer.PingWithStopCh(stopCh, time.Second)

	ctx, cancel := context.WithCancel(parentContext)
	if useAsync && batchSize == 0 {
		s.asyncPublish(stopCh, cancel)
		return ctx, nil
	}
	if useAsync && batchSize > 0 {
		s.asyncMultiPublish(stopCh, cancel, batchSize)
		return ctx, nil
	}
	s.publish(stopCh, cancel)
	return ctx, nil
}

func (s *publisher) asyncPublish(stopCh <-chan bool, cancel context.CancelFunc) {
	go func() {
		defer cancel()
		canceled := false
		errChan := make(chan error)
		for !canceled {
			select {
			case <-stopCh:
				canceled = true
			case err := <-errChan:
				s.logger.Error("publish async received", "topic", s.topic, "err", err.Error())
				if s.Producer.Ping() != nil {
					canceled = true
				}
			case data := <-s.inCh:
				if doneChan, err := PublishAsync(s.topic, data, s.Producer); err != nil {
					s.logger.Error("publish async", "topic", s.topic, "err", err.Error())
					if s.Producer.Ping() != nil {
						canceled = true
					}
				} else {
					CheckAsyncResult(doneChan, errChan, data)
				}
			}
		}
	}()
}

func (s *publisher) asyncMultiPublish(stopCh <-chan bool, cancel context.CancelFunc, batchSize int) chan<- []byte {
	inCh := make(chan []byte, 128)
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		defer cancel()
		canceled := false
		for !canceled {
			data := make([][]byte, 0)
			for i := 0; i < batchSize; {
				select {
				case <-stopCh: //收到ping失败的消息
					i = batchSize << 2
					canceled = true
				case in := <-inCh:
					i++
					data = append(data, in)
				case <-ticker.C: //接收超时
					i = batchSize << 2
				}
			}
			if len(data) > 0 {
				if _, err := PublishMultiAsync(s.topic, data, s.Producer); err != nil {
					s.logger.Error("publish multi async", "topic", s.topic, "err", err.Error())
					canceled = true
				}
			}
		}
	}()
	return inCh
}

func (s *publisher) publish(stopCh <-chan bool, cancel context.CancelFunc) {
	go func() {
		defer cancel()
		canceled := false
		for !canceled {
			select {
			case <-stopCh:
				canceled = true
			case in := <-s.inCh:
				if err := Publish(s.topic, in, s.Producer); err != nil {
					s.logger.Error("publish", "topic", s.topic, "err", err.Error())
					if s.Producer.Ping() != nil {
						canceled = true
					}
				}
			}
		}
	}()
}

func Publish(topic string, data []byte, producer *Producer) error {
	err := producer.Publish(topic, data)
	if err != nil {
		return err
	}
	return nil
}

func PublishAsync(topic string, data []byte, producer *Producer) (<-chan *nsq.ProducerTransaction, error) {
	responseChan := make(chan *nsq.ProducerTransaction, 1)
	err := producer.PublishAsync(topic, data, responseChan, nil)
	if err != nil {
		return nil, err
	}
	return responseChan, nil
}

func CheckAsyncResult(doneChan <-chan *nsq.ProducerTransaction, errChan chan<- error, data []byte) {
	go func() {
		select {
		case result := <-doneChan:
			if result.Error != nil {
				errChan <- errors.New(result.Error.Error() + ", DATA=" + string(data))
			}
		case <-time.After(time.Second * 30):
			errChan <- errors.New("wait result timeout with 30s, DATA=" + string(data))
		}
	}()
}

func PublishMultiAsync(topic string, data [][]byte, producer *Producer) (<-chan *nsq.ProducerTransaction, error) {

	doneChan := make(chan *nsq.ProducerTransaction, len(data))
	err := producer.MultiPublishAsync(topic, data, doneChan, nil)
	if err != nil {
		return nil, err
	}
	return doneChan, nil
}
