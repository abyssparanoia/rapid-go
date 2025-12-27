---
description: Domain error definition and handling patterns
globs:
  - "internal/domain/errors/**/*.go"
---

# Domain Errors Guidelines

## Error Definition Location

- Error factory functions: `internal/domain/errors/base.go`
- Error definitions: `internal/domain/errors/errors.go`
- Error categories: `internal/domain/errors/error_category.go`

## Error Structure

This project uses `github.com/abyssparanoia/goerr` package for error handling.

### Error Factory Functions (base.go)

```go
package errors

import "github.com/abyssparanoia/goerr"

func NewBadRequestError(errCode string, msg string) *goerr.Error {
    return goerr.New("%s", msg).
        WithCategory(ErrorCategoryBadRequest.String()).
        WithCode(errCode)
}

func NewUnauthorizedError(errCode string, msg string) *goerr.Error {
    return goerr.New("%s", msg).
        WithCategory(ErrorCategoryUnauthorized.String()).
        WithCode(errCode)
}

func NewNotFoundError(errCode string, msg string) *goerr.Error {
    return goerr.New("%s", msg).
        WithCategory(ErrorCategoryNotFound.String()).
        WithCode(errCode)
}

func NewInternalError(errCode string, msg string) *goerr.Error {
    return goerr.New("%s", msg).
        WithCategory(ErrorCategoryInternal.String()).
        WithCode(errCode)
}

// Other factory functions: NewForbiddenError, NewConflictError, NewCanceledError, NewServiceAvailableError
```

## Defining Domain Errors (errors.go)

```go
package errors

var (
    // General errors
    InternalErr               = NewInternalError("E100001", "An internal error has occurred")
    RequestInvalidArgumentErr = NewBadRequestError("E100002", "Request argument is invalid")
    InvalidIDTokenErr         = NewUnauthorizedError("E100003", "Invalid ID token")
    AssetInvalidErr           = NewBadRequestError("E100006", "Asset is invalid")
    AssetNotFoundErr          = NewNotFoundError("E100007", "Asset not found")

    // Tenant errors
    TenantNotFoundErr = NewNotFoundError("E200101", "Tenant not found")

    // Staff errors
    StaffNotFoundErr = NewNotFoundError("E200201", "Staff not found")
)
```

## Naming Conventions

| Error Type | Pattern | Example |
|------------|---------|---------|
| Not Found | `{Entity}NotFoundErr` | `StaffNotFoundErr`, `TenantNotFoundErr` |
| Already Exists | `{Entity}AlreadyExistsErr` | `StaffAlreadyExistsErr` |
| Invalid State | `{Entity}InvalidErr` | `AssetInvalidErr` |
| Permission | `{Action}ForbiddenErr` | `DeleteForbiddenErr` |

## Error Codes

Use format `E{category}{sequence}`:

- `E1xxxxx` - General/common errors
- `E2001xx` - Tenant errors
- `E2002xx` - Staff errors

```go
errCode: "E100001"  // Internal error
errCode: "E200101"  // Tenant not found
errCode: "E200201"  // Staff not found
```

## Using Errors in Repository

```go
func (r *staff) Get(ctx context.Context, query repository.GetStaffQuery) (*model.Staff, error) {
    // ...
    if err == sql.ErrNoRows {
        if query.OrFail {
            return nil, errors.StaffNotFoundErr.New().
                WithDetail("staff not found").
                WithValue("id", query.ID.String)
        }
        return nil, nil
    }
    return nil, errors.InternalErr.Wrap(err)
}
```

## Using Errors in Domain Model

```go
func ValidateAssetPath(
    assetType AssetType,
    path string,
) error {
    if !assetType.Valid() {
        return errors.AssetInvalidErr.New().
            WithDetail("asset_type is invalid").
            WithValue("asset_type", assetType.String())
    }
    if !strings.HasPrefix(path, assetType.String()) {
        return errors.AssetInvalidErr.New().
            WithDetail("path is invalid").
            WithValue("path", path)
    }
    return nil
}
```

## Using Errors in Usecase

```go
func (i *adminStaffInteractor) Get(ctx context.Context, param *input.AdminGetStaff) (*model.Staff, error) {
    staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
        ID:     null.StringFrom(param.StaffID),
        OrFail: true,
    })
    if err != nil {
        return nil, err  // Already wrapped domain error
    }

    // ...
}
```

## Error Mapping to gRPC Status

The error interceptor maps domain errors to gRPC status codes:

| Domain Error | gRPC Code | HTTP Status |
|--------------|-----------|-------------|
| `InternalErr` | `Internal` | 500 |
| `RequestInvalidArgumentErr` | `InvalidArgument` | 400 |
| `UnauthorizedErr` | `Unauthenticated` | 401 |
| `ForbiddenErr` | `PermissionDenied` | 403 |
| `*NotFoundErr` | `NotFound` | 404 |
| `*AlreadyExistsErr` | `AlreadyExists` | 409 |

## Error Checking

```go
import "github.com/abyssparanoia/goerr"

// Check specific error using goerr.Is
if goerr.Is(err, domainerrors.StaffNotFoundErr) {
    // Handle not found case
}

// Get error details
var goErr *goerr.Error
if errors.As(err, &goErr) {
    code := goErr.Code()
    category := goErr.Category()
}
```

## Adding Context to Errors

Use `WithDetail` and `WithValue` to add context:

```go
// Good - adds context with structured values
return errors.StaffNotFoundErr.New().
    WithDetail("staff not found").
    WithValue("staff_id", staffID).
    WithValue("tenant_id", tenantID)

// Bad - no context
return errors.StaffNotFoundErr.New()
```

## Best Practices

1. **One error per entity for NotFound** - `StaffNotFoundErr`, `TenantNotFoundErr`, not generic `NotFoundErr`
2. **Use `WithDetail` and `WithValue`** - Structured error context for debugging
3. **Use `.New()` to create new error instance** - Preserves stack trace
4. **Use `.Wrap(err)` for wrapping underlying errors** - Chain error causes
5. **Define business logic errors explicitly** - Don't use generic `RequestInvalidArgumentErr` for everything
6. **Keep error codes stable** - Clients may depend on them (e.g., `E200101`)
