package plugin

import (
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/pluginutil"
	"math"
	"os"

	"google.golang.org/grpc"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
)

// BackendPluginName is the name of the plugins that can be
// dispensed from the plugins server.
const BackendPluginName = "grpc-backend"

type ServeOpts struct {
	Factory logical.Factory
	Logger  log.Logger
}

// Serve is a helper function used to serve a grpc-backend plugins. This
// should be ran on the plugins's main process.
func Serve(opts *ServeOpts) error {
	logger := opts.Logger
	if logger == nil {
		logger = log.New(&log.LoggerOptions{
			Level:      log.Trace,
			Output:     os.Stderr,
			JSONFormat: true,
		})
	}

	// pluginMap is the map of plugins we can dispense.
	pluginSets := map[int]plugin.PluginSet{
		// Version 3 used to supports both protocols. We want to keep it around
		// since it's possible old plugins built against this version will still
		// work with gRPC. There is currently no difference between version 3
		// and version 4.
		3: plugin.PluginSet{
			"grpc-backend": &GRPCBackendPlugin{
				Factory:      opts.Factory,
				Logger:       logger,
				MetadataMode: false,
			},
		},
		4: plugin.PluginSet{
			"grpc-backend": &GRPCBackendPlugin{
				Factory:      opts.Factory,
				Logger:       logger,
				MetadataMode: false,
			},
		},
	}

	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		return err
	}

	serveOpts := &plugin.ServeConfig{
		HandshakeConfig:  handshakeConfig,
		VersionedPlugins: pluginSets,
		Logger:           logger,

		// A non-nil value here enables gRPC serving for this plugins...
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			opts = append(opts, grpc.MaxRecvMsgSize(math.MaxInt32))
			opts = append(opts, grpc.MaxSendMsgSize(math.MaxInt32))
			return plugin.DefaultGRPCServer(opts)
		},
	}

	plugin.Serve(serveOpts)

	return nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugins and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugins
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  4,
	MagicCookieKey:   "BACKEND_PLUGIN",
	MagicCookieValue: "8269d33a-b1c8-4f49-97d9-c8e5bed98e20",
}
