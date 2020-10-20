package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/hcl"
	"strings"
)

//配置key格式规范: /config/{APP_NAME},{PROFILE}/{DATA_KEY}
//配置支持继承
// /config/application/{DATA_KEY}
// /config/application,dev/{DATA_KEY}
// /config/application,{PROFILE}/{DATA_KEY}
// /config/{APP_NAME}/{DATA_KEY}
// /config/{APP_NAME},{PROFILE}/{DATA_KEY}
func (c *client) LoadConfig(out interface{}) error {
	if c.config.Config.DataKey == "" {
		c.config.Config.DataKey = "data"
	}
	cfg := c.config.Config
	app := c.config.Application
	keyPaths := []string{
		fmt.Sprintf("/config/application/%s", cfg.DataKey),
		fmt.Sprintf("/config/application,%s/%s", app.Profile, cfg.DataKey),
		fmt.Sprintf("/config/%s/%s", app.Name, cfg.DataKey),
		fmt.Sprintf("/config/%s,%s/%s", app.Name, app.Profile, cfg.DataKey),
	}
	options := c.queryOptions(nil)
	var _errors []string
	for _, path := range keyPaths {
		kvp, _, err := c.client.KV().Get(path, options)
		if nil != err {
			if strings.Contains(err.Error(), "Client.Timeout") {
				return err
			}
			_errors = append(_errors, fmt.Sprintf("%s => %s", path, err.Error()))
			continue
		}
		if nil != kvp {
			if hcl.Unmarshal(kvp.Value, out) != nil {
				return err
			}
		}
	}

	if len(_errors) == len(keyPaths) {
		return errors.New(strings.Join(_errors, ":"))
	}
	return nil
}
