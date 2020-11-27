package pluginutil

import (
	"context"
	"crypto/sha256"
	"fmt"
	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/wishicorp/sdk/version"
	"os"
	"os/exec"
	"path/filepath"
)

type BuiltinFactory func() (interface{}, error)

type PluginLookupUtil interface {
	LookupPlugin(context.Context, string) (*PluginRunner, error)
}

type PluginRunner struct {
	Name           string         `json:"name" structs:"name"`
	Command        string         `json:"command" structs:"command"`
	Args           []string       `json:"args" structs:"args"`
	Env            []string       `json:"env" structs:"env"`
	Sha256         []byte         `json:"sha256" structs:"sha256"`
	Builtin        bool           `json:"builtin" structs:"builtin"`
	BuiltinFactory BuiltinFactory `json:"-" structs:"-"`
}

// Run takes a wrapper RunnerUtil instance along with the go-plugins parameters and
// returns a configured plugins.Client with TLS Configured and a wrapping token set
// on PluginUnwrapTokenEnv for plugins process consumption.
func (r *PluginRunner) Run(ctx context.Context, pluginSets map[int]plugin.PluginSet, hs plugin.HandshakeConfig, env []string, logger log.Logger) (*plugin.Client, error) {
	return r.runCommon(ctx, pluginSets, hs, env, logger, false)
}

// RunMetadataMode returns a configured plugins.Client that will dispense a plugins
// in metadata mode. The PluginMetadataModeEnv is passed in as part of the Cmd to
func (r *PluginRunner) RunMetadataMode(ctx context.Context, pluginSets map[int]plugin.PluginSet, hs plugin.HandshakeConfig, env []string, logger log.Logger) (*plugin.Client, error) {
	return r.runCommon(ctx, pluginSets, hs, env, logger, true)

}

func (r *PluginRunner) runCommon(ctx context.Context, pluginSets map[int]plugin.PluginSet, hs plugin.HandshakeConfig, env []string, logger log.Logger, isMetadataMode bool) (*plugin.Client, error) {
	cmd := exec.Command(r.Command, r.Args...)

	// `env` should always go last to avoid overwriting internal values that might
	// have been provided externally.
	cmd.Env = append(cmd.Env, r.Env...)
	cmd.Env = append(cmd.Env, env...)

	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginVersionEnv, version.GetVersion().Version))

	if !isMetadataMode {
		// Add the metadata mode ENV and set it to false
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMetadataModeEnv, "false"))

	} else {
		logger = logger.With("metadata", "true")
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMetadataModeEnv, "true"))
	}

	secureConfig := &plugin.SecureConfig{
		Checksum: r.Sha256,
		Hash:     sha256.New(),
	}
	clientConfig := &plugin.ClientConfig{
		HandshakeConfig:  hs,
		VersionedPlugins: pluginSets,
		Cmd:              cmd,
		SecureConfig:     secureConfig,
		Logger:           logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
	}

	client := plugin.NewClient(clientConfig)

	return client, nil
}

// CtxCancelIfCanceled takes a context cancel func and a context. If the context is
// shutdown the cancelfunc is called. This is useful for merging two cancel
// functions.
func CtxCancelIfCanceled(f context.CancelFunc, ctxCanceler context.Context) chan struct{} {
	quitCh := make(chan struct{})
	go func() {
		select {
		case <-quitCh:
		case <-ctxCanceler.Done():
			f()
		}
	}()
	return quitCh
}

func LookBackendExecute(directory, name string) bool {
	pName := fmt.Sprintf("%s/backend.plugin", name)
	path := filepath.Join(directory, pName)
	exists, _ := PathExists(path)
	return exists
}

func PathExists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
