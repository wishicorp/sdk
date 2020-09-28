package pluginutil

import (
	"os"

	version "github.com/hashicorp/go-version"
)

var (
	// PluginMlockEnabled is the ENV name used to pass the configuration for
	// enabling mlock
	PluginMlockEnabled = "PLUGIN_MLOCK_ENABLED"

	// PluginVersionEnv is the ENV name used to pass the version of the
	// PLUGIN server to the plugins
	PluginVersionEnv = "PLUGIN_VERSION"

	// PluginMetadataModeEnv is an ENV name used to disable TLS communication
	// to boots mounting plugins.
	PluginMetadataModeEnv = "PLUGIN_METADATA_MODE"

	// PluginUnwrapTokenEnv is the ENV name used to pass unwrap tokens to the
	// plugins.
	PluginUnwrapTokenEnv = "PLUGIN_UNWRAP_TOKEN"

	// PluginCACertPEMEnv is an ENV name used for holding a CA PEM-encoded
	// string. Used for testing.
	PluginCACertPEMEnv = "PLUGIN_TESTING_PLUGIN_CA_PEM"
)

// OptionallyEnableMlock determines if mlock should be called, and if so enables
// mlock.
func OptionallyEnableMlock() error {
	//if os.Getenv(PluginMlockEnabled) == "true" {
	//	return mlock.LockMemory()
	//}

	return nil
}

// GRPCSupport defaults to returning true, unless PLUGIN_VERSION is missing or
// it fails to meet the version constraint.
func GRPCSupport() bool {
	verString := os.Getenv(PluginVersionEnv)
	// If the env var is empty, we fall back to netrpc for backward compatibility.
	if verString == "" {
		return false
	}
	if verString != "unknown" {
		ver, err := version.NewVersion(verString)
		if err != nil {
			return true
		}
		// Due to some regressions on 0.9.2 & 0.9.3 we now require version 0.9.4
		// to allow the plugins framework to default to gRPC.
		constraint, err := version.NewConstraint(">= 0.9.4")
		if err != nil {
			return true
		}
		return constraint.Check(ver)
	}
	return true
}

// InMetadataMode returns true if the plugins calling this function is running in metadata mode.
func InMetadataMode() bool {
	return os.Getenv(PluginMetadataModeEnv) == "true"
}
