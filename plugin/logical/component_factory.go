package logical

import (
	"context"
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/wishicorp/sdk/framework/worm"
	nsq_client "github.com/wishicorp/sdk/library/nsq-client"
	"github.com/wishicorp/sdk/library/rabbitmq"
	"github.com/wishicorp/sdk/library/redis"
	"github.com/wishicorp/sdk/library/rocket"
)

var ErrConfigEmpty = errors.New("config content empty")
var _ Component = (*component)(nil)
var _ ComponentFactory = (*componentFactory)(nil)

type componentFactory struct {
	cfg    *ComponentConfig
	logger hclog.Logger
}

type component struct {
	cfg    *ComponentConfig
	logger hclog.Logger
}

func NewComponent(cfg *ComponentConfig, logger hclog.Logger) *component {
	return &component{cfg: cfg, logger: logger}
}

func (c *component) FetchConfig(ctx context.Context, namespace, key, version string) (*ComponentConfig, error) {
	return c.cfg, nil
}
func (c *component) CreateFactory(ctx context.Context) (ComponentFactory, error) {
	return &componentFactory{cfg: c.cfg, logger: c.logger}, nil
}

func (c *componentFactory) NewRedisClient(ctx context.Context, namespace string) (RedisGroup, error) {
	if nil == c.cfg || nil == c.cfg.Redis || len(c.cfg.Redis) == 0 {
		return nil, ErrConfigEmpty
	}
	results := make(map[string]redis.RedisCli)
	for key, value := range c.cfg.Redis {
		prefix := value.KeyPrefix
		if namespace != "" {
			prefix = prefix + redis.RedisKeySep + namespace
		}
		cmd, err := redis.NewRedisCmd(value)
		if nil != err {
			return nil, err
		}
		results[key] = redis.NewRedisView(cmd, prefix, c.logger)
	}

	return results, nil
}

func (c *componentFactory) NewMySQLDao(ctx context.Context, namespace string) (DaoGroup, error) {
	if nil == c.cfg || nil == c.cfg.MySQL || len(c.cfg.MySQL) == 0 {
		return nil, ErrConfigEmpty
	}
	results := make(map[string]worm.BaseDao)
	namedLog := c.logger.ResetNamed(namespace)
	for key, value := range c.cfg.MySQL {
		conn, err := worm.NewConn(value, namedLog)
		if err != nil {
			return nil, err
		}
		results[key] = worm.NewBaseDao(conn)
	}

	return results, nil
}

func (c *componentFactory) NewRabbitMQ(ctx context.Context, namespace string) (RabbitGroup, error) {
	if nil == c.cfg || nil == c.cfg.RabbitMQ || len(c.cfg.RabbitMQ) == 0 {
		return nil, ErrConfigEmpty
	}
	group := RabbitGroup{}
	for name, config := range c.cfg.RabbitMQ {
		rb , err := rabbitmq.NewRabbitClient(config)
		if nil != err{
			return nil, err
		}
		group[name] = rb
	}
	return group, nil
}

func (c *componentFactory) NewRocketMQ(ctx context.Context, namespace string) (rocket.RocketMQ, error) {
	if nil == c.cfg || nil == c.cfg.RocketMQ {
		return nil, ErrConfigEmpty
	}
	return rocket.NewRocketMQ(c.cfg.RocketMQ), nil
}

func (c *componentFactory) NewNSQ(ctx context.Context, namespace string, topic string) (nsq_client.NSQClient, error) {
	if nil == c.cfg || nil == c.cfg.Nsq {
		return nil, ErrConfigEmpty
	}
	return nsq_client.NewNsqClient(c.cfg.Nsq, topic, c.logger.ResetNamed("nsq"))
}
