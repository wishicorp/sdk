package http_gateway

import (
	"context"
	"errors"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"github.com/wishicorp/sdk/plugin/framework"
	"github.com/wishicorp/sdk/plugin/logical"
)

//获取插件的操作名称
func (m *HttpGateway) getSchema(backend logical.Backend, request *logical.Request) (
	*logical.Schema, error) {

	result, err := backend.SchemaRequest(context.Background())
	if nil != err {
		return nil, err
	}

	var schema *logical.NamespaceSchema
	for _, n := range result.NamespaceSchemas {
		if n.Namespace == request.Namespace {
			schema = n
		}
	}
	if nil == schema {
		return nil, framework.ErrNamespaceNotExists
	}
	for operation, properties := range schema.Operations {
		if operation == request.Operation {
			return properties, nil
		}
	}
	return nil, framework.ErrOperationNotExists
}
func (m *HttpGateway) authorization(backend logical.Backend, request *logical.Request) (logical.Authorized, error) {
	if m.authMethod == nil {
		return logical.Authorized{}, nil
	}

	schema, err := m.getSchema(backend, request)
	if nil != err {
		return nil, err
	}
	if !schema.Authorized {
		return nil, nil
	}

	authBackend, has := m.pm.GetBackend(m.authMethod.Backend)
	if !has {
		return nil, logical.ErrAuthMethodNotFound
	}

	authReq := logical.Request{
		Operation:     logical.Operation(m.authMethod.Operation),
		Namespace:     m.authMethod.Namespace,
		Authorization: nil,
		Token:         request.Token,
		Data:          request.Data,
		Connection:    request.Connection,
	}
	authReply, err := authBackend.HandleRequest(m.ctx, &authReq)
	if nil != err {
		return nil, errors.New("auth " + err.Error())
	}
	if authReply.ResultCode != 0 {
		return nil, errors.New(authReply.ResultMsg)
	}
	if m.logger.IsTrace() {
		m.logger.Trace("auth reply", "reply", authReply)
	}
	var auth logical.Authorized
	if err := jsonutil.Swap(authReply.Data, &auth); err != nil {
		return nil, err
	}

	return auth, nil
}
