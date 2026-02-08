# Repository Patterns

Detailed code patterns for repository interfaces and implementations.

## Interface Definition

Location: `internal/domain/repository/{entity}.go`

```go
package repository

import (
    "context"

    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
    "github.com/aarondl/null/v9"
)

//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Example interface {
    Get(ctx context.Context, query GetExampleQuery) (*model.Example, error)
    List(ctx context.Context, query ListExamplesQuery) (model.Examples, error)
    Count(ctx context.Context, query ListExamplesQuery) (uint64, error)
    Create(ctx context.Context, example *model.Example) error
    Update(ctx context.Context, example *model.Example) error
    Delete(ctx context.Context, id string) error
}
```

### Optional Methods

```go
// BatchGet - for loading multiple entities by ID
BatchGet(ctx context.Context, ids []string, query BatchGetExamplesQuery) (model.ExampleMapByID, error)

// BatchCreate - for bulk inserts
BatchCreate(ctx context.Context, examples model.Examples) error
```

## Query Structs

### Get Query

```go
type GetExampleQuery struct {
    BaseGetOptions        // Embed base options
    ID       null.String  // Primary lookup
    TenantID null.String  // Scope filter
    AuthUID  null.String  // Alternative lookup
}
```

### List Query

```go
type ListExamplesQuery struct {
    BaseListOptions                              // Embed base options
    TenantID null.String                         // Required scope
    Status   nullable.Type[model.ExampleStatus]  // Optional enum filter
    SortKey  nullable.Type[model.ExampleSortKey] // Optional sort
}
```

### BatchGet Query

```go
type BatchGetExamplesQuery struct {
    BaseBatchGetOptions   // Embed base options
    TenantID null.String  // Scope filter
}
```

### Optional Field Types

| Use Case                   | Type                        |
| -------------------------- | --------------------------- |
| Optional string (ID, name) | `null.String`               |
| Optional number            | `null.Uint64`, `null.Int64` |
| Optional enum/custom type  | `nullable.Type[T]`          |

## Base Options (defined in base_options.go)

```go
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

type BaseBatchGetOptions struct {
    Preload    bool
    ForUpdate  bool
    SkipLocked bool
}
```

---

## Implementation

Location: `internal/infrastructure/{mysql|postgresql|spanner}/repository/{entity}.go`

### Struct and Constructor

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
    "github.com/aarondl/sqlboiler/v4/boil"
    "github.com/aarondl/sqlboiler/v4/queries/qm"
)

type example struct{}

func NewExample() repository.Example {
    return &example{}
}
```

### Get Method

```go
func (r *example) Get(
    ctx context.Context,
    query repository.GetExampleQuery,
) (*model.Example, error) {
    mods := []qm.QueryMod{}

    // Build WHERE conditions
    if query.ID.Valid {
        mods = append(mods, dbmodel.ExampleWhere.ID.EQ(query.ID.String))
    }
    if query.TenantID.Valid {
        mods = append(mods, dbmodel.ExampleWhere.TenantID.EQ(query.TenantID.String))
    }

    // Apply options
    mods = append(mods, r.buildPreload(query.Preload)...)
    mods = addForUpdateFromBaseGetOptions(mods, query.BaseGetOptions)

    dbExample, err := dbmodel.Examples(mods...).One(ctx, transactable.GetContextExecutor(ctx))
    if err != nil {
        if err == sql.ErrNoRows && !query.OrFail {
            return nil, nil  // Not found, OrFail=false
        } else if err == sql.ErrNoRows {
            return nil, errors.ExampleNotFoundErr.New().
                WithDetail("example is not found").
                WithValue("query", query)
        }
        return nil, errors.InternalErr.Wrap(err)
    }
    return marshaller.ExampleToModel(dbExample), nil
}
```

### List Method

```go
func (r *example) List(
    ctx context.Context,
    query repository.ListExamplesQuery,
) (model.Examples, error) {
    mods := []qm.QueryMod{}
    mods = append(mods, r.buildListQuery(query)...)

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
        case model.ExampleSortKeyCreatedAtAsc:
            mods = append(mods, qm.OrderBy("\"created_at\" ASC"))
        case model.ExampleSortKeyNameAsc:
            mods = append(mods, qm.OrderBy("\"name\" ASC"))
        case model.ExampleSortKeyNameDesc:
            mods = append(mods, qm.OrderBy("\"name\" DESC"))
        }
    }

    // Preload relations
    mods = append(mods, r.buildPreload(query.Preload)...)
    mods = addForUpdateFromBaseListOptions(mods, query.BaseListOptions)

    dbExamples, err := dbmodel.Examples(mods...).All(ctx, transactable.GetContextExecutor(ctx))
    if err != nil {
        return nil, errors.InternalErr.Wrap(err)
    }
    return marshaller.ExamplesToModel(dbExamples), nil
}
```

### Count Method

```go
func (r *example) Count(
    ctx context.Context,
    query repository.ListExamplesQuery,
) (uint64, error) {
    mods := []qm.QueryMod{}
    mods = append(mods, r.buildListQuery(query)...)

    ttl, err := dbmodel.Examples(mods...).Count(ctx, transactable.GetContextExecutor(ctx))
    if err != nil {
        return 0, errors.InternalErr.Wrap(err)
    }
    return uint64(ttl), nil
}
```

### Helper: buildListQuery

Reusable filter logic shared between List and Count:

```go
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

### Helper: buildPreload

```go
func (r *example) buildPreload(preload bool) []qm.QueryMod {
    if !preload {
        return nil
    }
    return []qm.QueryMod{
        qm.Load(dbmodel.ExampleRels.Tenant),
        // Add more relations as needed
    }
}
```

### Create Method

```go
func (r *example) Create(
    ctx context.Context,
    example *model.Example,
) error {
    dst := marshaller.ExampleToDBModel(example)
    if err := dst.Insert(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
        return errors.InternalErr.Wrap(err)
    }
    return nil
}
```

### BatchCreate Method

```go
func (r *example) BatchCreate(
    ctx context.Context,
    examples model.Examples,
) error {
    dsts := marshaller.ExamplesToDBModel(examples)
    if _, err := dsts.InsertAll(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
        return errors.InternalErr.Wrap(err)
    }
    return nil
}
```

### Update Method

```go
func (r *example) Update(
    ctx context.Context,
    example *model.Example,
) error {
    dst := marshaller.ExampleToDBModel(example)
    if _, err := dst.Update(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
        return errors.InternalErr.Wrap(err)
    }
    return nil
}
```

### Delete Method

```go
func (r *example) Delete(
    ctx context.Context,
    id string,
) error {
    dst := marshaller.ExampleToDBModel(&model.Example{ID: id}) //nolint:exhaustruct
    if _, err := dst.Delete(ctx, transactable.GetContextExecutor(ctx)); err != nil {
        return errors.InternalErr.Wrap(err)
    }
    return nil
}
```

### BatchGet Method

```go
func (r *example) BatchGet(
    ctx context.Context,
    ids []string,
    query repository.BatchGetExamplesQuery,
) (model.ExampleMapByID, error) {
    if len(ids) == 0 {
        return make(model.ExampleMapByID), nil
    }

    mods := []qm.QueryMod{
        dbmodel.ExampleWhere.ID.IN(ids),
    }
    if query.TenantID.Valid {
        mods = append(mods, dbmodel.ExampleWhere.TenantID.EQ(query.TenantID.String))
    }
    mods = append(mods, r.buildPreload(query.Preload)...)

    dbExamples, err := dbmodel.Examples(mods...).All(ctx, transactable.GetContextExecutor(ctx))
    if err != nil {
        return nil, errors.InternalErr.Wrap(err)
    }

    result := make(model.ExampleMapByID, len(dbExamples))
    for _, db := range dbExamples {
        result[db.ID] = marshaller.ExampleToModel(db)
    }
    return result, nil
}
```

## Transaction Context

Always use `transactable.GetContextExecutor(ctx)` to get the database executor:

```go
dbmodel.Examples(mods...).One(ctx, transactable.GetContextExecutor(ctx))
```

This ensures queries run within the transaction if one is active.

## Error Handling Pattern

```go
if err != nil {
    if err == sql.ErrNoRows && !query.OrFail {
        return nil, nil  // Return nil, nil when not found and OrFail=false
    } else if err == sql.ErrNoRows {
        return nil, errors.ExampleNotFoundErr.New().
            WithDetail("example is not found").
            WithValue("query", query)
    }
    return nil, errors.InternalErr.Wrap(err)  // Wrap other errors
}
```
