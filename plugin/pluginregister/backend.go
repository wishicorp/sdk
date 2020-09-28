package pluginregister

import (
	"github.com/wishicorp/sdk/plugin/logical"
	"sync/atomic"
)

type Backend struct {
	logical.Backend
	usedCount  uint64
	config     *logical.BackendConfig
	pluginType PluginType
}

func (p *Backend) UsedCount() uint64 {
	return atomic.LoadUint64(&p.usedCount)
}

func (p *Backend) Incr() {
	atomic.AddUint64(&p.usedCount, 1)
}

func (p *Backend) DeIncr() {
	value := atomic.LoadUint64(&p.usedCount)
	if value < 1 {
		return
	}
	atomic.StoreUint64(&p.usedCount, value-1)
}
