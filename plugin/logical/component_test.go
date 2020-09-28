package logical

import (
	"github.com/hashicorp/hcl"
	"github.com/wishicorp/sdk/framework/worm"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"testing"
)

var body = `
mysql "primary" {
  master = "root:123456@tcp(127.0.0.1)/virtual_coin?charset=utf8mb4&parseTime=true&loc=Local"
  use_master_slave = false
  show_sql = true
  slaves = [
    "root:123456@tcp(127.0.0.1)/virtual_coin?charset=utf8mb4&parseTime=true&loc=Local",
    "root:123456@tcp(127.0.0.1)/virtual_coin?charset=utf8mb4&parseTime=true&loc=Local"
    ]
}
mysql "secondary" {
  master = "root:123456@tcp(127.0.0.1)/virtual_coin?charset=utf8mb4&parseTime=true&loc=Local"
  use_master_slave = false
  show_sql = true
}
redis "primary" {
  addrs = ["127.0.0.1:6379"]
  use_cluster = false
}
rabbitmq "primary" {
  url = "amqp://guest:guest@localhost:5672/"
  auto_delete = false
  auto_ack = false
  no_wait = true
  exclusive = true
}
service {
    key = "value"
}

service {
    key = "foo"
}
`

type service struct {
	Key string
}
type Cfg struct {
	Mysql    map[string]worm.Config
	Services []service `hcl:"service"`
}

func TestComponent_FetchConfig(t *testing.T) {
	var c Cfg
	if err := hcl.Unmarshal([]byte(body), &c); err != nil {
		t.Fatal(err)
		return
	}
	t.Log(jsonutil.EncodeJSON(c.Mysql["primary"].Slaves[0]))
}
