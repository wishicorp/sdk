package consul

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"math/rand"
	"strings"
	"time"
)

const SvcTagSep = "__TAGS__"
const LanAddrKey = "lan_ipv4"
const WanAddrKey = "wan_ipv4"

type Service struct {
	ID             string
	Schema         string
	Name           string
	Address        string
	MatchBody      string
	CheckInterval  string
	Port           int
	Tags           []string
	HealthEndpoint string
	ServiceAddress map[string]api.ServiceAddress
}

func (s *Service) AddTags(tags ...string) {
	s.Tags = append(s.Tags, tags...)
}

func (c *client) GetServices(id, tag string) ([]*api.AgentService, error) {

	ss, _, err := c.client.Health().Service(id, tag, true, c.queryOptions(nil))
	if nil != err {
		return nil, err
	}
	if len(ss) == 0 {
		return nil, errors.New("service not found")
	}

	services := make([]*api.AgentService, 0)
	for e := range ss {
		services = append(services, ss[e].Service)
	}
	return services, nil
}

func (c *client) autoPullService(duration time.Duration) {
	go func() {
		ticker := time.NewTicker(duration)
		for {
			select {
			case <-ticker.C:
				for key, _ := range c.cachedServices {
					sArray := strings.Split(key, SvcTagSep)
					ss, err := c.GetServices(sArray[0], sArray[1])
					if nil == err && len(ss) > 0 {
						c.Lock()
						c.cachedServices[key] = ss
						c.Unlock()
					}
				}
			}
		}
	}()
}

func (c *client) GetService(id, tags string) (*api.AgentService, error) {

	key := id + SvcTagSep + tags

	var err error
	var ss []*api.AgentService

	if c.cachedServices[key] != nil && len(c.cachedServices[key]) > 0 {
		ss = c.cachedServices[key]
	} else {
		ss, err = c.GetServices(id, tags)
		if nil != err {
			return nil, err
		}
		c.Lock()
		c.cachedServices[key] = ss
		c.Unlock()
	}

	return ss[rand.Intn(len(ss))&0xffff], nil
}

func (c *client) GetServiceAddrPort(id string, useLan bool, tags string) (host string, port int, err error) {
	s, err := c.GetService(id, tags)
	if nil != err {
		return "", 0, err
	}
	var addr api.ServiceAddress
	var ok bool
	if useLan {
		addr, ok = s.TaggedAddresses[LanAddrKey]
	} else {
		addr, ok = s.TaggedAddresses[WanAddrKey]
	}

	if !ok {
		return "", 0, errors.New("service not found")
	}

	return addr.Address, addr.Port, nil
}
