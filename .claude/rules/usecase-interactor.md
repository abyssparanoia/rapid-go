---
description: Usecase layer interactor implementation patterns
globs:
  - "internal/usecase/**/*.go"
---

# Usecase Interactor Guidelines

## Naming Convention

- Interface: `{Actor}{Resource}Interactor` (e.g., `AdminUserInteractor`)
- Implementation: `{actor}{Resource}Interactor` (lowercase first letter)
- File: `{actor}_{resource}.go` for interface, `{actor}_{resource}_impl.go` for implementation

## Interface Definition

Location: `internal/usecase/{actor}_{resource}.go`

```go
package usecase

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type AdminExampleInteractor interface {
    Create(ctx context.Context, param *input.AdminCreateExample) (*model.Example, error)
    Get(ctx context.Context, param *input.AdminGetExample) (*model.Example, error)
    List(ctx context.Context, param *input.AdminListExamples) (*output.AdminListExamples, error)
    Update(ctx context.Context, param *input.AdminUpdateExample) (*model.Example, error)
    Delete(ctx context.Context, param *input.AdminDeleteExample) error
}
```

## Implementation Structure

Location: `internal/usecase/{actor}_{resource}_impl.go`

```go
package usecase

type adminExampleInteractor struct {
    transactable      repository.Transactable
    exampleRepository repository.Example
    // Add other dependencies
}

func NewAdminExampleInteractor(
    transactable repository.Transactable,
    exampleRepository repository.Example,
) AdminExampleInteractor {
    return &adminExampleInteractor{
        transactable:      transactable,
        exampleRepository: exampleRepository,
    }
}
```

## Input Structs

Location: `internal/usecase/input/{actor}_{resource}.go`

```go
package input

type AdminCreateExample struct {
    AdminID     string    `validate:"required"`
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

### Input Naming: `{Actor}{Action}{Resource}`

- `AdminCreateExample`
- `UserGetProfile`
- `AdminListUsers`

### Common Fields

- `AdminID` / `UserID` - Actor identifier from auth claims
- `TenantID` - Tenant context
- `RequestTime` - Current time from request context

### Validation Tags

- `validate:"required"` - Required field
- `validate:"required,max=256"` - Required with max length
- `validate:"required,min=1,max=100"` - For pagination

## Output Structs

Location: `internal/usecase/output/{actor}_{resource}.go`

```go
package output

type AdminListExamples struct {
    Examples   model.Examples
    TotalCount uint64
}
```

Only create output structs when:
- Returning multiple items (list with pagination)
- Returning computed values beyond the entity

For single entity returns, use `*model.Example` directly.

## Method Patterns

### Create

```go
func (i *adminExampleInteractor) Create(
    ctx context.Context,
    param *input.AdminCreateExample,
) (*model.Example, error) {
    // 1. Validate input
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

    // 4. Return with relations loaded
    return i.exampleRepository.Get(ctx, repository.GetExampleQuery{
        ID: null.StringFrom(example.ID),
        BaseGetOptions: repository.BaseGetOptions{
            OrFail:  true,
            Preload: true,
        },
    })
}
```

### Get

```go
func (i *adminExampleInteractor) Get(
    ctx context.Context,
    param *input.AdminGetExample,
) (*model.Example, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    return i.exampleRepository.Get(ctx, repository.GetExampleQuery{
        ID:       null.StringFrom(param.ExampleID),
        TenantID: null.StringFrom(param.TenantID),  // Scope to tenant
        BaseGetOptions: repository.BaseGetOptions{
            OrFail:  true,
            Preload: true,
        },
    })
}
```

### List with Pagination

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

    // Optional filters
    if param.Status != nil {
        query.Status = nullable.From(*param.Status)
    }

    examples, err := i.exampleRepository.List(ctx, query)
    if err != nil {
        return nil, err
    }

    totalCount, err := i.exampleRepository.Count(ctx, query)
    if err != nil {
        return nil, err
    }

    return &output.AdminListExamples{
        Examples:   examples,
        TotalCount: totalCount,
    }, nil
}
```

### Update

```go
func (i *adminExampleInteractor) Update(
    ctx context.Context,
    param *input.AdminUpdateExample,
) (*model.Example, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        // 1. Get with lock
        example, err := i.exampleRepository.Get(ctx, repository.GetExampleQuery{
            ID:       null.StringFrom(param.ExampleID),
            TenantID: null.StringFrom(param.TenantID),
            BaseGetOptions: repository.BaseGetOptions{
                OrFail:    true,
                ForUpdate: true,  // Lock for update
            },
        })
        if err != nil {
            return err
        }

        // 2. Apply updates via domain method
        example.Update(param.Name, param.Description, param.RequestTime)

        // 3. Persist
        return i.exampleRepository.Update(ctx, example)
    }); err != nil {
        return nil, err
    }

    // 4. Return updated entity with relations
    return i.exampleRepository.Get(ctx, repository.GetExampleQuery{
        ID: null.StringFrom(param.ExampleID),
        BaseGetOptions: repository.BaseGetOptions{
            OrFail:  true,
            Preload: true,
        },
    })
}
```

### Delete

```go
func (i *adminExampleInteractor) Delete(
    ctx context.Context,
    param *input.AdminDeleteExample,
) error {
    if err := param.Validate(); err != nil {
        return err
    }

    return i.transactable.RWTx(ctx, func(ctx context.Context) error {
        // Verify entity exists and belongs to tenant
        _, err := i.exampleRepository.Get(ctx, repository.GetExampleQuery{
            ID:       null.StringFrom(param.ExampleID),
            TenantID: null.StringFrom(param.TenantID),
            BaseGetOptions: repository.BaseGetOptions{
                OrFail:    true,
                ForUpdate: true,
            },
        })
        if err != nil {
            return err
        }

        return i.exampleRepository.Delete(ctx, param.ExampleID)
    })
}
```

## Transaction Rules

- Use `RWTx` for write operations (Create, Update, Delete)
- Use `ROTx` for read-only operations that need consistency
- Transaction boundary is always in the usecase layer
- Domain services assume transaction is already active

## Return Pattern Best Practices

### Always Enable Preload for Returned Entities

When returning domain entities from interactor methods, always set `Preload: true` in the repository query, even if there are currently no related entities:

```go
// Good - Always enable preload for returned entities
return i.exampleRepository.Get(ctx, repository.GetExampleQuery{
    ID: null.StringFrom(example.ID),
    BaseGetOptions: repository.BaseGetOptions{
        OrFail:  true,
        Preload: true,  // Always true for returned entities
    },
})
```

**Rationale**: When relations are added later, existing code will automatically include them. This prevents missing relation loading when the domain model evolves.

### Apply Asset Service Processing

When entities have asset URLs (images, files), apply the asset service batch processing before returning:

```go
// After fetching entities, apply asset URL processing
if err := i.assetService.BatchSetExampleURLs(ctx, examples, param.RequestTime); err != nil {
    return nil, err
}
```

**Rationale**: Same principle - ensures future asset fields are automatically processed.

## External Service Integration

When operations require synchronization with external services (IdP, email), include them within the transaction:

### Update with IdP Sync

```go
func (i *adminAdminInteractor) Update(
    ctx context.Context,
    param *input.AdminUpdateAdmin,
) (*model.Admin, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        admin, err := i.adminRepository.Get(ctx, repository.GetAdminQuery{
            ID: null.StringFrom(param.TargetAdminID),
            BaseGetOptions: repository.BaseGetOptions{
                OrFail:    true,
                ForUpdate: true,
            },
        })
        if err != nil {
            return err
        }

        // Use domain method for state change
        if param.Role.Valid {
            admin.UpdateRole(param.Role.Value(), param.RequestTime)
        }

        if err := i.adminRepository.Update(ctx, admin); err != nil {
            return err
        }

        // Sync to IdP within transaction
        if param.Role.Valid {
            if err := i.adminAuthenticationRepository.StoreClaims(
                ctx,
                admin.AuthUID,
                model.NewAdminClaims(
                    admin.AuthUID,
                    admin.Email,
                    null.StringFrom(admin.ID),
                    param.Role,
                ),
            ); err != nil {
                return err
            }
        }

        return nil
    }); err != nil {
        return nil, err
    }

    return i.adminRepository.Get(ctx, ...)
}
```

### Delete with IdP Cleanup

```go
func (i *adminAdminInteractor) Delete(
    ctx context.Context,
    param *input.AdminDeleteAdmin,
) error {
    if err := param.Validate(); err != nil {
        return err
    }

    // Authorization check
    if !param.AdminRole.IsRoot() {
        return errors.AdminForbiddenErr.Errorf("only root admin can delete")
    }

    return i.transactable.RWTx(ctx, func(ctx context.Context) error {
        admin, err := i.adminRepository.Get(ctx, repository.GetAdminQuery{
            ID: null.StringFrom(param.TargetAdminID),
            BaseGetOptions: repository.BaseGetOptions{
                OrFail:    true,
                ForUpdate: true,
            },
        })
        if err != nil {
            return err
        }

        // Delete from IdP first
        if err := i.adminAuthenticationRepository.DeleteUser(ctx, admin.AuthUID); err != nil {
            return err
        }

        // Then delete from database
        return i.adminRepository.Delete(ctx, param.TargetAdminID)
    })
}
```

## Optional Update Fields with nullable.Type

For optional update fields, use `nullable.Type[T]` instead of pointers:

### Input Struct

```go
type AdminUpdateAdmin struct {
    AdminID       string          `validate:"required"`
    AdminRole     model.AdminRole `validate:"required"`
    TargetAdminID string          `validate:"required"`
    Role          nullable.Type[model.AdminRole]  // Optional field
    RequestTime   time.Time       `validate:"required"`
}

func (p *AdminUpdateAdmin) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    // Validate optional field if present
    if p.Role.Valid && !p.Role.Value().Valid() {
        return errors.RequestInvalidArgumentErr.Errorf("invalid role: %s", p.Role.Value())
    }
    return nil
}
```

### Handler Usage

```go
func (h *Handler) UpdateAdmin(ctx context.Context, req *pb.UpdateAdminRequest) (*pb.UpdateAdminResponse, error) {
    claims, err := session_interceptor.RequireAdminSessionContext(ctx)
    if err != nil {
        return nil, err
    }

    param := input.NewAdminUpdateAdmin(
        claims.AdminID.String,
        claims.Role.Value(),
        req.AdminId,
        nullable.Type[model.AdminRole]{},  // Empty by default
        request_interceptor.GetRequestTime(ctx),
    )

    // Set optional field if provided in request
    if req.Role != nil {
        param.Role = nullable.TypeFrom(marshaller.AdminRoleToModel(*req.Role))
    }

    admin, err := h.adminInteractor.Update(ctx, param)
    // ...
}
```

### Usecase Usage

```go
// Check if optional field was provided
if param.Role.Valid {
    admin.UpdateRole(param.Role.Value(), param.RequestTime)
}
```

## Domain Method Usage (Domain Logic First)

**IMPORTANT**: Always use domain model methods for state changes instead of direct field assignment.

### Good Pattern

```go
// In usecase
if param.Role.Valid {
    admin.UpdateRole(param.Role.Value(), param.RequestTime)
}
```

### Anti-Pattern

```go
// Don't do this in usecase
if param.Role.Valid {
    admin.Role = param.Role.Value()
    admin.UpdatedAt = param.RequestTime
}
```

See `domain-model.md` for more details on domain method patterns.
