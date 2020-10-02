package compose

import (
	"context"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/hashicorp/go-hclog"
	"github.com/wishicorp/sdk/plugin/logical"
	"sync"
)

type Operator struct {
	Name         string
	Method       string
	Input        []*logical.Field
	Output       []*logical.Field
	FaultOnMiss  bool
	FaultOnError bool
}

type Compose struct {
	sync.Mutex
	Name          string
	Authorization []byte
	Token         string
	logger        hclog.Logger
	*linkedhashmap.Map
	Input  map[string][]*logical.Field
	Output map[string][]*logical.Field
}
type Config struct {
}

func NewCompose() *Compose {
	c := &Compose{}
	c.Map = linkedhashmap.New()
	c.Input = map[string][]*logical.Field{}
	c.Output = map[string][]*logical.Field{}
	return c
}
func (m *Compose) Setup(ctx context.Context, cfg *Config) error {
	return nil
}

func (m *Compose) Initialize(ctx context.Context, args *InitializationRequest) error {
	return nil
}

func (m *Compose) HandleRequest(context.Context, *Request) (*Response, error) {
	return nil, nil
}

func (m *Compose) Cleanup(context.Context) {
}

func (m *Compose) Logger(context.Context) hclog.Logger {
	return m.logger
}

func (m *Compose) Put(op *Operator) {
	m.Lock()
	defer m.Unlock()

	m.Map.Put(op.Name, op)
	m.Input[op.Name] = op.Input
	m.Output[op.Name] = op.Output
}
