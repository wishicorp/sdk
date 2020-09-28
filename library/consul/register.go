package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

func (c *client) Register(s *Service) error {
	c.Lock()
	defer c.Unlock()

	url := fmt.Sprintf("%s://%s:%d%s", s.Schema, s.Address, s.Port, s.HealthEndpoint)
	check := &api.AgentServiceCheck{
		CheckID:                        s.ID,
		HTTP:                           url,
		Interval:                       s.CheckInterval,
		Body:                           s.MatchBody,
		TLSSkipVerify:                  true,
		DeregisterCriticalServiceAfter: "5m",
	}

	a := &api.AgentServiceRegistration{
		ID:              s.ID,
		Name:            s.Name,
		Address:         s.Address,
		Port:            s.Port,
		TaggedAddresses: s.ServiceAddress,
		Tags:            append(s.Tags, c.config.Application.Profile),
		Check:           check,
	}
	err := c.client.Agent().ServiceRegister(a)
	return err
}
func (c *client) DeRegister(s *Service) error {
	c.Lock()
	defer c.Unlock()

	err := c.client.Agent().ServiceDeregister(s.ID)
	return err
}
