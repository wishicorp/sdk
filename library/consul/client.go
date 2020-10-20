//consul注册客户端
package consul

import (
	"crypto/tls"
	"github.com/hashicorp/consul/api"
	"net/http"
	"sync"
	"time"
)

type Client interface {
	//native registry client
	Client() *api.Client
	Register(s *Service) error
	DeRegister(s *Service) error

	//让服务进入维护模式，不可用
	Maintenance(id string, reason string) error
	DeMaintenance(id string) error

	//获取服务列表
	GetServices(id, tag string) ([]*api.AgentService, error)
	GetService(id, tag string) (*api.AgentService, error)
	//获取微服务路径
	GetServiceAddrPort(name string, useLan bool, tags string) (host string, port int, err error)
	//微服务客户端
	GetMicroHTTPClient(id string, useLan bool, tags string, header map[string][]string) (MicroHTTPClient, error)
	//丛配置中心加载配置文件,out对象必须为对象引用
	LoadConfig(out interface{}) error

	//创建一个session,ttl需大于15秒,behavior定义了session到期后的动作
	//如需深度定制session请获取native客户端创建
	NewSession(name string, ttl time.Duration, behavior SessionBehavior, opts *api.WriteOptions) (string, error)
	SessionInfo(id string, opts *api.QueryOptions) (*api.SessionEntry, error)
	//销毁session
	DestroySession(id string, opts *api.WriteOptions) error

	//对一个kv进行加锁
	//another标志key是否被其它session锁定
	//err==nil && false == another 加锁成功
	KVAcquire(key, session string, opts *api.QueryOptions) (success bool, err error)
	//释放一个session的锁
	KVRelease(key string, opts *api.QueryOptions) error
	//获取kv信息
	KVInfo(key string, opts *api.QueryOptions) (*api.KVPair, error)
	//检查或者设置key
	KVCas(p *api.KVPair, opts *api.WriteOptions) (bool, error)
	//获取kv并且序列化到out对象
	KVFire(key string, opts *api.QueryOptions, out interface{}) error

	Config() *Config

	queryOptions(*api.QueryOptions) *api.QueryOptions
	writeOptions(*api.WriteOptions) *api.WriteOptions
}

type Config struct {
	Datacenter  string
	ZoneAddress string //consul注册地址 127.0.0.1:8500
	Token       string
	Application struct {
		Name    string //应用名称
		Profile string //环境变量 dev test prod ...
	}
	Config struct {
		DataKey string //配置key
		Format  string //配置格式
	}
	TLSConfig api.TLSConfig
}

type client struct {
	sync.RWMutex
	client         *api.Client
	config         *Config
	cachedServices map[string][]*api.AgentService
}

func NewClient(c *Config) (Client, error) {
	cf := api.DefaultConfig()
	if c.Datacenter != "" {
		cf.Datacenter = c.Datacenter
	}
	cf.HttpClient = http.DefaultClient
	cf.HttpClient.Timeout = time.Second * 10
	cf.HttpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		TLSHandshakeTimeout: 30,
	}

	cf.Address = c.ZoneAddress
	cf.Token = c.Token
	cf.TLSConfig = c.TLSConfig
	cli, err := api.NewClient(cf)
	if nil != err {
		return nil, err
	}

	_client := &client{client: cli, config: c, cachedServices: map[string][]*api.AgentService{}}
	_client.autoPullService(time.Second * 30)
	return _client, err
}

func (c *client) Client() *api.Client {
	return c.client
}
func (c *client) Config() *Config {
	return c.config
}

func (c *client) queryOptions(opts *api.QueryOptions) *api.QueryOptions {
	q := &api.QueryOptions{RequireConsistent: true}
	if c.config.Datacenter != "" {
		q.Datacenter = c.config.Datacenter
	}
	q.Token = c.config.Token

	if opts != nil {
		if opts.Token != "" {
			q.Token = opts.Token
		}
		if opts.NodeMeta != nil {
			q.NodeMeta = opts.NodeMeta
		}
		q.AllowStale = opts.AllowStale
		q.Connect = opts.Connect
		q.Filter = opts.Filter
		q.UseCache = opts.UseCache
	}
	return q
}

func (c *client) writeOptions(opts *api.WriteOptions) *api.WriteOptions {
	w := api.WriteOptions{}
	if c.config.Datacenter != "" {
		w.Datacenter = c.config.Datacenter
	}
	if c.config.Token != "" {
		w.Token = c.config.Token
	}
	if opts != nil {
		if opts.Token != "" {
			w.Token = opts.Token
		}
	}
	return &w
}
