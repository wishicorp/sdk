package rocket

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/google/uuid"
	"os"
	"testing"
)

func mockConfig() *RktMQConfig {
	return &RktMQConfig{
		Broker:     os.Getenv("ROCKET_BROKER"),
		HttpBroker: os.Getenv("ROCKET_HTTP_BROKER"),
		AccessKey:  os.Getenv("ROCKET_ACCESS_KEY"),
		SecretKey:  os.Getenv("ROCKET_SECRET_KEY"),
		NameSpace:  os.Getenv("ROCKET_NAMESPACE"),
		Instance:   os.Getenv("ROCKET_INSTANCE"),
	}
}

func TestNewPushConsumer(t *testing.T) {

	r := NewRocketMQ(mockConfig())
	t.Log(mockConfig())
	c, err := r.NewConsumer("GID_iot_notice_charge")
	if nil != err {
		t.Fatal(err)
	}
	err = c.Subscribe("jdsh_test_msg", consumer.MessageSelector{},
		func(ctx context.Context, ext ...*primitive.MessageExt) (result consumer.ConsumeResult, err error) {
			t.Log("received: ", string(ext[0].MsgId))
			return 0, nil
		})
	if nil != err {
		t.Fatal(err)
	}
	err = c.Start()
	if nil != err {
		t.Fatal(err)
	}
	done := make(chan bool)
	<-done
}

func TestNewCommonProducer(t *testing.T) {
	r := NewRocketMQ(mockConfig())
	p, err := r.NewProducer()
	if err != nil {
		t.Fatal(err)
	}

	p.Start()
	msg :=
		fmt.Sprintf(`{"orderId": "%s", "ElectricAmount": 3.4, "ChargingPower": 4.5, "chargerId": 16357, "ChargeTime": 1587432168}`,
			uuid.New().String())
	for i := 0; i < 1; i++ {
		result, err := p.SendSync(context.Background(), primitive.NewMessage("jdsh_test_msg", []byte(msg)))
		if nil != err {
			t.Fatal(err)
		}
		t.Log(result)
	}

}
