package redis

import (
	"crypto/tls"
	"errors"
	"gopkg.in/redis.v5"
	"time"
)

var ErrRedisAddrsEmpty = errors.New("redis addrs is empty")

type Config struct {
	Addrs        []string      `hcl:"addrs"`      //集群地址
	Password     string        `hcl:"password"`   //密码
	KeyPrefix    string        `hcl:"key_prefix"` //key前缀
	DbIndex      int           `hcl:"db_index"`   //数据索引
	DialTimeout  time.Duration `hcl:"dial_timeout"`
	ReadTimeout  time.Duration `hcl:"read_timeout"`
	WriteTimeout time.Duration `hcl:"write_timeout"`
	ReadOnly     bool          `hcl:"read_only"`
	// PoolSize applies per cluster node and not for the whole cluster.
	PoolSize           int           `hcl:"pool_size"`
	PoolTimeout        time.Duration `hcl:"pool_timeout"`
	IdleTimeout        time.Duration `hcl:"idle_timeout"`
	IdleCheckFrequency time.Duration `hcl:"idle_check_frequency"`
	UseCluster         bool          `hcl:"use_cluster"`
	TLSConfig          *tls.Config   `hcl:"tls_config"`
}

type RedisCmd interface {
	redis.Cmdable
}

func NewRedisCmd(c *Config) (RedisCmd, error) {
	if len(c.Addrs) == 0 {
		return nil, ErrRedisAddrsEmpty
	}

	if c.UseCluster {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:              c.Addrs,
			ReadOnly:           c.ReadOnly,
			Password:           c.Password,
			DialTimeout:        c.DialTimeout,
			ReadTimeout:        c.ReadTimeout,
			WriteTimeout:       c.WriteTimeout,
			PoolSize:           c.PoolSize,
			PoolTimeout:        c.PoolTimeout,
			IdleTimeout:        c.IdleTimeout,
			IdleCheckFrequency: c.IdleCheckFrequency,
		}), nil
	}
	return redis.NewClient(&redis.Options{
		Addr:               c.Addrs[0],
		Password:           c.Password,
		DB:                 c.DbIndex,
		DialTimeout:        c.DialTimeout,
		ReadTimeout:        c.ReadTimeout,
		WriteTimeout:       c.WriteTimeout,
		PoolSize:           c.PoolSize,
		PoolTimeout:        c.PoolTimeout,
		IdleTimeout:        c.IdleTimeout,
		IdleCheckFrequency: c.IdleCheckFrequency,
		ReadOnly:           c.ReadOnly,
		TLSConfig:          c.TLSConfig,
	}), nil
}
