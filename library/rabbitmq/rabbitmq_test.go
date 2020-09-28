package rabbitmq

import (
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

func TestBroker(t *testing.T) {
	//
	c := &Config{
		Url:        "amqp://guest:guest@localhost:5672/",
		AutoDelete: false,
		AutoAck:    true,
		NoWait:     false,
		Exclusive:  false,
	}
	broker, err := NewRabbitClient(c)
	if nil != err {
		t.Log("on new broker", err)
	}
	ex := "test-broker-ex"
	queue := "test-broker-queue"
	routing := "test-broker-routing"

	broker.DeclareEx(ex, Fanout)
	broker.DeclareQueue(queue, amqp.Table{})
	broker.Bind(ex, queue, routing)

	go func() {
		broker.Subscribe(queue, func(d *amqp.Delivery) {
			t.Log("on Subscribe ", d.RoutingKey, string(d.Body), d.MessageId, d.CorrelationId)
		})
	}()
	go func() {
		broker.Subscribe(queue, func(d *amqp.Delivery) {
			t.Log("on Subscribe ", d.RoutingKey, string(d.Body), d.MessageId, d.CorrelationId)
		})
	}()
	err = broker.Publish(ex, routing, amqp.Publishing{
		MessageId:     uuid.New().String(),
		CorrelationId: uuid.New().String(),
		AppId:         uuid.New().String(),
		Type:          "text/plain",
		Body:          []byte("test message: " + time.Now().String()),
	})
	t.Log("on  Publish", err)
	broker.Close()
}
