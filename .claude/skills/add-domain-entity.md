---
name: add-domain-entity
description: REQUIRED Step 2 of CRUD workflow (after add-database-table). Use when creating domain models in internal/domain/model/, repository interfaces, marshallers, or repository implementations.
---

# Add Domain Entity

This skill guides you through creating domain layer components for a new entity.

## Prerequisites

- Database table already created (use **add-database-table** skill first)
- SQLBoiler model generated in `internal/infrastructure/{mysql|postgresql|spanner}/internal/dbmodel/`

## Step 1: Create Domain Model

Location: `internal/domain/model/{entity}.go`

```go
package model

import (
    "time"
    "github.com/abyssparanoia/rapid-go/internal/pkg/id"
    "github.com/volatiletech/null/v8"
)

type Example struct {
    ID          string
    TenantID    string
    Name        string
    Description string
    Status      ExampleStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time

    // Relations (read-only)
    ReadonlyReference *struct {
        Tenant *Tenant
    }
}

// Type aliases
type ExampleMapByID map[string]*Example
type Examples []*Example

// Constructor
func NewExample(
    tenantID string,
    name string,
    description string,
    t time.Time,
) *Example {
    return &Example{
        ID:                id.New(),
        TenantID:          tenantID,
        Name:              name,
        Description:       description,
        Status:            ExampleStatusDraft,
        CreatedAt:         t,
        UpdatedAt:         t,
        ReadonlyReference: nil,
    }
}

// Update method
func (e *Example) Update(
    name null.String,
    description null.String,
    t time.Time,
) *Example {
    if name.Valid {
        e.Name = name.String
    }
    if description.Valid {
        e.Description = description.String
    }
    e.UpdatedAt = t
    return e
}

// Helper methods
func (es Examples) IDs() []string {
    ids := make([]string, 0, len(es))
    for _, e := range es {
        ids = append(ids, e.ID)
    }
    return ids
}

func (es Examples) MapByID() ExampleMapByID {
    m := make(ExampleMapByID, len(es))
    for _, e := range es {
        m[e.ID] = e
    }
    return m
}

// Status type
type ExampleStatus string

const (
    ExampleStatusUnknown   ExampleStatus = "unknown"
    ExampleStatusDraft     ExampleStatus = "draft"
    ExampleStatusPublished ExampleStatus = "published"
    ExampleStatusArchived  ExampleStatus = "archived"
)

func (s ExampleStatus) String() string { return string(s) }
func (s ExampleStatus) Valid() bool {
    return s != ExampleStatusUnknown && s != ""
}

// Sort key
type ExampleSortKey string

const (
    ExampleSortKeyUnknown       ExampleSortKey = "unknown"
    ExampleSortKeyCreatedAtDesc ExampleSortKey = "created_at_desc"
    ExampleSortKeyNameAsc       ExampleSortKey = "name_asc"
)

func (k ExampleSortKey) Valid() bool {
    return k != ExampleSortKeyUnknown && k != ""
}
```

## Step 2: Add Domain Error

Location: `internal/domain/errors/errors.go`

Add not found error for the entity:

```go
var (
    // ... existing errors ...
    // NOTE: Pick a stable error code following the repo convention (see .claude/rules/domain-errors.md)
    ExampleNotFoundErr = NewNotFoundError("E2xxxxx", "Example not found")
)
```

## Step 3: Create Repository Interface

Location: `internal/domain/repository/{entity}.go`

```go
package repository

import (
    "context"
    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
    "github.com/volatiletech/null/v8"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Example interface {
    Get(ctx context.Context, query GetExampleQuery) (*model.Example, error)
    BatchGet(ctx context.Context, ids []string, query BatchGetExamplesQuery) (model.ExampleMapByID, error)
    List(ctx context.Context, query ListExamplesQuery) (model.Examples, error)
    Count(ctx context.Context, query ListExamplesQuery) (uint64, error)
    Create(ctx context.Context, example *model.Example) error
    Update(ctx context.Context, example *model.Example) error
    Delete(ctx context.Context, id string) error
}

type GetExampleQuery struct {
    ID       null.String
    TenantID null.String
    BaseGetOptions
}

type BatchGetExamplesQuery struct {
    TenantID null.String
    BaseBatchGetOptions
}

type ListExamplesQuery struct {
    TenantID null.String
    Status   nullable.Type[model.ExampleStatus]
    SortKey  nullable.Type[model.ExampleSortKey]
    BaseListOptions
}
```

## Step 4: Create Marshaller

Location: `internal/infrastructure/{mysql|postgresql|spanner}/internal/marshaller/{entity}.go`

### Basic Pattern

```go
package marshaller

import (
    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/internal/dbmodel"
)

func ExampleToModel(e *dbmodel.Example) *model.Example {
    m := &model.Example{
        ID:                e.ID,
        TenantID:          e.TenantID,
        Name:              e.Name,
        Description:       e.Description,
        Status:            model.ExampleStatus(e.Status),
        CreatedAt:         e.CreatedAt,
        UpdatedAt:         e.UpdatedAt,
        ReadonlyReference: nil,
    }

    // Populate ReadonlyReference if relations are loaded
    // IMPORTANT: Related entity's ReadonlyReference must remain nil
    if e.R != nil && e.R.Tenant != nil {
        m.ReadonlyReference = &struct{ Tenant *model.Tenant }{
            Tenant: TenantToModel(e.R.Tenant),  // Tenant.ReadonlyReference = nil
        }
    }
    return m
}

func ExamplesToModel(slice dbmodel.ExampleSlice) model.Examples {
    dsts := make(model.Examples, len(slice))
    for idx, e := range slice {
        dsts[idx] = ExampleToModel(e)
    }
    return dsts
}

func ExampleToDBModel(m *model.Example) *dbmodel.Example {
    return &dbmodel.Example{
        ID:          m.ID,
        TenantID:    m.TenantID,
        Name:        m.Name,
        Description: m.Description,
        Status:      m.Status.String(),
        CreatedAt:   m.CreatedAt,
        UpdatedAt:   m.UpdatedAt,
        R:           nil,  // ReadonlyReference is never written back
        L:           struct{}{},
    }
}
```

### Pattern for Nullable Timestamp Fields

For entities with optional timestamp fields (`null.Time`), use var declaration pattern:

```go
func InvitationToModel(i *dbmodel.Invitation) *model.Invitation {
    // Declare nullable fields first to prevent omission
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

### ReadonlyReference Rules

1. **Always initialize to nil** in the struct literal
2. **Populate only if relations are loaded** (`e.R != nil && e.R.Relation != nil`)
3. **Related entity's ReadonlyReference must be nil** - no recursive loading
4. **Never write back** - ReadonlyReference is ignored in ToDBModel

## Step 5: Create Repository Implementation

Location: `internal/infrastructure/{mysql|postgresql|spanner}/repository/{entity}.go`

```go
package repository

import (
    "context"
    "database/sql"

    "github.com/abyssparanoia/rapid-go/internal/domain/errors"
    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/domain/repository"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/internal/dbmodel"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/internal/marshaller"
    "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/transactable"
    "github.com/volatiletech/sqlboiler/v4/boil"
    "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type example struct{}

func NewExample() repository.Example {
    return &example{}
}

func (r *example) Get(ctx context.Context, query repository.GetExampleQuery) (*model.Example, error) {
    mods := []qm.QueryMod{}
    if query.ID.Valid {
        mods = append(mods, dbmodel.ExampleWhere.ID.EQ(query.ID.String))
    }
    if query.TenantID.Valid {
        mods = append(mods, dbmodel.ExampleWhere.TenantID.EQ(query.TenantID.String))
    }
    if query.ForUpdate {
        mods = append(mods, qm.For("UPDATE"))
    }
    if query.Preload {
        mods = r.addPreload(mods)
    }

    dbExample, err := dbmodel.Examples(mods...).One(ctx, transactable.GetContextExecutor(ctx))
    if err != nil {
        if err == sql.ErrNoRows && !query.OrFail {
            return nil, nil
        } else if err == sql.ErrNoRows {
            return nil, errors.ExampleNotFoundErr.Errorf("example not found")
        }
        return nil, errors.InternalErr.Wrap(err)
    }
    return marshaller.ExampleToModel(dbExample), nil
}

func (r *example) List(ctx context.Context, query repository.ListExamplesQuery) (model.Examples, error) {
    mods := r.buildListQuery(query)
    if query.Page.Valid && query.Limit.Valid {
        mods = append(mods,
            qm.Limit(int(query.Limit.Uint64)),
            qm.Offset(int(query.Limit.Uint64*(query.Page.Uint64-1))),
        )
    }
    if query.Preload {
        mods = r.addPreload(mods)
    }
    if query.SortKey.Valid && query.SortKey.Ptr().Valid() {
        switch query.SortKey.Value() {
        case model.ExampleSortKeyCreatedAtDesc:
            mods = append(mods, qm.OrderBy("\"created_at\" DESC"))
        case model.ExampleSortKeyNameAsc:
            mods = append(mods, qm.OrderBy("\"name\" ASC"))
        }
    }

    dbExamples, err := dbmodel.Examples(mods...).All(ctx, transactable.GetContextExecutor(ctx))
    if err != nil {
        return nil, errors.InternalErr.Wrap(err)
    }
    return marshaller.ExamplesToModel(dbExamples), nil
}

func (r *example) Count(ctx context.Context, query repository.ListExamplesQuery) (uint64, error) {
    mods := r.buildListQuery(query)
    ttl, err := dbmodel.Examples(mods...).Count(ctx, transactable.GetContextExecutor(ctx))
    if err != nil {
        return 0, errors.InternalErr.Wrap(err)
    }
    return uint64(ttl), nil
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

func (r *example) addPreload(mods []qm.QueryMod) []qm.QueryMod {
    mods = append(mods, qm.Load(dbmodel.ExampleRels.Tenant))
    return mods
}

func (r *example) BatchGet(ctx context.Context, ids []string, query repository.BatchGetExamplesQuery) (model.ExampleMapByID, error) {
    mods := []qm.QueryMod{dbmodel.ExampleWhere.ID.IN(ids)}
    if query.TenantID.Valid {
        mods = append(mods, dbmodel.ExampleWhere.TenantID.EQ(query.TenantID.String))
    }
    if query.Preload {
        mods = r.addPreload(mods)
    }

    dbExamples, err := dbmodel.Examples(mods...).All(ctx, transactable.GetContextExecutor(ctx))
    if err != nil {
        return nil, errors.InternalErr.Wrap(err)
    }

    result := make(map[string]*model.Example)
    for _, db := range dbExamples {
        result[db.ID] = marshaller.ExampleToModel(db)
    }
    return result, nil
}

func (r *example) Create(ctx context.Context, example *model.Example) error {
    dst := marshaller.ExampleToDBModel(example)
    if err := dst.Insert(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
        return errors.InternalErr.Wrap(err)
    }
    return nil
}

func (r *example) Update(ctx context.Context, example *model.Example) error {
    dst := marshaller.ExampleToDBModel(example)
    if _, err := dst.Update(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
        return errors.InternalErr.Wrap(err)
    }
    return nil
}

func (r *example) Delete(ctx context.Context, id string) error {
    dst := &dbmodel.Example{ID: id}
    if _, err := dst.Delete(ctx, transactable.GetContextExecutor(ctx)); err != nil {
        return errors.InternalErr.Wrap(err)
    }
    return nil
}
```

## Step 6: Generate Mocks

```bash
make generate.mock
```

## Checklist

- [ ] Domain model created with constructor and update method
- [ ] ReadonlyReference defined for relations (if any)
- [ ] Constructor sets ReadonlyReference to nil
- [ ] Status/enum types defined with Valid() method
- [ ] Domain error added
- [ ] Repository interface with go:generate directive
- [ ] Query structs defined with `nullable.Type[T]` for optional enum fields
- [ ] Marshaller created (both directions)
- [ ] Marshaller populates ReadonlyReference correctly (related entity's ReadonlyReference is nil)
- [ ] Marshaller uses var declaration pattern for nullable timestamp fields
- [ ] Repository implementation complete with addPreload helper
- [ ] Mocks generated

## Next Steps

After creating domain entity, use the **add-api-endpoint** skill to create:
- Usecase input/output
- Interactor
- Protocol Buffers definition
- gRPC handler
