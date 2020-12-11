package rocket

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliyunmq/mq-http-go-sdk"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"sync"
	"time"
)

type callFunc func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)

var _ PushConsumer = (*httpConsumer)(nil)

type httpConsumer struct {
	sync.Mutex
	MQClient   mq_http_sdk.MQClient
	instanceId string
	groupId    string
	handlers   map[string]*consumerHandler
	started    bool
}

func NewHttpConsumer(config *RktMQConfig, gid string) (c PushConsumer, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	// 设置HTTP接入域名（此处以公共云生产环境为例）
	endpoint := config.HttpBroker
	// AccessKey 阿里云身份验证，在阿里云服务器管理控制台创建
	accessKey := config.AccessKey
	// SecretKey 阿里云身份验证，在阿里云服务器管理控制台创建
	secretKey := config.SecretKey
	// 所属的 Topic
	// Topic所属实例ID，默认实例为空
	instanceId := config.NameSpace
	// 您在控制台创建的 Consumer ID(Group ID)
	client := mq_http_sdk.NewAliyunMQClient(endpoint, accessKey, secretKey, "")
	return &httpConsumer{MQClient: client, instanceId: instanceId, groupId: gid, handlers: map[string]*consumerHandler{}}, err
}

func (h *httpConsumer) Start() error {
	h.Lock()
	defer h.Unlock()
	h.started = true
	for _, handler := range h.handlers {
		go handler.start()
	}
	return nil
}

func (h *httpConsumer) Shutdown() error {
	return nil
}

func (h *httpConsumer) Subscribe(topic string, selector consumer.MessageSelector,
	f func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) error {
	h.Lock()
	defer h.Unlock()
	if h.started {
		return errors.New("must before start()")
	}
	mqConsumer := h.MQClient.GetConsumer(h.instanceId, topic, h.groupId, "")
	h.handlers[topic] = &consumerHandler{
		MQConsumer: mqConsumer,
		callFunc:   f,
	}
	return nil
}

func (h *httpConsumer) Unsubscribe(topic string) error {
	h.Lock()
	defer h.Unlock()
	delete(h.handlers, topic)
	return nil
}

type Message mq_http_sdk.ConsumeMessageEntry

func (m Message) toExtMessages(topic string) *primitive.MessageExt {
	result := &primitive.MessageExt{
		MsgId: m.MessageId,
		Message: primitive.Message{
			Topic: topic,
			Body:  []byte(m.MessageBody),
			Queue: &primitive.MessageQueue{Topic: topic},
		},
	}
	result.WithProperties(m.Properties)
	result.WithTag(m.MessageTag)
	return result
}

type consumerHandler struct {
	topic      string
	MQConsumer mq_http_sdk.MQConsumer
	callFunc   callFunc
}

func (c *consumerHandler) start() {
	respChan := make(chan mq_http_sdk.ConsumeMessageResponse)
	errChan := make(chan error)
	defer func() {
		close(respChan)
		close(errChan)
	}()
	go func() {
		for {
			select {
			case resp := <-respChan:
				// 处理业务逻辑
				for _, v := range resp.Messages {
					result, err := c.callFunc(context.Background(), Message(v).toExtMessages(c.topic))
					if nil == err && result == consumer.ConsumeSuccess {
						_ = c.MQConsumer.AckMessage([]string{v.ReceiptHandle})
					}
				}
			case <-errChan:
				time.Sleep(time.Second)
			case <-time.After(30 * time.Second):
			}
		}
	}()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			c.MQConsumer.ConsumeMessage(respChan, errChan,
				8, // 一次最多消费8条(最多可设置为16条)
				1, // 长轮询时间1秒（最多可设置为30秒）
			)
		}
		// 长轮询消费消息
		// 长轮询表示如果topic没有消息则请求会在服务端挂住1s，1s内如果有消息可以消费则立即返回
	}
}
