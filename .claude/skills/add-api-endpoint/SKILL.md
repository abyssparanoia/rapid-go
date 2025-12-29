---
name: add-api-endpoint
description: REQUIRED Step 3 of CRUD workflow (after add-domain-entity). Use when creating usecase interactors, proto definitions in schema/proto/, gRPC handlers, or registering in dependency.go.
---

# Add API Endpoint

This skill guides you through creating API layer components for a new entity.

## Prerequisites

- Domain entity already created (use **add-domain-entity** skill first)
- Repository interface and implementation ready

## Step 1: Create Usecase Input

Location: `internal/usecase/input/{actor}_{entity}.go`

```go
package input

import (
    "time"
    "github.com/abyssparanoia/rapid-go/internal/domain/errors"
    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/pkg/validation"
    "github.com/volatiletech/null/v8"
)

// Create
type AdminCreateExample struct {
    StaffID     string    `validate:"required"`
    TenantID    string    `validate:"required"`
    Name        string    `validate:"required,max=256"`
    Description string    `validate:"required"`
    RequestTime time.Time `validate:"required"`
}

func (p *AdminCreateExample) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    return nil
}

// Get
type AdminGetExample struct {
    StaffID     string    `validate:"required"`
    TenantID    string    `validate:"required"`
    ExampleID   string    `validate:"required"`
    RequestTime time.Time `validate:"required"`
}

func (p *AdminGetExample) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    return nil
}

// List
type AdminListExamples struct {
    StaffID     string `validate:"required"`
    TenantID    string `validate:"required"`
    Status      *model.ExampleStatus
    SortKey     *model.ExampleSortKey
    Page        uint64    `validate:"required,min=1"`
    Limit       uint64    `validate:"required,min=1,max=100"`
    RequestTime time.Time `validate:"required"`
}

func (p *AdminListExamples) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    return nil
}

// Update
type AdminUpdateExample struct {
    StaffID     string      `validate:"required"`
    TenantID    string      `validate:"required"`
    ExampleID   string      `validate:"required"`
    Name        null.String
    Description null.String
    RequestTime time.Time   `validate:"required"`
}

func (p *AdminUpdateExample) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    return nil
}

// Delete
type AdminDeleteExample struct {
    StaffID     string    `validate:"required"`
    TenantID    string    `validate:"required"`
    ExampleID   string    `validate:"required"`
    RequestTime time.Time `validate:"required"`
}

func (p *AdminDeleteExample) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    return nil
}
```

## Step 2: Create Usecase Output (if needed)

Location: `internal/usecase/output/{actor}_{entity}.go`

```go
package output

import "github.com/abyssparanoia/rapid-go/internal/domain/model"

type AdminListExamples struct {
    Examples   model.Examples
    TotalCount uint64
}
```

## Step 3: Create Interactor Interface

Location: `internal/usecase/{actor}_{entity}.go`

```go
package usecase

import (
    "context"
    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/usecase/input"
    "github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type AdminExampleInteractor interface {
    Create(ctx context.Context, param *input.AdminCreateExample) (*model.Example, error)
    Get(ctx context.Context, param *input.AdminGetExample) (*model.Example, error)
    List(ctx context.Context, param *input.AdminListExamples) (*output.AdminListExamples, error)
    Update(ctx context.Context, param *input.AdminUpdateExample) (*model.Example, error)
    Delete(ctx context.Context, param *input.AdminDeleteExample) error
}
```

## Step 4: Create Interactor Implementation

Location: `internal/usecase/{actor}_{entity}_impl.go`

```go
package usecase

import (
    "context"

    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/domain/repository"
    "github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
    "github.com/abyssparanoia/rapid-go/internal/usecase/input"
    "github.com/abyssparanoia/rapid-go/internal/usecase/output"
    "github.com/volatiletech/null/v8"
)

type adminExampleInteractor struct {
    transactable      repository.Transactable
    exampleRepository repository.Example
}

func NewAdminExampleInteractor(
    transactable repository.Transactable,
    exampleRepository repository.Example,
) AdminExampleInteractor {
    return &adminExampleInteractor{
        transactable:      transactable,
        exampleRepository: exampleRepository,
    }
}

func (i *adminExampleInteractor) Create(ctx context.Context, param *input.AdminCreateExample) (*model.Example, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    example := model.NewExample(param.TenantID, param.Name, param.Description, param.RequestTime)

    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        return i.exampleRepository.Create(ctx, example)
    }); err != nil {
        return nil, err
    }

    return i.exampleRepository.Get(ctx, repository.GetExampleQuery{
        ID:             null.StringFrom(example.ID),
        BaseGetOptions: repository.BaseGetOptions{OrFail: true, Preload: true},
    })
}

func (i *adminExampleInteractor) Get(ctx context.Context, param *input.AdminGetExample) (*model.Example, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    return i.exampleRepository.Get(ctx, repository.GetExampleQuery{
        ID:             null.StringFrom(param.ExampleID),
        TenantID:       null.StringFrom(param.TenantID),
        BaseGetOptions: repository.BaseGetOptions{OrFail: true, Preload: true},
    })
}

func (i *adminExampleInteractor) List(ctx context.Context, param *input.AdminListExamples) (*output.AdminListExamples, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    query := repository.ListExamplesQuery{
        TenantID:        null.StringFrom(param.TenantID),
        BaseListOptions: repository.BaseListOptions{Page: null.Uint64From(param.Page), Limit: null.Uint64From(param.Limit), Preload: true},
    }
    if param.Status != nil {
        query.Status = nullable.From(*param.Status)
    }
    if param.SortKey != nil {
        query.SortKey = nullable.From(*param.SortKey)
    }

    examples, err := i.exampleRepository.List(ctx, query)
    if err != nil {
        return nil, err
    }
    totalCount, err := i.exampleRepository.Count(ctx, query)
    if err != nil {
        return nil, err
    }

    return &output.AdminListExamples{Examples: examples, TotalCount: totalCount}, nil
}

func (i *adminExampleInteractor) Update(ctx context.Context, param *input.AdminUpdateExample) (*model.Example, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        example, err := i.exampleRepository.Get(ctx, repository.GetExampleQuery{
            ID:             null.StringFrom(param.ExampleID),
            TenantID:       null.StringFrom(param.TenantID),
            BaseGetOptions: repository.BaseGetOptions{OrFail: true, ForUpdate: true},
        })
        if err != nil {
            return err
        }
        example.Update(param.Name, param.Description, param.RequestTime)
        return i.exampleRepository.Update(ctx, example)
    }); err != nil {
        return nil, err
    }

    return i.exampleRepository.Get(ctx, repository.GetExampleQuery{
        ID:             null.StringFrom(param.ExampleID),
        BaseGetOptions: repository.BaseGetOptions{OrFail: true, Preload: true},
    })
}

func (i *adminExampleInteractor) Delete(ctx context.Context, param *input.AdminDeleteExample) error {
    if err := param.Validate(); err != nil {
        return err
    }

    return i.transactable.RWTx(ctx, func(ctx context.Context) error {
        _, err := i.exampleRepository.Get(ctx, repository.GetExampleQuery{
            ID:             null.StringFrom(param.ExampleID),
            TenantID:       null.StringFrom(param.TenantID),
            BaseGetOptions: repository.BaseGetOptions{OrFail: true, ForUpdate: true},
        })
        if err != nil {
            return err
        }
        return i.exampleRepository.Delete(ctx, param.ExampleID)
    })
}
```

## Step 5: Define Protocol Buffers

Location:

- Model / Enum: `schema/proto/rapid/{actor}_api/v1/model_{entity}.proto`
- Request/Response: `schema/proto/rapid/{actor}_api/v1/api_{entity}.proto`

```protobuf
syntax = "proto3";

package rapid.admin_api.v1;

import "google/api/annotations.proto";
import "google/protobuf/timestamp.proto";
import "protoc-gen-openapiv2/options/annotations.proto";

message Example {
  string id = 1;
  string tenant_id = 2;
  string name = 3;
  string description = 4;
  ExampleStatus status = 5;
  google.protobuf.Timestamp created_at = 6;
  google.protobuf.Timestamp updated_at = 7;
}

enum ExampleStatus {
  EXAMPLE_STATUS_UNSPECIFIED = 0;
  EXAMPLE_STATUS_DRAFT = 1;
  EXAMPLE_STATUS_PUBLISHED = 2;
  EXAMPLE_STATUS_ARCHIVED = 3;
}

enum ExampleSortKey {
  EXAMPLE_SORT_KEY_UNSPECIFIED = 0;
  EXAMPLE_SORT_KEY_CREATED_AT_DESC = 1;
  EXAMPLE_SORT_KEY_NAME_ASC = 2;
}

message CreateExampleRequest {
  string tenant_id = 1;
  string name = 2;
  string description = 3;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenant_id", "name", "description"] }
  };
}
message CreateExampleResponse { Example example = 1; }

message GetExampleRequest {
  string tenant_id = 1;
  string example_id = 2;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenant_id", "example_id"] }
  };
}
message GetExampleResponse { Example example = 1; }

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
}

message UpdateExampleRequest {
  string tenant_id = 1;
  string example_id = 2;
  optional string name = 3;
  optional string description = 4;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenant_id", "example_id"] }
  };
}
message UpdateExampleResponse { Example example = 1; }

message DeleteExampleRequest {
  string tenant_id = 1;
  string example_id = 2;

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: { required: ["tenant_id", "example_id"] }
  };
}
message DeleteExampleResponse {}
```

## Step 6: Add Service RPCs

Add to `schema/proto/rapid/{actor}_api/v1/api.proto`:

```protobuf
import "google/api/annotations.proto";
import "rapid/admin_api/v1/api_example.proto";

service AdminV1Service {
  // ... existing rpcs ...

  rpc CreateExample(CreateExampleRequest) returns (CreateExampleResponse) {
    option (google.api.http) = { post: "/admin/v1/tenants/{tenant_id}/examples" body: "*" };
  }
  rpc GetExample(GetExampleRequest) returns (GetExampleResponse) {
    option (google.api.http) = { get: "/admin/v1/tenants/{tenant_id}/examples/{example_id}" };
  }
  rpc ListExamples(ListExamplesRequest) returns (ListExamplesResponse) {
    option (google.api.http) = { get: "/admin/v1/tenants/{tenant_id}/examples" };
  }
  rpc UpdateExample(UpdateExampleRequest) returns (UpdateExampleResponse) {
    option (google.api.http) = { patch: "/admin/v1/tenants/{tenant_id}/examples/{example_id}" body: "*" };
  }
  rpc DeleteExample(DeleteExampleRequest) returns (DeleteExampleResponse) {
    option (google.api.http) = { delete: "/admin/v1/tenants/{tenant_id}/examples/{example_id}" };
  }
}
```

## Step 7: Generate Proto Code

```bash
make generate.buf
```

## Step 8: Create gRPC Handler Marshaller

Location: `internal/infrastructure/grpc/internal/handler/{actor}/marshaller/{entity}.go`

```go
package marshaller

import (
    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    pb "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
    "google.golang.org/protobuf/types/known/timestamppb"
)

func ExampleToPb(m *model.Example) *pb.Example {
    if m == nil { return nil }
    return &pb.Example{
        Id: m.ID, TenantId: m.TenantID, Name: m.Name, Description: m.Description,
        Status: ExampleStatusToPb(m.Status),
        CreatedAt: timestamppb.New(m.CreatedAt), UpdatedAt: timestamppb.New(m.UpdatedAt),
    }
}

func ExamplesToPb(models model.Examples) []*pb.Example {
    result := make([]*pb.Example, len(models))
    for i, m := range models { result[i] = ExampleToPb(m) }
    return result
}

func ExampleStatusToPb(s model.ExampleStatus) pb.ExampleStatus {
    switch s {
    case model.ExampleStatusDraft: return pb.ExampleStatus_EXAMPLE_STATUS_DRAFT
    case model.ExampleStatusPublished: return pb.ExampleStatus_EXAMPLE_STATUS_PUBLISHED
    case model.ExampleStatusArchived: return pb.ExampleStatus_EXAMPLE_STATUS_ARCHIVED
    default: return pb.ExampleStatus_EXAMPLE_STATUS_UNSPECIFIED
    }
}

func ExampleStatusToModel(s pb.ExampleStatus) model.ExampleStatus {
    switch s {
    case pb.ExampleStatus_EXAMPLE_STATUS_DRAFT: return model.ExampleStatusDraft
    case pb.ExampleStatus_EXAMPLE_STATUS_PUBLISHED: return model.ExampleStatusPublished
    case pb.ExampleStatus_EXAMPLE_STATUS_ARCHIVED: return model.ExampleStatusArchived
    default: return model.ExampleStatusUnknown
    }
}

func ExampleSortKeyToModel(s pb.ExampleSortKey) model.ExampleSortKey {
    switch s {
    case pb.ExampleSortKey_EXAMPLE_SORT_KEY_CREATED_AT_DESC: return model.ExampleSortKeyCreatedAtDesc
    case pb.ExampleSortKey_EXAMPLE_SORT_KEY_NAME_ASC: return model.ExampleSortKeyNameAsc
    default: return model.ExampleSortKeyUnknown
    }
}
```

## Step 9: Create gRPC Handler

Location: `internal/infrastructure/grpc/internal/handler/{actor}/{entity}.go`

```go
package admin

import (
    "context"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/admin/marshaller"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
    admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
    "github.com/abyssparanoia/rapid-go/internal/usecase/input"
    "github.com/volatiletech/null/v8"
)

func (h *AdminHandler) CreateExample(ctx context.Context, req *admin_apiv1.CreateExampleRequest) (*admin_apiv1.CreateExampleResponse, error) {
    // Admin API authorization is enforced by interceptors.
    // If you need staff claims, read them from session context:
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil { return nil, err }
    example, err := h.adminExampleInteractor.Create(ctx, &input.AdminCreateExample{
        StaffID: claims.StaffID.String, TenantID: req.TenantId,
        Name: req.Name, Description: req.Description,
        RequestTime: request_interceptor.GetRequestTime(ctx),
    })
    if err != nil { return nil, err }
    return &admin_apiv1.CreateExampleResponse{Example: marshaller.ExampleToPb(example)}, nil
}

func (h *AdminHandler) GetExample(ctx context.Context, req *admin_apiv1.GetExampleRequest) (*admin_apiv1.GetExampleResponse, error) {
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil { return nil, err }
    example, err := h.adminExampleInteractor.Get(ctx, &input.AdminGetExample{
        StaffID: claims.StaffID.String, TenantID: req.TenantId,
        ExampleID: req.ExampleId,
        RequestTime: request_interceptor.GetRequestTime(ctx),
    })
    if err != nil { return nil, err }
    return &admin_apiv1.GetExampleResponse{Example: marshaller.ExampleToPb(example)}, nil
}

func (h *AdminHandler) ListExamples(ctx context.Context, req *admin_apiv1.ListExamplesRequest) (*admin_apiv1.ListExamplesResponse, error) {
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil { return nil, err }
    param := &input.AdminListExamples{
        StaffID: claims.StaffID.String, TenantID: req.TenantId,
        Page: req.Page, Limit: req.Limit,
        RequestTime: request_interceptor.GetRequestTime(ctx),
    }
    if req.Status != nil {
        status := marshaller.ExampleStatusToModel(*req.Status)
        param.Status = &status
    }
    if req.SortKey != nil {
        sortKey := marshaller.ExampleSortKeyToModel(*req.SortKey)
        param.SortKey = &sortKey
    }
    result, err := h.adminExampleInteractor.List(ctx, param)
    if err != nil { return nil, err }
    return &admin_apiv1.ListExamplesResponse{
        Examples: marshaller.ExamplesToPb(result.Examples),
        TotalCount: result.TotalCount,
    }, nil
}

func (h *AdminHandler) UpdateExample(ctx context.Context, req *admin_apiv1.UpdateExampleRequest) (*admin_apiv1.UpdateExampleResponse, error) {
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil { return nil, err }
    param := &input.AdminUpdateExample{
        StaffID: claims.StaffID.String, TenantID: req.TenantId,
        ExampleID: req.ExampleId,
        RequestTime: request_interceptor.GetRequestTime(ctx),
    }
    if req.Name != nil { param.Name = null.StringFrom(*req.Name) }
    if req.Description != nil { param.Description = null.StringFrom(*req.Description) }
    example, err := h.adminExampleInteractor.Update(ctx, param)
    if err != nil { return nil, err }
    return &admin_apiv1.UpdateExampleResponse{Example: marshaller.ExampleToPb(example)}, nil
}

func (h *AdminHandler) DeleteExample(ctx context.Context, req *admin_apiv1.DeleteExampleRequest) (*admin_apiv1.DeleteExampleResponse, error) {
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil { return nil, err }
    err := h.adminExampleInteractor.Delete(ctx, &input.AdminDeleteExample{
        StaffID: claims.StaffID.String, TenantID: req.TenantId,
        ExampleID: req.ExampleId,
        RequestTime: request_interceptor.GetRequestTime(ctx),
    })
    if err != nil { return nil, err }
    return &admin_apiv1.DeleteExampleResponse{}, nil
}
```

## Step 10: Update Handler Struct

Update `internal/infrastructure/grpc/internal/handler/{actor}/handler.go` (this repo uses `AdminHandler` + `NewAdminHandler(...)` pattern):

```go
type AdminHandler struct {
    tenantInteractor usecase.AdminTenantInteractor
    staffInteractor  usecase.AdminStaffInteractor
    assetInteractor  usecase.AdminAssetInteractor
    exampleInteractor usecase.AdminExampleInteractor // Add
}

func NewAdminHandler(
    tenantInteractor usecase.AdminTenantInteractor,
    staffInteractor usecase.AdminStaffInteractor,
    assetInteractor usecase.AdminAssetInteractor,
    exampleInteractor usecase.AdminExampleInteractor, // Add
) admin_apiv1.AdminV1ServiceServer {
    return &AdminHandler{
        tenantInteractor: tenantInteractor,
        staffInteractor:  staffInteractor,
        assetInteractor:  assetInteractor,
        exampleInteractor: exampleInteractor,
    }
}
```

## Step 11: Update DI

Update `internal/infrastructure/dependency/dependency.go` to expose the new interactor, then pass it when building the gRPC server:

```go
// Add to Dependency struct
AdminExampleInteractor usecase.AdminExampleInteractor

// Add to Inject()
exampleRepository := db_repository.NewExample()
d.AdminExampleInteractor = usecase.NewAdminExampleInteractor(transactable, exampleRepository)
```

Then update `internal/infrastructure/grpc/run.go` to pass the new interactor to `admin.NewAdminHandler(...)`.

## Step 12: Generate Mocks & Test

```bash
make generate.mock
make test
```

## Checklist

- [ ] Usecase input structs with validation
- [ ] Usecase output struct (if List)
- [ ] Interactor interface with go:generate
- [ ] Interactor implementation
- [ ] Proto message definitions
- [ ] Proto service RPCs
- [ ] Proto code generated
- [ ] gRPC handler marshaller
- [ ] gRPC handler methods
- [ ] Handler struct updated
- [ ] DI configuration updated
- [ ] Mocks generated
- [ ] Tests passing
