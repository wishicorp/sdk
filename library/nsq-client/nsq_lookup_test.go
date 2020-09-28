package nsq_client

import (
	"github.com/nsqio/go-nsq"
	"testing"
)

func TestLookup(t *testing.T) {
	lookup := NewLookup(nsq.NewConfig(), []string{"http://127.0.0.1:4161/nodes"}, "logging")
	address := lookup.LookupNsqdAddress()
	t.Log(address)
}
