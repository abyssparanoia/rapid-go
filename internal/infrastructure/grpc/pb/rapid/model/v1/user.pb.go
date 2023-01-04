// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: rapid/model/v1/user.proto

package modelv1

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UserRole int32

const (
	UserRole_USER_ROLE_UNSPECIFIED UserRole = 0
	UserRole_USER_ROLE_NORMAL      UserRole = 1
	UserRole_USER_ROLE_ADMIN       UserRole = 2
)

// Enum value maps for UserRole.
var (
	UserRole_name = map[int32]string{
		0: "USER_ROLE_UNSPECIFIED",
		1: "USER_ROLE_NORMAL",
		2: "USER_ROLE_ADMIN",
	}
	UserRole_value = map[string]int32{
		"USER_ROLE_UNSPECIFIED": 0,
		"USER_ROLE_NORMAL":      1,
		"USER_ROLE_ADMIN":       2,
	}
)

func (x UserRole) Enum() *UserRole {
	p := new(UserRole)
	*p = x
	return p
}

func (x UserRole) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserRole) Descriptor() protoreflect.EnumDescriptor {
	return file_rapid_model_v1_user_proto_enumTypes[0].Descriptor()
}

func (UserRole) Type() protoreflect.EnumType {
	return &file_rapid_model_v1_user_proto_enumTypes[0]
}

func (x UserRole) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserRole.Descriptor instead.
func (UserRole) EnumDescriptor() ([]byte, []int) {
	return file_rapid_model_v1_user_proto_rawDescGZIP(), []int{0}
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Role        UserRole               `protobuf:"varint,2,opt,name=role,proto3,enum=rapid.model.v1.UserRole" json:"role,omitempty"`
	AuthUid     string                 `protobuf:"bytes,3,opt,name=auth_uid,json=authUid,proto3" json:"auth_uid,omitempty"`
	DisplayName string                 `protobuf:"bytes,4,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	ImageUrl    string                 `protobuf:"bytes,5,opt,name=image_url,json=imageUrl,proto3" json:"image_url,omitempty"`
	Email       string                 `protobuf:"bytes,6,opt,name=email,proto3" json:"email,omitempty"`
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt   *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// Types that are assignable to OneofTenant:
	//	*User_TenantId
	//	*User_Tenant
	OneofTenant isUser_OneofTenant `protobuf_oneof:"oneof_tenant"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rapid_model_v1_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_rapid_model_v1_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_rapid_model_v1_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetRole() UserRole {
	if x != nil {
		return x.Role
	}
	return UserRole_USER_ROLE_UNSPECIFIED
}

func (x *User) GetAuthUid() string {
	if x != nil {
		return x.AuthUid
	}
	return ""
}

func (x *User) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *User) GetImageUrl() string {
	if x != nil {
		return x.ImageUrl
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *User) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (m *User) GetOneofTenant() isUser_OneofTenant {
	if m != nil {
		return m.OneofTenant
	}
	return nil
}

func (x *User) GetTenantId() string {
	if x, ok := x.GetOneofTenant().(*User_TenantId); ok {
		return x.TenantId
	}
	return ""
}

func (x *User) GetTenant() *Tenant {
	if x, ok := x.GetOneofTenant().(*User_Tenant); ok {
		return x.Tenant
	}
	return nil
}

type isUser_OneofTenant interface {
	isUser_OneofTenant()
}

type User_TenantId struct {
	TenantId string `protobuf:"bytes,9,opt,name=tenant_id,json=tenantId,proto3,oneof"`
}

type User_Tenant struct {
	Tenant *Tenant `protobuf:"bytes,10,opt,name=tenant,proto3,oneof"`
}

func (*User_TenantId) isUser_OneofTenant() {}

func (*User_Tenant) isUser_OneofTenant() {}

var File_rapid_model_v1_user_proto protoreflect.FileDescriptor

var file_rapid_model_v1_user_proto_rawDesc = []byte{
	0x0a, 0x19, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x76, 0x31,
	0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x72, 0x61, 0x70,
	0x69, 0x64, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69,
	0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x61,
	0x70, 0x69, 0x64, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe7, 0x03, 0x0a, 0x04, 0x55, 0x73,
	0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x2c, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x18, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x76,
	0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65,
	0x12, 0x19, 0x0a, 0x08, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x61, 0x75, 0x74, 0x68, 0x55, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x64,
	0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x74, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x30, 0x0a, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x48, 0x00,
	0x52, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x3a, 0x59, 0x92, 0x41, 0x56, 0x0a, 0x54, 0xd2,
	0x01, 0x02, 0x69, 0x64, 0xd2, 0x01, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0xd2, 0x01, 0x08, 0x61, 0x75,
	0x74, 0x68, 0x5f, 0x75, 0x69, 0x64, 0xd2, 0x01, 0x0c, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0xd2, 0x01, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x75, 0x72,
	0x6c, 0xd2, 0x01, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0xd2, 0x01, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0xd2, 0x01, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x42, 0x0e, 0x0a, 0x0c, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x5f, 0x74, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x2a, 0x50, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x12,
	0x19, 0x0a, 0x15, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x55, 0x53,
	0x45, 0x52, 0x5f, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x4e, 0x4f, 0x52, 0x4d, 0x41, 0x4c, 0x10, 0x01,
	0x12, 0x13, 0x0a, 0x0f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x41, 0x44,
	0x4d, 0x49, 0x4e, 0x10, 0x02, 0x42, 0xd3, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x72, 0x61,
	0x70, 0x69, 0x64, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x76, 0x31, 0x42, 0x09, 0x55, 0x73,
	0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x58, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x62, 0x79, 0x73, 0x73, 0x70, 0x61, 0x72, 0x61, 0x6e,
	0x6f, 0x69, 0x61, 0x2f, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2d, 0x67, 0x6f, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75, 0x63,
	0x74, 0x75, 0x72, 0x65, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x72, 0x61, 0x70,
	0x69, 0x64, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x76, 0x31, 0x3b, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x52, 0x4d, 0x58, 0xaa, 0x02, 0x0e, 0x52, 0x61, 0x70, 0x69,
	0x64, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0e, 0x52, 0x61, 0x70,
	0x69, 0x64, 0x5c, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1a, 0x52, 0x61,
	0x70, 0x69, 0x64, 0x5c, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x10, 0x52, 0x61, 0x70, 0x69, 0x64,
	0x3a, 0x3a, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_rapid_model_v1_user_proto_rawDescOnce sync.Once
	file_rapid_model_v1_user_proto_rawDescData = file_rapid_model_v1_user_proto_rawDesc
)

func file_rapid_model_v1_user_proto_rawDescGZIP() []byte {
	file_rapid_model_v1_user_proto_rawDescOnce.Do(func() {
		file_rapid_model_v1_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_rapid_model_v1_user_proto_rawDescData)
	})
	return file_rapid_model_v1_user_proto_rawDescData
}

var file_rapid_model_v1_user_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_rapid_model_v1_user_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_rapid_model_v1_user_proto_goTypes = []interface{}{
	(UserRole)(0),                 // 0: rapid.model.v1.UserRole
	(*User)(nil),                  // 1: rapid.model.v1.User
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
	(*Tenant)(nil),                // 3: rapid.model.v1.Tenant
}
var file_rapid_model_v1_user_proto_depIdxs = []int32{
	0, // 0: rapid.model.v1.User.role:type_name -> rapid.model.v1.UserRole
	2, // 1: rapid.model.v1.User.created_at:type_name -> google.protobuf.Timestamp
	2, // 2: rapid.model.v1.User.updated_at:type_name -> google.protobuf.Timestamp
	3, // 3: rapid.model.v1.User.tenant:type_name -> rapid.model.v1.Tenant
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_rapid_model_v1_user_proto_init() }
func file_rapid_model_v1_user_proto_init() {
	if File_rapid_model_v1_user_proto != nil {
		return
	}
	file_rapid_model_v1_tenant_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_rapid_model_v1_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
	file_rapid_model_v1_user_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*User_TenantId)(nil),
		(*User_Tenant)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rapid_model_v1_user_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rapid_model_v1_user_proto_goTypes,
		DependencyIndexes: file_rapid_model_v1_user_proto_depIdxs,
		EnumInfos:         file_rapid_model_v1_user_proto_enumTypes,
		MessageInfos:      file_rapid_model_v1_user_proto_msgTypes,
	}.Build()
	File_rapid_model_v1_user_proto = out.File
	file_rapid_model_v1_user_proto_rawDesc = nil
	file_rapid_model_v1_user_proto_goTypes = nil
	file_rapid_model_v1_user_proto_depIdxs = nil
}
