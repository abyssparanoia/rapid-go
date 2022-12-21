// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: rapid/model/v1/pagination.proto

package modelv1

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

type Pagination struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CurrentPage uint64 `protobuf:"varint,1,opt,name=current_page,json=currentPage,proto3" json:"current_page,omitempty"`
	PrevPage    uint64 `protobuf:"varint,2,opt,name=prev_page,json=prevPage,proto3" json:"prev_page,omitempty"`
	NextPage    uint64 `protobuf:"varint,3,opt,name=next_page,json=nextPage,proto3" json:"next_page,omitempty"`
	TotalPage   uint64 `protobuf:"varint,4,opt,name=total_page,json=totalPage,proto3" json:"total_page,omitempty"`
	TotalCount  uint64 `protobuf:"varint,5,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	HasNext     bool   `protobuf:"varint,6,opt,name=has_next,json=hasNext,proto3" json:"has_next,omitempty"`
}

func (x *Pagination) Reset() {
	*x = Pagination{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rapid_model_v1_pagination_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pagination) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pagination) ProtoMessage() {}

func (x *Pagination) ProtoReflect() protoreflect.Message {
	mi := &file_rapid_model_v1_pagination_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pagination.ProtoReflect.Descriptor instead.
func (*Pagination) Descriptor() ([]byte, []int) {
	return file_rapid_model_v1_pagination_proto_rawDescGZIP(), []int{0}
}

func (x *Pagination) GetCurrentPage() uint64 {
	if x != nil {
		return x.CurrentPage
	}
	return 0
}

func (x *Pagination) GetPrevPage() uint64 {
	if x != nil {
		return x.PrevPage
	}
	return 0
}

func (x *Pagination) GetNextPage() uint64 {
	if x != nil {
		return x.NextPage
	}
	return 0
}

func (x *Pagination) GetTotalPage() uint64 {
	if x != nil {
		return x.TotalPage
	}
	return 0
}

func (x *Pagination) GetTotalCount() uint64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

func (x *Pagination) GetHasNext() bool {
	if x != nil {
		return x.HasNext
	}
	return false
}

var File_rapid_model_v1_pagination_proto protoreflect.FileDescriptor

var file_rapid_model_v1_pagination_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x76, 0x31,
	0x2f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0e, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x76,
	0x31, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70,
	0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x80, 0x02, 0x0a, 0x0a, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x21, 0x0a, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50,
	0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x72, 0x65, 0x76, 0x5f, 0x70, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x70, 0x72, 0x65, 0x76, 0x50, 0x61, 0x67, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x08, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x0a, 0x0b,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x19, 0x0a,
	0x08, 0x68, 0x61, 0x73, 0x5f, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x68, 0x61, 0x73, 0x4e, 0x65, 0x78, 0x74, 0x3a, 0x3a, 0x92, 0x41, 0x37, 0x0a, 0x35, 0xd2,
	0x01, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0xd2, 0x01,
	0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x70, 0x61, 0x67, 0x65, 0xd2, 0x01, 0x0b, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0xd2, 0x01, 0x08, 0x68, 0x61, 0x73, 0x5f,
	0x6e, 0x65, 0x78, 0x74, 0x42, 0xc9, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x72, 0x61, 0x70,
	0x69, 0x64, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x76, 0x31, 0x42, 0x0f, 0x50, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x48,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x62, 0x79, 0x73, 0x73,
	0x70, 0x61, 0x72, 0x61, 0x6e, 0x6f, 0x69, 0x61, 0x2f, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2d, 0x67,
	0x6f, 0x2f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70,
	0x62, 0x2f, 0x72, 0x61, 0x70, 0x69, 0x64, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x76, 0x31,
	0x3b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x52, 0x4d, 0x58, 0xaa, 0x02,
	0x0e, 0x52, 0x61, 0x70, 0x69, 0x64, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x0e, 0x52, 0x61, 0x70, 0x69, 0x64, 0x5c, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x5c, 0x56, 0x31,
	0xe2, 0x02, 0x1a, 0x52, 0x61, 0x70, 0x69, 0x64, 0x5c, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x5c, 0x56,
	0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x10,
	0x52, 0x61, 0x70, 0x69, 0x64, 0x3a, 0x3a, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x3a, 0x3a, 0x56, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rapid_model_v1_pagination_proto_rawDescOnce sync.Once
	file_rapid_model_v1_pagination_proto_rawDescData = file_rapid_model_v1_pagination_proto_rawDesc
)

func file_rapid_model_v1_pagination_proto_rawDescGZIP() []byte {
	file_rapid_model_v1_pagination_proto_rawDescOnce.Do(func() {
		file_rapid_model_v1_pagination_proto_rawDescData = protoimpl.X.CompressGZIP(file_rapid_model_v1_pagination_proto_rawDescData)
	})
	return file_rapid_model_v1_pagination_proto_rawDescData
}

var file_rapid_model_v1_pagination_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_rapid_model_v1_pagination_proto_goTypes = []interface{}{
	(*Pagination)(nil), // 0: rapid.model.v1.Pagination
}
var file_rapid_model_v1_pagination_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rapid_model_v1_pagination_proto_init() }
func file_rapid_model_v1_pagination_proto_init() {
	if File_rapid_model_v1_pagination_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rapid_model_v1_pagination_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pagination); i {
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
			RawDescriptor: file_rapid_model_v1_pagination_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rapid_model_v1_pagination_proto_goTypes,
		DependencyIndexes: file_rapid_model_v1_pagination_proto_depIdxs,
		MessageInfos:      file_rapid_model_v1_pagination_proto_msgTypes,
	}.Build()
	File_rapid_model_v1_pagination_proto = out.File
	file_rapid_model_v1_pagination_proto_rawDesc = nil
	file_rapid_model_v1_pagination_proto_goTypes = nil
	file_rapid_model_v1_pagination_proto_depIdxs = nil
}
