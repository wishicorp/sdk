package pluginregister

import (
	"context"
	"errors"
	"github.com/hashicorp/go-hclog"
	uuid "github.com/satori/go.uuid"
	"github.com/wishicorp/sdk/library/consul"
	"github.com/wishicorp/sdk/plugin"
	"github.com/wishicorp/sdk/plugin/logical"
	"sync"
)

var (
	PluginExists    = errors.New("plugins exists")
	PluginNotExists = errors.New("plugins not exists")
)

type PluginManager struct {
	rwMutex   sync.RWMutex
	logger    hclog.Logger
	directory string
	backends  map[string]*Backend
	registry  *PluginRegistry
	consul    consul.Client
	component logical.Component
}

//plugin管理器
func NewPluginManager(directory string, consul consul.Client,
	component logical.Component, logger hclog.Logger) *PluginManager {

	return &PluginManager{
		logger:    logger.Named("plugin-manager"),
		rwMutex:   sync.RWMutex{},
		directory: directory,
		backends:  map[string]*Backend{},
		consul:    consul,
		component: component,
		registry:  &PluginRegistry{directory: directory, entities: map[string]PluginEntry{}},
	}
}

func (p *PluginManager) AutoDiscovery() {
}

func (p *PluginManager) GetBackend(name string) (*Backend, bool) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	backend, has := p.backends[name]
	return backend, has
}

func (p *PluginManager) Deregister(name string) error {
	p.logger.Info("deregister", "name", name)

	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	backend, ok := p.backends[name]
	if !ok {
		return PluginNotExists
	}
	if backend.pluginType != PluginTypeLogical {
		return errors.New("plugin deregister unsupported")
	}

	delete(p.backends, name)
	p.registry.deregister(name)
	p.logger.Info("deregister", "name", name, "used", backend.UsedCount())

	for backend.UsedCount() > 0 {
	}
	backend.Cleanup(context.Background())
	return nil
}

func (p *PluginManager) Register(entry PluginEntry) error {
	p.logger.Info("register", "name", entry.GetName())
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	_, ok := p.backends[entry.GetName()]
	if ok {
		return PluginExists
	}

	if err := p.registry.register(entry); err != nil {
		return err
	}

	if err := p.start(entry); err != nil {
		p.registry.deregister(entry.GetName())
		return err
	}

	return nil
}

func (p *PluginManager) List() []map[string]string {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	var plugins []map[string]string
	for _, p2 := range p.registry.list() {
		plugins = append(plugins, toMapValue(p2))
	}
	return plugins
}

func (p *PluginManager) Trace(name string, on bool) error {
	p.logger.Info("trace", "name", name, "on", on)

	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	backend, ok := p.backends[name]
	if !ok {
		return PluginNotExists
	}
	if on {
		backend.Logger().SetLevel(hclog.Trace)
	} else {
		backend.Logger().SetLevel(hclog.Info)
	}
	return nil
}

func (p *PluginManager) Reload(entry *logicalEntry) error {
	if err := p.Deregister(entry.GetName()); err != nil {
		return err
	}
	return p.Register(entry)
}

func (p *PluginManager) Cleanup() {
	for key, backend := range p.backends {
		backend.Cleanup(context.Background())
		delete(p.backends, key)
	}
}

func (p *PluginManager) start(pluginEntry PluginEntry) error {
	name := pluginEntry.GetName()
	logger := p.logger.Named(name)
	conf := &logical.BackendConfig{
		Logger:        logger,
		BackendUUID:   uuid.NewV1().String(),
		Config:        map[string]string{"plugin_name": pluginEntry.GetName()},
		ConsulView:    logical.NewConsulView(name, "", p.consul),
		ComponentView: p.component,
	}

	var err error
	var backend logical.Backend
	var pluginType PluginType
	switch pluginEntry.(type) {
	case *logicalEntry:
		backend, err = p.newLogicalBackend(conf)
		pluginType = PluginTypeLogical
	case *builtinEntry:
		backend, err = p.newBuiltinBackend(conf)
		pluginType = PluginTypeBuiltin
	}
	if nil != err {
		return err
	}

	initRequest := &logical.InitializationRequest{}
	if err := backend.Initialize(context.Background(), initRequest); err != nil {
		return err
	}

	p.backends[name] = &Backend{Backend: backend, pluginType: pluginType, config: conf}

	return nil
}

func (p *PluginManager) newBuiltinBackend(conf *logical.BackendConfig) (logical.Backend, error) {
	return plugin.NewBackend(context.Background(), p.registry, conf, false)
}

func (p *PluginManager) newLogicalBackend(conf *logical.BackendConfig) (logical.Backend, error) {

	return plugin.NewBackend(context.Background(), p.registry, conf, false)
}

func toMapValue(plugin PluginEntry) map[string]string {
	value := map[string]string{
		"name": plugin.GetName(),
	}
	switch plugin.(type) {
	case *builtinEntry:
		value["type"] = "builtin"
	case *logicalEntry:
		value["type"] = "logical"
	default:
		value["type"] = "unknown"
	}
	return value
}
