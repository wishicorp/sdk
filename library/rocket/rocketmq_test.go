package rocket

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"testing"
)

func TestNewPushConsumer(t *testing.T) {

	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("GID_iot_notice_charge"),
		consumer.WithNameServer([]string{""}),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: "",
			SecretKey: "",
		}),
		consumer.WithInstance("ALIYUN"),
		consumer.WithNamespace(""),
	)
	t.Log(err)
	err = c.Subscribe("iot_notice_charge_report", consumer.MessageSelector{},
		func(ctx context.Context, ext ...*primitive.MessageExt) (result consumer.ConsumeResult, err error) {
			t.Log("received: ", ext)
			return 0, nil
		})
	t.Log(c, err)
	err = c.Start()
	t.Log(err)
	done := make(chan bool)
	<-done
}

func TestNewCommonProducer(t *testing.T) {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{""}),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: "",
			SecretKey: "",
		}),
		producer.WithInstanceName("ALIYUN"),
		producer.WithNamespace(""),
	)
	if err != nil {
		return
	}
	t.Log(err)
	p.Start()
	msg := "{\"orderId\": \"020259e2-7f1d-11ea-8bcd-80e65008316e\", \"ElectricAmount\": 3.4, \"ChargingPower\": 4.5, \"chargerId\": 16357, \"ChargeTime\": 1587432168}"

	result, err := p.SendSync(context.Background(), primitive.NewMessage("iot_notice_charge_report", []byte(msg)))
	t.Log(result, err)

}
