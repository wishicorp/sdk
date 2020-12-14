package consul

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	cli, err := NewClient(mockConfig())
	if nil != err {
		t.Fatal(err)
	}
	if session, err := cli.NewSession("bill-session", time.Hour, BehaviorRelease, nil); err != nil {
		t.Fatal(err)
	} else {
		t.Log(session)
	}
}
