// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: rapid/public_api/v1/authentication.proto

package public_apiv1

import (
	reflect "reflect"
	sync "sync"

	v1 "github.com/abyssparanoia/rapid-go/schema/proto/pb/rapid/model/v1"
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

type SignInRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SignInRequest) Reset() {
	*x = SignInRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rapid_public_api_v1_authentication_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignInRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignInRequest) ProtoMessage() {}

func (x *SignInRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rapid_public_api_v1_authentication_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignInRequest.ProtoReflect.Descriptor instead.
func (*SignInRequest) Descriptor() ([]byte, []int) {
	return file_rapid_public_api_v1_authentication_proto_rawDescGZIP(), []int{0}
}

type SignInResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *v1.User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *SignInResponse) Reset() {
	*x = SignInResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rapid_public_api_v1_authentication_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignInResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignInResponse) ProtoMessage() {}

func (x *SignInResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rapid_public_api_v1_authentication_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignInResponse.ProtoReflect.Descriptor instead.
func (*SignInResponse) Descriptor() ([]byte, []int) {
	return file_rapid_public_api_v1_authentication_proto_rawDescGZIP(), []int{1}
}

func (x *SignInResponse) GetUser() *v1.User {
	if x != nil {
		return x.User
	}
	return nil
}

var File_rapid_public_api_v1_authentication_proto protoreflect.FileDescriptor

var file_rapid_public_api_v1_authentication_proto_rawDesc = []byte{
	0x0a, 0x28, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x72, 0x61, 0x70, 0x69,
	0x64, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a,
	0x19, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x76, 0x31, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32,
	0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0f, 0x0a, 0x0d, 0x53, 0x69,
	0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x48, 0x0a, 0x0e, 0x53,
	0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x72, 0x61,
	0x70, 0x69, 0x64, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x3a, 0x0c, 0x92, 0x41, 0x09, 0x0a, 0x07, 0xd2, 0x01,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x42, 0xec, 0x01, 0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x2e, 0x72, 0x61,
	0x70, 0x69, 0x64, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x42, 0x13, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x52, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x62, 0x79, 0x73, 0x73, 0x70, 0x61, 0x72, 0x61, 0x6e, 0x6f,
	0x69, 0x61, 0x2f, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2d, 0x67, 0x6f, 0x2f, 0x73, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x2f, 0x72, 0x61, 0x70, 0x69,
	0x64, 0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x52,
	0x50, 0x58, 0xaa, 0x02, 0x12, 0x52, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x41, 0x70, 0x69, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x12, 0x52, 0x61, 0x70, 0x69, 0x64, 0x5c,
	0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x41, 0x70, 0x69, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1e, 0x52,
	0x61, 0x70, 0x69, 0x64, 0x5c, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x41, 0x70, 0x69, 0x5c, 0x56,
	0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x14,
	0x52, 0x61, 0x70, 0x69, 0x64, 0x3a, 0x3a, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x41, 0x70, 0x69,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rapid_public_api_v1_authentication_proto_rawDescOnce sync.Once
	file_rapid_public_api_v1_authentication_proto_rawDescData = file_rapid_public_api_v1_authentication_proto_rawDesc
)

func file_rapid_public_api_v1_authentication_proto_rawDescGZIP() []byte {
	file_rapid_public_api_v1_authentication_proto_rawDescOnce.Do(func() {
		file_rapid_public_api_v1_authentication_proto_rawDescData = protoimpl.X.CompressGZIP(file_rapid_public_api_v1_authentication_proto_rawDescData)
	})
	return file_rapid_public_api_v1_authentication_proto_rawDescData
}

var file_rapid_public_api_v1_authentication_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rapid_public_api_v1_authentication_proto_goTypes = []interface{}{
	(*SignInRequest)(nil),  // 0: rapid.public_api.v1.SignInRequest
	(*SignInResponse)(nil), // 1: rapid.public_api.v1.SignInResponse
	(*v1.User)(nil),        // 2: rapid.model.v1.User
}
var file_rapid_public_api_v1_authentication_proto_depIdxs = []int32{
	2, // 0: rapid.public_api.v1.SignInResponse.user:type_name -> rapid.model.v1.User
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rapid_public_api_v1_authentication_proto_init() }
func file_rapid_public_api_v1_authentication_proto_init() {
	if File_rapid_public_api_v1_authentication_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rapid_public_api_v1_authentication_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignInRequest); i {
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
		file_rapid_public_api_v1_authentication_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignInResponse); i {
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
			RawDescriptor: file_rapid_public_api_v1_authentication_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rapid_public_api_v1_authentication_proto_goTypes,
		DependencyIndexes: file_rapid_public_api_v1_authentication_proto_depIdxs,
		MessageInfos:      file_rapid_public_api_v1_authentication_proto_msgTypes,
	}.Build()
	File_rapid_public_api_v1_authentication_proto = out.File
	file_rapid_public_api_v1_authentication_proto_rawDesc = nil
	file_rapid_public_api_v1_authentication_proto_goTypes = nil
	file_rapid_public_api_v1_authentication_proto_depIdxs = nil
}