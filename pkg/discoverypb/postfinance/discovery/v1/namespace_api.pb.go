// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.2
// source: postfinance/discovery/v1/namespace_api.proto

package discoveryv1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type RegisterNamespaceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Export int32  `protobuf:"varint,2,opt,name=export,proto3" json:"export,omitempty"`
}

func (x *RegisterNamespaceRequest) Reset() {
	*x = RegisterNamespaceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_postfinance_discovery_v1_namespace_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterNamespaceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterNamespaceRequest) ProtoMessage() {}

func (x *RegisterNamespaceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_postfinance_discovery_v1_namespace_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterNamespaceRequest.ProtoReflect.Descriptor instead.
func (*RegisterNamespaceRequest) Descriptor() ([]byte, []int) {
	return file_postfinance_discovery_v1_namespace_api_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterNamespaceRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RegisterNamespaceRequest) GetExport() int32 {
	if x != nil {
		return x.Export
	}
	return 0
}

type RegisterNamespaceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace *Namespace `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
}

func (x *RegisterNamespaceResponse) Reset() {
	*x = RegisterNamespaceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_postfinance_discovery_v1_namespace_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterNamespaceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterNamespaceResponse) ProtoMessage() {}

func (x *RegisterNamespaceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_postfinance_discovery_v1_namespace_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterNamespaceResponse.ProtoReflect.Descriptor instead.
func (*RegisterNamespaceResponse) Descriptor() ([]byte, []int) {
	return file_postfinance_discovery_v1_namespace_api_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterNamespaceResponse) GetNamespace() *Namespace {
	if x != nil {
		return x.Namespace
	}
	return nil
}

type UnregisterNamespaceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *UnregisterNamespaceRequest) Reset() {
	*x = UnregisterNamespaceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_postfinance_discovery_v1_namespace_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnregisterNamespaceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnregisterNamespaceRequest) ProtoMessage() {}

func (x *UnregisterNamespaceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_postfinance_discovery_v1_namespace_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnregisterNamespaceRequest.ProtoReflect.Descriptor instead.
func (*UnregisterNamespaceRequest) Descriptor() ([]byte, []int) {
	return file_postfinance_discovery_v1_namespace_api_proto_rawDescGZIP(), []int{2}
}

func (x *UnregisterNamespaceRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type UnregisterNamespaceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UnregisterNamespaceResponse) Reset() {
	*x = UnregisterNamespaceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_postfinance_discovery_v1_namespace_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnregisterNamespaceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnregisterNamespaceResponse) ProtoMessage() {}

func (x *UnregisterNamespaceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_postfinance_discovery_v1_namespace_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnregisterNamespaceResponse.ProtoReflect.Descriptor instead.
func (*UnregisterNamespaceResponse) Descriptor() ([]byte, []int) {
	return file_postfinance_discovery_v1_namespace_api_proto_rawDescGZIP(), []int{3}
}

type ListNamespaceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListNamespaceRequest) Reset() {
	*x = ListNamespaceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_postfinance_discovery_v1_namespace_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListNamespaceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListNamespaceRequest) ProtoMessage() {}

func (x *ListNamespaceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_postfinance_discovery_v1_namespace_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListNamespaceRequest.ProtoReflect.Descriptor instead.
func (*ListNamespaceRequest) Descriptor() ([]byte, []int) {
	return file_postfinance_discovery_v1_namespace_api_proto_rawDescGZIP(), []int{4}
}

type ListNamespaceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespaces []*Namespace `protobuf:"bytes,1,rep,name=namespaces,proto3" json:"namespaces,omitempty"`
}

func (x *ListNamespaceResponse) Reset() {
	*x = ListNamespaceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_postfinance_discovery_v1_namespace_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListNamespaceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListNamespaceResponse) ProtoMessage() {}

func (x *ListNamespaceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_postfinance_discovery_v1_namespace_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListNamespaceResponse.ProtoReflect.Descriptor instead.
func (*ListNamespaceResponse) Descriptor() ([]byte, []int) {
	return file_postfinance_discovery_v1_namespace_api_proto_rawDescGZIP(), []int{5}
}

func (x *ListNamespaceResponse) GetNamespaces() []*Namespace {
	if x != nil {
		return x.Namespaces
	}
	return nil
}

var File_postfinance_discovery_v1_namespace_api_proto protoreflect.FileDescriptor

var file_postfinance_discovery_v1_namespace_api_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x70, 0x6f, 0x73, 0x74, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2f, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18,
	0x70, 0x6f, 0x73, 0x74, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x1a, 0x28, 0x70, 0x6f, 0x73, 0x74, 0x66, 0x69,
	0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f,
	0x76, 0x31, 0x2f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x46, 0x0a, 0x18, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x5e, 0x0a, 0x19, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x70, 0x6f, 0x73, 0x74, 0x66,
	0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x09, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x22, 0x30, 0x0a, 0x1a, 0x55, 0x6e, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x1d, 0x0a, 0x1b, 0x55, 0x6e,
	0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x0a, 0x14, 0x4c, 0x69, 0x73,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x5c, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0a, 0x6e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23,
	0x2e, 0x70, 0x6f, 0x73, 0x74, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x64, 0x69, 0x73,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x52, 0x0a, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x32,
	0xd7, 0x03, 0x0a, 0x0c, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x41, 0x50, 0x49,
	0x12, 0x97, 0x01, 0x0a, 0x11, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x32, 0x2e, 0x70, 0x6f, 0x73, 0x74, 0x66, 0x69, 0x6e,
	0x61, 0x6e, 0x63, 0x65, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x70, 0x6f, 0x73,
	0x74, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x22, 0x0e, 0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x3a, 0x01, 0x2a, 0x12, 0xa1, 0x01, 0x0a, 0x13, 0x55,
	0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x12, 0x34, 0x2e, 0x70, 0x6f, 0x73, 0x74, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65,
	0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x6e,
	0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x35, 0x2e, 0x70, 0x6f, 0x73, 0x74, 0x66,
	0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x2e, 0x76, 0x31, 0x2e, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x2a, 0x15, 0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x12, 0x88,
	0x01, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65,
	0x12, 0x2e, 0x2e, 0x70, 0x6f, 0x73, 0x74, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2f, 0x2e, 0x70, 0x6f, 0x73, 0x74, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e, 0x2f, 0x76, 0x31, 0x2f, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x42, 0x58, 0x0a, 0x1b, 0x63, 0x68, 0x2e,
	0x70, 0x6f, 0x73, 0x74, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x42, 0x11, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x41, 0x70, 0x69, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x24, 0x70,
	0x6f, 0x73, 0x74, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x3b, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_postfinance_discovery_v1_namespace_api_proto_rawDescOnce sync.Once
	file_postfinance_discovery_v1_namespace_api_proto_rawDescData = file_postfinance_discovery_v1_namespace_api_proto_rawDesc
)

func file_postfinance_discovery_v1_namespace_api_proto_rawDescGZIP() []byte {
	file_postfinance_discovery_v1_namespace_api_proto_rawDescOnce.Do(func() {
		file_postfinance_discovery_v1_namespace_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_postfinance_discovery_v1_namespace_api_proto_rawDescData)
	})
	return file_postfinance_discovery_v1_namespace_api_proto_rawDescData
}

var file_postfinance_discovery_v1_namespace_api_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_postfinance_discovery_v1_namespace_api_proto_goTypes = []interface{}{
	(*RegisterNamespaceRequest)(nil),    // 0: postfinance.discovery.v1.RegisterNamespaceRequest
	(*RegisterNamespaceResponse)(nil),   // 1: postfinance.discovery.v1.RegisterNamespaceResponse
	(*UnregisterNamespaceRequest)(nil),  // 2: postfinance.discovery.v1.UnregisterNamespaceRequest
	(*UnregisterNamespaceResponse)(nil), // 3: postfinance.discovery.v1.UnregisterNamespaceResponse
	(*ListNamespaceRequest)(nil),        // 4: postfinance.discovery.v1.ListNamespaceRequest
	(*ListNamespaceResponse)(nil),       // 5: postfinance.discovery.v1.ListNamespaceResponse
	(*Namespace)(nil),                   // 6: postfinance.discovery.v1.Namespace
}
var file_postfinance_discovery_v1_namespace_api_proto_depIdxs = []int32{
	6, // 0: postfinance.discovery.v1.RegisterNamespaceResponse.namespace:type_name -> postfinance.discovery.v1.Namespace
	6, // 1: postfinance.discovery.v1.ListNamespaceResponse.namespaces:type_name -> postfinance.discovery.v1.Namespace
	0, // 2: postfinance.discovery.v1.NamespaceAPI.RegisterNamespace:input_type -> postfinance.discovery.v1.RegisterNamespaceRequest
	2, // 3: postfinance.discovery.v1.NamespaceAPI.UnregisterNamespace:input_type -> postfinance.discovery.v1.UnregisterNamespaceRequest
	4, // 4: postfinance.discovery.v1.NamespaceAPI.ListNamespace:input_type -> postfinance.discovery.v1.ListNamespaceRequest
	1, // 5: postfinance.discovery.v1.NamespaceAPI.RegisterNamespace:output_type -> postfinance.discovery.v1.RegisterNamespaceResponse
	3, // 6: postfinance.discovery.v1.NamespaceAPI.UnregisterNamespace:output_type -> postfinance.discovery.v1.UnregisterNamespaceResponse
	5, // 7: postfinance.discovery.v1.NamespaceAPI.ListNamespace:output_type -> postfinance.discovery.v1.ListNamespaceResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_postfinance_discovery_v1_namespace_api_proto_init() }
func file_postfinance_discovery_v1_namespace_api_proto_init() {
	if File_postfinance_discovery_v1_namespace_api_proto != nil {
		return
	}
	file_postfinance_discovery_v1_namespace_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_postfinance_discovery_v1_namespace_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterNamespaceRequest); i {
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
		file_postfinance_discovery_v1_namespace_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterNamespaceResponse); i {
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
		file_postfinance_discovery_v1_namespace_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnregisterNamespaceRequest); i {
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
		file_postfinance_discovery_v1_namespace_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnregisterNamespaceResponse); i {
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
		file_postfinance_discovery_v1_namespace_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListNamespaceRequest); i {
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
		file_postfinance_discovery_v1_namespace_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListNamespaceResponse); i {
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
			RawDescriptor: file_postfinance_discovery_v1_namespace_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_postfinance_discovery_v1_namespace_api_proto_goTypes,
		DependencyIndexes: file_postfinance_discovery_v1_namespace_api_proto_depIdxs,
		MessageInfos:      file_postfinance_discovery_v1_namespace_api_proto_msgTypes,
	}.Build()
	File_postfinance_discovery_v1_namespace_api_proto = out.File
	file_postfinance_discovery_v1_namespace_api_proto_rawDesc = nil
	file_postfinance_discovery_v1_namespace_api_proto_goTypes = nil
	file_postfinance_discovery_v1_namespace_api_proto_depIdxs = nil
}
