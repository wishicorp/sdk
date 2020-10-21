package logical

import (
	"context"
	"github.com/hashicorp/consul/api"
	"github.com/wishicorp/sdk/library/consul"
	"time"
)

type Consul interface {
	Native(ctx context.Context) (*api.Client, error)
	Config(ctx context.Context) (*consul.Config, error)
	//sandbox 表示是否在 namespace 下获取配置
	//配置的路径为 /config/key{.version}
	GetConfig(ctx context.Context, key, version string, sandbox bool) ([]byte, error)

	//获取服务列表
	GetService(ctx context.Context, id, tag string) (*api.AgentService, error)

	//获取微服务路径
	GetServiceAddrPort(ctx context.Context, name string, useLan bool, tags string) (host string, port int, err error)
	//微服务客户端
	GetMicroHTTPClient(ctx context.Context, id string, useLan bool, tags string, header map[string][]string) (consul.MicroHTTPClient, error)

	//创建一个session,ttl需大于15秒,behavior定义了session到期后的动作
	//如需深度定制session请获取native客户端创建
	NewSession(ctx context.Context, name string, ttl time.Duration, behavior consul.SessionBehavior) (string, error)
	SessionInfo(ctx context.Context, id string) (*api.SessionEntry, error)
	//销毁session
	DestroySession(ctx context.Context, id string) error

	//对一个kv进行加锁
	//another标志key是否被其它session锁定
	//err==nil && false == another 加锁成功
	KVAcquire(ctx context.Context, key, session string) (success bool, err error)
	//释放一个session的锁
	KVRelease(ctx context.Context, key string) error
	//获取kv信息
	KVInfo(ctx context.Context, key string) (*api.KVPair, error)
	//检查或者设置key
	KVCas(ctx context.Context, p *api.KVPair) (bool, error)
	KVList(ctx context.Context, prefix string) (api.KVPairs, error)
	KVCreate(ctx context.Context, p *api.KVPair) error
}
