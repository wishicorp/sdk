package plugin

import (
	"context"
	"errors"
	"fmt"
	"github.com/wishicorp/sdk/helper/errwrap"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/pluginutil"
	"sync"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
)

// BackendPluginClient is a wrapper around backendPluginClient
// that also contains its plugins.Client instance. It's primarily
// used to cleanly kill the client on Cleanup()
type BackendPluginClient struct {
	client *plugin.Client
	sync.Mutex

	logical.Backend
}

// Cleanup calls the RPC client's Cleanup() func and also calls
// the go-plugins's client Kill() func
func (b *BackendPluginClient) Cleanup(ctx context.Context) {
	b.Logger().Trace("Cleanup", "client", "BackendPluginClient")
	b.Backend.Cleanup(ctx)
	b.client.Kill()
}

func NewBackend(ctx context.Context, sys pluginutil.PluginLookupUtil,
	conf *logical.BackendConfig, isMetadataMode bool) (logical.Backend, error) {

	name, ok := conf.Config["plugin_name"]
	if !ok {
		return nil, errors.New("config[plugin_name] not set")
	}

	// Look for plugins in the plugins catalog
	pluginRunner, err := sys.LookupPlugin(ctx, name)
	if err != nil {
		return nil, err
	}

	var backend logical.Backend
	if pluginRunner.Builtin {
		rawFactory, err := pluginRunner.BuiltinFactory()
		if err != nil {
			return nil, errwrap.Wrapf("error getting plugins type: {{err}}", err)
		}
		if factory, ok := rawFactory.(logical.Factory); !ok {
			return nil, fmt.Errorf("unsupported backend type: %q", name)
		} else {
			if backend, err = factory(ctx, conf); err != nil {
				return nil, err
			}
		}
	} else {
		backend, err = NewPluginClient(ctx, pluginRunner, conf.Logger, isMetadataMode)
		if err != nil {
			return nil, err
		}
		if err := backend.Setup(context.Background(), conf); err != nil {
			return nil, err
		}
	}

	return backend, nil
}

func NewPluginClient(ctx context.Context,
	pluginRunner *pluginutil.PluginRunner, logger log.Logger, isMetadataMode bool) (logical.Backend, error) {
	pluginSet := map[int]plugin.PluginSet{
		3: plugin.PluginSet{
			"grpc-backend": &GRPCBackendPlugin{
				MetadataMode: isMetadataMode,
			},
		},
		4: plugin.PluginSet{
			"grpc-backend": &GRPCBackendPlugin{
				MetadataMode: isMetadataMode,
			},
		},
	}

	namedLogger := logger.ResetNamed(pluginRunner.Name).Named("grpc-backend")
	var client *plugin.Client
	var err error
	if isMetadataMode {
		client, err = pluginRunner.RunMetadataMode(ctx, pluginSet, handshakeConfig, []string{}, namedLogger)
	} else {
		client, err = pluginRunner.Run(ctx, pluginSet, handshakeConfig, []string{}, namedLogger)
	}
	if err != nil {
		return nil, err
	}

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}
	// Request the plugins
	raw, err := rpcClient.Dispense("grpc-backend")
	if err != nil {
		return nil, err
	}

	var backend logical.Backend
	var transport string
	// We should have a logical grpc-backend type now. This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	switch raw.(type) {
	case *backendGRPCPluginClient:
		backend = raw.(*backendGRPCPluginClient)
		transport = "gRPC"
	default:
		return nil, errors.New("unsupported plugins client type")
	}

	//// Wrap the grpc-backend in a tracing middleware
	//if namedLogger.IsTrace() {
	//	backend = &backendTracingMiddleware{
	//		logger: namedLogger.With("transport", transport),
	//		next:   backend,
	//	}
	//}
	backend = &backendTracingMiddleware{
		logger: namedLogger.With("transport", transport),
		next:   backend,
	}
	return &BackendPluginClient{
		client:  client,
		Backend: backend,
	}, nil
}
