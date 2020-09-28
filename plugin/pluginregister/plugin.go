package pluginregister

import (
	"github.com/wishicorp/sdk/plugin/logical"
	"sync"
)

type PluginType string

const (
	PluginTypeUnknown PluginType = "unknown"
	PluginTypeBuiltin PluginType = "builtin"
	PluginTypeLogical PluginType = "logical"
)

type PluginEntry interface {
	GetName() string
}

type pluginEntry struct {
	sync.RWMutex
	name string
}

type builtinEntry struct {
	*pluginEntry
	Factory logical.Factory
}

func NewBuiltinEntry(name string, factory logical.Factory) *builtinEntry {
	return &builtinEntry{
		pluginEntry: &pluginEntry{
			name: name,
		},
		Factory: factory,
	}
}

type logicalEntry struct {
	version string
	*pluginEntry
	Args []string
}

func NewLogicalEntry(name string, args []string) *logicalEntry {
	return &logicalEntry{
		pluginEntry: &pluginEntry{
			name: name,
		},
		Args: args,
	}
}

func (p *pluginEntry) GetName() string {
	return p.name
}
