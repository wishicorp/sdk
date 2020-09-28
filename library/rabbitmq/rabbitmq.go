package rabbitmq

import (
	"crypto/tls"
	"github.com/streadway/amqp"
	"sync"
)

//rabbitmq连接配置
type Config struct {
	Url        string      `hcl:"url"`
	AutoDelete bool        `hcl:"auto_delete"`
	AutoAck    bool        `hcl:"auto_ack"`
	NoWait     bool        `hcl:"no_wait"`
	Exclusive  bool        `hcl:"exclusive"`
	TLSConfig  *tls.Config `hcl:"tls_config"`
}

//rabbitmq接口
type RabbitMQ interface {
	//定义交换机
	DeclareEx(exchange string, kind ExchangeKind) error
	DeclareQueue(queue string, args amqp.Table) (amqp.Queue, error)
	//绑定队列
	Bind(exchange, queue, routing string) error
	//订阅消息
	Subscribe(queue string, callback func(*amqp.Delivery)) error
	//发布消息
	Publish(exchange, routing string, message amqp.Publishing) error
	//关闭连接
	Close()
	NewDLXRabbitMQ() (DLXRabbitMQ, error)
}

// 交换机类型
type ExchangeKind string

const (
	Direct ExchangeKind = "direct"
	Fanout ExchangeKind = "fanout"
	Topic  ExchangeKind = "topic"
)

type rabbitmq struct {
	mutex   sync.Mutex
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *Config
}

//创建rabbitmq客户端
func NewRabbitClient(c *Config) (RabbitMQ, error) {
	var err error
	var conn *amqp.Connection

	if c.TLSConfig != nil {
		conn, err = amqp.DialTLS(c.Url, c.TLSConfig)
	} else {
		conn, err = amqp.Dial(c.Url)
	}
	if nil != err {
		return nil, err
	}
	return &rabbitmq{
		config: c,
		conn:   conn,
		mutex:  sync.Mutex{},
	}, nil
}

func (b *rabbitmq) NewDLXRabbitMQ() (DLXRabbitMQ, error) {
	return NewDLXRabbitClient(b.config)
}

func (b *rabbitmq) DeclareEx(exchange string, kind ExchangeKind) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	ch, _ := b.conn.Channel()

	err := ch.ExchangeDeclare(
		exchange,            // name
		string(kind),        // type
		true,                // durable
		b.config.AutoDelete, // auto-deleted
		false,               // internal
		b.config.NoWait,     // no-wait
		nil,                 // arguments
	)
	b.channel = ch
	return err
}

func (b *rabbitmq) DeclareQueue(queue string, args amqp.Table) (amqp.Queue, error) {
	return b.channel.QueueDeclare(
		queue,               // name
		true,                // durable
		b.config.AutoDelete, // delete when unused
		b.config.Exclusive,  // exclusive
		b.config.NoWait,     // no-wait
		args,                // arguments
	)
}
func (b *rabbitmq) Bind(exchange, queue, routing string) error {
	return b.channel.QueueBind(
		queue,    // queue name
		routing,  // routing key
		exchange, // exchange
		b.config.NoWait,
		nil,
	)
}

func (b *rabbitmq) Subscribe(queue string, callback func(*amqp.Delivery)) error {

	msgCh, err := b.channel.Consume(
		queue,              // queue
		"",                 // consumer
		b.config.AutoAck,   // auto-ack
		b.config.Exclusive, // exclusive
		false,              // no-local
		b.config.NoWait,    // no-wait
		nil,                // args
	)
	if nil != err {
		return err
	}
	go func() {
		for {
			select {
			case msg := <-msgCh:
				callback(&msg)
			}
		}
	}()
	return nil
}

func (b *rabbitmq) Publish(exchange, routing string, message amqp.Publishing) error {
	err := b.channel.Publish(
		exchange, // exchange
		routing,  // routing key
		false,    // mandatory
		false,    // immediate
		message,
	)
	return err
}

func (b *rabbitmq) Close() {
	if nil != b.channel {
		b.channel.Close()
	}
	if nil != b.conn {
		_ = b.conn.Close()
	}

}
