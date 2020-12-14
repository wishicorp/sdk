package consul

import (
	"github.com/hashicorp/consul/api"
	"testing"
)

func TestClient_KVInfo(t *testing.T) {
	centralConfig := Config{
		Token:       "",
	}
	cli, err := NewClient(&centralConfig)
	if nil != err {
		t.Fatal(err)
	}
	key := "/config/jdsh-pay,dev/1.0.0"
	if kvp, err := cli.KVInfo(key, nil); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(kvp.Value))
	}
}

func TestClient_KVCas(t *testing.T) {
	centralConfig := Config{}
	cli, err := NewClient(&centralConfig)
	if nil != err {
		t.Fatal(err)
	}
	value := "kvp test value"
	kvp := api.KVPair{
		Key:   "session/member-bill",
		Value: []byte(value),
	}
	if newKvp, err := cli.KVCas(&kvp, nil); err != nil {
		t.Fatal(err)
	} else {
		t.Log(newKvp)
	}
}
func TestClient_KVAcquire(t *testing.T) {
	centralConfig := Config{}
	cli, err := NewClient(&centralConfig)
	if nil != err {
		t.Fatal(err)
	}
	key := "session/member-bill"
	session := "53e473ca-5b89-fede-54d8-e1ef08a1de10"
	if kvp, err := cli.KVAcquire(key, session, nil); err != nil {
		t.Fatal(err)
	} else {
		t.Log( kvp)
	}
}

func TestClient_KVRelease(t *testing.T) {
	centralConfig := Config{}
	cli, err := NewClient(&centralConfig)
	if nil != err {
		t.Fatal(err)
	}
	key := "session/member-bill"
	if err := cli.KVRelease(key, nil); err != nil {
		t.Log(err)
	}
}
