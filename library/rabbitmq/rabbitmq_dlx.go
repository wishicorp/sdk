package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

//rabbitmq接口
type DLXRabbitMQ interface {
	//定义延迟交换机 ttl 为队列的内的消息国旗时间
	Declare(exchange, queue string, ttl time.Duration) error
	//订阅消息
	Subscribe(callback func(*amqp.Delivery))
	//发布消息
	Publish(message amqp.Publishing) error
	//关闭连接
	Close()
}
type dlxRabbitmq struct {
	*rabbitmq
	sync.Mutex
	msgTTL time.Duration

	dlxSendQueue    string
	dlxSendExchange string

	dlxRecvQueue    string
	dlxRecvExchange string

	sendCh *amqp.Channel
}

func (d *dlxRabbitmq) Declare(exchange, queue string, msgTTL time.Duration) error {
	d.Lock()
	defer d.Unlock()
	d.msgTTL = msgTTL

	d.dlxSendQueue = fmt.Sprintf("dlx-send-%s", queue)
	d.dlxSendExchange = fmt.Sprintf("dlx-send-%s", exchange)

	d.dlxRecvQueue = fmt.Sprintf("dlx-receive-%s", queue)
	d.dlxRecvExchange = fmt.Sprintf("dlx-receive-%s", exchange)

	//死信接收交换机和队列定义
	if err := d.DeclareEx(d.dlxRecvExchange, Fanout); err != nil {
		return err
	}
	if _, err := d.DeclareQueue(d.dlxRecvQueue, amqp.Table{}); err != nil {
		return err
	}
	if err := d.Bind(d.dlxRecvExchange, d.dlxRecvQueue, ""); err != nil {
		return err
	}

	//投递队列，此处默认使用direct模式，如果使用fanout模式需要定义exchange和绑定操作
	args := amqp.Table{
		"x-dead-letter-exchange": d.dlxRecvExchange,
		"x-message-ttl":          msgTTL.Milliseconds(),
	}
	if _, err := d.DeclareQueue(d.dlxSendQueue, args); err != nil {
		return err
	}

	ch, _ := d.conn.Channel()
	d.sendCh = ch

	return nil
}

func (d *dlxRabbitmq) Subscribe(callback func(*amqp.Delivery)) {
	d.rabbitmq.Subscribe(d.dlxRecvQueue, callback)
}

func (d *dlxRabbitmq) Publish(message amqp.Publishing) error {
	return d.sendCh.Publish("", d.dlxSendQueue, false, false, message)
}

//创建rabbitmq客户端
func NewDLXRabbitClient(c *Config) (DLXRabbitMQ, error) {
	conn, err := amqp.Dial(c.Url)
	if nil != err {
		return nil, err
	}
	mq := &rabbitmq{
		config: c,
		conn:   conn,
		mutex:  sync.Mutex{},
	}
	return &dlxRabbitmq{rabbitmq: mq}, nil
}
