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
- Test file: `{actor}_{resource}_impl_test.go` for unit tests

## Method Ordering

**All interface methods must be defined in the following order:**

1. **Get methods** - Single resource retrieval
2. **List methods** - Collection retrieval with pagination
3. **Create methods** - Resource creation
4. **Custom operations (no ID)** - Special operations without resource ID
5. **Update methods** - Resource modification
6. **Custom operations (with ID)** - Special operations with resource ID
7. **Delete methods** - Resource deletion

**Example ordering:**

```go
type AdminStaffInteractor interface {
    // Get
    Get(ctx context.Context, param *input.AdminGetStaff) (*model.Staff, error)

    // List
    List(ctx context.Context, param *input.AdminListStaffs) (*output.ListStaffs, error)

    // Create
    Create(ctx context.Context, param *input.AdminCreateStaff) (*model.Staff, error)

    // Custom (no ID)
    SendNotifications(ctx context.Context, param *input.AdminSendStaffNotifications) error

    // Update
    Update(ctx context.Context, param *input.AdminUpdateStaff) (*model.Staff, error)

    // Custom (with ID)
    SendNotification(ctx context.Context, param *input.AdminSendStaffNotification) error

    // Delete
    Delete(ctx context.Context, param *input.AdminDeleteStaff) error
}
```

**Implementation file methods must follow the same order.**

## Unit Testing Requirement

**ALL usecase interactor implementations MUST have corresponding unit tests.**

### Test File Structure

- **Location**: Same directory as implementation (`internal/usecase/`)
- **Naming**: `{actor}_{resource}_impl_test.go`
- **Pattern**: Table-driven tests using `map[string]testcaseFunc`

### Required Test Coverage

For each method in the interactor, implement tests covering:
- **invalid argument** - Validation error cases
- **not found** - Entity doesn't exist (for Get/Update/Delete)
- **success** - Happy path scenario

### Test Pattern

```go
func TestAdminStaffInteractor_Get(t *testing.T) {
    t.Parallel()

    type args struct {
        staffID string
    }

    type want struct {
        staff          *model.Staff
        expectedResult error
    }

    type testcase struct {
        args    args
        usecase AdminStaffInteractor
        want    want
    }

    type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

    tests := map[string]testcaseFunc{
        "invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
            // Setup test case with empty args
        },
        "not found": func(ctx context.Context, ctrl *gomock.Controller) testcase {
            // Setup test case with mocked repository returning NotFoundErr
        },
        "success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
            // Setup test case with mocked repository returning entity
        },
    }

    for name, tc := range tests {
        tc := tc
        t.Run(name, func(t *testing.T) {
            t.Parallel()
            ctx := t.Context()
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            tc := tc(ctx, ctrl)

            got, err := tc.usecase.Get(ctx, input.NewAdminGetStaff(tc.args.staffID))
            if tc.want.expectedResult == nil {
                require.NoError(t, err)
                require.Equal(t, tc.want.staff, got)
            } else {
                require.ErrorContains(t, err, tc.want.expectedResult.Error())
            }
        })
    }
}
```

### Test Utilities

- **Factory**: Use `factory.NewFactory()` to generate test data
- **Mocks**: Use `mock_repository`, `mock_service` packages
- **Test Transaction**: Use `mock_repository.TestMockTransactable()` for RWTx/ROTx
- **Parallel Execution**: Always use `t.Parallel()` for independent tests

### Verification

Run tests with:
```bash
make test
```

Ensure 100% coverage of usecase methods with meaningful test cases.

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
    assetService      service.Asset  // Always include this dependency
    // Add other dependencies
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

    // 4. Return with relations loaded (Preload: true is required)
    example, err := i.exampleRepository.Get(ctx, repository.GetExampleQuery{
        ID: null.StringFrom(example.ID),
        BaseGetOptions: repository.BaseGetOptions{
            OrFail:  true,
            Preload: true,
        },
    })
    if err != nil {
        return nil, err
    }

    // 5. Apply asset URL processing (call even if no assets exist)
    if err := i.assetService.BatchSetExampleURLs(ctx, model.Examples{example}, param.RequestTime); err != nil {
        return nil, err
    }

    return example, nil
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

    example, err := i.exampleRepository.Get(ctx, repository.GetExampleQuery{
        ID:       null.StringFrom(param.ExampleID),
        TenantID: null.StringFrom(param.TenantID),  // Scope to tenant
        BaseGetOptions: repository.BaseGetOptions{
            OrFail:  true,
            Preload: true,
        },
    })
    if err != nil {
        return nil, err
    }

    // Apply asset URL processing (call even if no assets exist)
    if err := i.assetService.BatchSetExampleURLs(ctx, model.Examples{example}, param.RequestTime); err != nil {
        return nil, err
    }

    return example, nil
}
```

### List with Pagination & SortKey (Unified Specification)

**IMPORTANT**: All List operations MUST include SortKey support. This is a unified specification across the codebase.

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
        SortKey: nullable.TypeFrom(param.SortKey),  // REQUIRED - Always include
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

**Key Points for List Operations:**

1. **SortKey is Mandatory**: Every List operation must accept and pass SortKey to repository
   - Input struct field: `SortKey model.XXXSortKey` (NON-nullable)
   - Repository query field: `SortKey nullable.TypeFrom(param.SortKey)`
   - Default value (CreatedAtDesc) is applied in input constructor

2. **Pagination Defaults**: Applied in input layer constructor, not validation
   - `page == 0` → `page = 1`
   - `limit == 0` → `limit = 30`

3. **Preload Required**: Always set `Preload: true` for returned entities
   - Ensures ReadonlyReference is populated for response marshalling

4. **Count Query**: Use same query struct (including filters and SortKey) for consistency

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

## Return Pattern Best Practices (Defensive Programming)

**IMPORTANT**: When returning entities, apply the following patterns **even if there are currently no targets**. This is defensive programming to prevent omissions when relations or asset fields are added later.

### Always Enable Preload for Returned Entities

**Always** set `Preload: true`. Always set it even if ReadonlyReference is currently empty.

```go
// Good - Set Preload: true even if no relations exist currently
return i.exampleRepository.Get(ctx, repository.GetExampleQuery{
    ID: null.StringFrom(example.ID),
    BaseGetOptions: repository.BaseGetOptions{
        OrFail:  true,
        Preload: true,  // Always true - prepare for future relation additions
    },
})

// Same applies to List
examples, err := i.exampleRepository.List(ctx, repository.ListExamplesQuery{
    TenantID: null.StringFrom(param.TenantID),
    BaseListOptions: repository.BaseListOptions{
        Page:    null.Uint64From(param.Page),
        Limit:   null.Uint64From(param.Limit),
        Preload: true,  // Always true
    },
})
```

**Rationale**: When relations are added later, existing code will automatically include them. This prevents missed updates as the domain model evolves.

### Always Apply Asset Service Processing

**Always** call `BatchSet{Entity}URLs`. Always call it even if the entity currently has no asset fields (such as profile images).

```go
// Good - Call BatchSet even if no asset fields exist
func (i *adminExampleInteractor) Get(
    ctx context.Context,
    param *input.AdminGetExample,
) (*model.Example, error) {
    // ...
    example, err := i.exampleRepository.Get(ctx, repository.GetExampleQuery{
        ID: null.StringFrom(param.ExampleID),
        BaseGetOptions: repository.BaseGetOptions{
            OrFail:  true,
            Preload: true,
        },
    })
    if err != nil {
        return nil, err
    }

    // Wrap single entity in model.Examples when calling
    if err := i.assetService.BatchSetExampleURLs(ctx, model.Examples{example}, param.RequestTime); err != nil {
        return nil, err
    }

    return example, nil
}

// For List operations
func (i *adminExampleInteractor) List(
    ctx context.Context,
    param *input.AdminListExamples,
) (*output.AdminListExamples, error) {
    // ...
    examples, err := i.exampleRepository.List(ctx, query)
    if err != nil {
        return nil, err
    }

    // Pass the slice directly
    if err := i.assetService.BatchSetExampleURLs(ctx, examples, param.RequestTime); err != nil {
        return nil, err
    }

    return &output.AdminListExamples{Examples: examples, TotalCount: totalCount}, nil
}
```

**Rationale**: When asset fields (such as image URLs) are added later, existing code will automatically set the URLs. Additionally, since `BatchSet` recursively sets URLs for related entities in ReadonlyReference, it also handles asset additions to related entities.

### Anti-Pattern: Conditional Processing

```go
// Bad - Only call when assets exist
if len(example.ProfileImagePath) > 0 {
    if err := i.assetService.BatchSetExampleURLs(...); err != nil { ... }
}

// Bad - Only set Preload when relations are needed
if needsRelations {
    query.BaseGetOptions.Preload = true
}
```

Avoid these patterns as they cause missed updates during future expansions.

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
