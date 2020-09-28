package logical

import (
	"crypto/tls"
)

// Connection represents the connection information for a request. This
// is present on the Request structure for credential backends.
type Connection struct {
	RemoteAddr string               `json:"remote_addr"`
	ConnState  *tls.ConnectionState `sentinel:""`
}
