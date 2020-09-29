package grpc_gateway

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/wishicorp/sdk/helper/jsonutil"
	proto "github.com/wishicorp/sdk/plugin/gateway/grpc-gateway/proto"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/pluginregister"
)

var _ proto.RpcGatewayServer = (*GRPCGatewayImpl)(nil)

type GRPCGatewayImpl struct {
	pm     *pluginregister.PluginManager
	logger hclog.Logger
}

func (m *GRPCGateway) NewGRPCGatewayImpl() *GRPCGatewayImpl {
	return &GRPCGatewayImpl{
		pm:     m.pm,
		logger: m.logger,
	}
}

func (m *GRPCGatewayImpl) Schemas(ctx context.Context, args *proto.SchemasArgs) (*proto.SchemasReply, error) {
	backends := make([]map[string]string, 0)
	protoSchemas := map[string]*proto.Schemas{}
	if args.Backend != "" {
		backends = append(backends, map[string]string{"name": args.Backend})
	} else {
		backends = m.pm.List()
	}

	for _, b := range backends {
		backend, has := m.pm.GetBackend(b["name"])
		if !has {
			continue
		}
		resp, err := backend.SchemaRequest(context.Background())
		if nil != err {
			return nil, err
		}
		if args.Namespace != "" {
			for _, schema := range resp.NamespaceSchemas {
				if schema.Namespace == args.Namespace {
					protoSchemas[args.Namespace] =
						m.toProtoNamespaceSchemas([]*logical.NamespaceSchema{schema})
				}
			}
			break
		} else {
			protoSchemas[b["name"]] = m.toProtoNamespaceSchemas(resp.NamespaceSchemas)
		}
	}

	return &proto.SchemasReply{SchemasMap: protoSchemas}, nil
}

//TODO 未完成实现，后期和httpGateway 合并实现
func (m *GRPCGatewayImpl) ExecRequest(ctx context.Context, args *proto.RequestArgs) (*proto.RequestReply, error) {
	return nil, logical.ErrNotImplementation
}

func (m *GRPCGatewayImpl) toProtoNamespaceSchemas(nss logical.NamespaceSchemas) *proto.Schemas {
	var schemas []*proto.Schema
	for _, schema := range nss {
		proSchema := proto.Schema{
			Namespace:   schema.Namespace,
			Description: schema.Description,
			Operations:  m.toProOperations(schema.Operations),
		}
		schemas = append(schemas, &proSchema)
	}
	return &proto.Schemas{Schemas: schemas}
}

func (m *GRPCGatewayImpl) toProOperations(operation map[logical.Operation]*logical.Schema) map[string]*proto.Operation {
	operations := map[string]*proto.Operation{}
	for l, schema := range operation {
		var op proto.Operation
		if err := jsonutil.Swap(schema, &op); err != nil {
			m.logger.Error("trans operation", "err", err)
		}
		operations[string(l)] = &op
	}
	return operations
}
