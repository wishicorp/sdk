package consul

import "testing"

type exampleConfig struct {
	App interface{} `hcl:"mysql" yaml:"mysql"`
}

func TestClient_LoadConfig(t *testing.T) {
	centralConfig := Config{
		Application: struct {
			Name    string
			Profile string
		}{Name: "jdsh-pay", Profile: "dev"},
		Token: "57c5d69a-5f19-469b-0543-12a487eecc66",
		Config: struct {
			DataKey string
			Format  string
		}{DataKey: "1.0.0", Format: "hcl"},
	}
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
