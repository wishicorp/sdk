package consul

import "testing"

type exampleConfig struct {
	App string `hcl:"app"`
}

func TestClient_LoadConfig(t *testing.T) {
	centralConfig := Config{}
	cli, err := NewClient(&centralConfig)
	if nil != err {
		t.Fatal(err)
	}

	//载入配置
	var cfg exampleConfig
	if err := cli.LoadConfig(&cfg); err != nil {
		t.Fatal(err)
	}
	t.Log(cfg)
}
