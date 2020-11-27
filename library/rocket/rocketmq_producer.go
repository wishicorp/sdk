package rocket

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"sync"
)

type Producer interface {
	rocketmq.Producer
	// SendMessageSync send a message with sync
	Start() error
	Shutdown() error
	SendSync(ctx context.Context, mq ...*primitive.Message) (*primitive.SendResult, error)
	SendAsync(ctx context.Context, mq func(ctx context.Context, result *primitive.SendResult, err error),
		msg ...*primitive.Message) error
	SendOneWay(ctx context.Context, mq ...*primitive.Message) error
}

type commonProducer struct {
	sync.Mutex
	rocketmq.Producer
}

func NewCommonProducer(config *RktMQConfig) (Producer, error) {
	p, _ := rocketmq.NewProducer(
		producer.WithNameServer([]string{config.Broker}),
		producer.WithRetry(3),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: config.AccessKey,
			SecretKey: config.SecretKey,
		}),
		producer.WithInstanceName(config.Instance),
		producer.WithNamespace(config.NameSpace),
	)
	return &commonProducer{Producer: p}, nil
}

//func NewOrderlyProducer(config *RktMQConfig) (Producer, error) {
//	c := &rocketmq.{
//		ClientConfig:  defaultConfig(config),
//		ProducerModel: rocketmq.OrderlyProducer,
//	}
//	producer, err := rocketmq.NewProducer(c)
//	if err != nil {
//		return nil, err
//	}
//	err = producer.Start()
//	if err != nil {
//		return nil, err
//	}
//	return &commonProducer{Producer: producer}, nil
//}
//
//func NewTransProducer(config *RktMQConfig) (Producer, error) {
//	c := &rocketmq.ProducerConfig{
//		ClientConfig:  defaultConfig(config),
//		ProducerModel: rocketmq.TransProducer,
//	}
//	producer, err := rocketmq.NewProducer(c)
//	if err != nil {
//		return nil, err
//	}
//	err = producer.Start()
//	if err != nil {
//		return nil, err
//	}
//	return &commonProducer{Producer: producer}, nil
//}

//func NewTransactionProducer(config *RktMQConfig,
//	producerModel rocketmq.ProducerModel,
//	listener rocketmq.TransactionLocalListener, arg interface{}) (rocketmq.TransactionProducer, error) {
//	c := &rocketmq.ProducerConfig{
//		ClientConfig:  defaultConfig(config),
//		ProducerModel: producerModel,
//	}
//	producer, err := rocketmq.NewTransactionProducer(c, listener, arg)
//	if err != nil {
//		return nil, err
//	}
//	err = producer.Start()
//	if err != nil {
//		return nil, err
//	}
//	return producer, nil
//}

func (p *commonProducer) Start() error {
	p.Lock()
	defer p.Unlock()
	return p.Producer.Start()
}

func (p *commonProducer) Shutdown() error {
	p.Lock()
	defer p.Unlock()
	return p.Producer.Shutdown()
}

func (p *commonProducer) SendSync(ctx context.Context, msg ...*primitive.Message) (*primitive.SendResult, error) {
	p.Lock()
	defer p.Unlock()
	return p.Producer.SendSync(ctx, msg...)
}
func (p *commonProducer) SendAsync(ctx context.Context, mq func(ctx context.Context, result *primitive.SendResult, err error),
	msg ...*primitive.Message) error {
	p.Lock()
	defer p.Unlock()
	return p.Producer.SendAsync(ctx, mq, msg...)
}
func (p *commonProducer) SendOneWay(ctx context.Context, msg ...*primitive.Message) error {
	p.Lock()
	defer p.Unlock()
	return p.Producer.SendOneWay(ctx, msg...)
}
