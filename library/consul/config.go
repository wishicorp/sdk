package consul

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl"
	"gopkg.in/yaml.v2"
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
		c.config.Config.DataKey = "1.0.0"
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
	var succCount = 0
	for _, path := range keyPaths {
		kvp, _, err := c.client.KV().Get(path, options)
		if nil != err {
			return err
		}
		if kvp == nil {
			continue
		}
		if err := c.decodeConfig(path, kvp.Value, out); err != nil {
			return err
		}
		succCount++
	}
	if 0 == succCount{
		return fmt.Errorf("config %v is empty", keyPaths)
	}
	return nil
}

func (c *client)decodeConfig(key string,value []byte, out interface{})(err error)  {
	var format string = strings.ToLower(c.config.Config.Format)
	switch format {
	case "hcl":
		err = hcl.Unmarshal(value, out)
	case "json":
		err = json.Unmarshal(value, out)
	default:
		err = yaml.Unmarshal(value, out)
	}
	if err != nil {
		return fmt.Errorf("unmarshal:%s => format:%s err:%s", key, format, err.Error())
	}
	return nil
}