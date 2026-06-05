---
description: gRPC handler implementation patterns
globs:
  - "internal/infrastructure/grpc/**/*.go"
---

# gRPC Handler Guidelines

## Directory Structure

```
internal/infrastructure/grpc/
├── internal/
│   ├── handler/
│   │   ├── admin/             # Admin API handlers
│   │   │   ├── handler.go     # Handler struct & constructor
│   │   │   ├── tenant.go      # Resource handlers (example: tenant)
│   │   │   ├── staff.go
│   │   │   ├── asset.go
│   │   │   └── marshaller/    # Proto <-> Domain converters
│   │   ├── debug/             # Debug API handlers
│   │   └── public/            # Public API handlers
│   └── interceptor/
│       ├── authorization_interceptor/
│       ├── request_interceptor/
│       └── session_interceptor/
└── pb/                        # Generated protobuf code
```

## Handler Method Ordering

**All handler methods must be implemented in the following order:**

1. **Get methods** - Single resource retrieval
2. **List methods** - Collection retrieval with pagination
3. **Create methods** - Resource creation
4. **Custom operations (no ID)** - Special operations without resource ID
5. **Update methods** - Resource modification
6. **Custom operations (with ID)** - Special operations with resource ID
7. **Delete methods** - Resource deletion

**Example ordering in handler file:**

```go
// Get
func (h *AdminHandler) GetStaff(ctx context.Context, req *pb.GetStaffRequest) (*pb.GetStaffResponse, error) {...}

// List
func (h *AdminHandler) ListStaffs(ctx context.Context, req *pb.ListStaffsRequest) (*pb.ListStaffsResponse, error) {...}

// Create
func (h *AdminHandler) CreateStaff(ctx context.Context, req *pb.CreateStaffRequest) (*pb.CreateStaffResponse, error) {...}

// Custom (no ID)
func (h *AdminHandler) SendStaffNotifications(ctx context.Context, req *pb.SendStaffNotificationsRequest) (*pb.SendStaffNotificationsResponse, error) {...}

// Update
func (h *AdminHandler) UpdateStaff(ctx context.Context, req *pb.UpdateStaffRequest) (*pb.UpdateStaffResponse, error) {...}

// Custom (with ID)
func (h *AdminHandler) SendStaffNotification(ctx context.Context, req *pb.SendStaffNotificationRequest) (*pb.SendStaffNotificationResponse, error) {...}

// Delete
func (h *AdminHandler) DeleteStaff(ctx context.Context, req *pb.DeleteStaffRequest) (*pb.DeleteStaffResponse, error) {...}
```

**Important**: This ordering ensures consistency with proto definitions and usecase interfaces.

## Handler Structure

Location: `internal/infrastructure/grpc/internal/handler/{admin|public|debug}/handler.go`

```go
package admin

type Handler struct {
    tenantInteractor usecase.AdminTenantInteractor
    staffInteractor  usecase.AdminStaffInteractor
    assetInteractor  usecase.AdminAssetInteractor
}

func NewHandler(
    tenantInteractor usecase.AdminTenantInteractor,
    staffInteractor usecase.AdminStaffInteractor,
    assetInteractor usecase.AdminAssetInteractor,
) admin_apiv1.AdminV1ServiceServer { /* ... */ }
```

### Handlers Must Not Hold Repositories Directly

Handler structs may only depend on `usecase.*Interactor` (and `pb.Unimplemented*Server`). They must NOT take a `repository.*` field — even when the repository call is "just" for authorization context (e.g. resolving a vehicle's `device_group_id` for Layer 3 device-group permission).

Why:
- Layered architecture: `infrastructure/grpc` is meant to talk to `usecase`, not jump over it into `domain/repository`.
- Cross-cutting concerns (validation, asset URL hydration, tenant ownership checks) live in the interactor. A handler bypass tends to drift toward duplicating those checks.
- DI surface bloat: every "auth-only" repository field added to `Dependency` becomes a permanent leak that future handlers copy.

```go
// BAD - handler holds repository for authorization lookup
type TenantHandler struct {
    vehicleInteractor usecase.TenantVehicleInteractor
    vehicleRepository repository.Vehicle // forbidden, even for auth lookup
}

func (h *TenantHandler) ListVehicleLocationLogs(ctx context.Context, req ...) (...) {
    vehicle, err := h.vehicleRepository.Get(ctx, ...) // bypasses usecase
    // Layer 3 auth using vehicle.DeviceGroupID ...
}

// GOOD - add a Get method to the interactor; handler routes through usecase
type TenantHandler struct {
    vehicleInteractor    usecase.TenantVehicleInteractor
    vehicleMapInteractor usecase.TenantVehicleMapInteractor
}

func (h *TenantHandler) ListVehicleLocationLogs(ctx context.Context, req ...) (...) {
    vehicle, err := h.vehicleInteractor.Get(ctx, input.NewTenantGetVehicle(
        perm.TenantID, req.GetVehicleId(), request_interceptor.GetRequestTime(ctx),
    ))
    // Layer 3 auth using vehicle.DeviceGroupID ...
    result, err := h.vehicleMapInteractor.ListLocationLogs(ctx, ...)
}
```

If no suitable `Get` exists, add one to the matching interactor (with tenant-ownership validation inside). Do not introduce a one-off "auth helper" that takes a repository.

The same rule applies to `Dependency` in `internal/infrastructure/dependency/dependency.go`: only `usecase.*Interactor` fields are exposed to the handler layer; raw repositories stay private.

## Handler Method Pattern

Location: `internal/infrastructure/grpc/internal/handler/{actor}/{resource}.go`

```go
package admin

func (h *AdminHandler) CreateTenant(
    ctx context.Context,
    req *admin_apiv1.CreateTenantRequest,
) (*admin_apiv1.CreateTenantResponse, error) {
    // 1. Extract request time from context
    requestTime := request_interceptor.GetRequestTime(ctx)

    // 2. Call interactor with input struct
    tenant, err := h.tenantInteractor.Create(ctx, input.NewAdminCreateTenant(req.GetName(), requestTime))
    if err != nil {
        return nil, err  // Error interceptor handles conversion
    }

    // 3. Convert domain model to proto and return
    return &admin_apiv1.CreateTenantResponse{
        Tenant: marshaller.TenantToPB(tenant),
    }, nil
}
```

## Context Helpers

```go
// Admin API authorization is handled by interceptors.
// If a handler needs staff claims, read them from session context:
claims, err := session_interceptor.RequireStaffSessionContext(ctx)
if err != nil { return nil, err }

// Get request timestamp
requestTime := request_interceptor.GetRequestTime(ctx)
```

## Handling Optional Proto Fields

```go
func (h *Handler) ListExamples(
    ctx context.Context,
    req *pb.ListExamplesRequest,
) (*pb.ListExamplesResponse, error) {
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil { return nil, err }
    requestTime := request_interceptor.GetRequestTime(ctx)

    param := &input.AdminListExamples{
        StaffID:     claims.StaffID.String,
        TenantID:    req.TenantId,
        Page:        req.Page,
        Limit:       req.Limit,
        RequestTime: requestTime,
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

    return &pb.ListExamplesResponse{
        Examples:   marshaller.ExamplesToPb(result.Examples),
        TotalCount: result.TotalCount,
    }, nil
}
```

## Handling Optional Proto String Fields

For optional proto `string` fields (e.g., `optional string display_name`), use `null.StringFromPtr()` to convert directly from `*string` to `null.String`:

```go
func (h *StaffHandler) UpdateMe(
    ctx context.Context,
    req *staff_apiv1.UpdateMeRequest,
) (*staff_apiv1.UpdateMeResponse, error) {
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil {
        return nil, err
    }
    requestTime := request_interceptor.GetRequestTime(ctx)

    // Good - Direct call to constructor with StringFromPtr
    staff, err := h.meInteractor.Update(
        ctx,
        input.NewStaffUpdateMe(
            claims.TenantID.String,
            claims.StaffID.String,
            null.StringFromPtr(req.DisplayName),    // Handles nil pointer correctly
            null.StringFromPtr(req.ImageAssetId),
            requestTime,
        ),
    )
    if err != nil {
        return nil, err
    }

    return &staff_apiv1.UpdateMeResponse{
        Staff: marshaller.StaffToPB(staff),
    }, nil
}
```

### Anti-Pattern: Param Mutation

```go
// Bad - Don't create param first and then mutate fields
param := input.NewStaffUpdateMe(
    claims.TenantID.String,
    claims.StaffID.String,
    null.String{},
    null.String{},
    requestTime,
)

if req.DisplayName != nil {
    param.DisplayName = null.StringFrom(*req.DisplayName)
}
if req.ImageAssetId != nil {
    param.ImageAssetID = null.StringFrom(*req.ImageAssetId)
}

staff, err := h.meInteractor.Update(ctx, param)
```

### Rules

1. **Use `null.StringFromPtr()`** for optional proto `string` fields
   - `null.StringFromPtr(req.Field)` handles `nil` → `null.String{Valid: false}`
   - `null.StringFromPtr(&"value")` → `null.String{Valid: true, String: "value"}`

2. **Use `nullable.TypeFromPtr()`** for optional proto enum fields
   - Similar pattern for custom domain types

3. **Call constructor directly** in handler
   - Pass all arguments inline instead of mutating after construction
   - More concise and reduces risk of forgetting to set fields

## Marshaller (Proto <-> Domain)

Location: `internal/infrastructure/grpc/internal/handler/{actor}/marshaller/{resource}.go`

### File Organization

Each resource should have its own marshaller file:

```
internal/infrastructure/grpc/internal/handler/{actor}/marshaller/
├── admin.go              # Admin resource marshaller
├── admin_invitation.go   # AdminInvitation resource marshaller
├── tenant.go             # Tenant resource marshaller
└── common.go             # Shared enum conversions (if any)
```

**Do not** combine multiple resource marshallers into a single file. Keep marshaller files focused on one domain entity.

### Domain to Proto

#### Partial Pattern Marshallers

各エンティティに2つのmarshaller関数を定義：

```go
package marshaller

// CRUD直接レスポンス用（Fullメッセージ）
func ExampleToPb(m *model.Example) *pb.Example {
    if m == nil {
        return nil
    }

    var tenant *pb.TenantPartial
    if m.ReadonlyReference != nil && m.ReadonlyReference.Tenant != nil {
        tenant = TenantPartialToPb(m.ReadonlyReference.Tenant)
    }

    return &pb.Example{
        Id:        m.ID,
        Tenant:    tenant,
        Name:      m.Name,
        Status:    ExampleStatusToPb(m.Status),
        CreatedAt: timestamppb.New(m.CreatedAt),
        UpdatedAt: timestamppb.New(m.UpdatedAt),
    }
}

// 他リソースへの埋め込み用（Partialメッセージ）
func ExamplePartialToPb(m *model.Example) *pb.ExamplePartial {
    if m == nil {
        return nil
    }

    var tenant *pb.TenantPartial
    if m.ReadonlyReference != nil && m.ReadonlyReference.Tenant != nil {
        tenant = TenantPartialToPb(m.ReadonlyReference.Tenant)
    }

    return &pb.ExamplePartial{
        Id:     m.ID,
        Tenant: tenant,
        Name:   m.Name,
        Status: ExampleStatusToPb(m.Status),
    }
}

func ExamplesToPb(models model.Examples) []*pb.Example {
    result := make([]*pb.Example, len(models))
    for i, m := range models {
        result[i] = ExampleToPb(m)
    }
    return result
}

func ExamplesPartialToPb(models model.Examples) []*pb.ExamplePartial {
    result := make([]*pb.ExamplePartial, len(models))
    for i, m := range models {
        result[i] = ExamplePartialToPb(m)
    }
    return result
}
```

#### Preload必須

Partialパターンでは親参照が必須フィールド。レスポンス返却時は必ず`Preload: true`を指定：

```go
// REQUIRED - Preload must be true for responses
example, err := i.exampleRepository.Get(ctx, repository.GetExampleQuery{
    ID: null.StringFrom(param.ExampleID),
    BaseGetOptions: repository.BaseGetOptions{
        OrFail:  true,
        Preload: true,  // Partialパターンでは必須
    },
})
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
```

### Marshaller Return Pattern (Avoiding Field Mapping Omissions)

To avoid field mapping omissions, use a variable declaration pattern instead of inline struct initialization when dealing with optional/nullable fields:

```go
func AdminInvitationToPb(m *model.AdminInvitation) *pb.AdminInvitation {
    if m == nil {
        return nil
    }

    // Declare variables for optional/nullable fields first
    var acceptedAt *timestamppb.Timestamp
    if m.AcceptedAt.Valid {
        acceptedAt = timestamppb.New(m.AcceptedAt.Time)
    }

    var rejectedAt *timestamppb.Timestamp
    if m.RejectedAt.Valid {
        rejectedAt = timestamppb.New(m.RejectedAt.Time)
    }

    // Build proto message with all fields explicitly listed
    result := &pb.AdminInvitation{
        Id:             m.ID,
        InvitationCode: m.InvitationCode,
        Status:         AdminInvitationStatusToPb(m.Status),
        AcceptedAt:     acceptedAt,
        RejectedAt:     rejectedAt,
        CreatedAt:      timestamppb.New(m.CreatedAt),
        UpdatedAt:      timestamppb.New(m.UpdatedAt),
    }

    return result
}
```

**Why this pattern:**
- Explicit variable declarations make it easier to verify all fields are mapped
- Code reviewers can more easily spot missing field mappings
- Prevents accidental omission of nullable field handling

## Error Handling

- Return errors directly from interactors
- Error interceptor converts domain errors to gRPC status codes
- Domain errors (`errors.XxxErr`) map to appropriate gRPC codes

```go
// In handler - just return error
if err != nil {
    return nil, err
}

// Error interceptor handles conversion:
// errors.NotFoundErr -> codes.NotFound
// errors.InvalidArgumentErr -> codes.InvalidArgument
// errors.InternalErr -> codes.Internal
```

## DI Registration

Location: `internal/infrastructure/dependency/dependency.go`

```go
// Add handler to Dependency struct
AdminHandler *admin.Handler

// In Inject() method
d.AdminHandler = admin.NewHandler(
    d.AdminExampleInteractor,
    // ... other interactors
)
```
