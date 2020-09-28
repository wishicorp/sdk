package logical

import (
	"context"
	"errors"
	"github.com/wishicorp/sdk/framework/worm"
	nsq_client "github.com/wishicorp/sdk/library/nsq-client"
	"github.com/wishicorp/sdk/library/rabbitmq"
	"github.com/wishicorp/sdk/library/redis"
	"github.com/wishicorp/sdk/library/rocket"
)

const GroupCfgPrimaryKey = "primary"

var ErrConfigGroupPrimaryNotSet = errors.New("config group primary not set")

/**
多数据源采用kv多形式存储, primary 为必须定义的数据源配置
插件的backend初始化就可以以map的方式获取dao对象
mysql "primary" {
  master = ""
  slaves = []
  use_master_slave = false
  show_sql = true
}
*/
type RedisCfgGroup map[string]*redis.Config

/**
同mysql的方式
redis "primary" {
  addrs = ["127.0.0.1:6379"]
  use_cluster = false
  key_predix = "various:dev"
}
*/
type MySQLCfgGroup map[string]*worm.Config
type RabbitMQCfgGroup map[string]*rabbitmq.Config
type ComponentConfig struct {
	Redis    RedisCfgGroup       `json:"redis" hcl:"redis"`
	MySQL    MySQLCfgGroup       `json:"mysql" hcl:"mysql"`
	RabbitMQ RabbitMQCfgGroup    `json:"rabbitmq" hcl:"rabbitmq"`
	RocketMQ *rocket.RktMQConfig `json:"rocketmq" hcl:"rocketmq"`
	Nsq      *nsq_client.Config  `json:"nsq" hcl:"nsq"`
	/**
	扩展参数定义
	*/
	Extras map[string]map[string]interface{} `json:"extras"  hcl:"extras"`
}

func (c *ComponentConfig) Validate() error {
	_, ok := c.Redis[GroupCfgPrimaryKey]
	if !ok {
		return ErrConfigGroupPrimaryNotSet
	}
	_, ok = c.MySQL[GroupCfgPrimaryKey]
	if !ok {
		return ErrConfigGroupPrimaryNotSet
	}
	return nil
}

func (m RedisCfgGroup) Primary() (*redis.Config, bool) {
	p, ok := m[GroupCfgPrimaryKey]
	return p, ok
}
func (m MySQLCfgGroup) Primary() (*worm.Config, bool) {
	p, ok := m[GroupCfgPrimaryKey]
	return p, ok
}
func (m RabbitMQCfgGroup) Primary() (*rabbitmq.Config, bool) {
	p, ok := m[GroupCfgPrimaryKey]
	return p, ok
}
type RedisGroup map[string]redis.RedisCli
type DaoGroup map[string]worm.BaseDao
type RabbitGroup map[string]rabbitmq.RabbitMQ
type RocketGroup map[string]rocket.Producer

func (m RedisGroup) Primary() redis.RedisCli {
	return m[GroupCfgPrimaryKey]
}
func (m DaoGroup) Primary() worm.BaseDao {
	return m[GroupCfgPrimaryKey]
}
func (m RabbitGroup) Primary() rabbitmq.RabbitMQ {
	return m[GroupCfgPrimaryKey]
}

type ComponentFactory interface {
	/**
	namespace 用于隔离多个插件key之间的冲突
	最终组成: 不带namespace  "various:dev:Authorization:062ad5fc-c643-44e6-a916-b29f8bc4384c"
	         带namespace   "various:dev:${namespace}:Authorization:062ad5fc-c643-44e6-a916-b29f8bc4384c"
	*/
	NewRedisClient(ctx context.Context, namespace string) (RedisGroup, error)
	NewMySQLDao(ctx context.Context, namespace string) (DaoGroup, error)
	NewRabbitMQ(ctx context.Context, namespace string) (RabbitGroup, error)
	NewRocketMQ(ctx context.Context, namespace string) (rocket.RocketMQ, error)
	NewNSQ(ctx context.Context, topic string, namespace string) (nsq_client.NSQClient, error)
}

type Component interface {
	FetchConfig(ctx context.Context, namespace, key, version string) (*ComponentConfig, error)
	CreateFactory(ctx context.Context) (ComponentFactory, error)
}

const ConfigExample = `
mysql "primary" {
  master = "root:123456@tcp(127.0.0.1)/virtual_coin?charset=utf8mb4&parseTime=true&loc=Local"
  slaves = [
    "root:123456@tcp(127.0.0.1)/virtual_coin?charset=utf8mb4&parseTime=true&loc=Local",
    "root:123456@tcp(127.0.0.1)/virtual_coin?charset=utf8mb4&parseTime=true&loc=Local"
    ]
  use_master_slave = false
  show_sql = true
}

redis "primary" {
  addrs = ["127.0.0.1:6379"]
  use_cluster = false
}

rabbitmq "primary" {
  url = "amqp://guest:guest@localhost:5672/"
  auto_delete = false
  auto_ack = false
  no_wait = true
  exclusive = true
}
rocketmq{
  broker = ""
  access_key =  ""
  secret_key = ""
  instance = "ALIYUN"
  name_space = ""
}
extras {
    key = {
       foo = "bar"
    }
}
`
