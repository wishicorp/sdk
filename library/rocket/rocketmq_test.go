package rocket

import (
	"context"
	"encoding/json"
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
var cfg = `{"broker":"http://MQ_INST_1401091344273301_BcunFXU4.cn-hangzhou.mq-internal.aliyuncs.com:8080","http_broker":"http://1401091344273301.mqrest.cn-hangzhou.aliyuncs.com","access_key":"LTAI4GCELPQryWEhBC5gWShZ","secret_key":"gYc8zpjYfk8FYiRgkeGFlicQpg6vln","name_space":"MQ_INST_1401091344273301_BcunFXU4","instance":"ALIYUN"}`
var c RktMQConfig

func init()  {
	_ = json.Unmarshal([]byte(cfg), &c)
}
func TestNewPushConsumer(t *testing.T) {

	r := NewRocketMQ(&c)
	t.Log(&c)
	c, err := r.NewConsumer("GID_jdsh_test_msg")
	if nil != err {
		t.Fatal(err)
	}
	err = c.Subscribe("jdsh_coupon_combo_order", consumer.MessageSelector{},
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
	r := NewRocketMQ(&c)
	p, err := r.NewProducer()
	if err != nil {
		t.Fatal(err)
	}

	p.Start()
	msg :=
		fmt.Sprintf(`{"orderId": "%s", "ElectricAmount": 3.4, "ChargingPower": 4.5, "chargerId": 16357, "ChargeTime": 1587432168}`,
			uuid.New().String())
	for i := 0; i < 1; i++ {
		result, err := p.SendSync(context.Background(), primitive.NewMessage("jdsh_coupon_combo_pay", []byte(msg)))
		if nil != err {
			t.Fatal(err)
		}
		t.Log(result)
	}
}
