package proto

import (
	"encoding/json"
	"errors"
	"github.com/wishicorp/sdk/helper/errutil"
	"github.com/wishicorp/sdk/plugin/logical"
)

const (
	ErrTypeUnknown uint32 = iota
	ErrTypeUserError
	ErrTypeInternalError
	ErrTypeCodedError
	ErrTypeStatusBadRequest
	ErrTypeUnsupportedOperation
	ErrTypeUnsupportedPath
	ErrTypeInvalidRequest
	ErrTypePermissionDenied
	ErrTypeMultiAuthzPending
)

func ProtoErrToErr(e *ProtoError) error {
	if e == nil {
		return nil
	}

	var err error
	switch e.ErrType {
	case ErrTypeUnknown:
		err = errors.New(e.ErrMsg)
	case ErrTypeUserError:
		err = errutil.UserError{Err: e.ErrMsg}
	case ErrTypeInternalError:
		err = errutil.InternalError{Err: e.ErrMsg}
	case ErrTypeCodedError:
		err = logical.CodedError(int(e.ErrCode), e.ErrMsg)
	case ErrTypeStatusBadRequest:
		err = &logical.StatusBadRequest{Err: e.ErrMsg}
	case ErrTypeUnsupportedOperation:
		err = logical.ErrUnsupportedOperation
	case ErrTypeUnsupportedPath:
		err = logical.ErrUnsupportedPath
	case ErrTypeInvalidRequest:
		err = logical.ErrInvalidRequest
	case ErrTypePermissionDenied:
		err = logical.ErrPermissionDenied
	case ErrTypeMultiAuthzPending:
		err = logical.ErrMultiAuthzPending
	}

	return err
}

func ErrToString(e error) string {
	if e == nil {
		return ""
	}

	return e.Error()
}
func ErrToProtoErr(e error) *ProtoError {
	if e == nil {
		return nil
	}
	pbErr := &ProtoError{
		ErrMsg:  e.Error(),
		ErrType: ErrTypeUnknown,
	}

	switch e.(type) {
	case errutil.UserError:
		pbErr.ErrType = ErrTypeUserError
	case errutil.InternalError:
		pbErr.ErrType = ErrTypeInternalError
	case logical.HTTPCodedError:
		pbErr.ErrType = ErrTypeCodedError
		pbErr.ErrCode = int64(e.(logical.HTTPCodedError).Code())
	case *logical.StatusBadRequest:
		pbErr.ErrType = ErrTypeStatusBadRequest
	}

	switch {
	case e == logical.ErrUnsupportedOperation:
		pbErr.ErrType = ErrTypeUnsupportedOperation
	case e == logical.ErrUnsupportedPath:
		pbErr.ErrType = ErrTypeUnsupportedPath
	case e == logical.ErrInvalidRequest:
		pbErr.ErrType = ErrTypeInvalidRequest
	case e == logical.ErrPermissionDenied:
		pbErr.ErrType = ErrTypePermissionDenied
	case e == logical.ErrMultiAuthzPending:
		pbErr.ErrType = ErrTypeMultiAuthzPending
	}

	return pbErr
}

func LogicalRequestToProtoRequest(r *logical.Request) (*Request, error) {

	if r == nil {
		return nil, errors.New("request is null")
	}

	headers := map[string]*Header{}
	for k, v := range r.Headers {
		headers[k] = &Header{Header: v}
	}

	return &Request{
		Id:            r.ID,
		Operation:     string(r.Operation),
		Namespace:     r.Namespace,
		Data:          r.Data,
		Token:         r.Token,
		Authorization: r.Authorization,
		Headers:       headers,
	}, nil
}

func ProtoRequestToLogicalRequest(r *Request) (*logical.Request, error) {
	if r == nil {
		return nil, nil
	}

	var headers map[string][]string
	if len(r.Headers) > 0 {
		headers = make(map[string][]string, len(r.Headers))
		for k, v := range r.Headers {
			headers[k] = v.Header
		}
	}

	return &logical.Request{
		ID:            r.Id,
		Operation:     logical.Operation(r.Operation),
		Namespace:     r.Namespace,
		Token:         r.Token,
		Authorization: r.Authorization,
		Data:          r.Data,
		Headers:       headers,
		Connection:    ProtoConnectionToLogicalConnection(r.Connection),
	}, nil
}

func LogicalConnectionToProtoConnection(c *logical.Connection) *Connection {
	if c == nil {
		return nil
	}

	return &Connection{
		RemoteAddr: c.RemoteAddr,
	}
}

func ProtoConnectionToLogicalConnection(c *Connection) *logical.Connection {
	if c == nil {
		return nil
	}

	return &logical.Connection{
		RemoteAddr: c.RemoteAddr,
	}
}

func LogicalResponseToProtoResponse(r *logical.Response) (*HandlerResponse, error) {
	if r == nil {
		return nil, nil
	}

	buf, err := json.Marshal(r.Data)
	if err != nil {
		return nil, err
	}

	headers := map[string]*Header{}
	for k, v := range r.Headers {
		headers[k] = &Header{Header: v}
	}

	return &HandlerResponse{
		ResultCode: r.ResultCode,
		ResultMsg:  r.ResultMsg,
		Data:       string(buf),
		Headers:    headers,
	}, nil
}
func ProtoResponseToLogicalResponse(r *HandlerResponse) (*logical.Response, error) {
	if r == nil {
		return nil, nil
	}

	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(r.Data), &data)
	if err != nil {
		return nil, err
	}

	var headers map[string][]string
	if len(r.Headers) > 0 {
		headers = make(map[string][]string, len(r.Headers))
		for k, v := range r.Headers {
			headers[k] = v.Header
		}
	}

	return &logical.Response{
		ResultCode: r.ResultCode,
		ResultMsg:  r.ResultMsg,
		Data:       data,
		Headers:    headers,
	}, nil
}

func LogicalStorageEntryToProtoStorageEntry(e *logical.StorageEntry) *StorageEntry {
	if e == nil {
		return nil
	}

	return &StorageEntry{
		Key:   e.Key,
		Value: e.Value,
	}
}

func ProtoStorageEntryToLogicalStorageEntry(e *StorageEntry) *logical.StorageEntry {
	if e == nil {
		return nil
	}

	return &logical.StorageEntry{
		Key:   e.Key,
		Value: e.Value,
	}
}

func protoFieldToLogicalField(fields []*Field) []*logical.Field {
	if nil == fields || len(fields) == 0 {
		return []*logical.Field{}
	}
	var outFields []*logical.Field
	for _, field := range fields {
		of := logical.Field{
			Field:      field.Field,
			Name:       field.Name,
			Kind:       field.Kind,
			Required:   field.Required,
			Deprecated: field.Deprecated,
			IsList:     field.IsList,
			Example:    field.Example,
		}
		if field.Reference != nil {
			of.Reference = protoFieldToLogicalField(field.Reference)
		}
		outFields = append(outFields, &of)
	}
	return outFields
}
func ProtoNamespaceSchemasToLigicalNamespaceSchemas(ns *SchemaRequestReply) *logical.SchemaResponse {
	var schemas logical.NamespaceSchemas
	for _, schema := range ns.NamespaceSchemas {
		operations := map[logical.Operation]*logical.Schema{}
		for key, opt := range schema.Operations {
			operations[logical.Operation(key)] = &logical.Schema{
				Description: opt.Description,
				Authorized:  opt.Authorized,
				Deprecated:  opt.Deprecated,

				Input:  protoFieldToLogicalField(opt.Input),
				Output: protoFieldToLogicalField(opt.Output),
			}
		}

		sc := logical.NamespaceSchema{
			Namespace:   schema.Namespace,
			Description: schema.Description,
			Operations:  operations,
		}
		schemas = append(schemas, &sc)
	}
	response := &logical.SchemaResponse{
		NamespaceSchemas: schemas,
	}
	return response
}

func logicalFieldToProtoField(fields []*logical.Field) []*Field {
	if nil == fields || len(fields) == 0 {
		return []*Field{}
	}
	var outFields []*Field
	for _, field := range fields {
		of := Field{
			Field:      field.Field,
			Name:       field.Name,
			Kind:       field.Kind,
			Required:   field.Required,
			Deprecated: field.Deprecated,
			Example:    field.Example,
			IsList:     field.IsList,
		}
		if field.Reference != nil {
			of.Reference = logicalFieldToProtoField(field.Reference)
		}
		outFields = append(outFields, &of)
	}
	return outFields
}

func LogicalNamespaceSchemasToProtoNamespaceSchemas(ns *logical.SchemaResponse) *SchemaRequestReply {
	var schemas []*NamespaceSchema
	for _, schema := range ns.NamespaceSchemas {
		operations := map[string]*Schema{}
		for key, opt := range schema.Operations {
			operations[string(key)] = &Schema{
				Description: opt.Description,
				Authorized:  opt.Authorized,
				Deprecated:  opt.Deprecated,
				Input:       logicalFieldToProtoField(opt.Input),
				Output:      logicalFieldToProtoField(opt.Output),
			}
		}

		sc := NamespaceSchema{
			Namespace:   schema.Namespace,
			Description: schema.Description,
			Operations:  operations,
		}
		schemas = append(schemas, &sc)
	}
	response := &SchemaRequestReply{
		NamespaceSchemas: schemas,
	}
	return response
}
