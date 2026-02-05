---
description: Protocol Buffers definition conventions (rapid-go)
globs:
  - "schema/proto/**/*.proto"
---

# Protocol Buffers Guidelines

This document summarizes proto definition rules for this repository's gRPC / gRPC-Gateway / OpenAPI generation.

## Directory Structure

Split definitions by audience under `schema/proto/rapid/`.

```
schema/proto/rapid/
├── admin_api/v1/
│   ├── api.proto                # Service definition (RPC + HTTP annotations only)
│   ├── api_asset.proto          # Request/Response messages (per resource)
│   ├── api_staff.proto
│   ├── api_tenant.proto
│   ├── model_asset.proto        # Domain models / enums
│   ├── model_staff.proto
│   ├── model_tenant.proto
│   └── model_pagination.proto
├── public_api/v1/
│   ├── api.proto
│   └── api_health_check.proto
└── debug_api/v1/
    └── api.proto
```

## Package Naming

- `rapid.admin_api.v1`
- `rapid.public_api.v1`
- `rapid.debug_api.v1`

## File Organization

### `api.proto` - Service Definition

- Keep **only** service definitions and HTTP annotations here. Put Request/Response and models into `api_*.proto` / `model_*.proto`.
- `api.proto` should import `api_*.proto` and reference their message types in RPC signatures.

Example (admin):

```protobuf
syntax = "proto3";

package rapid.admin_api.v1;

import "google/api/annotations.proto";
import "rapid/admin_api/v1/api_asset.proto";
import "rapid/admin_api/v1/api_staff.proto";
import "rapid/admin_api/v1/api_tenant.proto";

service AdminV1Service {
  rpc GetTenant(GetTenantRequest) returns (GetTenantResponse) {
    option (google.api.http) = {get: "/admin/v1/tenants/{tenant_id}"};
  }

  rpc ListTenants(ListTenantsRequest) returns (ListTenantsResponse) {
    option (google.api.http) = {get: "/admin/v1/tenants"};
  }

  rpc CreateTenant(CreateTenantRequest) returns (CreateTenantResponse) {
    option (google.api.http) = {
      post: "/admin/v1/tenants"
      body: "*"
    };
  }

  rpc UpdateTenant(UpdateTenantRequest) returns (UpdateTenantResponse) {
    option (google.api.http) = {
      patch: "/admin/v1/tenants/{tenant_id}"
      body: "*"
    };
  }

  rpc DeleteTenant(DeleteTenantRequest) returns (DeleteTenantResponse) {
    option (google.api.http) = {delete: "/admin/v1/tenants/{tenant_id}"};
  }
}
```

### `api_{resource}.proto` - Request/Response Definitions

- Define Request/Response messages here (and supporting types such as list `Pagination` if needed).
- Mark required fields using `protoc-gen-openapiv2` schema annotations (this repository's convention).

Example (admin tenant):

```protobuf
syntax = "proto3";

package rapid.admin_api.v1;

import "protoc-gen-openapiv2/options/annotations.proto";
import "rapid/admin_api/v1/model_pagination.proto";
import "rapid/admin_api/v1/model_tenant.proto";

message GetTenantRequest {
  string tenant_id = 1;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenant_id"] }
  };
}

message GetTenantResponse {
  Tenant tenant = 1;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenant"] }
  };
}

message ListTenantsRequest {
  uint64 page = 1;
  uint64 limit = 2;
}

message ListTenantsResponse {
  repeated Tenant tenants = 1;
  Pagination pagination = 2;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenants", "pagination"] }
  };
}

message UpdateTenantRequest {
  string tenant_id = 1;
  optional string name = 2;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenant_id"] }
  };
}
```

### `model_{resource}.proto` - Model / Enum Definitions

- Place resource messages and enums here.
- If you use `google.protobuf.Timestamp`, import `google/protobuf/timestamp.proto`.
- Each resource requires **both Full and Partial message definitions** (see Partial Pattern section below).

Example (admin staff):

```protobuf
syntax = "proto3";

package rapid.admin_api.v1;

import "google/protobuf/timestamp.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "rapid/admin_api/v1/model_tenant.proto";

// Full - for direct CRUD responses (with timestamps)
message Staff {
  string id = 1;
  TenantPartial tenant = 2;
  StaffRole role = 3;
  string auth_uid = 4;
  string display_name = 5;
  string image_url = 6;
  string email = 7;
  google.protobuf.Timestamp created_at = 8;
  google.protobuf.Timestamp updated_at = 9;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: {
      required: [
        "id",
        "tenant",
        "role",
        "auth_uid",
        "display_name",
        "image_url",
        "email",
        "created_at",
        "updated_at"
      ]
    }
  };
}

// Partial - for embedding in other resources (no timestamps)
message StaffPartial {
  string id = 1;
  TenantPartial tenant = 2;
  StaffRole role = 3;
  string auth_uid = 4;
  string display_name = 5;
  string image_url = 6;
  string email = 7;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: {
      required: ["id", "tenant", "role", "auth_uid", "display_name", "image_url", "email"]
    }
  };
}

enum StaffRole {
  STAFF_ROLE_UNSPECIFIED = 0;
  STAFF_ROLE_NORMAL = 1;
  STAFF_ROLE_ADMIN = 2;
}
```

## Partial Pattern

全リソースに対して`XXXPartial`メッセージを定義する。

### 定義ルール

| Message Type | 用途 | Timestamps |
|--------------|------|------------|
| `{Entity}` | CRUD直接レスポンス | あり (created_at, updated_at) |
| `{Entity}Partial` | 他リソースへの埋め込み | なし |

### Partialの構造

```protobuf
// Full - CRUD直接レスポンス用（timestamps含む）
message Example {
  string id = 1;
  TenantPartial tenant = 2;  // 親参照はPartial
  string name = 3;
  ExampleStatus status = 4;
  google.protobuf.Timestamp created_at = 5;
  google.protobuf.Timestamp updated_at = 6;
}

// Partial - 他リソースへの埋め込み用（timestamps除外）
message ExamplePartial {
  string id = 1;
  TenantPartial tenant = 2;  // 親参照もPartial
  string name = 3;
  ExampleStatus status = 4;
  // created_at, updated_at は含まない
}
```

### Request vs Response

- **Response**: `{Parent}Partial`を使用（`{parent}_id`ではなく`{Parent}Partial {parent}`）
- **Request**: 引き続き`string {parent}_id`を使用

```protobuf
// Response - TenantPartialを使用
message Staff {
  TenantPartial tenant = 2;  // string tenant_idではない
}

// Request - tenant_id stringを使用
message CreateStaffRequest {
  string tenant_id = 1;  // IDのまま
}
```

### ネストしたPartial

PartialがPartialを含む場合も同様にPartialを使用：

```protobuf
message StaffPartial {
  string id = 1;
  TenantPartial tenant = 2;  // PartialがPartialを含む
  // ...
}
```

### Field Number Convention（更新）

| Range | Purpose |
|-------|---------|
| 1-99 | 通常フィールド（Partial埋め込み含む） |
| 100+ | 予約（使用しない） |

**Note**: Partialパターンでは`optional {Parent} {parent} = 101`は不要。親参照は常に field 1-99 の範囲で必須フィールドとして定義される。

## Naming Conventions

### Service Names

- `{Audience}V1Service` (e.g. `AdminV1Service`, `PublicV1Service`, `DebugV1Service`)

### Message Names

- Resource: PascalCase like `Tenant`, `Staff`, `Asset`
- Request/Response: `{Action}{Resource}Request` / `{Action}{Resource}Response`

### Field Names

- snake_case (e.g. `tenant_id`, `display_name`, `created_at`)

### Enum Names

- Always start enums with `*_UNSPECIFIED = 0`
- SCREAMING_SNAKE_CASE

## HTTP Annotations (gRPC-Gateway)

Existing APIs in this repository use the following path prefixes by audience:

- admin: `/admin/v1/...`
- debug: `/debug/v1/...`
- public: `/v1/...` (NOT `/public/v1`)

Basic CRUD patterns:

```protobuf
// Create
post: "/admin/v1/tenants"
body: "*"

// Get
get: "/admin/v1/tenants/{tenant_id}"

// List
get: "/admin/v1/tenants"

// Update (partial)
patch: "/admin/v1/tenants/{tenant_id}"
body: "*"

// Delete
delete: "/admin/v1/tenants/{tenant_id}"
```

Special endpoint patterns:

- Operation-style endpoints using `/-/` (e.g. `/admin/v1/assets/-/presigned_url`, `/debug/v1/staffs/-/id_token`)

## Optional Fields

- For partial updates (PATCH), use `optional` for fields that may be updated (e.g. `optional string name = 2;`).

## OpenAPI (protoc-gen-openapiv2) Required

- Explicitly specify required fields for Request/Response using `openapiv2_schema.required` in most cases.
- For list requests, you may omit `required` depending on existing definitions—follow existing patterns in this repo.

## Code Generation

```bash
make generate.buf
```

Generated files location: `internal/infrastructure/grpc/pb/`

## Import Paths

```protobuf
import "google/api/annotations.proto";
import "google/protobuf/timestamp.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "rapid/admin_api/v1/api_tenant.proto";
import "rapid/admin_api/v1/model_tenant.proto";
```

