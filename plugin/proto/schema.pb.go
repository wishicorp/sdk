// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.4
// source: schema.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type SchemaRequestArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SchemaRequestArgs) Reset() {
	*x = SchemaRequestArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schema_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SchemaRequestArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SchemaRequestArgs) ProtoMessage() {}

func (x *SchemaRequestArgs) ProtoReflect() protoreflect.Message {
	mi := &file_schema_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SchemaRequestArgs.ProtoReflect.Descriptor instead.
func (*SchemaRequestArgs) Descriptor() ([]byte, []int) {
	return file_schema_proto_rawDescGZIP(), []int{0}
}

//接口属性列
//type Field struct {
//    Field      string `json:"field"`
//    Name       string `json:"name"`
//    Type       string `json:"type"`
//    Required   bool   `json:"required"`
//    Deprecated bool   `json:"deprecated"`
//}
type Field struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field      string   `protobuf:"bytes,1,opt,name=field,proto3" json:"field,omitempty"`
	Name       string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Kind       string   `protobuf:"bytes,3,opt,name=kind,proto3" json:"kind,omitempty"`
	Required   bool     `protobuf:"varint,4,opt,name=required,proto3" json:"required,omitempty"`
	Deprecated bool     `protobuf:"varint,5,opt,name=deprecated,proto3" json:"deprecated,omitempty"`
	IsList     bool     `protobuf:"varint,6,opt,name=isList,proto3" json:"isList,omitempty"`
	Example    string   `protobuf:"bytes,8,opt,name=example,proto3" json:"example,omitempty"`
	Reference  []*Field `protobuf:"bytes,9,rep,name=reference,proto3" json:"reference,omitempty"`
}

func (x *Field) Reset() {
	*x = Field{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schema_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Field) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Field) ProtoMessage() {}

func (x *Field) ProtoReflect() protoreflect.Message {
	mi := &file_schema_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Field.ProtoReflect.Descriptor instead.
func (*Field) Descriptor() ([]byte, []int) {
	return file_schema_proto_rawDescGZIP(), []int{1}
}

func (x *Field) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *Field) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Field) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *Field) GetRequired() bool {
	if x != nil {
		return x.Required
	}
	return false
}

func (x *Field) GetDeprecated() bool {
	if x != nil {
		return x.Deprecated
	}
	return false
}

func (x *Field) GetIsList() bool {
	if x != nil {
		return x.IsList
	}
	return false
}

func (x *Field) GetExample() string {
	if x != nil {
		return x.Example
	}
	return ""
}

func (x *Field) GetReference() []*Field {
	if x != nil {
		return x.Reference
	}
	return nil
}

////接口属性
//type Schema struct {
//	Description string   `json:"description"`
//	Authorized  bool     `json:"authorized"`
//	Deprecated  bool     `json:"deprecated"`
//	Input       []*Field `json:"input,omitempty"`
//	Output      []*Field `json:"output,omitempty"`
//}
type Schema struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description string   `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	Authorized  bool     `protobuf:"varint,2,opt,name=authorized,proto3" json:"authorized,omitempty"`
	Deprecated  bool     `protobuf:"varint,3,opt,name=deprecated,proto3" json:"deprecated,omitempty"`
	Input       []*Field `protobuf:"bytes,4,rep,name=input,proto3" json:"input,omitempty"`
	Output      []*Field `protobuf:"bytes,5,rep,name=output,proto3" json:"output,omitempty"`
}

func (x *Schema) Reset() {
	*x = Schema{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schema_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Schema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Schema) ProtoMessage() {}

func (x *Schema) ProtoReflect() protoreflect.Message {
	mi := &file_schema_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Schema.ProtoReflect.Descriptor instead.
func (*Schema) Descriptor() ([]byte, []int) {
	return file_schema_proto_rawDescGZIP(), []int{2}
}

func (x *Schema) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Schema) GetAuthorized() bool {
	if x != nil {
		return x.Authorized
	}
	return false
}

func (x *Schema) GetDeprecated() bool {
	if x != nil {
		return x.Deprecated
	}
	return false
}

func (x *Schema) GetInput() []*Field {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *Schema) GetOutput() []*Field {
	if x != nil {
		return x.Output
	}
	return nil
}

//type NamespaceSchema struct {
//	Namespace   string `json:"namespace"`
//	Description string `json:"description"`
//	Operations  map[Operation]*Schema
//}
type NamespaceSchema struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace   string             `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Description string             `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Operations  map[string]*Schema `protobuf:"bytes,3,rep,name=operations,proto3" json:"operations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *NamespaceSchema) Reset() {
	*x = NamespaceSchema{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schema_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NamespaceSchema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NamespaceSchema) ProtoMessage() {}

func (x *NamespaceSchema) ProtoReflect() protoreflect.Message {
	mi := &file_schema_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NamespaceSchema.ProtoReflect.Descriptor instead.
func (*NamespaceSchema) Descriptor() ([]byte, []int) {
	return file_schema_proto_rawDescGZIP(), []int{3}
}

func (x *NamespaceSchema) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *NamespaceSchema) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *NamespaceSchema) GetOperations() map[string]*Schema {
	if x != nil {
		return x.Operations
	}
	return nil
}

type SchemaRequestReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NamespaceSchemas []*NamespaceSchema `protobuf:"bytes,1,rep,name=Namespaces,proto3" json:"Namespaces,omitempty"`
	Err              string             `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *SchemaRequestReply) Reset() {
	*x = SchemaRequestReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schema_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SchemaRequestReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SchemaRequestReply) ProtoMessage() {}

func (x *SchemaRequestReply) ProtoReflect() protoreflect.Message {
	mi := &file_schema_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SchemaRequestReply.ProtoReflect.Descriptor instead.
func (*SchemaRequestReply) Descriptor() ([]byte, []int) {
	return file_schema_proto_rawDescGZIP(), []int{4}
}

func (x *SchemaRequestReply) GetNamespaceSchemas() []*NamespaceSchema {
	if x != nil {
		return x.NamespaceSchemas
	}
	return nil
}

func (x *SchemaRequestReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

var File_schema_proto protoreflect.FileDescriptor

var file_schema_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x13, 0x0a, 0x11, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x72, 0x67, 0x73, 0x22, 0xdf, 0x01, 0x0a, 0x05, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x12, 0x1e,
	0x0a, 0x0a, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0a, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x69, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x69, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x12, 0x2a, 0x0a, 0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x09, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x52, 0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x22, 0xb4, 0x01, 0x0a,
	0x06, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x61,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x65, 0x70,
	0x72, 0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x64,
	0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x12, 0x22, 0x0a, 0x05, 0x69, 0x6e, 0x70,
	0x75, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x24, 0x0a,
	0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x06, 0x6f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x22, 0xe7, 0x01, 0x0a, 0x0f, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x46, 0x0a, 0x0a, 0x6f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x53, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x0a, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a,
	0x4c, 0x0a, 0x0f, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x23, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x6a, 0x0a,
	0x12, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x42, 0x0a, 0x10, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65,
	0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x53,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x10, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65,
	0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x42, 0x0b, 0x48, 0x02, 0x5a, 0x07, 0x2e,
	0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_schema_proto_rawDescOnce sync.Once
	file_schema_proto_rawDescData = file_schema_proto_rawDesc
)

func file_schema_proto_rawDescGZIP() []byte {
	file_schema_proto_rawDescOnce.Do(func() {
		file_schema_proto_rawDescData = protoimpl.X.CompressGZIP(file_schema_proto_rawDescData)
	})
	return file_schema_proto_rawDescData
}

var file_schema_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_schema_proto_goTypes = []interface{}{
	(*SchemaRequestArgs)(nil),  // 0: proto.SchemaRequestArgs
	(*Field)(nil),              // 1: proto.Field
	(*Schema)(nil),             // 2: proto.Schema
	(*NamespaceSchema)(nil),    // 3: proto.NamespaceSchema
	(*SchemaRequestReply)(nil), // 4: proto.SchemaRequestReply
	nil,                        // 5: proto.NamespaceSchema.OperationsEntry
}
var file_schema_proto_depIdxs = []int32{
	1, // 0: proto.Field.reference:type_name -> proto.Field
	1, // 1: proto.Schema.input:type_name -> proto.Field
	1, // 2: proto.Schema.output:type_name -> proto.Field
	5, // 3: proto.NamespaceSchema.operations:type_name -> proto.NamespaceSchema.OperationsEntry
	3, // 4: proto.SchemaRequestReply.Namespaces:type_name -> proto.NamespaceSchema
	2, // 5: proto.NamespaceSchema.OperationsEntry.value:type_name -> proto.Schema
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_schema_proto_init() }
func file_schema_proto_init() {
	if File_schema_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_schema_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SchemaRequestArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_schema_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Field); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_schema_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Schema); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_schema_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NamespaceSchema); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_schema_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SchemaRequestReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_schema_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_schema_proto_goTypes,
		DependencyIndexes: file_schema_proto_depIdxs,
		MessageInfos:      file_schema_proto_msgTypes,
	}.Build()
	File_schema_proto = out.File
	file_schema_proto_rawDesc = nil
	file_schema_proto_goTypes = nil
	file_schema_proto_depIdxs = nil
}
