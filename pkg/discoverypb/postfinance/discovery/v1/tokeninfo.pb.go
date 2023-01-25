// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: postfinance/discovery/v1/tokeninfo.proto

package discoveryv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// TokenInfo represents a machine token.
type TokenInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id is the id of the token.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// namespaces defines which namespaces the token has access to.
	Namespaces []string `protobuf:"bytes,2,rep,name=namespaces,proto3" json:"namespaces,omitempty"`
	// expires_at shows the expiry time
	ExpiresAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
}

func (x *TokenInfo) Reset() {
	*x = TokenInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_postfinance_discovery_v1_tokeninfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenInfo) ProtoMessage() {}

func (x *TokenInfo) ProtoReflect() protoreflect.Message {
	mi := &file_postfinance_discovery_v1_tokeninfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenInfo.ProtoReflect.Descriptor instead.
func (*TokenInfo) Descriptor() ([]byte, []int) {
	return file_postfinance_discovery_v1_tokeninfo_proto_rawDescGZIP(), []int{0}
}

func (x *TokenInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TokenInfo) GetNamespaces() []string {
	if x != nil {
		return x.Namespaces
	}
	return nil
}

func (x *TokenInfo) GetExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiresAt
	}
	return nil
}

var File_postfinance_discovery_v1_tokeninfo_proto protoreflect.FileDescriptor

var file_postfinance_discovery_v1_tokeninfo_proto_rawDesc = []byte{
	0x0a, 0x28, 0x70, 0x6f, 0x73, 0x74, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2f, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x70, 0x6f, 0x73, 0x74,
	0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x76, 0x0a, 0x09, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x61, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74, 0x42, 0x55, 0x0a,
	0x1b, 0x63, 0x68, 0x2e, 0x70, 0x6f, 0x73, 0x74, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2e,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x42, 0x0e, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x69, 0x6e, 0x66, 0x6f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x24,
	0x70, 0x6f, 0x73, 0x74, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2f, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x3b, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_postfinance_discovery_v1_tokeninfo_proto_rawDescOnce sync.Once
	file_postfinance_discovery_v1_tokeninfo_proto_rawDescData = file_postfinance_discovery_v1_tokeninfo_proto_rawDesc
)

func file_postfinance_discovery_v1_tokeninfo_proto_rawDescGZIP() []byte {
	file_postfinance_discovery_v1_tokeninfo_proto_rawDescOnce.Do(func() {
		file_postfinance_discovery_v1_tokeninfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_postfinance_discovery_v1_tokeninfo_proto_rawDescData)
	})
	return file_postfinance_discovery_v1_tokeninfo_proto_rawDescData
}

var file_postfinance_discovery_v1_tokeninfo_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_postfinance_discovery_v1_tokeninfo_proto_goTypes = []interface{}{
	(*TokenInfo)(nil),             // 0: postfinance.discovery.v1.TokenInfo
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_postfinance_discovery_v1_tokeninfo_proto_depIdxs = []int32{
	1, // 0: postfinance.discovery.v1.TokenInfo.expires_at:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_postfinance_discovery_v1_tokeninfo_proto_init() }
func file_postfinance_discovery_v1_tokeninfo_proto_init() {
	if File_postfinance_discovery_v1_tokeninfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_postfinance_discovery_v1_tokeninfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenInfo); i {
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
			RawDescriptor: file_postfinance_discovery_v1_tokeninfo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_postfinance_discovery_v1_tokeninfo_proto_goTypes,
		DependencyIndexes: file_postfinance_discovery_v1_tokeninfo_proto_depIdxs,
		MessageInfos:      file_postfinance_discovery_v1_tokeninfo_proto_msgTypes,
	}.Build()
	File_postfinance_discovery_v1_tokeninfo_proto = out.File
	file_postfinance_discovery_v1_tokeninfo_proto_rawDesc = nil
	file_postfinance_discovery_v1_tokeninfo_proto_goTypes = nil
	file_postfinance_discovery_v1_tokeninfo_proto_depIdxs = nil
}
