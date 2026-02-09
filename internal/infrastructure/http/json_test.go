package http

import (
	"testing"

	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestCustomJSONPb_Int64AsNumber(t *testing.T) {
	tests := map[string]func(t *testing.T){
		"Int64AsNumber false - int64 fields are serialized as strings": func(t *testing.T) {
			marshaler := &CustomJSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				Int64AsNumber: false, // Default behavior
			}

			msg := &admin_apiv1.Pagination{
				CurrentPage: 1,
				PrevPage:    0,
				NextPage:    2,
				TotalPage:   412,
				TotalCount:  12345,
				HasNext:     true,
			}

			result, err := marshaler.Marshal(msg)
			require.NoError(t, err)

			// int64 fields should be strings
			assert.Contains(t, string(result), `"current_page":"1"`)
			assert.Contains(t, string(result), `"total_count":"12345"`)
			assert.Contains(t, string(result), `"total_page":"412"`)
		},
		"Int64AsNumber true - int64 fields are serialized as numbers": func(t *testing.T) {
			marshaler := &CustomJSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				Int64AsNumber: true, // Convert to numbers
			}

			msg := &admin_apiv1.Pagination{
				CurrentPage: 1,
				PrevPage:    0,
				NextPage:    2,
				TotalPage:   412,
				TotalCount:  12345,
				HasNext:     true,
			}

			result, err := marshaler.Marshal(msg)
			require.NoError(t, err)

			// int64 fields should be numbers
			assert.Contains(t, string(result), `"current_page":1`)
			assert.Contains(t, string(result), `"total_count":12345`)
			assert.Contains(t, string(result), `"total_page":412`)
		},
		"Int64AsNumber true - nested messages with int64 fields": func(t *testing.T) {
			marshaler := &CustomJSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				Int64AsNumber: true,
			}

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

			result, err := marshaler.Marshal(msg)
			require.NoError(t, err)

			// Nested pagination fields should be numbers
			assert.Contains(t, string(result), `"current_page":2`)
			assert.Contains(t, string(result), `"total_count":9999`)
			assert.Contains(t, string(result), `"total_page":200`)
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}
