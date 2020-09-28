//rocket 对producer和consumer进行了封装实现了线程安全
package rocket

type RktMQConfig struct {
	Broker    string `yaml:"broker"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	NameSpace string `yaml:"name_space"`
	Instance  string `yaml:"instance"`
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
	return NewPushConsumer(r.cfg, gid)
}
