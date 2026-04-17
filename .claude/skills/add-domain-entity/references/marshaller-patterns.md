# Marshaller Patterns

Detailed code patterns for converting between DB models and domain models.

## File Location

`internal/infrastructure/{mysql|postgresql|spanner}/internal/marshaller/{entity}.go`

Each entity should have its own marshaller file.

## Basic Pattern: ToModel

```go
package marshaller

import (
    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/internal/dbmodel"
    "github.com/aarondl/null/v8"
)

func ExampleToModel(e *dbmodel.Example) *model.Example {
    m := &model.Example{
        ID:          e.ID,
        TenantID:    e.TenantID,
        Name:        e.Name,
        Description: e.Description,
        Status:      model.NewExampleStatus(e.Status),  // Use constructor for enums
        CreatedAt:   e.CreatedAt,
        UpdatedAt:   e.UpdatedAt,

        // Initialize computed fields
        ImageURL: null.String{},

        // Initialize ReadonlyReference
        ReadonlyReference: nil,
    }

    // Populate ReadonlyReference if relations are loaded
    if e.R != nil {
        var tenant *model.Tenant
        if e.R.Tenant != nil {
            tenant = TenantToModel(e.R.Tenant)
        }
        if tenant != nil {
            m.ReadonlyReference = &struct {
                Tenant *model.Tenant
            }{
                Tenant: tenant,
            }
        }
    }

    return m
}
```

## Slice Conversion

```go
func ExamplesToModel(slice dbmodel.ExampleSlice) model.Examples {
    dsts := make(model.Examples, len(slice))
    for idx, e := range slice {
        dsts[idx] = ExampleToModel(e)
    }
    return dsts
}
```

## Basic Pattern: ToDBModel

```go
func ExampleToDBModel(m *model.Example) *dbmodel.Example {
    return &dbmodel.Example{
        ID:          m.ID,
        TenantID:    m.TenantID,
        Name:        m.Name,
        Description: m.Description,
        Status:      m.Status.String(),  // Convert enum to string
        CreatedAt:   m.CreatedAt,
        UpdatedAt:   m.UpdatedAt,
        // ReadonlyReference and computed fields are NOT converted back
        R: nil,
        L: struct{}{},
    }
}

func ExamplesToDBModel(m model.Examples) dbmodel.ExampleSlice {
    dsts := make(dbmodel.ExampleSlice, len(m))
    for idx, e := range m {
        dsts[idx] = ExampleToDBModel(e)
    }
    return dsts
}
```

## Nullable Timestamp Pattern

For entities with optional timestamp fields (`null.Time`), use var declaration pattern to prevent field omission:

```go
func InvitationToModel(i *dbmodel.Invitation) *model.Invitation {
    // Declare nullable fields first to ensure all are handled
    var acceptedAt null.Time
    if i.AcceptedAt.Valid {
        acceptedAt = null.TimeFrom(i.AcceptedAt.Time)
    }

    var rejectedAt null.Time
    if i.RejectedAt.Valid {
        rejectedAt = null.TimeFrom(i.RejectedAt.Time)
    }

    var invalidatedAt null.Time
    if i.InvalidatedAt.Valid {
        invalidatedAt = null.TimeFrom(i.InvalidatedAt.Time)
    }

    // Build result with all fields explicitly listed
    result := &model.Invitation{
        ID:             i.ID,
        InvitationCode: i.InvitationCode,
        Status:         model.NewInvitationStatus(i.Status),
        Email:          i.Email,
        ExpiresAt:      i.ExpiresAt,
        AcceptedAt:     acceptedAt,
        RejectedAt:     rejectedAt,
        InvalidatedAt:  invalidatedAt,
        CreatedAt:      i.CreatedAt,
        UpdatedAt:      i.UpdatedAt,
        ReadonlyReference: nil,
    }

    return result
}
```

### Why Use This Pattern

- Explicit variable declarations make field mapping visible
- Prevents accidentally omitting nullable field handling
- Code reviewers can easily verify all fields are mapped

## ReadonlyReference Rules

### Rule 1: Always Initialize to nil

```go
m := &model.Example{
    // ... fields ...
    ReadonlyReference: nil,  // Always start as nil
}
```

### Rule 2: Populate Only if Relations are Loaded

```go
if e.R != nil {
    var tenant *model.Tenant
    if e.R.Tenant != nil {
        tenant = TenantToModel(e.R.Tenant)
    }
    if tenant != nil {
        m.ReadonlyReference = &struct {
            Tenant *model.Tenant
        }{
            Tenant: tenant,
        }
    }
}
```

### Rule 3: Related Entity's ReadonlyReference Must Be nil

When converting a related entity, its own `ReadonlyReference` remains nil. This prevents recursive/circular loading.

```go
// TenantToModel - ReadonlyReference is always nil
func TenantToModel(t *dbmodel.Tenant) *model.Tenant {
    return &model.Tenant{
        ID:        t.ID,
        Name:      t.Name,
        CreatedAt: t.CreatedAt,
        UpdatedAt: t.UpdatedAt,
        ReadonlyReference: nil,  // No recursive loading
    }
}
```

### Rule 4: Never Write Back

ReadonlyReference is ignored in ToDBModel:

```go
func ExampleToDBModel(m *model.Example) *dbmodel.Example {
    return &dbmodel.Example{
        // ... only persist direct fields ...
        R: nil,      // Relations not written
        L: struct{}{},
    }
}
```

## Multiple Relations

```go
func OrderToModel(o *dbmodel.Order) *model.Order {
    m := &model.Order{
        ID:        o.ID,
        CustomerID: o.CustomerID,
        Status:    model.NewOrderStatus(o.Status),
        CreatedAt: o.CreatedAt,
        UpdatedAt: o.UpdatedAt,
        ReadonlyReference: nil,
    }

    if o.R != nil {
        var customer *model.Customer
        if o.R.Customer != nil {
            customer = CustomerToModel(o.R.Customer)
        }

        var items model.OrderItems
        if o.R.OrderItems != nil {
            items = OrderItemsToModel(o.R.OrderItems)
        }

        // Only set if at least one relation is loaded
        if customer != nil || len(items) > 0 {
            m.ReadonlyReference = &struct {
                Customer *model.Customer
                Items    model.OrderItems
            }{
                Customer: customer,
                Items:    items,
            }
        }
    }

    return m
}
```

## Enum Conversion

Use constructor functions for enums to handle unknown values safely:

```go
// In domain model (model/example_status.go)
func NewExampleStatus(str string) ExampleStatus {
    switch str {
    case ExampleStatusDraft.String():
        return ExampleStatusDraft
    case ExampleStatusPublished.String():
        return ExampleStatusPublished
    case ExampleStatusArchived.String():
        return ExampleStatusArchived
    default:
        return ExampleStatusUnknown
    }
}

// In marshaller
Status: model.NewExampleStatus(e.Status),  // Safe conversion
```

## Computed Fields

Computed fields (like signed URLs) are initialized empty in the marshaller and set later by domain services:

```go
m := &model.Example{
    // ... fields ...
    ImageURL: null.String{},  // Set empty, service will populate
}
```

## Complete Example

```go
package marshaller

import (
    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/internal/dbmodel"
    "github.com/aarondl/null/v8"
)

func StaffToModel(e *dbmodel.Staff) *model.Staff {
    m := &model.Staff{
        ID:          e.ID,
        TenantID:    e.TenantID,
        Role:        model.NewStaffRole(e.Role),
        AuthUID:     e.AuthUID,
        DisplayName: e.DisplayName,
        ImagePath:   e.ImagePath,
        Email:       e.Email,
        CreatedAt:   e.CreatedAt,
        UpdatedAt:   e.UpdatedAt,

        ImageURL:          null.String{},
        ReadonlyReference: nil,
    }

    if e.R != nil {
        var tenant *model.Tenant
        if e.R.Tenant != nil {
            tenant = TenantToModel(e.R.Tenant)
        }
        if tenant != nil {
            m.ReadonlyReference = &struct {
                Tenant *model.Tenant
            }{
                Tenant: tenant,
            }
        }
    }

    return m
}

func StaffsToModel(slice dbmodel.StaffSlice) model.Staffs {
    dsts := make(model.Staffs, len(slice))
    for idx, e := range slice {
        dsts[idx] = StaffToModel(e)
    }
    return dsts
}

func StaffToDBModel(m *model.Staff) *dbmodel.Staff {
    return &dbmodel.Staff{
        ID:          m.ID,
        TenantID:    m.TenantID,
        Role:        m.Role.String(),
        AuthUID:     m.AuthUID,
        DisplayName: m.DisplayName,
        ImagePath:   m.ImagePath,
        Email:       m.Email,
        CreatedAt:   m.CreatedAt,
        UpdatedAt:   m.UpdatedAt,
        R:           nil,
        L:           struct{}{},
    }
}

func StaffsToDBModel(m model.Staffs) dbmodel.StaffSlice {
    dsts := make(dbmodel.StaffSlice, len(m))
    for idx, e := range m {
        dsts[idx] = StaffToDBModel(e)
    }
    return dsts
}
```
