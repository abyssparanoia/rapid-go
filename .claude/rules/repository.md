---
description: Repository interface, implementation, and marshaller patterns
globs:
  - "internal/domain/repository/**/*.go"
  - "internal/infrastructure/mysql/repository/**/*.go"
  - "internal/infrastructure/mysql/internal/marshaller/**/*.go"
  - "internal/infrastructure/postgresql/repository/**/*.go"
  - "internal/infrastructure/postgresql/internal/marshaller/**/*.go"
  - "internal/infrastructure/spanner/repository/**/*.go"
  - "internal/infrastructure/spanner/internal/marshaller/**/*.go"
---

# Repository Guidelines

## Interface Definition (Domain Layer)

Location: `internal/domain/repository/{entity}.go`

```go
package repository

import (
    "context"

    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/volatiletech/null/v8"
)

//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Staff interface {
    Get(
        ctx context.Context,
        query GetStaffQuery,
    ) (*model.Staff, error)
    List(
        ctx context.Context,
        query ListStaffQuery,
    ) (model.Staffs, error)
    Count(
        ctx context.Context,
        query ListStaffQuery,
    ) (uint64, error)
    Create(
        ctx context.Context,
        staff *model.Staff,
    ) error
    BatchCreate(
        ctx context.Context,
        staffs model.Staffs,
    ) error
    Update(
        ctx context.Context,
        staff *model.Staff,
    ) error
    Delete(
        ctx context.Context,
        id string,
    ) error
}
```

### Always include `//go:generate` directive for mock generation

Note: Use `go tool` syntax for go generate directive.

## Query Structs

```go
type GetStaffQuery struct {
    BaseGetOptions
    ID      null.String
    AuthUID null.String
}

type ListStaffQuery struct {
    BaseListOptions
    TenantID null.String
}
```

### Optional Fields: Use `nullable.Type[T]` Instead of Pointers

For optional filter fields with custom types (enums, domain types), always use `nullable.Type[T]` instead of pointers:

```go
// Good - Use nullable.Type for optional enum/custom type fields
type ListExamplesQuery struct {
    Status  nullable.Type[model.ExampleStatus]
    SortKey nullable.Type[model.ExampleSortKey]
    Role    nullable.Type[model.AdminRole]
}

// Bad - Avoid pointers for optional filter fields
type ListExamplesQuery struct {
    Status  *model.ExampleStatus   // Don't use pointers
    SortKey *model.ExampleSortKey  // Don't use pointers
    Role    *model.AdminRole       // Don't use pointers
}
```

**Why `nullable.Type[T]`:**
- Consistent with codebase conventions
- Provides `.Valid` and `.Value()` methods for safer access
- Works seamlessly with validation patterns in repository implementations
- Avoids nil pointer dereference risks

**When to use which:**
| Type | Use Case |
|------|----------|
| `null.String` | Optional string fields (IDs, names) |
| `null.Uint64` | Optional numeric fields (pagination) |
| `nullable.Type[T]` | Optional enum/custom type fields (status, role, sort key) |

### Base Options (defined in `base_options.go`)

```go
package repository

import "github.com/volatiletech/null/v8"

type BaseGetOptions struct {
    OrFail     bool  // Return error if not found (vs nil)
    Preload    bool  // Load relations
    ForUpdate  bool  // SELECT FOR UPDATE
    SkipLocked bool  // SKIP LOCKED
}

type BaseListOptions struct {
    Page       null.Uint64
    Limit      null.Uint64
    Preload    bool
    ForUpdate  bool
    SkipLocked bool
}
```

## Implementation (Infrastructure Layer)

Location: `internal/infrastructure/{mysql|postgresql|spanner}/repository/{entity}.go`

```go
package repository

type example struct{}

func NewExample() repository.Example {
    return &example{}
}
```

### Get Method Pattern

```go
func (r *example) Get(ctx context.Context, query repository.GetExampleQuery) (*model.Example, error) {
    mods := []qm.QueryMod{}

    // Build query conditions
    if query.ID.Valid {
        mods = append(mods, dbmodel.ExampleWhere.ID.EQ(query.ID.String))
    }
    if query.TenantID.Valid {
        mods = append(mods, dbmodel.ExampleWhere.TenantID.EQ(query.TenantID.String))
    }

    // Handle options
    if query.ForUpdate {
        mods = append(mods, qm.For("UPDATE"))
    }
    if query.Preload {
        mods = r.addPreload(mods)
    }

    dbEntity, err := dbmodel.Examples(mods...).One(ctx, transactable.GetContextExecutor(ctx))
    if err != nil {
        if err == sql.ErrNoRows && !query.OrFail {
            return nil, nil  // Not found, but OrFail=false
        } else if err == sql.ErrNoRows {
            return nil, errors.ExampleNotFoundErr.Errorf("example not found")
        }
        return nil, errors.InternalErr.Wrap(err)
    }
    return marshaller.ExampleToModel(dbEntity), nil
}
```

### List Method Pattern

```go
func (r *example) List(ctx context.Context, query repository.ListExamplesQuery) (model.Examples, error) {
    mods := r.buildListQuery(query)

    // Pagination
    if query.Page.Valid && query.Limit.Valid {
        mods = append(mods,
            qm.Limit(int(query.Limit.Uint64)),
            qm.Offset(int(query.Limit.Uint64*(query.Page.Uint64-1))),
        )
    }

    // Sorting
    if query.SortKey.Valid && query.SortKey.Ptr().Valid() {
        switch query.SortKey.Value() {
        case model.ExampleSortKeyCreatedAtDesc:
            mods = append(mods, qm.OrderBy("\"created_at\" DESC"))
        case model.ExampleSortKeyNameAsc:
            mods = append(mods, qm.OrderBy("\"name\" ASC"))
        }
    }

    // Preload
    if query.Preload {
        mods = r.addPreload(mods)
    }

    dbEntities, err := dbmodel.Examples(mods...).All(ctx, transactable.GetContextExecutor(ctx))
    if err != nil {
        return nil, errors.InternalErr.Wrap(err)
    }
    return marshaller.ExamplesToModel(dbEntities), nil
}

func (r *example) buildListQuery(query repository.ListExamplesQuery) []qm.QueryMod {
    mods := []qm.QueryMod{}
    if query.TenantID.Valid {
        mods = append(mods, dbmodel.ExampleWhere.TenantID.EQ(query.TenantID.String))
    }
    if query.Status.Valid && query.Status.Value().Valid() {
        mods = append(mods, dbmodel.ExampleWhere.Status.EQ(query.Status.Value().String()))
    }
    return mods
}
```

### Preload Helper

```go
func (r *example) addPreload(mods []qm.QueryMod) []qm.QueryMod {
    mods = append(mods, qm.Load(dbmodel.ExampleRels.Tenant))
    // Add more relations as needed
    return mods
}
```

## Marshaller (Infrastructure Layer)

Location: `internal/infrastructure/{mysql|postgresql|spanner}/internal/marshaller/{entity}.go`

### DB Model → Domain Model

```go
func ExampleToModel(e *dbmodel.Example) *model.Example {
    m := &model.Example{
        ID:        e.ID,
        TenantID:  e.TenantID,
        Name:      e.Name,
        Status:    model.ExampleStatus(e.Status),
        CreatedAt: e.CreatedAt,
        UpdatedAt: e.UpdatedAt,
        ReadonlyReference: nil,
    }

    // Handle relations
    if e.R != nil && e.R.Tenant != nil {
        m.ReadonlyReference = &struct{ Tenant *model.Tenant }{
            Tenant: TenantToModel(e.R.Tenant),
        }
    }
    return m
}
```

### ReadonlyReference Marshalling Rules

**IMPORTANT**: When populating `ReadonlyReference`, the related entity's own `ReadonlyReference` must remain `nil`.

```go
// Correct - Related entity's ReadonlyReference is nil
func ExampleToModel(e *dbmodel.Example) *model.Example {
    m := &model.Example{
        ID:                e.ID,
        ReadonlyReference: nil,
    }

    if e.R != nil && e.R.Tenant != nil {
        // TenantToModel returns Tenant with ReadonlyReference = nil
        m.ReadonlyReference = &struct{ Tenant *model.Tenant }{
            Tenant: TenantToModel(e.R.Tenant),
        }
    }
    return m
}

// TenantToModel - ReadonlyReference is always nil (no recursive loading)
func TenantToModel(s *dbmodel.Tenant) *model.Tenant {
    return &model.Tenant{
        ID:                s.ID,
        Name:              s.Name,
        ReadonlyReference: nil,  // Always nil - no recursive loading
    }
}
```

**Why no recursive ReadonlyReference:**
- Prevents circular dependencies (A → B → A)
- Reduces memory usage and query complexity
- Related entities are for display purposes only, not for navigation

### Nullable Timestamp Fields

For optional timestamp fields (`null.Time`), use var declaration pattern to prevent field mapping omissions:

```go
func InvitationToModel(i *dbmodel.Invitation) *model.Invitation {
    // Declare nullable fields first
    var acceptedAt null.Time
    if i.AcceptedAt.Valid {
        acceptedAt = null.TimeFrom(i.AcceptedAt.Time)
    }

    var rejectedAt null.Time
    if i.RejectedAt.Valid {
        rejectedAt = null.TimeFrom(i.RejectedAt.Time)
    }

    return &model.Invitation{
        ID:         i.ID,
        Status:     model.InvitationStatus(i.Status),
        AcceptedAt: acceptedAt,
        RejectedAt: rejectedAt,
        CreatedAt:  i.CreatedAt,
        UpdatedAt:  i.UpdatedAt,
        ReadonlyReference: nil,
    }
}
```

**Why var declaration:**
- Ensures all fields are explicitly handled
- Prevents accidentally omitting fields in struct literal
- Makes nullable field handling visible and reviewable

### Slice Version

```go
func ExamplesToModel(slice dbmodel.ExampleSlice) model.Examples {
    dsts := make(model.Examples, len(slice))
    for idx, e := range slice {
        dsts[idx] = ExampleToModel(e)
    }
    return dsts
}
```

### Domain Model → DB Model

```go
func ExampleToDBModel(m *model.Example) *dbmodel.Example {
    return &dbmodel.Example{
        ID:        m.ID,
        TenantID:  m.TenantID,
        Name:      m.Name,
        Status:    m.Status.String(),
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
        R:         nil,  // Relations not written
        L:         struct{}{},
    }
}
```

**Note**: `ReadonlyReference` is never converted back to DB model - it's read-only.

## Transaction Context

Always use `transactable.GetContextExecutor(ctx)` to get the database executor:

```go
dbmodel.Examples(mods...).One(ctx, transactable.GetContextExecutor(ctx))
```

This ensures queries run within the transaction if one is active.
