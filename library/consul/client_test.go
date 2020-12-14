package consul

import (
	"github.com/hashicorp/consul/api"
	"os"
	"testing"
)

func TestClient_Client(t *testing.T) {

}

func mockConfig() *Config {
	return &Config{
		ZoneAddress: os.Getenv("CONSUL_ADDR"),
		Token:       os.Getenv("CONSUL_TOKEN"),
		Application: struct {
			Name    string
			Profile string
		}{Name: os.Getenv("APP_NAME"), Profile: os.Getenv("APP_PROFILE")},
		Config: struct {
			DataKey string
			Format  string
		}{DataKey: os.Getenv("CONFIG_DATA_KEY"), Format: os.Getenv("CONFIG_FORMAT")},
		TLSConfig: api.TLSConfig{},
	}
}
