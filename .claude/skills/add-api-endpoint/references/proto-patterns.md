# Protocol Buffers Patterns Reference

Detailed code patterns for proto definitions.

## File Organization

Split into separate files by purpose:

```
schema/proto/rapid/{actor}_api/v1/
├── api.proto              # Service definition (RPCs + HTTP annotations)
├── api_{entity}.proto     # Request/Response messages
└── model_{entity}.proto   # Model messages + enums
```

## Model Definition

Location: `schema/proto/rapid/{actor}_api/v1/model_{entity}.proto`

```protobuf
syntax = "proto3";

package rapid.admin_api.v1;

import "google/protobuf/timestamp.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "rapid/admin_api/v1/model_tenant.proto";

// Main entity message
message Example {
  string id = 1;
  string tenant_id = 2;
  string name = 3;
  string description = 4;
  ExampleStatus status = 5;
  google.protobuf.Timestamp created_at = 6;
  google.protobuf.Timestamp updated_at = 7;

  // Readonly relation (field numbers 101+)
  optional Tenant tenant = 101;
}

// Status enum
enum ExampleStatus {
  EXAMPLE_STATUS_UNSPECIFIED = 0;
  EXAMPLE_STATUS_DRAFT = 1;
  EXAMPLE_STATUS_PUBLISHED = 2;
  EXAMPLE_STATUS_ARCHIVED = 3;
}

// Sort key enum
enum ExampleSortKey {
  EXAMPLE_SORT_KEY_UNSPECIFIED = 0;
  EXAMPLE_SORT_KEY_CREATED_AT_DESC = 1;
  EXAMPLE_SORT_KEY_NAME_ASC = 2;
}
```

## Request/Response Messages

Location: `schema/proto/rapid/{actor}_api/v1/api_{entity}.proto`

```protobuf
syntax = "proto3";

package rapid.admin_api.v1;

import "protoc-gen-openapiv2/options/annotations.proto";
import "rapid/admin_api/v1/model_{entity}.proto";
import "rapid/admin_api/v1/model_pagination.proto";

// Create
message CreateExampleRequest {
  string tenant_id = 1;
  string name = 2;
  string description = 3;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenant_id", "name", "description"] }
  };
}

message CreateExampleResponse {
  Example example = 1;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["example"] }
  };
}

// Get
message GetExampleRequest {
  string tenant_id = 1;
  string example_id = 2;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenant_id", "example_id"] }
  };
}

message GetExampleResponse {
  Example example = 1;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["example"] }
  };
}

// List
message ListExamplesRequest {
  string tenant_id = 1;
  optional ExampleStatus status = 2;
  optional ExampleSortKey sort_key = 3;
  uint64 page = 4;
  uint64 limit = 5;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenant_id", "page", "limit"] }
  };
}

message ListExamplesResponse {
  repeated Example examples = 1;
  uint64 total_count = 2;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["examples", "total_count"] }
  };
}

// Update
message UpdateExampleRequest {
  string tenant_id = 1;
  string example_id = 2;
  optional string name = 3;
  optional string description = 4;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenant_id", "example_id"] }
  };
}

message UpdateExampleResponse {
  Example example = 1;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["example"] }
  };
}

// Delete
message DeleteExampleRequest {
  string tenant_id = 1;
  string example_id = 2;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenant_id", "example_id"] }
  };
}

message DeleteExampleResponse {}
```

## Service Definition

Location: `schema/proto/rapid/{actor}_api/v1/api.proto`

Add RPCs to existing service:

```protobuf
syntax = "proto3";

package rapid.admin_api.v1;

import "google/api/annotations.proto";
import "rapid/admin_api/v1/api_example.proto";

service AdminV1Service {
  // ... existing rpcs ...

  // Example CRUD
  rpc CreateExample(CreateExampleRequest) returns (CreateExampleResponse) {
    option (google.api.http) = {
      post: "/admin/v1/tenants/{tenant_id}/examples"
      body: "*"
    };
  }

  rpc GetExample(GetExampleRequest) returns (GetExampleResponse) {
    option (google.api.http) = {
      get: "/admin/v1/tenants/{tenant_id}/examples/{example_id}"
    };
  }

  rpc ListExamples(ListExamplesRequest) returns (ListExamplesResponse) {
    option (google.api.http) = {
      get: "/admin/v1/tenants/{tenant_id}/examples"
    };
  }

  rpc UpdateExample(UpdateExampleRequest) returns (UpdateExampleResponse) {
    option (google.api.http) = {
      patch: "/admin/v1/tenants/{tenant_id}/examples/{example_id}"
      body: "*"
    };
  }

  rpc DeleteExample(DeleteExampleRequest) returns (DeleteExampleResponse) {
    option (google.api.http) = {
      delete: "/admin/v1/tenants/{tenant_id}/examples/{example_id}"
    };
  }
}
```

## HTTP Path Conventions

| Actor | Path Prefix |
|-------|-------------|
| Admin | `/admin/v1/...` |
| Debug | `/debug/v1/...` |
| Public | `/v1/...` (no `public` prefix) |

### CRUD Patterns

```
POST   /admin/v1/tenants/{tenant_id}/examples        # Create
GET    /admin/v1/tenants/{tenant_id}/examples/{id}   # Get
GET    /admin/v1/tenants/{tenant_id}/examples        # List
PATCH  /admin/v1/tenants/{tenant_id}/examples/{id}   # Update
DELETE /admin/v1/tenants/{tenant_id}/examples/{id}   # Delete
```

### Special Operations

Use `/-/` for operation-style endpoints:

```
POST /admin/v1/assets/-/presigned_url
POST /debug/v1/staffs/-/id_token
```

## Key Conventions

### Enum Naming

```protobuf
enum ExampleStatus {
  EXAMPLE_STATUS_UNSPECIFIED = 0;  // Always first
  EXAMPLE_STATUS_DRAFT = 1;
  EXAMPLE_STATUS_PUBLISHED = 2;
}
```

- Always start with `{TYPE}_UNSPECIFIED = 0`
- Use SCREAMING_SNAKE_CASE
- Prefix with type name

### Optional Fields

Use `optional` for partial updates and optional filters:

```protobuf
message UpdateExampleRequest {
  optional string name = 3;        // Partial update
}

message ListExamplesRequest {
  optional ExampleStatus status = 2;  // Optional filter
}
```

### Readonly Relations

Use field numbers 101+ for readonly relations:

```protobuf
message Example {
  string id = 1;
  // ... regular fields 1-99 ...

  // Readonly relations 101+
  optional Tenant tenant = 101;
  optional Staff created_by = 102;
}
```

### Required Fields Annotation

Always add `openapiv2_schema` for request/response messages:

```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
  json_schema: { required: ["tenant_id", "name"] }
};
```

## Import Paths

```protobuf
import "google/api/annotations.proto";
import "google/protobuf/timestamp.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "rapid/admin_api/v1/model_tenant.proto";
```
