package plugin

import (
	"context"
	"errors"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/pluginutil"
	"github.com/wishicorp/sdk/plugin/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"sync/atomic"
)

var ErrPluginShutdown = errors.New("plugins is shut down")
var ErrClientInMetadataMode = errors.New("plugins client can not perform action while in metadata mode")

// Validate backendGRPCPluginClient satisfies the logical.Backend interface
var _ logical.Backend = &backendGRPCPluginClient{}

// backendPluginClient implements logical.Backend and is the
// go-plugins client.
type backendGRPCPluginClient struct {
	broker       *plugin.GRPCBroker
	client       proto.BackendClient
	metadataMode bool

	logger log.Logger

	// This is used to signal to the Cleanup function that it can proceed
	// because we have a defined server
	cleanupCh chan struct{}

	// server is the grpc server used for serving storage and sysview requests.
	server *atomic.Value

	// clientConn is the underlying grpc connection to the server, we store it
	// so it can be cleaned up.
	clientConn *grpc.ClientConn
	doneCtx    context.Context
}

func (b *backendGRPCPluginClient) Name() string {
	reply, err := b.client.Name(context.Background(), &proto.Empty{})
	if err != nil {
		return ""
	}
	return reply.GetName()
}

func (b *backendGRPCPluginClient) SchemaRequest(ctx context.Context) (*logical.SchemaReply, error) {
	if b.metadataMode {
		return nil, ErrClientInMetadataMode
	}
	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	resp, err := b.client.SchemaRequest(ctx, &proto.Empty{})
	if err != nil {
		if b.doneCtx.Err() != nil {
			return nil, ErrPluginShutdown
		}
		return nil, err
	}

	return proto.ProtoNamespaceSchemasToLigicalNamespaceSchemas(resp), nil
}

func (b *backendGRPCPluginClient) Setup(ctx context.Context, config *logical.BackendConfig) error {
	consul := &GRPCConsulServer{
		impl: config.ConsulView,
	}
	component := &GRPCComponentServer{
		impl: config.ComponentView,
	}
	// Register the server in this closure.
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		opts = append(opts, grpc.MaxRecvMsgSize(math.MaxInt32))
		opts = append(opts, grpc.MaxSendMsgSize(math.MaxInt32))

		s := grpc.NewServer(opts...)
		b.server.Store(s)
		proto.RegisterConsulServer(s, consul)
		proto.RegisterComponentServer(s, component)
		close(b.cleanupCh)
		return s
	}
	brokerID := b.broker.NextId()
	go b.broker.AcceptAndServe(brokerID, serverFunc)

	args := &proto.SetupArgs{
		BrokerId:    brokerID,
		Config:      config.Config,
		BackendUUID: config.BackendUUID,
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	reply, err := b.client.Setup(ctx, args)
	if err != nil {
		return err
	}
	if reply.Err != "" {
		return errors.New(reply.Err)
	}

	b.logger = config.Logger

	return nil
}

func (b *backendGRPCPluginClient) Initialize(ctx context.Context, args *logical.InitializationRequest) error {
	if b.metadataMode {
		return ErrClientInMetadataMode
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	reply, err := b.client.Initialize(ctx, &proto.InitializeArgs{Params: args.Params}, largeMsgGRPCCallOpts...)
	if err != nil {
		if b.doneCtx.Err() != nil {
			return ErrPluginShutdown
		}

		grpcStatus, ok := status.FromError(err)
		if ok && grpcStatus.Code() == codes.Unimplemented {
			return nil
		}

		return err
	}
	if reply.Err != nil {
		return proto.ProtoErrToErr(reply.Err)
	}

	return nil
}

func (b *backendGRPCPluginClient) HandleExistenceCheck(ctx context.Context, req *logical.Request) (bool, bool, error) {
	if b.metadataMode {
		return false, false, ErrClientInMetadataMode
	}

	protoReq, err := proto.LogicalRequestToProtoRequest(req)
	if err != nil {
		return false, false, err
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()
	reply, err := b.client.HandleExistenceCheck(ctx, &proto.HandleExistenceCheckArgs{
		Request: protoReq,
	}, largeMsgGRPCCallOpts...)

	if err != nil {
		if b.doneCtx.Err() != nil {
			return false, false, ErrPluginShutdown
		}
		return false, false, err
	}
	if reply.Err != nil {
		return false, false, proto.ProtoErrToErr(reply.Err)
	}

	return reply.CheckFound, reply.Exists, nil
}

func (b *backendGRPCPluginClient) HandleRequest(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	if b.metadataMode {
		return nil, ErrClientInMetadataMode
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	protoReq, err := proto.LogicalRequestToProtoRequest(req)
	if err != nil {
		return nil, err
	}

	reply, err := b.client.HandleRequest(ctx, &proto.HandleRequestArgs{
		Request: protoReq,
	}, largeMsgGRPCCallOpts...)
	if err != nil {
		if b.doneCtx.Err() != nil {
			return nil, ErrPluginShutdown
		}

		return nil, err
	}
	resp, err := proto.ProtoResponseToLogicalResponse(reply.Response)
	if err != nil {
		return nil, err
	}
	if reply.Err != nil {
		return resp, proto.ProtoErrToErr(reply.Err)
	}

	return resp, nil
}

func (b *backendGRPCPluginClient) Logger() log.Logger {
	return b.logger
}

func (b *backendGRPCPluginClient) Version(ctx context.Context) string {

	reply, err := b.client.Version(b.doneCtx, &proto.Empty{})

	if err != nil {
		return "unknown"
	}

	return reply.Version
}

func (b *backendGRPCPluginClient) Cleanup(ctx context.Context) {
	b.logger.Trace("cleanup", "context", ctx)

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	b.client.Cleanup(ctx, &proto.Empty{})

	// This will block until Setup has run the function to create a new server
	// in b.server. If we stop here before it has a chance to actually start
	// listening, when it starts listening it will immediately error out and
	// exit, which is fine. Overall this ensures that we do not miss stopping
	// the server if it ends up being created after Cleanup is called.
	<-b.cleanupCh
	server := b.server.Load()
	if server != nil {
		server.(*grpc.Server).GracefulStop()
	}
	b.clientConn.Close()
}

func (b *backendGRPCPluginClient) Type() logical.BackendType {
	reply, err := b.client.Type(b.doneCtx, &proto.Empty{})
	if err != nil {
		return logical.TypeUnknown
	}

	return logical.BackendType(reply.Type)
}
