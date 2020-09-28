package plugin

import (
	"context"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"github.com/wishicorp/sdk/plugin/logical"
	"time"

	log "github.com/hashicorp/go-hclog"
)

// backendPluginClient implements logical.Backend and is the
// go-plugins client.
type backendTracingMiddleware struct {
	logger log.Logger

	next logical.Backend
}

func (b *backendTracingMiddleware) Version(ctx context.Context) string {
	return b.next.Version(ctx)
}

// Validate the backendTracingMiddle object satisfies the grpc-backend interface
var _ logical.Backend = (*backendTracingMiddleware)(nil)

func (b *backendTracingMiddleware) SchemaRequest(ctx context.Context) (*logical.SchemaResponse, error) {
	defer func(then time.Time) {
		b.logger.Trace("SchemaRequest", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("SchemaRequest", "status", "started")
	return b.next.SchemaRequest(ctx)
}

func (b *backendTracingMiddleware) Initialize(ctx context.Context, req *logical.InitializationRequest) (err error) {
	defer func(then time.Time) {
		b.logger.Trace("initialize", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("initialize", "status", "started")
	return b.next.Initialize(ctx, req)
}

func (b *backendTracingMiddleware) HandleExistenceCheck(ctx context.Context, req *logical.Request) (found bool, exists bool, err error) {
	defer func(then time.Time) {
		b.logger.Trace("handle existence check", "namespace", req.Namespace, "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("handle existence check", "namespace", req.Namespace, "request", req, "status", "started", "req")
	return b.next.HandleExistenceCheck(ctx, req)
}
func (b *backendTracingMiddleware) HandleRequest(ctx context.Context, req *logical.Request) (resp *logical.Response, err error) {

	if b.next.Logger().IsTrace() {
		defer func(then time.Time) {
			b.logger.Trace("handle request", "path", req.Namespace,
				"status", "finished", "err", err, "took", time.Since(then),
				"resp", jsonutil.EncodeToString(resp))
		}(time.Now())

		b.logger.Trace("handle request",
			"path", req.Namespace, "status", "started", "request",
			jsonutil.EncodeToString(req), "req", req)
	}
	resp, err = b.next.HandleRequest(ctx, req)

	return resp, err
}

func (b *backendTracingMiddleware) Logger() log.Logger {
	return b.next.Logger()
}

func (b *backendTracingMiddleware) Cleanup(ctx context.Context) {
	defer func(then time.Time) {
		b.logger.Trace("cleanup", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("cleanup", "status", "started")
	b.next.Cleanup(ctx)
}

func (b *backendTracingMiddleware) Setup(ctx context.Context, config *logical.BackendConfig) (err error) {
	defer func(then time.Time) {
		b.logger.Trace("setup", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("setup", "status", "started")
	return b.next.Setup(ctx, config)
}

func (b *backendTracingMiddleware) Type() logical.BackendType {
	defer func(then time.Time) {
		b.logger.Trace("type", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("type", "status", "started")
	return b.next.Type()
}
