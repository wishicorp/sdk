package consul

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	centralConfig := Config{}
	cli, err := NewClient(&centralConfig)
	if nil != err {
		t.Fatal(err)
	}
	if session, err := cli.NewSession("bill-session", time.Hour, BehaviorRelease, nil); err != nil {
		t.Fatal(err)
	} else {
		t.Log(session)
	}
}
