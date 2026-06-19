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
    "github.com/aarondl/null/v8"
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
    SortKey  nullable.Type[model.StaffSortKey]
}
```

### SortKey in List Queries (Unified Specification)

**All List queries must include a SortKey field** using `nullable.Type[model.XXXSortKey]`:

```go
type ListTenantsQuery struct {
    BaseListOptions
    SortKey nullable.Type[model.TenantSortKey]
}

type ListStaffsQuery struct {
    BaseListOptions
    TenantID null.String
    SortKey  nullable.Type[model.StaffSortKey]
}
```

**Key points:**

- SortKey is **always** `nullable.Type[model.XXXSortKey]` in repository queries
- Input layer resolves nullable to non-nullable with default (`CreatedAtDesc`)
- Repository implementation checks `query.SortKey.Valid && query.SortKey.Value().Valid()` before applying sort

### Optional Fields: `null/v8` for Primitives, `nullable.Type[T]` for Custom Types

**Primitive and time fields use `null/v8` (`null.Int64`, `null.Bool`, `null.Float64`,
`null.Time`, `null.String`). Reserve `nullable.Type[T]` for custom domain types that `null/v8`
does NOT provide** — enums (`model.ExampleStatus`, `model.AdminRole`, sort keys), domain structs,
and `civil.Date`. Never use `nullable.Type[int64/bool/float64/time.Time/string]` — there is a
matching `null/v8` type for each. Optional fields must not use bare pointers.

```go
// Good - primitives via null/v8; custom types via nullable.Type[T]
type ListExamplesQuery struct {
    IsActive null.Bool                          // primitive → null/v8
    MinScore null.Int64                         // primitive → null/v8
    Status   nullable.Type[model.ExampleStatus] // custom enum → nullable.Type
    SortKey  nullable.Type[model.ExampleSortKey]
}

// Bad - nullable.Type used for a primitive (use null.Bool / null.Int64)
type ListExamplesQuery struct {
    IsActive nullable.Type[bool]  // WRONG → null.Bool
    MinScore nullable.Type[int64] // WRONG → null.Int64
}

// Bad - bare pointers for optional filter fields
type ListExamplesQuery struct {
    Status *model.ExampleStatus // WRONG → nullable.Type[model.ExampleStatus]
}
```

In repository impls, read primitive `null/v8` fields via the typed accessor (`.Bool`, `.Int64`,
`.Float64`), and custom `nullable.Type[T]` via `.Value()`:

```go
if query.IsActive.Valid {
    mods = append(mods, dbmodel.ExampleWhere.IsActive.EQ(query.IsActive.Bool))
}
if query.Status.Valid && query.Status.Value().Valid() {
    mods = append(mods, dbmodel.ExampleWhere.Status.EQ(query.Status.Value().String()))
}
```

**Why split this way:**

- `null/v8` is the canonical optional type for primitives/time across the codebase (DB models,
  domain models, inputs all use it); `nullable.Type[T]` exists only to cover types `null/v8`
  cannot represent (generics over custom types).
- Both provide `.Valid` for safe access and avoid nil-pointer dereference.

**When to use which:**
| Type | Use Case |
|------|----------|
| `null.String` | Optional string fields (IDs, names) |
| `null.Int64` | Optional integer fields |
| `null.Float64` | Optional decimal/rate fields |
| `null.Bool` | Optional boolean fields |
| `null.Time` | Optional timestamp fields |
| `null.Uint64` | Optional pagination fields |
| `nullable.Type[T]` | Optional **custom-type** fields only: enums (status, role, sort key), domain structs, `civil.Date` |

### Base Options (defined in `base_options.go`)

```go
package repository

import "github.com/aarondl/null/v8"

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

### IncludeDeleted Option on Get Queries

For entities with soft-delete (`deleted_at` column), add `IncludeDeleted bool` directly to the entity's `Get{Entity}Query` struct (not in `BaseGetOptions`, since soft-delete is entity-specific):

```go
type GetPaymentMethodQuery struct {
    BaseGetOptions
    ID             null.String
    ExternalID     null.String
    OrganizationID null.String
    IncludeDeleted bool  // When true, includes soft-deleted rows
}
```

**When to use `IncludeDeleted: true`:**
- Idempotency checks in webhook handlers — a soft-deleted row means the event was already processed (or the record was detached). Skip rather than re-create.
- Restore workflows — look up a soft-deleted record to restore it.

**Repository implementation pattern:**

```go
func (r *example) Get(ctx context.Context, query repository.GetExampleQuery) (*model.Example, error) {
    mods := []qm.QueryMod{}

    // Only filter out soft-deleted rows when IncludeDeleted is false
    if !query.IncludeDeleted {
        mods = append(mods, dbmodel.ExampleWhere.DeletedAt.IsNull())
    }

    // ... rest of query building
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

**IMPORTANT**: Sorting must be applied **BEFORE** pagination for correct SQL semantics.

#### MySQL Implementation

```go
func (r *example) List(ctx context.Context, query repository.ListExamplesQuery) (model.Examples, error) {
    mods := r.buildListQuery(query)

    // Sorting (BEFORE pagination)
    if query.SortKey.Valid && query.SortKey.Value().Valid() {
        switch query.SortKey.Value() {
        case model.ExampleSortKeyCreatedAtDesc:
            mods = append(mods, qm.OrderBy("`created_at` DESC"))
        case model.ExampleSortKeyCreatedAtAsc:
            mods = append(mods, qm.OrderBy("`created_at` ASC"))
        case model.ExampleSortKeyNameAsc:
            mods = append(mods, qm.OrderBy("`name` ASC"))
        case model.ExampleSortKeyNameDesc:
            mods = append(mods, qm.OrderBy("`name` DESC"))
        case model.ExampleSortKeyUnknown:
            return nil, errors.InternalErr.Errorf("invalid sort key: %s", query.SortKey.Value())
        }
    }

    // Pagination (AFTER sorting)
    if query.Page.Valid && query.Limit.Valid {
        mods = append(mods,
            qm.Limit(int(query.Limit.Uint64)),
            qm.Offset(int(query.Limit.Uint64*(query.Page.Uint64-1))),
        )
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

#### PostgreSQL Implementation

PostgreSQL uses double quotes instead of backticks for identifiers. Always reference `dbmodel.{Table}Columns.{Field}` constants instead of hardcoding column name strings — this prevents silent breaks when columns are renamed.

```go
func (r *example) List(ctx context.Context, query repository.ListExamplesQuery) (model.Examples, error) {
    mods := r.buildListQuery(query)

    // Sorting (BEFORE pagination)
    // IMPORTANT: Use dbmodel column constants, never hardcode string literals
    if query.SortKey.Valid && query.SortKey.Value().Valid() {
        switch query.SortKey.Value() {
        case model.ExampleSortKeyCreatedAtDesc:
            mods = append(mods, qm.OrderBy("\""+dbmodel.ExampleColumns.CreatedAt+"\" DESC"))
        case model.ExampleSortKeyCreatedAtAsc:
            mods = append(mods, qm.OrderBy("\""+dbmodel.ExampleColumns.CreatedAt+"\" ASC"))
        case model.ExampleSortKeyNameAsc:
            mods = append(mods, qm.OrderBy("\""+dbmodel.ExampleColumns.Name+"\" ASC"))
        case model.ExampleSortKeyNameDesc:
            mods = append(mods, qm.OrderBy("\""+dbmodel.ExampleColumns.Name+"\" DESC"))
        case model.ExampleSortKeyUnknown:
            return nil, errors.InternalErr.Errorf("invalid sort key: %s", query.SortKey.Value())
        }
    }

    // Pagination (AFTER sorting)
    if query.Page.Valid && query.Limit.Valid {
        mods = append(mods,
            qm.Limit(int(query.Limit.Uint64)),
            qm.Offset(int(query.Limit.Uint64*(query.Page.Uint64-1))),
        )
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
```

#### Key Differences Between Databases

| Database   | Identifier Quoting | Example                 |
| ---------- | ------------------ | ----------------------- |
| MySQL      | Backticks          | `` `created_at` DESC `` |
| PostgreSQL | Double quotes      | `"created_at" DESC`     |
| Spanner    | Backticks          | `` `created_at` DESC `` |

#### Column Name References — Use Generated Constants

When passing raw column name strings to `qm.OrderBy` / `qm.Where`, **always reference the SQLBoiler-generated `dbmodel.{Entity}Columns` constants** instead of hardcoding the column name.

Why:
- If the column is renamed or removed, hardcoded strings fail silently — only constant references surface as a build error.
- Hardcoded strings can only be tracked by grep, which makes rename refactors fragile.
- Constants enable IDE autocomplete and prevent typos.

```go
// GOOD - reference dbmodel.ExampleColumns
mods = append(mods, qm.OrderBy(fmt.Sprintf("`%s` ASC, `%s` ASC",
    dbmodel.ExampleColumns.CreatedAt,
    dbmodel.ExampleColumns.ID,
)))

// GOOD - composite WHERE clauses follow the same pattern (e.g. tuple-compare cursor pagination)
mods = append(mods, qm.Where(
    fmt.Sprintf("(`%s` > ? OR (`%s` = ? AND `%s` > ?))",
        dbmodel.ExampleColumns.CreatedAt,
        dbmodel.ExampleColumns.CreatedAt,
        dbmodel.ExampleColumns.ID,
    ),
    cursorTime, cursorTime, cursorID,
))

// BAD - hardcoded raw strings
mods = append(mods, qm.OrderBy("`created_at` ASC, `id` ASC"))
mods = append(mods, qm.Where("(`created_at` > ? OR (`created_at` = ? AND `id` > ?))", ...))
```

For JOINed queries that need the `table.column` form, use `dbmodel.{Entity}TableColumns` instead (it yields strings like `examples.created_at`).

Prefer the typed builder API (`dbmodel.{Entity}Where.{Field}.EQ(...)`, `.GT(...)`, etc.) whenever it covers the predicate. Reach for `qm.Where` with raw SQL only for cases the builder cannot express — e.g. tuple-compare cursor pagination expanded as `(a > ? OR (a = ? AND b > ?))`.

#### Unknown SortKey Handling

**IMPORTANT**: Always return an error for Unknown SortKey values. Do not silently skip sorting.

```go
case model.ExampleSortKeyUnknown:
    return nil, errors.InternalErr.Errorf("invalid sort key: %s", query.SortKey.Value())
```

**Anti-pattern** (do not use):

```go
case model.ExampleSortKeyUnknown:
    // No sorting applied for unknown  // WRONG - should return error
```

### Preload Helper

```go
func (r *example) addPreload(mods []qm.QueryMod) []qm.QueryMod {
    mods = append(mods, qm.Load(dbmodel.ExampleRels.Tenant))
    // Add more relations as needed
    return mods
}
```

### Owned Children vs Reference Relations — Unconditional Preload

`Preload bool` in `BaseGetOptions` / `BaseListOptions` controls loading of **reference** relations (e.g., `ReadonlyReference.Tenant`). It does NOT apply to **fully-owned child entities**.

| Relation Type | Load Condition | Example |
|---|---|---|
| Reference (lookup) | Conditional — only when `Preload: true` | `InvoiceRels.Organization` |
| Owned 1:N child | **Unconditional — always load** | `InvoiceRels.InvoiceItems` |
| Owned 1:1 child | **Unconditional — always load** | `InvoiceRels.InvoiceStripe` |

**Fully-owned children** (entities that cannot exist without their parent and are part of the same aggregate) **must always be loaded** regardless of the `Preload` flag. Never use a boolean parameter to gate their loading:

```go
// BAD - conditional preload of owned child
func (r *invoice) buildPreload(preloadStripe bool) []qm.QueryMod {
    mods := []qm.QueryMod{qm.Load(dbmodel.InvoiceRels.InvoiceItems)}
    if preloadStripe {
        mods = append(mods, qm.Load(dbmodel.InvoiceRels.InvoiceStripe))
    }
    return mods
}

// GOOD - always load owned children
func (r *invoice) buildPreload() []qm.QueryMod {
    return []qm.QueryMod{
        qm.Load(dbmodel.InvoiceRels.InvoiceItems),
        qm.Load(dbmodel.InvoiceRels.InvoiceStripe),
        qm.Load(fmt.Sprintf("%s.%s", dbmodel.InvoiceRels.InvoiceItems, dbmodel.InvoiceItemRels.InvoiceItemStripe)),
    }
}
```

**Why**: The API handler may not surface the owned children today, but the domain model is incomplete without them. Conditional loading causes nil-dereference bugs and makes behavior dependent on call-site flags.

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

### Struct Literal Return Rule

All conversion functions must return the struct using a struct literal with all fields explicitly listed. Do **not** initialize an empty struct and then assign fields one by one.

```go
// GOOD - struct literal return
func ExampleToModel(e *dbmodel.Example) *model.Example {
    return &model.Example{
        ID:        e.ID,
        TenantID:  e.TenantID,
        Name:      e.Name,
        Status:    model.ExampleStatus(e.Status),
        CreatedAt: e.CreatedAt,
        UpdatedAt: e.UpdatedAt,
        ReadonlyReference: nil,
    }
}

// BAD - field-by-field assignment on empty struct
func ExampleToModel(e *dbmodel.Example) *model.Example {
    m := &model.Example{}
    m.ID = e.ID
    m.TenantID = e.TenantID
    m.Name = e.Name
    return m
}
```

**Exception**: When conditional field assignment is required (e.g., `ReadonlyReference` populated only when `e.R != nil`, or nullable timestamps), use the `var` declaration pattern described above, then build the struct literal with the prepared variables.

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
