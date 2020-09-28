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
)

var ErrServerInMetadataMode = errors.New("plugins server can not perform action while in metadata mode")
var _ proto.BackendServer = (*backendGRPCPluginServer)(nil)

type backendGRPCPluginServer struct {
	broker  *plugin.GRPCBroker
	backend logical.Backend

	factory logical.Factory

	brokeredClient *grpc.ClientConn

	logger log.Logger
}

func (b *backendGRPCPluginServer) SchemaRequest(ctx context.Context, empty *proto.Empty) (*proto.SchemaRequestReply, error) {
	resp, err := b.backend.SchemaRequest(ctx)
	if err != nil {
		return nil, err
	}
	return proto.LogicalNamespaceSchemasToProtoNamespaceSchemas(resp), nil
}

func (b *backendGRPCPluginServer) Setup(ctx context.Context, args *proto.SetupArgs) (*proto.SetupReply, error) {
	// Dial for storage
	brokeredClient, err := b.broker.Dial(args.BrokerId)
	if err != nil {
		return &proto.SetupReply{}, err
	}

	b.brokeredClient = brokeredClient
	consul := newGRPCConsulClient(brokeredClient)
	component := newGRPCComponentClient(brokeredClient, b.logger)
	config := &logical.BackendConfig{
		ConsulView:    consul,
		ComponentView: component,
		Logger:        b.logger,
		Config:        args.Config,
		BackendUUID:   args.BackendUUID,
	}

	// Call the underlying grpc-backend factory after shims have been created
	// to set b.grpc-backend
	backend, err := b.factory(ctx, config)
	if err != nil {
		return &proto.SetupReply{
			Err: proto.ErrToString(err),
		}, nil
	}
	b.backend = backend
	return &proto.SetupReply{}, nil
}

func (b *backendGRPCPluginServer) Initialize(ctx context.Context, args *proto.InitializeArgs) (*proto.InitializeReply, error) {
	if pluginutil.InMetadataMode() {
		return &proto.InitializeReply{}, ErrServerInMetadataMode
	}

	req := &logical.InitializationRequest{
		Params: args.Params,
	}

	respErr := b.backend.Initialize(ctx, req)

	return &proto.InitializeReply{
		Err: proto.ErrToProtoErr(respErr),
	}, nil
}

func (b *backendGRPCPluginServer) HandleExistenceCheck(ctx context.Context, args *proto.HandleExistenceCheckArgs) (*proto.HandleExistenceCheckReply, error) {
	if pluginutil.InMetadataMode() {
		return &proto.HandleExistenceCheckReply{}, ErrServerInMetadataMode
	}

	logicalReq, err := proto.ProtoRequestToLogicalRequest(args.Request)
	if err != nil {
		return &proto.HandleExistenceCheckReply{}, err
	}
	//logicalReq.Storage = newGRPCStorageClient(b.brokeredClient)
	//logicalReq.Redis = newGRPCRedisClient(b.brokeredClient)
	checkFound, exists, err := b.backend.HandleExistenceCheck(ctx, logicalReq)
	return &proto.HandleExistenceCheckReply{
		CheckFound: checkFound,
		Exists:     exists,
		Err:        proto.ErrToProtoErr(err),
	}, err
}

func (b *backendGRPCPluginServer) HandleRequest(ctx context.Context, args *proto.HandleRequestArgs) (*proto.HandleRequestReply, error) {
	if pluginutil.InMetadataMode() {
		return &proto.HandleRequestReply{}, ErrServerInMetadataMode
	}

	logicalReq, err := proto.ProtoRequestToLogicalRequest(args.Request)
	if err != nil {
		return &proto.HandleRequestReply{}, err
	}

	//logicalReq.Storage = newGRPCStorageClient(b.brokeredClient)
	//logicalReq.Redis = newGRPCRedisClient(b.brokeredClient)
	resp, respErr := b.backend.HandleRequest(ctx, logicalReq)
	if respErr != nil {
		return nil, respErr
	}

	pbResp, err := proto.LogicalResponseToProtoResponse(resp)
	if err != nil {
		return &proto.HandleRequestReply{}, err
	}

	return &proto.HandleRequestReply{
		Response: pbResp,
		Err:      proto.ErrToProtoErr(respErr),
	}, nil
}

func (b *backendGRPCPluginServer) Cleanup(ctx context.Context, empty *proto.Empty) (*proto.Empty, error) {
	b.backend.Cleanup(ctx)
	// Close rpc clients
	b.brokeredClient.Close()
	return &proto.Empty{}, nil
}

func (b *backendGRPCPluginServer) Type(ctx context.Context, _ *proto.Empty) (*proto.TypeReply, error) {
	return &proto.TypeReply{
		Type: uint32(b.backend.Type()),
	}, nil
}

func (b backendGRPCPluginServer) Version(ctx context.Context, _ *proto.Empty) (*proto.VersionReply, error) {
	version := b.backend.Version(ctx)

	return &proto.VersionReply{Version: version}, nil
}
