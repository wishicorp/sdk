package http_gateway

import (
	"context"
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
func (m *HttpGateway) authorization(backend logical.Backend, request *logical.Request) (authResp *logical.Response, err error) {
	defer func() {
		if err != nil {
			m.logger.Error("authorization", "request", request, "err", err)
		}
	}()
	var schema *logical.Schema
	schema, err = m.getSchema(backend, request)
	if nil != err {
		return nil, err
	}
	if !schema.Authorized {
		return &logical.Response{
			ResultCode: 0,
			ResultMsg:  "",
			Content:    &logical.Content{Data: &logical.Authorized{}},
		}, nil
	}
	authBackend, has := m.pm.GetBackend(m.authMethod.Backend)
	if !has {
		err = logical.ErrAuthMethodNotFound
		return nil, err
	}

	authReq := logical.Request{
		Operation:     logical.Operation(m.authMethod.Operation),
		Namespace:     m.authMethod.Namespace,
		Authorization: nil,
		Token:         request.Token,
		Data:          request.Data,
		Connection:    request.Connection,
	}
	authResp, err = authBackend.HandleRequest(m.ctx, &authReq)
	if nil != err {
		return nil, err
	}
	if m.logger.IsTrace() {
		m.logger.Trace("auth reply", "reply", authResp)
	}
	return authResp, nil
}
