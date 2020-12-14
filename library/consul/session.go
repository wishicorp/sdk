package consul

import (
	"github.com/hashicorp/consul/api"
	"time"
)

type SessionBehavior string

const (
	BehaviorDelete  SessionBehavior = "delete"
	BehaviorRelease SessionBehavior = "release"
)

func (c *client) NewSession(name string, ttl time.Duration, behavior SessionBehavior, opts *api.WriteOptions) (string, error) {
	id, _, err := c.client.Session().Create(&api.SessionEntry{
		Name: name, TTL: ttl.String(), Behavior: string(behavior),
	}, c.writeOptions(opts))
	if nil != err {
		return "", err
	}
	return id, nil
}

func (c *client) SessionInfo(id string, opts *api.QueryOptions) (*api.SessionEntry, error) {
	entry, _, err := c.client.Session().Info(id, c.queryOptions(opts))
	if nil != err {
		return nil, err
	}
	return entry, nil
}

func (c *client) DestroySession(id string, opts *api.WriteOptions) error {
	_, err := c.client.Session().Destroy(id, c.writeOptions(opts))
	return err
}

func (c *client) CreateLocker(key string, opts *api.WriteOptions) (*api.Lock, error) {
	lock, err := c.client.LockKey(key)
	return lock, err
}
