package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/wishicorp/sdk/helper/threadutil"
	"strconv"
	"strings"
	"time"
)

const retryKeyPrefix = "retry"

type PutNext bool
type NotificationHandler func(p *Payload, err error) PutNext

func DefaultRetryPolicies() []time.Duration {
	return []time.Duration{
		time.Minute * 1,
		time.Minute * 5,
		time.Minute * 10,
		time.Minute * 30,
		time.Minute * 60,
		time.Minute * 120,
	}
}

type PayloadDecoder interface {
	Decode(payload *Payload) string
}
type PayloadEncoder interface {
	Encode(date string) *Payload
}
type Payload struct {
	count int64
	Value string
	cache RedisCli
}

type Notification interface {
	PutNotification(p *Payload)
	Subscribe(handler NotificationHandler) error
}

type notification struct {
	key      string
	node     string
	cache    RedisCli
	policies []time.Duration
}

func (p *Payload) String() string {
	return encodePayload(p)
}

func NewNotification(node, key string, cache RedisCli, policies []time.Duration) Notification {
	if "" == node || "" == key || nil == cache {
		return nil
	}
	n := &notification{
		key:      key,
		node:     node,
		cache:    cache,
		policies: policies,
	}
	if nil == policies {
		n.policies = DefaultRetryPolicies()
	}
	return n
}

func (n *notification) PutNotification(p *Payload) {
	setKey := fmt.Sprintf("%s-%s:%s", retryKeyPrefix, n.key, p.String())
	_ = n.cache.Set(context.Background(), setKey, []byte(p.Value), n.policies[p.count].String())
}

func (n *notification) lock(p *Payload) bool {
	setKey := fmt.Sprintf("%s-%s:%s", retryKeyPrefix, n.key, p.Value)
	value := fmt.Sprintf("%s-%d", n.node, threadutil.GetRoutineID())
	ret, _ := n.cache.SetNX(context.Background(), setKey, []byte(value), "0")
	if ret {
		return ret
	}
	rval, _ := n.cache.Get(context.Background(), setKey)
	return value == string(rval)
}

func (n *notification) unlock(p *Payload) {
	setKey := fmt.Sprintf("%s-%s:%s", retryKeyPrefix, n.key, p.Value)
	p.cache.Del(context.Background(), setKey)
}

func (n *notification) Subscribe(handler NotificationHandler) error {
	skey := fmt.Sprintf("%s-%s", retryKeyPrefix, n.key)
	space := fmt.Sprintf("__key*__:%s:%s*", n.cache.KeyPrefix(), skey)
	psub, err := n.cache.PSubscribe(space)
	if nil != err {
		return err
	}
	go func() {
		for {
			message, err := psub.ReceiveMessage()
			if nil != err {
				continue
			}
			if message.Payload == "expired" {

				payload, err := decodePayload(message.Channel)
				if nil != err {
					continue
				}
				payload.count += 1
				payload.cache = n.cache

				if !n.lock(payload) {
					continue
				}

				putNext := handler(payload, err)

				if putNext && payload.count < int64(len(n.policies)) {
					n.PutNotification(payload)
				}

				if !putNext || (putNext && payload.count >= int64(len(n.policies))) {
					n.unlock(payload)
				}

			}
		}
	}()
	return nil
}
func decodePayload(s string) (*Payload, error) {
	ss := strings.Split(s, ":")
	if len(ss) < 2 {
		return nil, errors.New("illegal payload: " + s)
	}
	ss = ss[1:]
	retryS := ss[len(ss)-1]
	rss := strings.Split(retryS, "||")
	if len(rss) < 2 {
		return nil, errors.New("illegal payload: " + s)
	}

	payload := Payload{}
	payload.Value = rss[0]
	payload.count, _ = strconv.ParseInt(rss[1], 10, 32)
	return &payload, nil
}

func encodePayload(p *Payload) string {
	return fmt.Sprintf("%s||%d", p.Value, p.count)
}
