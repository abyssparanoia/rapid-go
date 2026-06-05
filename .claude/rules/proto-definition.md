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

## RPC Method Ordering

**All RPC methods must be defined in the following order across all proto files:**

1. **Get methods** - Single resource retrieval (e.g., `GetStaff`)
2. **List methods** - Collection retrieval with pagination (e.g., `ListStaffs`)
3. **Create methods** - Resource creation (e.g., `CreateStaff`)
4. **Custom operations (no ID)** - Special operations without resource ID (e.g., `rpc SendNotifications` with path `/staffs:send_notifications`)
5. **Update methods** - Resource modification (e.g., `UpdateStaff`)
6. **Custom operations (with ID)** - Special operations with resource ID (e.g., `rpc SendNotification` with path `/staffs/{staff_id}:send_notification`)
7. **Delete methods** - Resource deletion (e.g., `DeleteStaff`)

**Example ordering in api.proto:**

```protobuf
service AdminV1Service {
  // Get
  rpc GetStaff(GetStaffRequest) returns (GetStaffResponse) {...}

  // List
  rpc ListStaffs(ListStaffsRequest) returns (ListStaffsResponse) {...}

  // Create
  rpc CreateStaff(CreateStaffRequest) returns (CreateStaffResponse) {...}

  // Custom (no ID)
  rpc SendStaffNotifications(SendStaffNotificationsRequest) returns (SendStaffNotificationsResponse) {...}

  // Update
  rpc UpdateStaff(UpdateStaffRequest) returns (UpdateStaffResponse) {...}

  // Custom (with ID)
  rpc SendStaffNotification(SendStaffNotificationRequest) returns (SendStaffNotificationResponse) {...}

  // Delete
  rpc DeleteStaff(DeleteStaffRequest) returns (DeleteStaffResponse) {...}
}
```

**Important**: This ordering applies to both `api.proto` service definitions and `api_{resource}.proto` message definitions.

## File Organization

### `api.proto` - Service Definition

- Keep **only** service definitions and HTTP annotations here. Put Request/Response and models into `api_*.proto` / `model_*.proto`.
- `api.proto` should import `api_*.proto` and reference their message types in RPC signatures.
- **API仕様コメントはrpc定義に記述する**。認証要否・処理フロー・サーバー側処理などの仕様はすべて `api.proto` のrpc定義上のコメントに集約する。

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
- **API仕様コメントをRequest/Responseメッセージに書かない**。認証要否・処理フローなどの仕様は `api.proto` のrpc定義に記述する。

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

#### Shared Enums Belong in `model_{resource}.proto`

ある enum を **複数の Request メッセージ** で参照する場合 (例: `ListXxx` と `ExportXxxCSV` 両方が同じ `DailyMetricType` を必須フィールドに持つケース)、その enum は `model_{resource}.proto` の **トップレベル** に置く。`ListXxxRequest` の中にネストさせて、別 Request から `ListXxxRequest.SomeEnum` の修飾名で参照するのは禁止。

理由:
- enum は意味的に「画面表示用 List API の都合」ではなく「ドメイン値型」。`ListXxxRequest` の中に隠れていると、新しい RPC が同じ enum を再利用しづらい。
- ネスト + 修飾参照は **構造的カップリング** を生む — `ListXxxRequest` の名前変更や enum リネームが他 Request を silent に巻き込む。
- 1 ヶ所だけ修飾名で参照するパターンは proto レベルでは合法だが、Go 側の生成型名 (`ListXxxRequest_SomeEnum_VALUE`) が冗長になり、テスト・marshaller の取り回しも悪化する。

```protobuf
// BAD - SomeEnum を ListXxxRequest にネストし、ExportXxxCSVRequest から修飾参照
message ListXxxRequest {
  enum SomeEnum { SOME_ENUM_UNSPECIFIED = 0; ... }
  SomeEnum field = 5;
}
message ExportXxxCSVRequest {
  ListXxxRequest.SomeEnum field = 5; // 構造的カップリング
}

// GOOD - SomeEnum をトップレベルに置き、両 Request から直接参照
// model_xxx.proto
enum SomeEnum { SOME_ENUM_UNSPECIFIED = 0; ... }

// api_xxx.proto
message ListXxxRequest {
  SomeEnum field = 5;
}
message ExportXxxCSVRequest {
  SomeEnum field = 5;
}
```

**例外**: `ListXxxSortKey` のように **その List API 固有** で他 RPC で再利用する見込みがない enum は、List Request にネストする (既存規約: SortKey Enum Naming Convention セクション参照)。判断基準は「他 Request から参照する可能性があるか」。


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

### Partial間の依存ルール

**Partial間では他のPartialを参照しない**のが原則。唯一の例外は**親エンティティ（TenantPartial）**のみ。

親以外のエンティティへの参照は `string {entity}_id` を使用する：

```protobuf
// Good - 親参照のみPartialを許可
message StaffPartial {
  string id = 1;
  TenantPartial tenant = 2;  // 親参照 → Partial OK
  // ...
}

// Good - 親以外はIDで参照
message StaffInvitationPartial {
  string id = 1;
  TenantPartial tenant = 2;           // 親参照 → Partial OK
  optional string invited_by_staff_id = 6;  // 親以外 → string ID
  // ...
}

// Bad - 親以外のPartialを参照
message StaffInvitationPartial {
  optional StaffPartial invited_by = 6;  // NG: Partial間の依存
}
```

**Note**: Full メッセージ（CRUD直接レスポンス）では `StaffPartial` などの参照は許可される。この制約は Partial メッセージのみに適用される。

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

Path naming rule:

- **Use `snake_case` for both the resource segment and the custom verb in URL paths.** Do not mix hyphens and underscores in the same path. Existing endpoints all use snake_case (e.g. `bots:resolve_code`, `staff_invitations:accept`, `license_plate_ocr_sessions/{...}:recognize`, `vehicle_operating_statuses`); a new endpoint such as `location-logs:export-csv` must be written `location_logs:export_csv`.

## Optional Fields

- For partial updates (PATCH), use `optional` for fields that may be updated (e.g. `optional string name = 2;`).

## List Request Patterns (Pagination & SortKey)

All List operations follow a unified specification for pagination and sorting.

### Required Fields in ListXXXRequest

Only truly required fields (e.g., `tenant_id` for tenant-scoped lists) should be marked as required in OpenAPI schema. Pagination fields (`page`, `limit`) and `sort_key` are **optional** - defaults are applied in the input layer constructor.

```protobuf
message ListStaffsRequest {
  string tenant_id = 1;
  uint64 page = 2;
  uint64 limit = 3;
  optional ListStaffsSortKey sort_key = 4;

  enum ListStaffsSortKey {
    LIST_STAFFS_SORT_KEY_UNSPECIFIED = 0;
    LIST_STAFFS_SORT_KEY_CREATED_AT_DESC = 1;
    LIST_STAFFS_SORT_KEY_CREATED_AT_ASC = 2;
    LIST_STAFFS_SORT_KEY_DISPLAY_NAME_ASC = 3;
    LIST_STAFFS_SORT_KEY_DISPLAY_NAME_DESC = 4;
  }

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: {
      required: ["tenant_id"]  // Only tenant_id is required
    }
  };
}
```

### SortKey Enum Naming Convention

- **Enum name**: `List{Entity}sSortKey` (nested inside `List{Entity}sRequest`)
- **Enum values**: `LIST_{ENTITY}S_SORT_KEY_{FIELD}_{DIRECTION}`
  - Always include `UNSPECIFIED = 0` as first value
  - Direction: `ASC` or `DESC`
  - Common fields: `CREATED_AT`, `UPDATED_AT`, entity-specific fields (e.g., `NAME`, `DISPLAY_NAME`)

### SortKey Enum Positioning

**IMPORTANT**: The enum definition must appear **BEFORE** the field that uses it, not after.

```protobuf
message ListStaffsRequest {
  string tenant_id = 1;
  uint64 page = 2;
  uint64 limit = 3;

  // CORRECT - Enum defined BEFORE the field that uses it
  enum ListStaffsSortKey {
    LIST_STAFFS_SORT_KEY_UNSPECIFIED = 0;
    LIST_STAFFS_SORT_KEY_CREATED_AT_DESC = 1;
    LIST_STAFFS_SORT_KEY_CREATED_AT_ASC = 2;
    LIST_STAFFS_SORT_KEY_DISPLAY_NAME_ASC = 3;
    LIST_STAFFS_SORT_KEY_DISPLAY_NAME_DESC = 4;
  }

  optional ListStaffsSortKey sort_key = 4;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: {
      required: ["tenant_id"]
    }
  };
}
```

**Anti-pattern** (do not use):
```protobuf
message ListStaffsRequest {
  // ...
  optional ListStaffsSortKey sort_key = 4;  // Field declared first

  enum ListStaffsSortKey {  // Enum defined after - WRONG
    // ...
  }
}
```

### Default Values

- **page**: Default `1` if unspecified (applied in input constructor)
- **limit**: Default `30` if unspecified (applied in input constructor)
- **sort_key**: Default `CreatedAtDesc` if unspecified (applied in input constructor)

**Important**: Do NOT use proto3 default values or validation tags. Defaults are handled in the Go input layer constructor.

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

