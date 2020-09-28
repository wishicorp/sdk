package pluginregister

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/pluginutil"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type PluginRegistry struct {
	mutex     sync.Mutex
	directory string
	entities  map[string]PluginEntry
}

func (p *PluginRegistry) register(entry PluginEntry) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	_, ok := p.entities[entry.GetName()]
	if ok {
		return PluginExists
	}
	p.entities[entry.GetName()] = entry
	return nil
}

func (p *PluginRegistry) deregister(name string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	delete(p.entities, name)
}

func (p *PluginRegistry) LookupPlugin(ctx context.Context, name string) (*pluginutil.PluginRunner, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	directory := p.directory

	entry, ok := p.entities[name]
	if !ok {
		return nil, PluginNotExists
	}

	switch value := entry.(type) {
	case *builtinEntry:
		return &pluginutil.PluginRunner{
			Name:           name,
			Builtin:        true,
			BuiltinFactory: toFunc(value.Factory),
		}, nil
	case *logicalEntry:
		pName := fmt.Sprintf("%s/%s", value.GetName(), logical.BackendExecuteName)
		path := filepath.Join(directory, pName)
		if !pluginutil.LookBackendExecute(directory, name) {
			return nil, fmt.Errorf("execute file not exists[%s]", pName)
		}

		return &pluginutil.PluginRunner{
			Name:    name,
			Args:    value.Args,
			Command: path,
			Sha256:  checksum(path),
		}, nil
	default:
		return nil, PluginNotExists
	}
}

func (p *PluginRegistry) list() []PluginEntry {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var entities []PluginEntry
	for _, entry := range p.entities {
		entities = append(entities, entry)
	}
	return entities
}

func toFunc(ifc interface{}) func() (interface{}, error) {
	return func() (interface{}, error) {
		return ifc, nil
	}
}

func checksum(filePath string) []byte {
	hash := sha256.New()
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	_, err = io.Copy(hash, file)
	if err != nil {
		return nil
	}

	sum := hash.Sum(nil)
	return sum
}
