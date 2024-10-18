// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: rapid/admin_api/v1/api_asset.proto

package admin_apiv1

import (
	reflect "reflect"
	sync "sync"

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

type AssetType int32

const (
	AssetType_ASSET_TYPE_UNSPECIFIED AssetType = 0
	AssetType_ASSET_TYPE_USER_IMAGE  AssetType = 1
)

// Enum value maps for AssetType.
var (
	AssetType_name = map[int32]string{
		0: "ASSET_TYPE_UNSPECIFIED",
		1: "ASSET_TYPE_USER_IMAGE",
	}
	AssetType_value = map[string]int32{
		"ASSET_TYPE_UNSPECIFIED": 0,
		"ASSET_TYPE_USER_IMAGE":  1,
	}
)

func (x AssetType) Enum() *AssetType {
	p := new(AssetType)
	*p = x
	return p
}

func (x AssetType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AssetType) Descriptor() protoreflect.EnumDescriptor {
	return file_rapid_admin_api_v1_api_asset_proto_enumTypes[0].Descriptor()
}

func (AssetType) Type() protoreflect.EnumType {
	return &file_rapid_admin_api_v1_api_asset_proto_enumTypes[0]
}

func (x AssetType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AssetType.Descriptor instead.
func (AssetType) EnumDescriptor() ([]byte, []int) {
	return file_rapid_admin_api_v1_api_asset_proto_rawDescGZIP(), []int{0}
}

type CreateAssetPresignedURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AssetType   AssetType `protobuf:"varint,1,opt,name=asset_type,json=assetType,proto3,enum=rapid.admin_api.v1.AssetType" json:"asset_type,omitempty"`
	ContentType string    `protobuf:"bytes,2,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
}

func (x *CreateAssetPresignedURLRequest) Reset() {
	*x = CreateAssetPresignedURLRequest{}
	mi := &file_rapid_admin_api_v1_api_asset_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateAssetPresignedURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAssetPresignedURLRequest) ProtoMessage() {}

func (x *CreateAssetPresignedURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rapid_admin_api_v1_api_asset_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAssetPresignedURLRequest.ProtoReflect.Descriptor instead.
func (*CreateAssetPresignedURLRequest) Descriptor() ([]byte, []int) {
	return file_rapid_admin_api_v1_api_asset_proto_rawDescGZIP(), []int{0}
}

func (x *CreateAssetPresignedURLRequest) GetAssetType() AssetType {
	if x != nil {
		return x.AssetType
	}
	return AssetType_ASSET_TYPE_UNSPECIFIED
}

func (x *CreateAssetPresignedURLRequest) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

type CreateAssetPresignedURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AssetId      string `protobuf:"bytes,1,opt,name=asset_id,json=assetId,proto3" json:"asset_id,omitempty"`
	PresignedUrl string `protobuf:"bytes,2,opt,name=presigned_url,json=presignedUrl,proto3" json:"presigned_url,omitempty"`
}

func (x *CreateAssetPresignedURLResponse) Reset() {
	*x = CreateAssetPresignedURLResponse{}
	mi := &file_rapid_admin_api_v1_api_asset_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateAssetPresignedURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAssetPresignedURLResponse) ProtoMessage() {}

func (x *CreateAssetPresignedURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rapid_admin_api_v1_api_asset_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAssetPresignedURLResponse.ProtoReflect.Descriptor instead.
func (*CreateAssetPresignedURLResponse) Descriptor() ([]byte, []int) {
	return file_rapid_admin_api_v1_api_asset_proto_rawDescGZIP(), []int{1}
}

func (x *CreateAssetPresignedURLResponse) GetAssetId() string {
	if x != nil {
		return x.AssetId
	}
	return ""
}

func (x *CreateAssetPresignedURLResponse) GetPresignedUrl() string {
	if x != nil {
		return x.PresignedUrl
	}
	return ""
}

var File_rapid_admin_api_v1_api_asset_proto protoreflect.FileDescriptor

var file_rapid_admin_api_v1_api_asset_proto_rawDesc = []byte{
	0x0a, 0x22, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x70, 0x69, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa4, 0x01, 0x0a, 0x1e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74, 0x50, 0x72, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3c, 0x0a, 0x0a, 0x61,
	0x73, 0x73, 0x65, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1d, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09,
	0x61, 0x73, 0x73, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x3a, 0x21, 0x92, 0x41,
	0x1e, 0x0a, 0x1c, 0xd2, 0x01, 0x0a, 0x61, 0x73, 0x73, 0x65, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0xd2, 0x01, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x22,
	0x83, 0x01, 0x0a, 0x1f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74, 0x50,
	0x72, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x73, 0x73, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x73, 0x73, 0x65, 0x74, 0x49, 0x64, 0x12, 0x23,
	0x0a, 0x0d, 0x70, 0x72, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x75, 0x72, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64,
	0x55, 0x72, 0x6c, 0x3a, 0x20, 0x92, 0x41, 0x1d, 0x0a, 0x1b, 0xd2, 0x01, 0x08, 0x61, 0x73, 0x73,
	0x65, 0x74, 0x5f, 0x69, 0x64, 0xd2, 0x01, 0x0d, 0x70, 0x72, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x5f, 0x75, 0x72, 0x6c, 0x2a, 0x42, 0x0a, 0x09, 0x41, 0x73, 0x73, 0x65, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x41, 0x53, 0x53, 0x45, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x19,
	0x0a, 0x15, 0x41, 0x53, 0x53, 0x45, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x53, 0x45,
	0x52, 0x5f, 0x49, 0x4d, 0x41, 0x47, 0x45, 0x10, 0x01, 0x42, 0xef, 0x01, 0x0a, 0x16, 0x63, 0x6f,
	0x6d, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x42, 0x0d, 0x41, 0x70, 0x69, 0x41, 0x73, 0x73, 0x65, 0x74, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x60, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x61, 0x62, 0x79, 0x73, 0x73, 0x70, 0x61, 0x72, 0x61, 0x6e, 0x6f, 0x69, 0x61, 0x2f,
	0x72, 0x61, 0x70, 0x69, 0x64, 0x2d, 0x67, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65,
	0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x5f, 0x61, 0x70, 0x69, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x52, 0x41, 0x58, 0xaa, 0x02, 0x11,
	0x52, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x41, 0x70, 0x69, 0x2e, 0x56,
	0x31, 0xca, 0x02, 0x11, 0x52, 0x61, 0x70, 0x69, 0x64, 0x5c, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x41,
	0x70, 0x69, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1d, 0x52, 0x61, 0x70, 0x69, 0x64, 0x5c, 0x41, 0x64,
	0x6d, 0x69, 0x6e, 0x41, 0x70, 0x69, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x52, 0x61, 0x70, 0x69, 0x64, 0x3a, 0x3a, 0x41,
	0x64, 0x6d, 0x69, 0x6e, 0x41, 0x70, 0x69, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_rapid_admin_api_v1_api_asset_proto_rawDescOnce sync.Once
	file_rapid_admin_api_v1_api_asset_proto_rawDescData = file_rapid_admin_api_v1_api_asset_proto_rawDesc
)

func file_rapid_admin_api_v1_api_asset_proto_rawDescGZIP() []byte {
	file_rapid_admin_api_v1_api_asset_proto_rawDescOnce.Do(func() {
		file_rapid_admin_api_v1_api_asset_proto_rawDescData = protoimpl.X.CompressGZIP(file_rapid_admin_api_v1_api_asset_proto_rawDescData)
	})
	return file_rapid_admin_api_v1_api_asset_proto_rawDescData
}

var file_rapid_admin_api_v1_api_asset_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_rapid_admin_api_v1_api_asset_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rapid_admin_api_v1_api_asset_proto_goTypes = []any{
	(AssetType)(0),                          // 0: rapid.admin_api.v1.AssetType
	(*CreateAssetPresignedURLRequest)(nil),  // 1: rapid.admin_api.v1.CreateAssetPresignedURLRequest
	(*CreateAssetPresignedURLResponse)(nil), // 2: rapid.admin_api.v1.CreateAssetPresignedURLResponse
}
var file_rapid_admin_api_v1_api_asset_proto_depIdxs = []int32{
	0, // 0: rapid.admin_api.v1.CreateAssetPresignedURLRequest.asset_type:type_name -> rapid.admin_api.v1.AssetType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rapid_admin_api_v1_api_asset_proto_init() }
func file_rapid_admin_api_v1_api_asset_proto_init() {
	if File_rapid_admin_api_v1_api_asset_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rapid_admin_api_v1_api_asset_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rapid_admin_api_v1_api_asset_proto_goTypes,
		DependencyIndexes: file_rapid_admin_api_v1_api_asset_proto_depIdxs,
		EnumInfos:         file_rapid_admin_api_v1_api_asset_proto_enumTypes,
		MessageInfos:      file_rapid_admin_api_v1_api_asset_proto_msgTypes,
	}.Build()
	File_rapid_admin_api_v1_api_asset_proto = out.File
	file_rapid_admin_api_v1_api_asset_proto_rawDesc = nil
	file_rapid_admin_api_v1_api_asset_proto_goTypes = nil
	file_rapid_admin_api_v1_api_asset_proto_depIdxs = nil
}
