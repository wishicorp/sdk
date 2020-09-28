package plugin

import (
	"context"
	"errors"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/proto"

	"google.golang.org/grpc"
)

var _ logical.Storage = (*GRPCStorageClient)(nil)

var ErrGRPCStorageNotImplemented = errors.New("grpc storage not implemented")

func newGRPCStorageClient(conn *grpc.ClientConn) *GRPCStorageClient {
	return &GRPCStorageClient{
		client: proto.NewStorageClient(conn),
	}
}

// GRPCStorageClient is an implementation of logical.Storage that communicates
// over RPC.
type GRPCStorageClient struct {
	client proto.StorageClient
}

func (s *GRPCStorageClient) List(ctx context.Context, prefix string) ([]*logical.StorageEntry, error) {
	reply, err := s.client.List(ctx, &proto.StorageListArgs{
		Prefix: prefix,
	}, largeMsgGRPCCallOpts...)
	if err != nil {
		return []*logical.StorageEntry{}, err
	}
	if reply.Err != "" {
		return nil, errors.New(reply.Err)
	}
	var entities []*logical.StorageEntry

	for _, entity := range reply.Entities {
		entities = append(entities, proto.ProtoStorageEntryToLogicalStorageEntry(entity))
	}

	return entities, nil
}

func (s *GRPCStorageClient) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	reply, err := s.client.Get(ctx, &proto.StorageGetArgs{
		Key: key,
	}, largeMsgGRPCCallOpts...)
	if err != nil {
		return nil, err
	}
	if reply.Err != "" {
		return nil, errors.New(reply.Err)
	}
	return proto.ProtoStorageEntryToLogicalStorageEntry(reply.Entry), nil
}

func (s *GRPCStorageClient) Put(ctx context.Context, entry *logical.StorageEntry) error {
	reply, err := s.client.Put(ctx, &proto.StoragePutArgs{
		Entry: proto.LogicalStorageEntryToProtoStorageEntry(entry),
	}, largeMsgGRPCCallOpts...)
	if err != nil {
		return err
	}
	if reply.Err != "" {
		return errors.New(reply.Err)
	}
	return nil
}

func (s *GRPCStorageClient) Delete(ctx context.Context, key string) error {
	reply, err := s.client.Delete(ctx, &proto.StorageDeleteArgs{
		Key: key,
	})
	if err != nil {
		return err
	}
	if reply.Err != "" {
		return errors.New(reply.Err)
	}
	return nil
}

// StorageServer is a net/rpc compatible structure for serving
type GRPCStorageServer struct {
	impl logical.Storage
}

func (s *GRPCStorageServer) List(ctx context.Context, args *proto.StorageListArgs) (*proto.StorageListReply, error) {
	if nil == s.impl {
		return nil, ErrGRPCStorageNotImplemented
	}

	reply, err := s.impl.List(ctx, args.Prefix)
	var entities []*proto.StorageEntry
	for _, entry := range reply {
		entities = append(entities, proto.LogicalStorageEntryToProtoStorageEntry(entry))
	}
	return &proto.StorageListReply{
		Entities: entities,
		Err:      proto.ErrToString(err),
	}, nil
}

func (s *GRPCStorageServer) Get(ctx context.Context, args *proto.StorageGetArgs) (*proto.StorageGetReply, error) {
	if nil == s.impl {
		return nil, ErrGRPCStorageNotImplemented
	}

	storageEntry, err := s.impl.Get(ctx, args.Key)
	return &proto.StorageGetReply{
		Entry: proto.LogicalStorageEntryToProtoStorageEntry(storageEntry),
		Err:   proto.ErrToString(err),
	}, nil
}

func (s *GRPCStorageServer) Put(ctx context.Context, args *proto.StoragePutArgs) (*proto.StoragePutReply, error) {
	if nil == s.impl {
		return nil, ErrGRPCStorageNotImplemented
	}

	err := s.impl.Put(ctx, proto.ProtoStorageEntryToLogicalStorageEntry(args.Entry))
	return &proto.StoragePutReply{
		Err: proto.ErrToString(err),
	}, nil
}

func (s *GRPCStorageServer) Delete(ctx context.Context, args *proto.StorageDeleteArgs) (*proto.StorageDeleteReply, error) {
	if nil == s.impl {
		return nil, ErrGRPCStorageNotImplemented
	}

	err := s.impl.Delete(ctx, args.Key)
	return &proto.StorageDeleteReply{
		Err: proto.ErrToString(err),
	}, nil
}

// NOOPStorage is used to deny access to the storage interface while running a
// backend plugin in metadata mode.
type NOOPStorage struct{}

func (s *NOOPStorage) List(_ context.Context, prefix string) ([]*logical.StorageEntry, error) {
	return []*logical.StorageEntry{}, nil
}

func (s *NOOPStorage) Get(_ context.Context, key string) (*logical.StorageEntry, error) {
	return nil, nil
}

func (s *NOOPStorage) Put(_ context.Context, entry *logical.StorageEntry) error {
	return nil
}

func (s *NOOPStorage) Delete(_ context.Context, key string) error {
	return nil
}
