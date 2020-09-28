package plugin

import (
	"context"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/proto"
	"sync/atomic"

	"google.golang.org/grpc"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
)

var _ plugin.Plugin = (*GRPCBackendPlugin)(nil)
var _ plugin.GRPCPlugin = (*GRPCBackendPlugin)(nil)

// GRPCBackendPlugin is the plugins.Plugin implementation that only supports GRPC
// transport
type GRPCBackendPlugin struct {
	Factory      logical.Factory
	MetadataMode bool
	Logger       log.Logger
	plugin.NetRPCUnsupportedPlugin
}

func (b GRPCBackendPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterBackendServer(s, &backendGRPCPluginServer{
		broker:  broker,
		factory: b.Factory,
		// We pass the logger down into the grpc-backend so go-plugins will forward
		// logs for us.
		logger: b.Logger,
	})
	return nil
}

func (b *GRPCBackendPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	ret := &backendGRPCPluginClient{
		client:       proto.NewBackendClient(c),
		clientConn:   c,
		broker:       broker,
		cleanupCh:    make(chan struct{}),
		doneCtx:      ctx,
		metadataMode: b.MetadataMode,
	}

	// Create the value and set the type
	ret.server = new(atomic.Value)
	ret.server.Store((*grpc.Server)(nil))

	return ret, nil
}
