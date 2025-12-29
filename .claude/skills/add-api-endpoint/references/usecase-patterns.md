# Usecase Patterns Reference

Detailed code patterns for usecase layer implementation.

## Input Structs

Location: `internal/usecase/input/{actor}_{entity}.go`

### Create Input

```go
package input

import (
    "time"
    "github.com/abyssparanoia/rapid-go/internal/domain/errors"
    "github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

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
```

### Get Input

```go
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
```

### List Input

```go
import "github.com/abyssparanoia/rapid-go/internal/domain/model"

type AdminListExamples struct {
    StaffID     string `validate:"required"`
    TenantID    string `validate:"required"`
    Status      *model.ExampleStatus   // Optional filter
    SortKey     *model.ExampleSortKey  // Optional sort
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
```

### Update Input

```go
import "github.com/volatiletech/null/v8"

type AdminUpdateExample struct {
    StaffID     string      `validate:"required"`
    TenantID    string      `validate:"required"`
    ExampleID   string      `validate:"required"`
    Name        null.String // Optional update field
    Description null.String // Optional update field
    RequestTime time.Time   `validate:"required"`
}

func (p *AdminUpdateExample) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    return nil
}
```

### Delete Input

```go
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

## Output Structs

Location: `internal/usecase/output/{actor}_{entity}.go`

Only needed for List operations:

```go
package output

import "github.com/abyssparanoia/rapid-go/internal/domain/model"

type AdminListExamples struct {
    Examples   model.Examples
    TotalCount uint64
}
```

## Interactor Interface

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

## Interactor Implementation

Location: `internal/usecase/{actor}_{entity}_impl.go`

### Constructor

```go
package usecase

import (
    "context"

    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/domain/repository"
    "github.com/abyssparanoia/rapid-go/internal/domain/service"
    "github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
    "github.com/abyssparanoia/rapid-go/internal/usecase/input"
    "github.com/abyssparanoia/rapid-go/internal/usecase/output"
    "github.com/volatiletech/null/v8"
)

type adminExampleInteractor struct {
    transactable      repository.Transactable
    exampleRepository repository.Example
    assetService      service.Asset  // Always include for URL processing
}

func NewAdminExampleInteractor(
    transactable repository.Transactable,
    exampleRepository repository.Example,
    assetService service.Asset,
) AdminExampleInteractor {
    return &adminExampleInteractor{
        transactable:      transactable,
        exampleRepository: exampleRepository,
        assetService:      assetService,
    }
}
```

### Create Method

```go
func (i *adminExampleInteractor) Create(
    ctx context.Context,
    param *input.AdminCreateExample,
) (*model.Example, error) {
    // 1. Validate
    if err := param.Validate(); err != nil {
        return nil, err
    }

    // 2. Create domain entity
    example := model.NewExample(
        param.TenantID,
        param.Name,
        param.Description,
        param.RequestTime,
    )

    // 3. Persist in transaction
    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        return i.exampleRepository.Create(ctx, example)
    }); err != nil {
        return nil, err
    }

    // 4. Get with relations
    example, err := i.exampleRepository.Get(ctx, repository.GetExampleQuery{
        ID:             null.StringFrom(example.ID),
        BaseGetOptions: repository.BaseGetOptions{OrFail: true, Preload: true},
    })
    if err != nil {
        return nil, err
    }

    // 5. Apply asset URL processing (always call even if no assets)
    if err := i.assetService.BatchSetExampleURLs(ctx, model.Examples{example}, param.RequestTime); err != nil {
        return nil, err
    }

    return example, nil
}
```

### Get Method

```go
func (i *adminExampleInteractor) Get(
    ctx context.Context,
    param *input.AdminGetExample,
) (*model.Example, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    example, err := i.exampleRepository.Get(ctx, repository.GetExampleQuery{
        ID:             null.StringFrom(param.ExampleID),
        TenantID:       null.StringFrom(param.TenantID),
        BaseGetOptions: repository.BaseGetOptions{OrFail: true, Preload: true},
    })
    if err != nil {
        return nil, err
    }

    // Apply asset URL processing (always call)
    if err := i.assetService.BatchSetExampleURLs(ctx, model.Examples{example}, param.RequestTime); err != nil {
        return nil, err
    }

    return example, nil
}
```

### List Method

```go
func (i *adminExampleInteractor) List(
    ctx context.Context,
    param *input.AdminListExamples,
) (*output.AdminListExamples, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    query := repository.ListExamplesQuery{
        TenantID: null.StringFrom(param.TenantID),
        BaseListOptions: repository.BaseListOptions{
            Page:    null.Uint64From(param.Page),
            Limit:   null.Uint64From(param.Limit),
            Preload: true,
        },
    }

    // Handle optional filters
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

    // Apply asset URL processing (always call)
    if err := i.assetService.BatchSetExampleURLs(ctx, examples, param.RequestTime); err != nil {
        return nil, err
    }

    return &output.AdminListExamples{
        Examples:   examples,
        TotalCount: totalCount,
    }, nil
}
```

### Update Method

```go
func (i *adminExampleInteractor) Update(
    ctx context.Context,
    param *input.AdminUpdateExample,
) (*model.Example, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        // Get with lock
        example, err := i.exampleRepository.Get(ctx, repository.GetExampleQuery{
            ID:             null.StringFrom(param.ExampleID),
            TenantID:       null.StringFrom(param.TenantID),
            BaseGetOptions: repository.BaseGetOptions{OrFail: true, ForUpdate: true},
        })
        if err != nil {
            return err
        }

        // Apply updates via domain method
        example.Update(param.Name, param.Description, param.RequestTime)

        return i.exampleRepository.Update(ctx, example)
    }); err != nil {
        return nil, err
    }

    // Get with relations
    example, err := i.exampleRepository.Get(ctx, repository.GetExampleQuery{
        ID:             null.StringFrom(param.ExampleID),
        BaseGetOptions: repository.BaseGetOptions{OrFail: true, Preload: true},
    })
    if err != nil {
        return nil, err
    }

    // Apply asset URL processing (always call)
    if err := i.assetService.BatchSetExampleURLs(ctx, model.Examples{example}, param.RequestTime); err != nil {
        return nil, err
    }

    return example, nil
}
```

### Delete Method

```go
func (i *adminExampleInteractor) Delete(
    ctx context.Context,
    param *input.AdminDeleteExample,
) error {
    if err := param.Validate(); err != nil {
        return err
    }

    return i.transactable.RWTx(ctx, func(ctx context.Context) error {
        // Verify exists and belongs to tenant
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

## Key Patterns

### Always Use Preload

Set `Preload: true` for all returned entities, even if ReadonlyReference is currently empty:

```go
BaseGetOptions: repository.BaseGetOptions{OrFail: true, Preload: true}
```

### Always Call Asset Service

Call `BatchSet{Entity}URLs` even if entity has no asset fields currently:

```go
if err := i.assetService.BatchSetExampleURLs(ctx, model.Examples{example}, param.RequestTime); err != nil {
    return nil, err
}
```

### Transaction Boundaries

- **RWTx**: For Create, Update, Delete operations
- **ROTx**: For read-only operations needing consistency (rarely needed)
- Always use domain methods for state changes within transactions

### ForUpdate Lock

Use `ForUpdate: true` when fetching entity for Update or Delete:

```go
BaseGetOptions: repository.BaseGetOptions{OrFail: true, ForUpdate: true}
```
