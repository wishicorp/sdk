package compose

import "testing"

func TestNewCompose(t *testing.T) {
	c := NewCompose()
	iter := c.Iterator()
	c.Put(&Operator{Name: "account"})
	t.Log(iter.Next())
	t.Log(iter.Value())
}
