# Handler Patterns Reference

Detailed code patterns for gRPC handlers and DI configuration.

## Marshaller

Location: `internal/infrastructure/grpc/internal/handler/{actor}/marshaller/{entity}.go`

### Domain to Proto

```go
package marshaller

import (
    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    pb "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
    "google.golang.org/protobuf/types/known/timestamppb"
)

func ExampleToPb(m *model.Example) *pb.Example {
    if m == nil {
        return nil
    }

    result := &pb.Example{
        Id:          m.ID,
        TenantId:    m.TenantID,
        Name:        m.Name,
        Description: m.Description,
        Status:      ExampleStatusToPb(m.Status),
        CreatedAt:   timestamppb.New(m.CreatedAt),
        UpdatedAt:   timestamppb.New(m.UpdatedAt),
    }

    // Handle ReadonlyReference
    if m.ReadonlyReference != nil && m.ReadonlyReference.Tenant != nil {
        result.Tenant = TenantToPb(m.ReadonlyReference.Tenant)
    }

    return result
}

func ExamplesToPb(models model.Examples) []*pb.Example {
    result := make([]*pb.Example, len(models))
    for i, m := range models {
        result[i] = ExampleToPb(m)
    }
    return result
}
```

### Enum Conversions

```go
func ExampleStatusToPb(s model.ExampleStatus) pb.ExampleStatus {
    switch s {
    case model.ExampleStatusDraft:
        return pb.ExampleStatus_EXAMPLE_STATUS_DRAFT
    case model.ExampleStatusPublished:
        return pb.ExampleStatus_EXAMPLE_STATUS_PUBLISHED
    case model.ExampleStatusArchived:
        return pb.ExampleStatus_EXAMPLE_STATUS_ARCHIVED
    default:
        return pb.ExampleStatus_EXAMPLE_STATUS_UNSPECIFIED
    }
}

func ExampleStatusToModel(s pb.ExampleStatus) model.ExampleStatus {
    switch s {
    case pb.ExampleStatus_EXAMPLE_STATUS_DRAFT:
        return model.ExampleStatusDraft
    case pb.ExampleStatus_EXAMPLE_STATUS_PUBLISHED:
        return model.ExampleStatusPublished
    case pb.ExampleStatus_EXAMPLE_STATUS_ARCHIVED:
        return model.ExampleStatusArchived
    default:
        return model.ExampleStatusUnknown
    }
}

func ExampleSortKeyToModel(s pb.ExampleSortKey) model.ExampleSortKey {
    switch s {
    case pb.ExampleSortKey_EXAMPLE_SORT_KEY_CREATED_AT_DESC:
        return model.ExampleSortKeyCreatedAtDesc
    case pb.ExampleSortKey_EXAMPLE_SORT_KEY_NAME_ASC:
        return model.ExampleSortKeyNameAsc
    default:
        return model.ExampleSortKeyUnknown
    }
}
```

### Nullable Timestamp Fields

Use var declaration pattern for optional timestamp fields:

```go
func InvitationToPb(m *model.Invitation) *pb.Invitation {
    if m == nil {
        return nil
    }

    // Declare nullable fields first
    var acceptedAt *timestamppb.Timestamp
    if m.AcceptedAt.Valid {
        acceptedAt = timestamppb.New(m.AcceptedAt.Time)
    }

    var rejectedAt *timestamppb.Timestamp
    if m.RejectedAt.Valid {
        rejectedAt = timestamppb.New(m.RejectedAt.Time)
    }

    return &pb.Invitation{
        Id:         m.ID,
        Status:     InvitationStatusToPb(m.Status),
        AcceptedAt: acceptedAt,
        RejectedAt: rejectedAt,
        CreatedAt:  timestamppb.New(m.CreatedAt),
        UpdatedAt:  timestamppb.New(m.UpdatedAt),
    }
}
```

## Handler Methods

Location: `internal/infrastructure/grpc/internal/handler/{actor}/{entity}.go`

### Create Handler

```go
package admin

import (
    "context"

    "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/admin/marshaller"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
    admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
    "github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *Handler) CreateExample(
    ctx context.Context,
    req *admin_apiv1.CreateExampleRequest,
) (*admin_apiv1.CreateExampleResponse, error) {
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil {
        return nil, err
    }

    example, err := h.adminExampleInteractor.Create(ctx, &input.AdminCreateExample{
        StaffID:     claims.StaffID.String,
        TenantID:    req.TenantId,
        Name:        req.Name,
        Description: req.Description,
        RequestTime: request_interceptor.GetRequestTime(ctx),
    })
    if err != nil {
        return nil, err
    }

    return &admin_apiv1.CreateExampleResponse{
        Example: marshaller.ExampleToPb(example),
    }, nil
}
```

### Get Handler

```go
func (h *Handler) GetExample(
    ctx context.Context,
    req *admin_apiv1.GetExampleRequest,
) (*admin_apiv1.GetExampleResponse, error) {
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil {
        return nil, err
    }

    example, err := h.adminExampleInteractor.Get(ctx, &input.AdminGetExample{
        StaffID:     claims.StaffID.String,
        TenantID:    req.TenantId,
        ExampleID:   req.ExampleId,
        RequestTime: request_interceptor.GetRequestTime(ctx),
    })
    if err != nil {
        return nil, err
    }

    return &admin_apiv1.GetExampleResponse{
        Example: marshaller.ExampleToPb(example),
    }, nil
}
```

### List Handler

```go
func (h *Handler) ListExamples(
    ctx context.Context,
    req *admin_apiv1.ListExamplesRequest,
) (*admin_apiv1.ListExamplesResponse, error) {
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil {
        return nil, err
    }

    param := &input.AdminListExamples{
        StaffID:     claims.StaffID.String,
        TenantID:    req.TenantId,
        Page:        req.Page,
        Limit:       req.Limit,
        RequestTime: request_interceptor.GetRequestTime(ctx),
    }

    // Handle optional enum fields
    if req.Status != nil {
        status := marshaller.ExampleStatusToModel(*req.Status)
        param.Status = &status
    }
    if req.SortKey != nil {
        sortKey := marshaller.ExampleSortKeyToModel(*req.SortKey)
        param.SortKey = &sortKey
    }

    result, err := h.adminExampleInteractor.List(ctx, param)
    if err != nil {
        return nil, err
    }

    return &admin_apiv1.ListExamplesResponse{
        Examples:   marshaller.ExamplesToPb(result.Examples),
        TotalCount: result.TotalCount,
    }, nil
}
```

### Update Handler

```go
import "github.com/aarondl/null/v9"

func (h *Handler) UpdateExample(
    ctx context.Context,
    req *admin_apiv1.UpdateExampleRequest,
) (*admin_apiv1.UpdateExampleResponse, error) {
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil {
        return nil, err
    }

    param := &input.AdminUpdateExample{
        StaffID:     claims.StaffID.String,
        TenantID:    req.TenantId,
        ExampleID:   req.ExampleId,
        RequestTime: request_interceptor.GetRequestTime(ctx),
    }

    // Handle optional update fields
    if req.Name != nil {
        param.Name = null.StringFrom(*req.Name)
    }
    if req.Description != nil {
        param.Description = null.StringFrom(*req.Description)
    }

    example, err := h.adminExampleInteractor.Update(ctx, param)
    if err != nil {
        return nil, err
    }

    return &admin_apiv1.UpdateExampleResponse{
        Example: marshaller.ExampleToPb(example),
    }, nil
}
```

### Delete Handler

```go
func (h *Handler) DeleteExample(
    ctx context.Context,
    req *admin_apiv1.DeleteExampleRequest,
) (*admin_apiv1.DeleteExampleResponse, error) {
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil {
        return nil, err
    }

    if err := h.adminExampleInteractor.Delete(ctx, &input.AdminDeleteExample{
        StaffID:     claims.StaffID.String,
        TenantID:    req.TenantId,
        ExampleID:   req.ExampleId,
        RequestTime: request_interceptor.GetRequestTime(ctx),
    }); err != nil {
        return nil, err
    }

    return &admin_apiv1.DeleteExampleResponse{}, nil
}
```

## Handler Struct

Location: `internal/infrastructure/grpc/internal/handler/{actor}/handler.go`

```go
package admin

import (
    admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
    "github.com/abyssparanoia/rapid-go/internal/usecase"
)

type Handler struct {
    tenantInteractor  usecase.AdminTenantInteractor
    staffInteractor   usecase.AdminStaffInteractor
    assetInteractor   usecase.AdminAssetInteractor
    exampleInteractor usecase.AdminExampleInteractor  // Add new field
}

func NewHandler(
    tenantInteractor usecase.AdminTenantInteractor,
    staffInteractor usecase.AdminStaffInteractor,
    assetInteractor usecase.AdminAssetInteractor,
    exampleInteractor usecase.AdminExampleInteractor,  // Add new parameter
) admin_apiv1.AdminV1ServiceServer {
    return &Handler{
        tenantInteractor:  tenantInteractor,
        staffInteractor:   staffInteractor,
        assetInteractor:   assetInteractor,
        exampleInteractor: exampleInteractor,
    }
}
```

## DI Registration

Location: `internal/infrastructure/dependency/dependency.go`

### Add to Dependency Struct

```go
type Dependency struct {
    // ...existing fields...

    // Admin Interactors
    AdminExampleInteractor usecase.AdminExampleInteractor  // Add
}
```

### Add to Inject Method

```go
func (d *Dependency) Inject(ctx context.Context, e *environment.Environment) {
    // ...existing initialization...

    // Repository
    exampleRepository := db_repository.NewExample()

    // Interactor
    d.AdminExampleInteractor = usecase.NewAdminExampleInteractor(
        transactable,
        exampleRepository,
        assetService,
    )
}
```

### Update gRPC Server

Location: `internal/infrastructure/grpc/run.go`

```go
adminHandler := admin.NewHandler(
    dep.AdminTenantInteractor,
    dep.AdminStaffInteractor,
    dep.AdminAssetInteractor,
    dep.AdminExampleInteractor,  // Add
)
```

## Key Patterns

### Context Helpers

```go
// Get authenticated staff claims
claims, err := session_interceptor.RequireStaffSessionContext(ctx)
if err != nil {
    return nil, err
}

// Get request timestamp
requestTime := request_interceptor.GetRequestTime(ctx)
```

### Error Handling

Return errors directly from interactors - error interceptor handles conversion:

```go
example, err := h.adminExampleInteractor.Get(ctx, param)
if err != nil {
    return nil, err  // Error interceptor converts to gRPC status
}
```

### Optional Proto Fields

Handle `optional` proto fields by checking for nil:

```go
// Enum filter
if req.Status != nil {
    status := marshaller.ExampleStatusToModel(*req.Status)
    param.Status = &status
}

// String update
if req.Name != nil {
    param.Name = null.StringFrom(*req.Name)
}
```
