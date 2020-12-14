package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
)

func (c *client) KVFire(key string, opts *api.QueryOptions, out interface{}) error {
	kvp, _, err := c.client.KV().Get(key, c.queryOptions(opts))
	if nil != err {
		return err
	}
	if nil == kvp {
		return fmt.Errorf("key[%s] not exists", key)
	}
	return c.decodeConfig(key, kvp.Value, out)
}

func (c *client) KVInfo(key string, opts *api.QueryOptions) (*api.KVPair, error) {
	kvp, _, err := c.client.KV().Get(key, c.queryOptions(opts))
	if nil != err {
		return nil, err
	}
	if nil == kvp {
		return nil, fmt.Errorf("key[%s] not exists", key)
	}
	return kvp, nil
}

func (c *client) KVRelease(key string, opts *api.QueryOptions) error {
	kvp, err := c.KVInfo(key, opts)
	if nil != err {
		return err
	}

	ret, _, err := c.client.KV().Release(kvp, c.writeOptions(&api.WriteOptions{Token: opts.Token}))
	if nil != err {
		return err
	}
	if !ret {
		return errors.New("release session failure")
	}
	return nil
}

func (c *client) KVAcquire(key, session string, opts *api.QueryOptions) (success bool, err error) {

	kvp, err := c.KVInfo(key, opts)
	if nil != err {
		return false, err
	}

	kvp.Session = session
	success, _, err = c.client.KV().Acquire(kvp, c.writeOptions(&api.WriteOptions{Token: opts.Token}))
	return success, err
}
func (c *client) KVCas(p *api.KVPair, opts *api.WriteOptions) (bool, error) {
	ret, _, err := c.client.KV().CAS(p, c.writeOptions(opts))
	return ret, err
}
