package grpc_gateway

import (
	"context"
	"github.com/wishicorp/sdk/plugin/gateway/grpc-gateway/proto"
	"github.com/wishicorp/sdk/plugin/logical"
)

var _ proto.RpcGatewayServer = (*GRPCGateway)(nil)

//TODO 未完成实现，后期和httpGateway 合并实现
func (m *GRPCGateway) Schemas(ctx context.Context, args *proto.SchemasArgs) (*proto.SchemasReply, error) {
	return nil, logical.ErrNotImplementation
}

func (m *GRPCGateway) ExecRequest(ctx context.Context, args *proto.RequestArgs) (*proto.RequestReply, error) {
	return nil, logical.ErrNotImplementation
}

func (m *GRPCGateway) executeGetSchemas() map[string]interface{} {
	schemas := map[string]interface{}{}
	backends := m.pm.List()
	for _, b := range backends {
		backend, has := m.pm.GetBackend(b["name"])
		if !has {
			continue
		}
		resp, err := backend.SchemaRequest(context.Background())
		if nil != err {
			continue
		}
		schemas[b["name"]] = resp.NamespaceSchemas
	}
	return schemas
}
