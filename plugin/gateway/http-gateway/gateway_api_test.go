package http_gateway

import (
	"strings"
	"testing"
)

func TestGetRemoteAddr(t *testing.T) {
	t.Log(strings.HasPrefix("127.0.0.1:6889", "127.0.0.1"))
}
