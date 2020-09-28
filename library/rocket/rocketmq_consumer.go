package rocket

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"sync"
)

type PushConsumer interface {
	rocketmq.PushConsumer
	Subscribe(topic string, selector consumer.MessageSelector,
		f func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) error
	Start() error
}
type pushConsumer struct {
	sync.Mutex
	rocketmq.PushConsumer
}

func NewPushConsumer(config *RktMQConfig, gid string) (PushConsumer, error) {

	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(gid),
		consumer.WithNameServer([]string{config.Broker}),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: config.AccessKey,
			SecretKey: config.SecretKey,
		}),
		consumer.WithInstance(config.Instance),
		consumer.WithNamespace(config.NameSpace),
	)
	if nil != err {
		return nil, err
	}
	return &pushConsumer{PushConsumer: c}, nil
}

//func NewPushConsumer(config *RktMQConfig,
//	messageModel rocket.MessageModel,
//	consumerModel rocket.ConsumerModel)(PushConsumer, error) {
//	c := &rocket.PushConsumerConfig{
//		ClientConfig: defaultConfig(config),
//		Model:         messageModel,
//		ConsumerModel: consumerModel,
//	}
//	consumer, err := rocket.NewPushConsumer(c)
//	if err != nil {
//		return nil,err
//	}
//	return &pushConsumer{PushConsumer: consumer},nil
//}

func (p *pushConsumer) Subscribe(topic string, selector consumer.MessageSelector,
	f func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) error {
	p.Lock()
	defer p.Unlock()
	return p.PushConsumer.Subscribe(topic, selector, f)
}

func (p *pushConsumer) Start() error {
	p.Lock()
	defer p.Unlock()
	return p.PushConsumer.Start()
}
