package plugin

import (
	"google.golang.org/grpc"
	"math"
)

var largeMsgGRPCCallOpts []grpc.CallOption = []grpc.CallOption{
	grpc.MaxCallSendMsgSize(math.MaxInt32),
	grpc.MaxCallRecvMsgSize(math.MaxInt32),
}

// This is the implementation of plugins.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugins can be served over
// gRPC.
//type GRPCBackendPlugin2 struct {
//	logical.Backend
//	Factory logical.Factory
//	Logger  log.Logger
//	plugins.NetRPCUnsupportedPlugin
//}
//
//func (p *GRPCBackendPlugin2) GRPCServer(broker *plugins.GRPCBroker, s *grpc.Server) error {
//	proto.RegisterBackendServer(s, &GRPCBackendServer{
//		Impl:   p.Backend,
//		broker: broker,
//	})
//	return nil
//}
//
//func (p *GRPCBackendPlugin2) GRPCClient(ctx context.Context, broker *plugins.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
//	return &GRPCBackendClient{
//		client: proto.NewBackendClient(c),
//		broker: broker,
//	}, nil
//}
//
//var _ plugins.GRPCPlugin = &GRPCBackendPlugin2{}
