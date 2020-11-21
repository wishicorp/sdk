package rabbitmq

import (
	"github.com/streadway/amqp"
	"testing"
	"time"
)

func TestNewDLXRabbitClient(t *testing.T) {
	//
	c := &Config{
		Url:        "amqp://guest:guest@localhost:5672/",
		AutoDelete: false,
		AutoAck:    true,
		NoWait:     false,
		Exclusive:  false,
	}
	broker, err := NewDLXRabbitClient(c)
	if nil != err {
		t.Log("on new broker", err)
	}
	defer broker.Close()

	ex := "exchange"
	queue := "queue"

	if err := broker.Declare(ex, queue, time.Second*5); err != nil {
		t.Fatal(err)
		return
	}
	//total := 5

	go func() {
		broker.Subscribe(func(d *amqp.Delivery) {
			t.Log(time.Now().String()+" received delay msg 1: ", d.RoutingKey, string(d.Body), d.MessageId, d.CorrelationId)
			//不确认重入队列
		})
	}()

	go func() {
		broker.Subscribe(func(d *amqp.Delivery) {
			t.Log(time.Now().String()+" received delay msg 2: ", d.RoutingKey, string(d.Body), d.MessageId, d.CorrelationId)
			//确认消息
		})
	}()
	//for i := 0; i < total; i++ {
	//	err = broker.Publish(amqp.Publishing{
	//		Body: []byte("test delayed message: " + time.Now().String()),
	//	})
	//	t.Log("on  Publish", err)
	//}
	done := make(chan bool)
	<-done
}
