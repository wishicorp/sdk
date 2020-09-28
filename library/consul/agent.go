package consul

func (c *client) Maintenance(id string, reason string) error {
	if err := c.client.Agent().EnableServiceMaintenance(id, reason); err != nil {
		return err
	}
	return nil
}

func (c *client) DeMaintenance(id string) error {
	if err := c.client.Agent().DisableServiceMaintenance(id); err != nil {
		return err
	}
	return nil
}
