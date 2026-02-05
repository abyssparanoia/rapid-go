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

## Asset Service Pattern

The Asset service is responsible for generating presigned URLs for assets (images, files, etc.) stored in external storage (GCS, S3).

### Interface Method Naming Convention

AssetService methods follow strict naming and signature conventions:

```go
type Asset interface {
    CreatePresignedURL(
        ctx context.Context,
        assetType model.AssetType,
        contentType model.ContentType,
        requestTime time.Time,
    ) (*AssetCreatePresignedURLResult, error)

    GetWithValidate(
        ctx context.Context,
        assetType model.AssetType,
        assetID string,
    ) (string, error)

    // BatchSet methods - Always use plural slice type
    BatchSetTenantURLs(ctx context.Context, tenants model.Tenants) error
    BatchSetStaffURLs(ctx context.Context, staffs model.Staffs) error
    BatchSet{Entity}URLs(ctx context.Context, {entities} model.{Entity}s) error
}
```

### Method Signature Rules

**CRITICAL RULES - Must be followed for all new resource types:**

1. **Method name**: `BatchSet{Entity}URLs` where `{Entity}` is singular form
   - ✅ `BatchSetStaffURLs` (Staff → Staffs)
   - ✅ `BatchSetTenantURLs` (Tenant → Tenants)
   - ❌ NOT `BatchSetStaffURL` (missing 's' at end)

2. **Parameter type**: ALWAYS use plural slice type `model.{Entity}s`
   - ✅ `tenants model.Tenants` (type alias for `[]*Tenant`)
   - ✅ `staffs model.Staffs` (type alias for `[]*Staff`)
   - ❌ NOT `tenants []*model.Tenant` (use type alias)

3. **Return type**: `error` only (modifies entities in-place)

### Implementation Pattern

```go
func (s *assetService) BatchSetStaffURLs(
    ctx context.Context,
    staffs model.Staffs,
) error {
    // Generate URLs for each staff's ImagePath
    for _, staff := range staffs {
        if staff.ImagePath != "" {
            url, err := s.assetRepository.GetPresignedURL(ctx, staff.ImagePath)
            if err != nil {
                return err
            }
            staff.ImageURL = null.StringFrom(url)
        }

        // Recursively set URLs for ReadonlyReference relations
        if staff.ReadonlyReference != nil && staff.ReadonlyReference.Tenant != nil {
            if err := s.BatchSetTenantURLs(ctx, model.Tenants{staff.ReadonlyReference.Tenant}); err != nil {
                return err
            }
        }
    }
    return nil
}
```

### Usage in Usecase Layer

**ALWAYS call BatchSet methods when returning entities, even if no asset fields currently exist:**

```go
// Single entity - wrap in slice
staff, err := i.staffRepository.Get(ctx, query)
if err != nil {
    return nil, err
}
// MUST call even if Staff has no image field yet (defensive programming)
if err := i.assetService.BatchSetStaffURLs(ctx, model.Staffs{staff}); err != nil {
    return nil, err
}
return staff, nil

// Multiple entities - pass slice directly
staffs, err := i.staffRepository.List(ctx, query)
if err != nil {
    return nil, err
}
// MUST call even if Staff has no image field yet (defensive programming)
if err := i.assetService.BatchSetStaffURLs(ctx, staffs); err != nil {
    return nil, err
}
return output.NewAdminListStaffs(staffs, pagination), nil
```

### When to Add New BatchSet Methods

Add a new `BatchSet{Entity}URLs` method when:
1. **New resource with asset field** - Entity has `ImagePath`, `FilePath`, etc.
2. **Future-proofing** - Even if no asset field exists yet, add the method for defensive programming
3. **ReadonlyReference contains assets** - Entity has relations that may contain assets

### Adding New Resource Support

When adding support for a new resource type (e.g., `Product`):

**Step 1: Add interface method**
```go
type Asset interface {
    // ... existing methods
    BatchSetProductURLs(ctx context.Context, products model.Products) error
}
```

**Step 2: Implement method**
```go
func (s *assetService) BatchSetProductURLs(
    ctx context.Context,
    products model.Products,
) error {
    for _, product := range products {
        // Set URLs for product's asset fields
        if product.ImagePath != "" {
            url, err := s.assetRepository.GetPresignedURL(ctx, product.ImagePath)
            if err != nil {
                return err
            }
            product.ImageURL = null.StringFrom(url)
        }

        // Set URLs for ReadonlyReference relations
        if product.ReadonlyReference != nil {
            if product.ReadonlyReference.Category != nil {
                if err := s.BatchSetCategoryURLs(ctx, model.Categories{product.ReadonlyReference.Category}); err != nil {
                    return err
                }
            }
        }
    }
    return nil
}
```

**Step 3: Call in usecase**
```go
// Always call when returning Product entities
if err := i.assetService.BatchSetProductURLs(ctx, model.Products{product}); err != nil {
    return nil, err
}
```

### Why This Pattern

1. **Defensive Programming**: Adding BatchSet methods even when no assets exist prevents missed updates when asset fields are added later
2. **Recursive URL Setting**: BatchSet methods automatically handle ReadonlyReference relations
3. **Consistent Interface**: Uniform method signatures make the codebase predictable
4. **Type Safety**: Using type aliases (`model.Staffs`) instead of raw slices provides better type checking

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
