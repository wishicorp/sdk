//rocket 对producer和consumer进行了封装实现了线程安全
package rocket

type RktMQConfig struct {
	Broker     string `yaml:"broker" hcl:"broker" json:"broker"`
	HttpBroker string `yaml:"http_broker" hcl:"http_broker" json:"http_broker"`
	AccessKey  string `yaml:"access_key" hcl:"access_key" json:"access_key"`
	SecretKey  string `yaml:"secret_key" hcl:"secret_key" json:"secret_key"`
	NameSpace  string `yaml:"name_space" hcl:"name_space" json:"name_space"`
	Instance   string `yaml:"instance" hcl:"instance" json:"instance"`
}
type RocketMQ interface {
	NewProducer() (Producer, error)
	NewConsumer(gid string) (PushConsumer, error)
}

var _ RocketMQ = (*rocketMQ)(nil)

type rocketMQ struct {
	cfg *RktMQConfig
}

func NewRocketMQ(cfg *RktMQConfig) RocketMQ {
	return &rocketMQ{cfg: cfg}
}

func (r *rocketMQ) NewProducer() (Producer, error) {
	return NewCommonProducer(r.cfg)
}

func (r *rocketMQ) NewConsumer(gid string) (PushConsumer, error) {
	if r.cfg.HttpBroker != ""{
		return NewHttpConsumer(r.cfg, gid)
	}
	return NewPushConsumer(r.cfg, gid)
}
