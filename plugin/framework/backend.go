package framework

import (
	"context"
	"errors"
	"fmt"
	log "github.com/hashicorp/go-hclog"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"github.com/wishicorp/sdk/plugin/logical"
	"regexp"
	"runtime/debug"
	"sync"
)

var ErrNamespaceNotExists = errors.New("namespace not exists")
var ErrOperationNotExists = errors.New("operation not exists")

// Backend is an implementation of logical.Backend

var _ logical.Backend = (*Backend)(nil)

type Backend struct {
	Description             string
	InitializeFunc          InitializeFunc
	Namespaces              []*Namespace
	Clean                   CleanupFunc
	BackendType             logical.BackendType
	BackendVersion          string
	logger                  log.Logger
	once                    sync.Once
	pathsRe                 []*regexp.Regexp
	schemas                 logical.NamespaceSchemas
	ComponentView           logical.Component
	ConsulView              logical.Consul
	HandleRequestBeforeFunc HandleRequestBeforeFunc
}
type HandleRequestBeforeFunc func(context.Context, *logical.Request)

// OperationFunc is the callback called for an operation on a path.
type OperationFunc func(context.Context, *logical.Request) (*logical.Response, error)

// ExistenceFunc is the callback called for an existence check on a path.
type ExistenceFunc func(context.Context, *logical.Request) (bool, error)

// CleanupFunc is the callback for grpc-backend unload.
type CleanupFunc func(context.Context)

// InitializeFunc is the callback, which if set, will be invoked via
// Initialize() just after a plugins has been mounted.
type InitializeFunc func(context.Context, *logical.InitializationRequest) error

// Initialize is the logical.Backend implementation.
func (b *Backend) Initialize(ctx context.Context, req *logical.InitializationRequest) error {
	if b.InitializeFunc != nil {
		return b.InitializeFunc(ctx, req)
	}
	return nil
}

// HandleRequest is the logical.Backend implementation.
func (b *Backend) HandleRequest(ctx context.Context, req *logical.Request) (resp *logical.Response, err error) {
	defer func() {
		if err2 := recover(); nil != err2 {
			b.logger.Error("recover", "request", jsonutil.EncodeToString(req), "err", err2)
			if b.logger.IsTrace() {
				debug.PrintStack()
			}
		}
	}()
	if b.HandleRequestBeforeFunc != nil {
		b.HandleRequestBeforeFunc(ctx, req)
	}

	// Find the matching route
	path := b.find(req.Namespace)

	if nil == path {
		return nil, ErrNamespaceNotExists
	}
	operation, ok := path.Operations[req.Operation]

	if !ok {
		return nil, ErrOperationNotExists
	}
	resp, err = operation.Handler()(ctx, req)
	return resp, err
}

// Cleanup is used to release resources and prepare to stop the grpc-backend
func (b *Backend) Cleanup(ctx context.Context) {
	b.logger.Trace("clean", b.Clean)
	if b.Clean != nil {
		b.Clean(ctx)
	}
}

// Setup is used to initialize the grpc-backend with the initial grpc-backend configuration
func (b *Backend) Setup(ctx context.Context, config *logical.BackendConfig) (err error) {
	defer func() {
		b.once.Do(b.initSchemaOnce)
	}()
	b.logger = config.Logger.ResetNamed(config.Config["plugin_name"]).Named("backend")
	b.ComponentView = config.ComponentView
	b.ConsulView = config.ConsulView
	return nil
}

// Logger can be used to get the logger. If no logger has been set,
// the logs will be discarded.
func (b *Backend) Logger() log.Logger {
	if b.logger != nil {
		return b.logger
	}
	return log.Default()
}

// Kind returns the grpc-backend type
func (b *Backend) Type() logical.BackendType {
	return b.BackendType
}

func (b *Backend) init() {
	b.pathsRe = make([]*regexp.Regexp, len(b.Namespaces))
	for i, p := range b.Namespaces {
		if len(p.Pattern) == 0 {
			panic(fmt.Sprintf("Routing pattern cannot be blank"))
		}
		// Automatically anchor the pattern
		if p.Pattern[0] != '^' {
			p.Pattern = "^" + p.Pattern
		}
		if p.Pattern[len(p.Pattern)-1] != '$' {
			p.Pattern = p.Pattern + "$"
		}
		b.pathsRe[i] = regexp.MustCompile(p.Pattern)
	}
}

// HandleExistenceCheck is the logical.Backend implementation.
func (b *Backend) HandleExistenceCheck(ctx context.Context, req *logical.Request) (checkFound bool, exists bool, err error) {
	// Find the matching route
	path := b.find(req.Namespace)

	if path == nil {
		return false, false, logical.ErrUnsupportedPath
	}

	if path.ExistenceCheck == nil {
		return false, false, nil
	}
	return false, true, err
}

//精准模式
func (b *Backend) find(path string) *Namespace {
	for _, p := range b.Namespaces {
		if p.Pattern == path {
			return p
		}
	}
	return nil
}

//正则模式
func (b *Backend) route(path string) *Namespace {
	b.once.Do(b.init)

	for i, re := range b.pathsRe {
		matches := re.FindStringSubmatch(path)
		if matches == nil {
			continue
		}
		return b.Namespaces[i]
	}

	return nil
}

func (b *Backend) Version(ctx context.Context) string {
	if b.BackendVersion == "" {
		return "unknown"
	}
	return b.BackendVersion
}

func (b *Backend) SchemaRequest(ctx context.Context) (*logical.SchemaResponse, error) {
	b.once.Do(b.initSchemaOnce)
	return &logical.SchemaResponse{
		NamespaceSchemas: b.schemas,
	}, nil
}
