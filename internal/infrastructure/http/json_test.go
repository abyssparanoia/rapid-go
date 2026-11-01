package http

import (
	"bytes"
	"testing"

	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestCustomJSONPb_Int64AsNumber(t *testing.T) {
	t.Parallel()

	newMarshaler := func(int64AsNumber bool) *CustomJSONPb {
		return &CustomJSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			Int64AsNumber: int64AsNumber,
		}
	}
	newPagination := func() *admin_apiv1.Pagination {
		return &admin_apiv1.Pagination{
			CurrentPage: 1,
			PrevPage:    0,
			NextPage:    2,
			TotalPage:   412,
			TotalCount:  12345,
			HasNext:     true,
		}
	}
	// newSample builds a message with the field shapes the admin API protos do not have:
	// map<string, int64>, repeated int64, google.protobuf.Int64Value, double and a numeric-looking string.
	newSample := func(t *testing.T) *dynamicpb.Message {
		t.Helper()
		const pkg = "rapid.http.jsontest"
		field := func(name string, number int32, typ descriptorpb.FieldDescriptorProto_Type, jsonName string) *descriptorpb.FieldDescriptorProto {
			return &descriptorpb.FieldDescriptorProto{
				Name:     proto.String(name),
				Number:   proto.Int32(number),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:     typ.Enum(),
				JsonName: proto.String(jsonName),
			}
		}
		counts := field("counts", 1, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, "counts")
		counts.Label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum()
		counts.TypeName = proto.String("." + pkg + ".Sample.CountsEntry")
		values := field("values", 2, descriptorpb.FieldDescriptorProto_TYPE_INT64, "values")
		values.Label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum()
		wrappedTotal := field("wrapped_total", 3, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, "wrappedTotal")
		wrappedTotal.TypeName = proto.String(".google.protobuf.Int64Value")

		file, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
			Name:       proto.String("rapid_http_json_test.proto"),
			Package:    proto.String(pkg),
			Syntax:     proto.String("proto3"),
			Dependency: []string{"google/protobuf/wrappers.proto"},
			MessageType: []*descriptorpb.DescriptorProto{{
				Name: proto.String("Sample"),
				Field: []*descriptorpb.FieldDescriptorProto{
					counts,
					values,
					wrappedTotal,
					field("label", 4, descriptorpb.FieldDescriptorProto_TYPE_STRING, "label"),
					field("ratio", 5, descriptorpb.FieldDescriptorProto_TYPE_DOUBLE, "ratio"),
				},
				NestedType: []*descriptorpb.DescriptorProto{{
					Name: proto.String("CountsEntry"),
					Field: []*descriptorpb.FieldDescriptorProto{
						field("key", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING, "key"),
						field("value", 2, descriptorpb.FieldDescriptorProto_TYPE_INT64, "value"),
					},
					Options: &descriptorpb.MessageOptions{MapEntry: proto.Bool(true)},
				}},
			}},
		}, protoregistry.GlobalFiles)
		require.NoError(t, err)

		md := file.Messages().ByName("Sample")
		msg := dynamicpb.NewMessage(md)
		countsMap := msg.Mutable(md.Fields().ByName("counts")).Map()
		countsMap.Set(protoreflect.ValueOfString("a").MapKey(), protoreflect.ValueOfInt64(1))
		countsMap.Set(protoreflect.ValueOfString("b").MapKey(), protoreflect.ValueOfInt64(-9007199254740993)) // |v| > 2^53
		valuesList := msg.Mutable(md.Fields().ByName("values")).List()
		valuesList.Append(protoreflect.ValueOfInt64(3))
		valuesList.Append(protoreflect.ValueOfInt64(9007199254740993)) // > 2^53
		msg.Set(md.Fields().ByName("wrapped_total"), protoreflect.ValueOfMessage(wrapperspb.Int64(42).ProtoReflect()))
		msg.Set(md.Fields().ByName("label"), protoreflect.ValueOfString("123"))
		msg.Set(md.Fields().ByName("ratio"), protoreflect.ValueOfFloat64(0.1))
		return msg
	}

	tests := map[string]func(t *testing.T){
		"Int64AsNumber false - int64 fields are serialized as strings": func(t *testing.T) {
			result, err := newMarshaler(false).Marshal(newPagination())
			require.NoError(t, err)

			assert.Contains(t, string(result), `"current_page":"1"`)
			assert.Contains(t, string(result), `"total_count":"12345"`)
			assert.Contains(t, string(result), `"total_page":"412"`)
		},
		"Int64AsNumber true - int64 fields are serialized as numbers": func(t *testing.T) {
			result, err := newMarshaler(true).Marshal(newPagination())
			require.NoError(t, err)

			assert.Contains(t, string(result), `"current_page":1`)
			assert.Contains(t, string(result), `"total_count":12345`)
			assert.Contains(t, string(result), `"total_page":412`)
		},
		"Int64AsNumber true - nested messages with int64 fields": func(t *testing.T) {
			msg := &admin_apiv1.ListTenantsResponse{
				Tenants: []*admin_apiv1.Tenant{
					{Id: "tenant1", Name: "Test Tenant 1"},
				},
				Pagination: &admin_apiv1.Pagination{
					CurrentPage: 2,
					PrevPage:    1,
					NextPage:    3,
					TotalPage:   200,
					TotalCount:  9999,
					HasNext:     true,
				},
			}

			result, err := newMarshaler(true).Marshal(msg)
			require.NoError(t, err)

			assert.Contains(t, string(result), `"current_page":2`)
			assert.Contains(t, string(result), `"total_count":9999`)
			assert.Contains(t, string(result), `"total_page":200`)
		},
		"Int64AsNumber true - camelCase field names when UseProtoNames is false": func(t *testing.T) {
			marshaler := newMarshaler(true)
			marshaler.UseProtoNames = false

			result, err := marshaler.Marshal(newPagination())
			require.NoError(t, err)

			assert.Contains(t, string(result), `"currentPage":1`)
			assert.Contains(t, string(result), `"totalCount":12345`)
		},
		"Int64AsNumber true - map values, repeated values and Int64Value are numbers, other fields untouched": func(t *testing.T) {
			result, err := newMarshaler(true).Marshal(newSample(t))
			require.NoError(t, err)

			assert.Contains(t, string(result), `"counts":{"a":1,"b":-9007199254740993}`)
			assert.Contains(t, string(result), `"values":[3,9007199254740993]`)
			assert.Contains(t, string(result), `"wrapped_total":42`)
			assert.Contains(t, string(result), `"label":"123"`)
			assert.Contains(t, string(result), `"ratio":0.1`)
		},
		"Int64AsNumber true - top-level Int64Value is a bare number": func(t *testing.T) {
			result, err := newMarshaler(true).Marshal(wrapperspb.Int64(7))
			require.NoError(t, err)

			assert.Equal(t, `7`, string(result))
		},
		"Int64AsNumber true - Indent is preserved": func(t *testing.T) {
			marshaler := newMarshaler(true)
			marshaler.Indent = "  "

			result, err := marshaler.Marshal(newPagination())
			require.NoError(t, err)

			assert.Contains(t, string(result), "\n  \"total_count\": 12345")
		},
		"Int64AsNumber true - NewEncoder applies the conversion": func(t *testing.T) {
			var buf bytes.Buffer
			require.NoError(t, newMarshaler(true).NewEncoder(&buf).Encode(newPagination()))

			assert.Contains(t, buf.String(), `"total_count":12345`)
			assert.True(t, bytes.HasSuffix(buf.Bytes(), []byte("\n")))
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tc(t)
		})
	}
}

func TestCustomJSONPb_Unmarshal(t *testing.T) {
	t.Parallel()

	marshaler := &CustomJSONPb{
		UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
	}

	tests := map[string]func(t *testing.T){
		"map - decodes string keys and values": func(t *testing.T) {
			var got map[string]string
			require.NoError(t, marshaler.Unmarshal([]byte(`{"a":"1","b":"2"}`), &got))

			assert.Equal(t, map[string]string{"a": "1", "b": "2"}, got)
		},
		"map - rejects a duplicate key": func(t *testing.T) {
			var got map[string]string
			err := marshaler.Unmarshal([]byte(`{"a":"1","a":"2"}`), &got)

			assert.ErrorContains(t, err, `duplicate map key "a"`)
		},
		"map - rejects two spellings of the same int32 key": func(t *testing.T) {
			var got map[int32]string
			err := marshaler.Unmarshal([]byte(`{"1":"x","01":"y"}`), &got)

			assert.ErrorContains(t, err, "duplicate map key")
		},
		"enum - decodes an in-range number": func(t *testing.T) {
			var got admin_apiv1.StaffRole
			require.NoError(t, marshaler.Unmarshal([]byte(`1`), &got))

			assert.Equal(t, admin_apiv1.StaffRole(1), got)
		},
		"enum - rejects a fractional number": func(t *testing.T) {
			var got admin_apiv1.StaffRole
			err := marshaler.Unmarshal([]byte(`1.5`), &got)

			assert.ErrorContains(t, err, "is not a valid value for enum")
		},
		"enum - rejects a number outside int32": func(t *testing.T) {
			var got admin_apiv1.StaffRole
			err := marshaler.Unmarshal([]byte(`4294967297`), &got)

			assert.ErrorContains(t, err, "is not a valid value for enum")
		},
		"enum - rejects a symbolic name": func(t *testing.T) {
			var got admin_apiv1.StaffRole
			err := marshaler.Unmarshal([]byte(`"STAFF_ROLE_ADMIN"`), &got)

			assert.ErrorContains(t, err, "symbolic enum")
		},
		"proto message - decodes via protojson": func(t *testing.T) {
			var got admin_apiv1.Pagination
			require.NoError(t, marshaler.Unmarshal([]byte(`{"current_page":"3","total_count":12,"unknown":true}`), &got))

			assert.Equal(t, uint64(3), got.GetCurrentPage())
			assert.Equal(t, uint64(12), got.GetTotalCount())
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tc(t)
		})
	}
}
