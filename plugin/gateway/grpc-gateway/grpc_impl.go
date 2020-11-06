package grpc_gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"github.com/wishicorp/sdk/plugin/gateway"
	gwproto "github.com/wishicorp/sdk/plugin/gateway/grpc-gateway/proto"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/pluginregister"
)

var _ gwproto.RpcGatewayServer = (*GRPCGatewayImpl)(nil)

type GRPCGatewayImpl struct {
	pm          *pluginregister.PluginManager
	logger      hclog.Logger
	authMethod  *gateway.Method
	authEnabled bool
}

func (m *GRPCGateway) NewGRPCGatewayImpl() *GRPCGatewayImpl {
	return &GRPCGatewayImpl{
		pm:          m.pm,
		logger:      m.logger,
		authMethod:  m.authMethod,
		authEnabled: m.authEnabled,
	}
}

func (m *GRPCGatewayImpl) Schemas(ctx context.Context, args *gwproto.SchemasArgs) (*gwproto.SchemasReply, error) {
	backends := make([]map[string]string, 0)
	protoSchemas := map[string]*gwproto.Schemas{}
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

	return &gwproto.SchemasReply{SchemasMap: protoSchemas}, nil
}

//TODO 未完成实现，后期和httpGateway 合并实现
func (m *GRPCGatewayImpl) ExecRequest(ctx context.Context, args *gwproto.RequestArgs) (gwReply *gwproto.RequestReply, err error) {
	defer func() {
		if nil != err{
			m.logger.Error("exec request",
				"args", jsonutil.EncodeToString(args),
				"err", err)
		}else {
			if m.logger.IsTrace(){
				m.logger.Trace("exec request",
					"args", jsonutil.EncodeToString(args),
					"reply", jsonutil.EncodeToString(gwReply))
			}
		}
	}()
	if args.Data == nil{
		return nil, errors.New("args[data] is null")
	}
	backend, has := m.pm.GetBackend(args.Backend)
	if !has {
		return nil, pluginregister.PluginNotExists
	}

	backend.Incr()
	defer backend.DeIncr()
	req := &logical.Request{
		Operation: logical.Operation(args.Operation),
		Namespace: args.Namespace,
		Token:     args.Token,
		Data: map[string][]byte{"data": args.Data},
	}

	if m.authEnabled && m.authMethod != nil {
		var authReply *logical.Response
		authReply, err = m.authorization(backend, req)
		if err != nil {
			return nil, fmt.Errorf("auth: %s", err.Error())
		}
		if authReply.ResultCode != 0 {
			return &gwproto.RequestReply{
				Result: &gwproto.Result{
					ResultCode: int32(authReply.ResultCode),
					ResultMsg:  authReply.ResultMsg,
				},
			}, nil
		}
		if err = jsonutil.Swap(authReply.Content.Data, &req.Authorized); err != nil {
			return nil, err
		}
	}
	var reply *logical.Response
	reply, err = backend.HandleRequest(ctx, req)
	if nil != err {
		return nil, err
	}
	var data []byte
	data, err = jsonutil.EncodeJSON(reply.Content)
	gwReply = &gwproto.RequestReply{
		Result: &gwproto.Result{
			ResultCode: int32(reply.ResultCode),
			ResultMsg:  reply.ResultMsg,
			Data:       data,
		},
	}
	return gwReply, err
}

func (m *GRPCGatewayImpl) toProtoNamespaceSchemas(nss logical.NamespaceSchemas) *gwproto.Schemas {
	var schemas []*gwproto.Schema
	for _, schema := range nss {
		proSchema := gwproto.Schema{
			Namespace:   schema.Namespace,
			Description: schema.Description,
			Operations:  m.toProOperations(schema.Operations),
		}
		schemas = append(schemas, &proSchema)
	}
	return &gwproto.Schemas{Schemas: schemas}
}

func (m *GRPCGatewayImpl) toProOperations(operation map[logical.Operation]*logical.Schema) map[string]*gwproto.Operation {
	operations := map[string]*gwproto.Operation{}
	for l, schema := range operation {
		var op gwproto.Operation
		if err := jsonutil.Swap(schema, &op); err != nil {
			m.logger.Error("trans operation", "err", err)
		}
		operations[string(l)] = &op
	}
	return operations
}
