// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: rapid/admin_api/v1/api.proto

package admin_apiv1

import (
	reflect "reflect"
	unsafe "unsafe"

	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_rapid_admin_api_v1_api_proto protoreflect.FileDescriptor

var file_rapid_admin_api_v1_api_proto_rawDesc = string([]byte{
	0x0a, 0x1c, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12,
	0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x22, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x70, 0x69, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x22, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x70, 0x69, 0x5f, 0x73, 0x74, 0x61,
	0x66, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x70, 0x69,
	0x5f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xd5, 0x07,
	0x0a, 0x0e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x56, 0x31, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0xaf, 0x01, 0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74,
	0x50, 0x72, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x12, 0x32, 0x2e, 0x72,
	0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74, 0x50, 0x72, 0x65,
	0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x33, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x73, 0x73, 0x65,
	0x74, 0x50, 0x72, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x25, 0x3a, 0x01, 0x2a,
	0x22, 0x20, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x73, 0x73, 0x65,
	0x74, 0x73, 0x2f, 0x2d, 0x2f, 0x70, 0x72, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x75,
	0x72, 0x6c, 0x12, 0x7f, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12,
	0x24, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x25, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x1f, 0x12, 0x1d, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f,
	0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x73, 0x2f, 0x7b, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x7d, 0x12, 0x79, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x73, 0x12, 0x26, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x72, 0x61, 0x70,
	0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x12, 0x11, 0x2f, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x73, 0x12, 0x7f,
	0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x27,
	0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x3a, 0x01, 0x2a, 0x22, 0x11, 0x2f, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x73, 0x12,
	0x8b, 0x01, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x12, 0x27, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x72, 0x61, 0x70, 0x69,
	0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x28, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x22, 0x3a, 0x01, 0x2a, 0x32, 0x1d,
	0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x73, 0x2f, 0x7b, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0x88, 0x01,
	0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x27,
	0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x25, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1f, 0x2a, 0x1d, 0x2f, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x73, 0x2f, 0x7b, 0x74, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0x7b, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x53, 0x74, 0x61, 0x66, 0x66, 0x12, 0x26, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x66, 0x66, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x27, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x66, 0x66,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15,
	0x3a, 0x01, 0x2a, 0x22, 0x10, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x73,
	0x74, 0x61, 0x66, 0x66, 0x73, 0x42, 0xea, 0x01, 0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x72, 0x61,
	0x70, 0x69, 0x64, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x42, 0x08, 0x41, 0x70, 0x69, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x60, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x62, 0x79, 0x73, 0x73, 0x70, 0x61,
	0x72, 0x61, 0x6e, 0x6f, 0x69, 0x61, 0x2f, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2d, 0x67, 0x6f, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74,
	0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f,
	0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x31, 0x3b, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x76, 0x31, 0xa2, 0x02,
	0x03, 0x52, 0x41, 0x58, 0xaa, 0x02, 0x11, 0x52, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x41, 0x64, 0x6d,
	0x69, 0x6e, 0x41, 0x70, 0x69, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x11, 0x52, 0x61, 0x70, 0x69, 0x64,
	0x5c, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x41, 0x70, 0x69, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1d, 0x52,
	0x61, 0x70, 0x69, 0x64, 0x5c, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x41, 0x70, 0x69, 0x5c, 0x56, 0x31,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x52,
	0x61, 0x70, 0x69, 0x64, 0x3a, 0x3a, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x41, 0x70, 0x69, 0x3a, 0x3a,
	0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var file_rapid_admin_api_v1_api_proto_goTypes = []any{
	(*CreateAssetPresignedURLRequest)(nil),  // 0: rapid.admin_api.v1.CreateAssetPresignedURLRequest
	(*GetTenantRequest)(nil),                // 1: rapid.admin_api.v1.GetTenantRequest
	(*ListTenantsRequest)(nil),              // 2: rapid.admin_api.v1.ListTenantsRequest
	(*CreateTenantRequest)(nil),             // 3: rapid.admin_api.v1.CreateTenantRequest
	(*UpdateTenantRequest)(nil),             // 4: rapid.admin_api.v1.UpdateTenantRequest
	(*DeleteTenantRequest)(nil),             // 5: rapid.admin_api.v1.DeleteTenantRequest
	(*CreateStaffRequest)(nil),              // 6: rapid.admin_api.v1.CreateStaffRequest
	(*CreateAssetPresignedURLResponse)(nil), // 7: rapid.admin_api.v1.CreateAssetPresignedURLResponse
	(*GetTenantResponse)(nil),               // 8: rapid.admin_api.v1.GetTenantResponse
	(*ListTenantsResponse)(nil),             // 9: rapid.admin_api.v1.ListTenantsResponse
	(*CreateTenantResponse)(nil),            // 10: rapid.admin_api.v1.CreateTenantResponse
	(*UpdateTenantResponse)(nil),            // 11: rapid.admin_api.v1.UpdateTenantResponse
	(*DeleteTenantResponse)(nil),            // 12: rapid.admin_api.v1.DeleteTenantResponse
	(*CreateStaffResponse)(nil),             // 13: rapid.admin_api.v1.CreateStaffResponse
}
var file_rapid_admin_api_v1_api_proto_depIdxs = []int32{
	0,  // 0: rapid.admin_api.v1.AdminV1Service.CreateAssetPresignedURL:input_type -> rapid.admin_api.v1.CreateAssetPresignedURLRequest
	1,  // 1: rapid.admin_api.v1.AdminV1Service.GetTenant:input_type -> rapid.admin_api.v1.GetTenantRequest
	2,  // 2: rapid.admin_api.v1.AdminV1Service.ListTenants:input_type -> rapid.admin_api.v1.ListTenantsRequest
	3,  // 3: rapid.admin_api.v1.AdminV1Service.CreateTenant:input_type -> rapid.admin_api.v1.CreateTenantRequest
	4,  // 4: rapid.admin_api.v1.AdminV1Service.UpdateTenant:input_type -> rapid.admin_api.v1.UpdateTenantRequest
	5,  // 5: rapid.admin_api.v1.AdminV1Service.DeleteTenant:input_type -> rapid.admin_api.v1.DeleteTenantRequest
	6,  // 6: rapid.admin_api.v1.AdminV1Service.CreateStaff:input_type -> rapid.admin_api.v1.CreateStaffRequest
	7,  // 7: rapid.admin_api.v1.AdminV1Service.CreateAssetPresignedURL:output_type -> rapid.admin_api.v1.CreateAssetPresignedURLResponse
	8,  // 8: rapid.admin_api.v1.AdminV1Service.GetTenant:output_type -> rapid.admin_api.v1.GetTenantResponse
	9,  // 9: rapid.admin_api.v1.AdminV1Service.ListTenants:output_type -> rapid.admin_api.v1.ListTenantsResponse
	10, // 10: rapid.admin_api.v1.AdminV1Service.CreateTenant:output_type -> rapid.admin_api.v1.CreateTenantResponse
	11, // 11: rapid.admin_api.v1.AdminV1Service.UpdateTenant:output_type -> rapid.admin_api.v1.UpdateTenantResponse
	12, // 12: rapid.admin_api.v1.AdminV1Service.DeleteTenant:output_type -> rapid.admin_api.v1.DeleteTenantResponse
	13, // 13: rapid.admin_api.v1.AdminV1Service.CreateStaff:output_type -> rapid.admin_api.v1.CreateStaffResponse
	7,  // [7:14] is the sub-list for method output_type
	0,  // [0:7] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_rapid_admin_api_v1_api_proto_init() }
func file_rapid_admin_api_v1_api_proto_init() {
	if File_rapid_admin_api_v1_api_proto != nil {
		return
	}
	file_rapid_admin_api_v1_api_asset_proto_init()
	file_rapid_admin_api_v1_api_staff_proto_init()
	file_rapid_admin_api_v1_api_tenant_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_rapid_admin_api_v1_api_proto_rawDesc), len(file_rapid_admin_api_v1_api_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rapid_admin_api_v1_api_proto_goTypes,
		DependencyIndexes: file_rapid_admin_api_v1_api_proto_depIdxs,
	}.Build()
	File_rapid_admin_api_v1_api_proto = out.File
	file_rapid_admin_api_v1_api_proto_goTypes = nil
	file_rapid_admin_api_v1_api_proto_depIdxs = nil
}
