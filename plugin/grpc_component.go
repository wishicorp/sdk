package plugin

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/proto"
	"google.golang.org/grpc"
	"sync"
)

var _ logical.Component = (*GRPCComponentClient)(nil)

func newGRPCComponentClient(conn *grpc.ClientConn, logger hclog.Logger) *GRPCComponentClient {
	return &GRPCComponentClient{
		client: proto.NewComponentClient(conn),
		once:   sync.Once{},
		logger: logger,
		cfg:    new(logical.ComponentConfig),
	}
}

// GRPCComponentClient is an implementation of logical.Component that communicates
// over RPC.
type GRPCComponentClient struct {
	client    proto.ComponentClient
	component logical.Component
	logger    hclog.Logger
	once      sync.Once
	cfg       *logical.ComponentConfig
}

func (s *GRPCComponentClient) CreateFactory(ctx context.Context) (logical.ComponentFactory, error) {
	s.once.Do(s.loadConfig)
	return logical.NewComponent(s.cfg, s.logger).CreateFactory(ctx)
}

func (s *GRPCComponentClient) loadConfig() {
	resp, _ := s.FetchConfig(context.Background(), "", "", "")
	s.cfg = resp
}

func (s *GRPCComponentClient) FetchConfig(ctx context.Context, namespace, key, version string) (*logical.ComponentConfig, error) {
	req := proto.FetchConfigArgs{
		Namespace: namespace,
		Key:       key,
		Version:   version,
	}
	resp, err := s.client.FetchConfig(ctx, &req, largeMsgGRPCCallOpts...)
	if nil != err {
		return nil, err
	}
	var config logical.ComponentConfig
	jsonutil.DecodeJSON(resp.Value, &config)
	return &config, err
}

type GRPCComponentServer struct {
	impl logical.Component
}

func (s *GRPCComponentServer) CreateFactory(ctx context.Context, empty *proto.CEmpty) (*proto.CEmpty, error) {
	return empty, nil
}

func (s *GRPCComponentServer) FetchConfig(ctx context.Context, c *proto.FetchConfigArgs) (*proto.FetchConfigReply, error) {
	resp, err := s.impl.FetchConfig(ctx, c.Namespace, c.Key, c.Version)
	if nil != err {
		return nil, err
	}
	values, err := jsonutil.EncodeJSON(resp)
	if nil != err {
		return nil, err
	}
	return &proto.FetchConfigReply{Value: values}, nil
}
