---
description: Domain service patterns for complex business logic
globs:
  - "internal/domain/service/**/*.go"
---

# Domain Service Guidelines

## When to Use Domain Services

Use domain services when:

- Business logic spans multiple entities
- Logic doesn't naturally belong to a single entity
- Complex calculations or validations are needed
- External service coordination is required (via interfaces)

## Service Structure

Domain services consist of:
- Interface definition (`{service}.go`) - defines the contract
- Implementation (`{service}_impl.go`) - implements the interface
- Tests (`{service}_impl_test.go`) - unit tests

### Interface Definition (staff.go)

```go
package service

import (
    "context"
    "time"

    "github.com/abyssparanoia/rapid-go/internal/domain/model"
)

//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_service
type Staff interface {
    Create(
        ctx context.Context,
        param StaffCreateParam,
    ) (*model.Staff, error)
}

type StaffCreateParam struct {
    TenantID    string
    Email       string
    Password    string
    StaffRole   model.StaffRole
    DisplayName string
    ImagePath   string
    RequestTime time.Time
}
```

### Implementation (staff_impl.go)

```go
package service

import (
    "context"

    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/domain/repository"
)

type staffService struct {
    staffRepository               repository.Staff
    staffAuthenticationRepository repository.StaffAuthentication
}

func NewStaff(
    staffRepository repository.Staff,
    staffAuthenticationRepository repository.StaffAuthentication,
) Staff {
    return &staffService{
        staffRepository,
        staffAuthenticationRepository,
    }
}
```

## Param/Result Pattern

### Input Parameters

```go
// Use Param suffix for input structs
type DoSomethingParam struct {
    ExampleID string
    Value     string
    Options   DoSomethingOptions
}

type DoSomethingOptions struct {
    SkipValidation bool
    NotifyUsers    bool
}
```

### Output Results

```go
// Use Result suffix for output structs
type DoSomethingResult struct {
    Example *model.Example
    Count   int
    Summary string
}
```

## Method Signatures

```go
func (s *ExampleService) DoSomething(
    ctx context.Context,
    param *DoSomethingParam,
) (*DoSomethingResult, error) {
    // 1. Validate input
    if param.ExampleID == "" {
        return nil, errors.RequestInvalidArgumentErr.Errorf("example_id is required")
    }

    // 2. Fetch required data
    example, err := s.exampleRepo.Get(ctx, repository.GetExampleQuery{
        ID:     null.StringFrom(param.ExampleID),
        OrFail: true,
    })
    if err != nil {
        return nil, err
    }

    // 3. Execute business logic
    result := s.processExample(example, param)

    // 4. Return result
    return &DoSomethingResult{
        Example: example,
        Count:   result.count,
    }, nil
}
```

## Domain Service vs Interactor

| Aspect | Domain Service | Interactor |
|--------|---------------|------------|
| Location | `domain/service/` | `usecase/` |
| Transaction | Assumes context has TX | Manages TX boundaries |
| Focus | Pure business logic | Application workflow |
| Dependencies | Only domain layer | Domain + infrastructure interfaces |

## Best Practices

1. **Keep services focused** - Each service handles one domain concept
2. **No transaction management** - Services assume transaction context exists
3. **Use Param/Result** - Clear input/output contracts
4. **Return domain errors** - Use errors from `domain/errors`
5. **No infrastructure dependencies** - Only repository interfaces

## Example: Publishing Workflow

```go
type PublishExampleParam struct {
    ExampleID string
    PublisherID string
}

type PublishExampleResult struct {
    Example       *model.Example
    Notifications []model.Notification
}

func (s *ExampleService) Publish(
    ctx context.Context,
    param *PublishExampleParam,
) (*PublishExampleResult, error) {
    // Get example with lock
    example, err := s.exampleRepo.Get(ctx, repository.GetExampleQuery{
        ID:        null.StringFrom(param.ExampleID),
        OrFail:    true,
        ForUpdate: true,
    })
    if err != nil {
        return nil, err
    }

    // Validate status transition
    if example.Status != model.ExampleStatusDraft {
        return nil, errors.ExampleStatusInvalidErr.Errorf(
            "cannot publish: current status=%s", example.Status,
        )
    }

    // Update status
    example.Status = model.ExampleStatusPublished
    example.UpdatedAt = time.Now()

    if err := s.exampleRepo.Update(ctx, example); err != nil {
        return nil, err
    }

    return &PublishExampleResult{
        Example: example,
    }, nil
}
```
