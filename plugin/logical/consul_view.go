package logical

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/wishicorp/sdk/library/consul"
	"time"
)

var _ Consul = (*ConsulView)(nil)

const ConfigPath = "config"
const SessionPath = "session"

type ConsulView struct {
	consul    consul.Client
	path      string
	namespace string
	profile   string
}

func NewConsulView(namespace, profile string, consul consul.Client) *ConsulView {
	return &ConsulView{
		consul:    consul,
		path:      ConfigPath,
		namespace: namespace,
		profile:   profile,
	}
}
func (c *ConsulView) Config(ctx context.Context) (*consul.Config, error) {
	return c.consul.Config(), nil
}

func (c *ConsulView) Native(ctx context.Context) (*api.Client, error) {
	return c.consul.Client(), nil
}

func (c *ConsulView) GetServiceAddrPort(ctx context.Context, name string, useLan bool, tags string) (host string, port int, err error) {
	return c.consul.GetServiceAddrPort(name, useLan, tags)
}

func (c *ConsulView) GetMicroHTTPClient(ctx context.Context, id string, useLan bool, tags string, header map[string][]string) (consul.MicroHTTPClient, error) {
	return c.consul.GetMicroHTTPClient(id, useLan, tags, header)
}

func (c *ConsulView) GetConfig(ctx context.Context, key, version string, sandbox bool) ([]byte, error) {
	var path string
	if !sandbox {
		path = fmt.Sprintf("/%s/%s", c.path, key)
	} else {
		path = c.expendKey(key)
	}
	if version != "" {
		path = path + "/" + version
	}
	if kvp, err := c.consul.KVInfo(path, &api.QueryOptions{}); err != nil {
		return nil, err
	} else {
		return kvp.Value, nil
	}
}

func (c *ConsulView) GetService(ctx context.Context, id, tag string) (*api.AgentService, error) {
	return c.consul.GetService(id, tag)
}

func (c *ConsulView) NewSession(ctx context.Context,
	name string, ttl time.Duration, behavior consul.SessionBehavior) (string, error) {
	return c.consul.NewSession(name, ttl, behavior, &api.WriteOptions{})
}

func (c *ConsulView) SessionInfo(ctx context.Context, id string) (*api.SessionEntry, error) {
	return c.consul.SessionInfo(id, &api.QueryOptions{})
}

func (c *ConsulView) DestroySession(ctx context.Context, id string) error {
	return c.consul.DestroySession(id, &api.WriteOptions{})
}

func (c *ConsulView) KVAcquire(ctx context.Context, key, session string) (success bool, err error) {
	return c.consul.KVAcquire(c.ExpendSessionKey(key), session, &api.QueryOptions{})
}

func (c *ConsulView) KVRelease(ctx context.Context, key string) error {
	return c.consul.KVRelease(c.ExpendSessionKey(key), &api.QueryOptions{})
}

func (c *ConsulView) KVInfo(ctx context.Context, key string) (*api.KVPair, error) {
	return c.consul.KVInfo(c.ExpendSessionKey(key), &api.QueryOptions{})
}

func (c *ConsulView) KVCas(ctx context.Context, p *api.KVPair) (bool, error) {
	return c.consul.KVCas(p, &api.WriteOptions{})
}

func (c *ConsulView) expendKey(key string) string {
	if c.profile != "" {
		return fmt.Sprintf("/%s/%s,%s/%s", c.path, c.namespace, c.profile, key)
	}
	return fmt.Sprintf("/%s/%s/%s", c.path, c.namespace, key)
}

func (c *ConsulView) KVList(ctx context.Context, prefix string) (api.KVPairs, error) {
	kvps, _, err := c.consul.Client().KV().List(prefix, &api.QueryOptions{})
	return kvps, err
}

func (c *ConsulView) KVCreate(ctx context.Context, p *api.KVPair) error {
	_, err := c.consul.Client().KV().Put(p, &api.WriteOptions{})
	return err
}

func (c *ConsulView) ExpendSessionKey(key string) string {
	if c.profile != "" {
		return fmt.Sprintf("/%s/%s,%s/%s", SessionPath, c.namespace, c.profile, key)
	}
	return fmt.Sprintf("/%s/%s/%s", SessionPath, c.namespace, key)
}
