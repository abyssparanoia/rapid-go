// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: rapid/public_api/v1/tenant.proto

package public_apiv1

import (
	reflect "reflect"
	sync "sync"

	v1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/model/v1"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PublicGetTenantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantId string `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
}

func (x *PublicGetTenantRequest) Reset() {
	*x = PublicGetTenantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rapid_public_api_v1_tenant_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublicGetTenantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublicGetTenantRequest) ProtoMessage() {}

func (x *PublicGetTenantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rapid_public_api_v1_tenant_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublicGetTenantRequest.ProtoReflect.Descriptor instead.
func (*PublicGetTenantRequest) Descriptor() ([]byte, []int) {
	return file_rapid_public_api_v1_tenant_proto_rawDescGZIP(), []int{0}
}

func (x *PublicGetTenantRequest) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

type PublicGetTenantResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tenant *v1.Tenant `protobuf:"bytes,1,opt,name=tenant,proto3" json:"tenant,omitempty"`
}

func (x *PublicGetTenantResponse) Reset() {
	*x = PublicGetTenantResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rapid_public_api_v1_tenant_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublicGetTenantResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublicGetTenantResponse) ProtoMessage() {}

func (x *PublicGetTenantResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rapid_public_api_v1_tenant_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublicGetTenantResponse.ProtoReflect.Descriptor instead.
func (*PublicGetTenantResponse) Descriptor() ([]byte, []int) {
	return file_rapid_public_api_v1_tenant_proto_rawDescGZIP(), []int{1}
}

func (x *PublicGetTenantResponse) GetTenant() *v1.Tenant {
	if x != nil {
		return x.Tenant
	}
	return nil
}

var File_rapid_public_api_v1_tenant_proto protoreflect.FileDescriptor

var file_rapid_public_api_v1_tenant_proto_rawDesc = []byte{
	0x0a, 0x20, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x13, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e,
	0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x48, 0x0a, 0x16, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x47, 0x65,
	0x74, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b,
	0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x3a, 0x11, 0x92, 0x41, 0x0e,
	0x0a, 0x0c, 0xd2, 0x01, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x22, 0x59,
	0x0a, 0x17, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x06, 0x74, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x72, 0x61, 0x70, 0x69,
	0x64, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x52, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x3a, 0x0e, 0x92, 0x41, 0x0b, 0x0a, 0x09,
	0xd2, 0x01, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x42, 0xf4, 0x01, 0x0a, 0x17, 0x63, 0x6f,
	0x6d, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x42, 0x0b, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x62, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x61, 0x62, 0x79, 0x73, 0x73, 0x70, 0x61, 0x72, 0x61, 0x6e, 0x6f, 0x69, 0x61, 0x2f, 0x72,
	0x61, 0x70, 0x69, 0x64, 0x2d, 0x67, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x52, 0x50, 0x58, 0xaa, 0x02,
	0x12, 0x52, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x41, 0x70, 0x69,
	0x2e, 0x56, 0x31, 0xca, 0x02, 0x12, 0x52, 0x61, 0x70, 0x69, 0x64, 0x5c, 0x50, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x41, 0x70, 0x69, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1e, 0x52, 0x61, 0x70, 0x69, 0x64,
	0x5c, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x41, 0x70, 0x69, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x14, 0x52, 0x61, 0x70, 0x69,
	0x64, 0x3a, 0x3a, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x41, 0x70, 0x69, 0x3a, 0x3a, 0x56, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rapid_public_api_v1_tenant_proto_rawDescOnce sync.Once
	file_rapid_public_api_v1_tenant_proto_rawDescData = file_rapid_public_api_v1_tenant_proto_rawDesc
)

func file_rapid_public_api_v1_tenant_proto_rawDescGZIP() []byte {
	file_rapid_public_api_v1_tenant_proto_rawDescOnce.Do(func() {
		file_rapid_public_api_v1_tenant_proto_rawDescData = protoimpl.X.CompressGZIP(file_rapid_public_api_v1_tenant_proto_rawDescData)
	})
	return file_rapid_public_api_v1_tenant_proto_rawDescData
}

var file_rapid_public_api_v1_tenant_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rapid_public_api_v1_tenant_proto_goTypes = []interface{}{
	(*PublicGetTenantRequest)(nil),  // 0: rapid.public_api.v1.PublicGetTenantRequest
	(*PublicGetTenantResponse)(nil), // 1: rapid.public_api.v1.PublicGetTenantResponse
	(*v1.Tenant)(nil),               // 2: rapid.model.v1.Tenant
}
var file_rapid_public_api_v1_tenant_proto_depIdxs = []int32{
	2, // 0: rapid.public_api.v1.PublicGetTenantResponse.tenant:type_name -> rapid.model.v1.Tenant
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rapid_public_api_v1_tenant_proto_init() }
func file_rapid_public_api_v1_tenant_proto_init() {
	if File_rapid_public_api_v1_tenant_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rapid_public_api_v1_tenant_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublicGetTenantRequest); i {
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
		file_rapid_public_api_v1_tenant_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublicGetTenantResponse); i {
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
			RawDescriptor: file_rapid_public_api_v1_tenant_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rapid_public_api_v1_tenant_proto_goTypes,
		DependencyIndexes: file_rapid_public_api_v1_tenant_proto_depIdxs,
		MessageInfos:      file_rapid_public_api_v1_tenant_proto_msgTypes,
	}.Build()
	File_rapid_public_api_v1_tenant_proto = out.File
	file_rapid_public_api_v1_tenant_proto_rawDesc = nil
	file_rapid_public_api_v1_tenant_proto_goTypes = nil
	file_rapid_public_api_v1_tenant_proto_depIdxs = nil
}
