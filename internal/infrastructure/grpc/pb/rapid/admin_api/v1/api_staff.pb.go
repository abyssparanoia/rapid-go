// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        (unknown)
// source: rapid/admin_api/v1/api_staff.proto

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

type CreateStaffRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TenantId      string                 `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	Email         string                 `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	DisplayName   string                 `protobuf:"bytes,3,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	Role          StaffRole              `protobuf:"varint,4,opt,name=role,proto3,enum=rapid.admin_api.v1.StaffRole" json:"role,omitempty"`
	ImageAssetId  string                 `protobuf:"bytes,5,opt,name=image_asset_id,json=imageAssetId,proto3" json:"image_asset_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateStaffRequest) Reset() {
	*x = CreateStaffRequest{}
	mi := &file_rapid_admin_api_v1_api_staff_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateStaffRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateStaffRequest) ProtoMessage() {}

func (x *CreateStaffRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rapid_admin_api_v1_api_staff_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateStaffRequest.ProtoReflect.Descriptor instead.
func (*CreateStaffRequest) Descriptor() ([]byte, []int) {
	return file_rapid_admin_api_v1_api_staff_proto_rawDescGZIP(), []int{0}
}

func (x *CreateStaffRequest) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *CreateStaffRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *CreateStaffRequest) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *CreateStaffRequest) GetRole() StaffRole {
	if x != nil {
		return x.Role
	}
	return StaffRole_STAFF_ROLE_UNSPECIFIED
}

func (x *CreateStaffRequest) GetImageAssetId() string {
	if x != nil {
		return x.ImageAssetId
	}
	return ""
}

type CreateStaffResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Staff         *Staff                 `protobuf:"bytes,1,opt,name=staff,proto3" json:"staff,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateStaffResponse) Reset() {
	*x = CreateStaffResponse{}
	mi := &file_rapid_admin_api_v1_api_staff_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateStaffResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateStaffResponse) ProtoMessage() {}

func (x *CreateStaffResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rapid_admin_api_v1_api_staff_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateStaffResponse.ProtoReflect.Descriptor instead.
func (*CreateStaffResponse) Descriptor() ([]byte, []int) {
	return file_rapid_admin_api_v1_api_staff_proto_rawDescGZIP(), []int{1}
}

func (x *CreateStaffResponse) GetStaff() *Staff {
	if x != nil {
		return x.Staff
	}
	return nil
}

var File_rapid_admin_api_v1_api_staff_proto protoreflect.FileDescriptor

var file_rapid_admin_api_v1_api_staff_proto_rawDesc = []byte{
	0x0a, 0x22, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x70, 0x69, 0x5f, 0x73, 0x74, 0x61, 0x66, 0x66, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x5f, 0x73, 0x74, 0x61, 0x66, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x85,
	0x02, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x66, 0x66, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x69, 0x73, 0x70,
	0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x31, 0x0a, 0x04, 0x72,
	0x6f, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x72, 0x61, 0x70, 0x69,
	0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x74, 0x61, 0x66, 0x66, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x24,
	0x0a, 0x0e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x41, 0x73, 0x73,
	0x65, 0x74, 0x49, 0x64, 0x3a, 0x40, 0x92, 0x41, 0x3d, 0x0a, 0x3b, 0xd2, 0x01, 0x09, 0x74, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0xd2, 0x01, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0xd2,
	0x01, 0x0c, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0xd2, 0x01,
	0x04, 0x72, 0x6f, 0x6c, 0x65, 0xd2, 0x01, 0x0e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x61, 0x73,
	0x73, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x22, 0x55, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x53, 0x74, 0x61, 0x66, 0x66, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x66, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x72,
	0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x74, 0x61, 0x66, 0x66, 0x52, 0x05, 0x73, 0x74, 0x61, 0x66, 0x66, 0x3a, 0x0d,
	0x92, 0x41, 0x0a, 0x0a, 0x08, 0xd2, 0x01, 0x05, 0x73, 0x74, 0x61, 0x66, 0x66, 0x42, 0xef, 0x01,
	0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x42, 0x0d, 0x41, 0x70, 0x69, 0x53, 0x74, 0x61,
	0x66, 0x66, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x60, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x62, 0x79, 0x73, 0x73, 0x70, 0x61, 0x72, 0x61, 0x6e,
	0x6f, 0x69, 0x61, 0x2f, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2d, 0x67, 0x6f, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75, 0x63,
	0x74, 0x75, 0x72, 0x65, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x72, 0x61, 0x70,
	0x69, 0x64, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x52, 0x41,
	0x58, 0xaa, 0x02, 0x11, 0x52, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x41,
	0x70, 0x69, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x11, 0x52, 0x61, 0x70, 0x69, 0x64, 0x5c, 0x41, 0x64,
	0x6d, 0x69, 0x6e, 0x41, 0x70, 0x69, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1d, 0x52, 0x61, 0x70, 0x69,
	0x64, 0x5c, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x41, 0x70, 0x69, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x52, 0x61, 0x70, 0x69,
	0x64, 0x3a, 0x3a, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x41, 0x70, 0x69, 0x3a, 0x3a, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rapid_admin_api_v1_api_staff_proto_rawDescOnce sync.Once
	file_rapid_admin_api_v1_api_staff_proto_rawDescData = file_rapid_admin_api_v1_api_staff_proto_rawDesc
)

func file_rapid_admin_api_v1_api_staff_proto_rawDescGZIP() []byte {
	file_rapid_admin_api_v1_api_staff_proto_rawDescOnce.Do(func() {
		file_rapid_admin_api_v1_api_staff_proto_rawDescData = protoimpl.X.CompressGZIP(file_rapid_admin_api_v1_api_staff_proto_rawDescData)
	})
	return file_rapid_admin_api_v1_api_staff_proto_rawDescData
}

var file_rapid_admin_api_v1_api_staff_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rapid_admin_api_v1_api_staff_proto_goTypes = []any{
	(*CreateStaffRequest)(nil),  // 0: rapid.admin_api.v1.CreateStaffRequest
	(*CreateStaffResponse)(nil), // 1: rapid.admin_api.v1.CreateStaffResponse
	(StaffRole)(0),              // 2: rapid.admin_api.v1.StaffRole
	(*Staff)(nil),               // 3: rapid.admin_api.v1.Staff
}
var file_rapid_admin_api_v1_api_staff_proto_depIdxs = []int32{
	2, // 0: rapid.admin_api.v1.CreateStaffRequest.role:type_name -> rapid.admin_api.v1.StaffRole
	3, // 1: rapid.admin_api.v1.CreateStaffResponse.staff:type_name -> rapid.admin_api.v1.Staff
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_rapid_admin_api_v1_api_staff_proto_init() }
func file_rapid_admin_api_v1_api_staff_proto_init() {
	if File_rapid_admin_api_v1_api_staff_proto != nil {
		return
	}
	file_rapid_admin_api_v1_model_staff_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rapid_admin_api_v1_api_staff_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rapid_admin_api_v1_api_staff_proto_goTypes,
		DependencyIndexes: file_rapid_admin_api_v1_api_staff_proto_depIdxs,
		MessageInfos:      file_rapid_admin_api_v1_api_staff_proto_msgTypes,
	}.Build()
	File_rapid_admin_api_v1_api_staff_proto = out.File
	file_rapid_admin_api_v1_api_staff_proto_rawDesc = nil
	file_rapid_admin_api_v1_api_staff_proto_goTypes = nil
	file_rapid_admin_api_v1_api_staff_proto_depIdxs = nil
}
