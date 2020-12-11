package rocket

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"testing"
)

func TestNewPushConsumer(t *testing.T) {

	r := NewRocketMQ(&RktMQConfig{
	})
	c, err := r.NewConsumer("GID_iot_notice_charge")
	if nil != err{
		t.Fatal(err)
	}
	err = c.Subscribe("jdsh_test_msg", consumer.MessageSelector{},
		func(ctx context.Context, ext ...*primitive.MessageExt) (result consumer.ConsumeResult, err error) {
			t.Log("received: ", ctx, ext[0])
			return 0, nil
		})
	if nil != err{
		t.Fatal(err)
	}
	err = c.Start()
	if nil != err{
		t.Fatal(err)
	}
	done := make(chan bool)
	<-done
}

func TestNewCommonProducer(t *testing.T) {
	r := NewRocketMQ(&RktMQConfig{
	})
	p, err := r.NewProducer()
	if err != nil{
		t.Fatal(err)
	}

	p.Start()
	msg := "{" +
		"\"orderId\": \"020259e2-7f1d-11ea-8bcd-80e65008316e\", " +
		"\"ElectricAmount\": 3.4," +
		" \"ChargingPower\": 4.5," +
		" \"chargerId\": 16357," +
		" \"ChargeTime\": 1587432168" +
		"}"
	for i := 0; i < 100; i++ {
		result, err := p.SendSync(context.Background(), primitive.NewMessage("jdsh_test_msg", []byte(msg)))
		t.Log(result, err)
	}


}
