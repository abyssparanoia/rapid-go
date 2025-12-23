---
description: External service integration patterns (Cognito, Firebase, S3, etc.)
globs:
  - "internal/domain/repository/*authentication*.go"
  - "internal/infrastructure/cognito/**/*.go"
  - "internal/infrastructure/firebase/**/*.go"
  - "internal/usecase/**/*_impl.go"
---

# External Service Integration Guidelines

## Overview

External services (IdP, storage, email) are integrated through repository interfaces defined in the domain layer, with implementations in the infrastructure layer.

## Supported Identity Providers

This project supports multiple IdP backends:

- **AWS Cognito**: `internal/infrastructure/cognito/`
- **Firebase Auth**: `internal/infrastructure/firebase/`

The repository interface pattern allows swapping IdP implementations without changing business logic.

## Authentication Repository Pattern (Cognito / Firebase)

### Interface Definition

Location: `internal/domain/repository/admin_authentication.go`

```go
package repository

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository

type AdminAuthentication interface {
    // GetClaims retrieves user claims from IdP
    GetClaims(
        ctx context.Context,
        authUID string,
    ) (*model.AdminClaims, error)

    // StoreClaims updates user claims in IdP
    StoreClaims(
        ctx context.Context,
        authUID string,
        claims *model.AdminClaims,
    ) error

    // DeleteUser removes user from IdP
    DeleteUser(
        ctx context.Context,
        authUID string,
    ) error
}
```

### Claims Model

```go
type AdminClaims struct {
    AuthUID string
    Email   string
    AdminID null.String                    // null until user entity created
    Role    nullable.Type[model.AdminRole] // null until role assigned
}

func NewAdminClaims(
    authUID string,
    email string,
    adminID null.String,
    role nullable.Type[model.AdminRole],
) *AdminClaims {
    return &AdminClaims{
        AuthUID: authUID,
        Email:   email,
        AdminID: adminID,
        Role:    role,
    }
}
```

### Implementation (Cognito)

Location: `internal/infrastructure/cognito/repository/admin_authentication.go`

```go
package repository

type adminAuthentication struct {
    cli        *cognitoidentityprovider.Client
    userPoolID string
}

func NewAdminAuthentication(
    cli *cognitoidentityprovider.Client,
    userPoolID string,
) repository.AdminAuthentication {
    return &adminAuthentication{
        cli:        cli,
        userPoolID: userPoolID,
    }
}

func (r *adminAuthentication) StoreClaims(
    ctx context.Context,
    authUID string,
    claims *model.AdminClaims,
) error {
    attrs := []types.AttributeType{}

    if claims.AdminID.Valid {
        attrs = append(attrs, types.AttributeType{
            Name:  aws.String("custom:admin_id"),
            Value: aws.String(claims.AdminID.String),
        })
    }
    if claims.Role.Valid {
        attrs = append(attrs, types.AttributeType{
            Name:  aws.String("custom:admin_role"),
            Value: aws.String(claims.Role.Value().String()),
        })
    }

    req := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
        UserPoolId:     aws.String(r.userPoolID),
        Username:       aws.String(authUID),
        UserAttributes: attrs,
    }

    _, err := r.cli.AdminUpdateUserAttributes(ctx, req)
    if err != nil {
        return errors.InternalErr.Wrap(err)
    }
    return nil
}

func (r *adminAuthentication) DeleteUser(
    ctx context.Context,
    authUID string,
) error {
    req := &cognitoidentityprovider.AdminDeleteUserInput{
        UserPoolId: aws.String(r.userPoolID),
        Username:   aws.String(authUID),
    }

    _, err := r.cli.AdminDeleteUser(ctx, req)
    if err != nil {
        return errors.InternalErr.Wrap(err)
    }
    return nil
}
```

## Sync Points with IdP

### 1. On User Creation (Accept Invitation)

Store initial claims when user entity is created:

```go
func (i *interactor) Accept(ctx context.Context, param *input.AcceptInvitation) (*model.Admin, error) {
    // ... within transaction

    // Create user entity
    admin = model.NewAdmin(...)
    if err := i.adminRepository.Create(ctx, admin); err != nil {
        return err
    }

    // Sync claims to IdP
    if err := i.adminAuthenticationRepository.StoreClaims(
        ctx,
        param.AuthUID,
        model.NewAdminClaims(
            param.AuthUID,
            admin.Email,
            null.StringFrom(admin.ID),      // Set admin_id
            nullable.TypeFrom(admin.Role),  // Set role
        ),
    ); err != nil {
        return err
    }
    // ...
}
```

### 2. On Role Update

Sync claims when user role changes:

```go
func (i *interactor) Update(ctx context.Context, param *input.UpdateAdmin) (*model.Admin, error) {
    // ... within transaction

    // Update via domain method
    admin.UpdateRole(param.Role.Value(), param.RequestTime)

    if err := i.adminRepository.Update(ctx, admin); err != nil {
        return err
    }

    // Sync updated role to IdP
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
    // ...
}
```

### 3. On User Deletion

Delete user from IdP when user is deleted:

```go
func (i *interactor) Delete(ctx context.Context, param *input.DeleteAdmin) error {
    // ... within transaction

    // Get user to retrieve AuthUID
    admin, err := i.adminRepository.Get(ctx, ...)
    if err != nil {
        return err
    }

    // Delete from IdP first (before DB deletion)
    if err := i.adminAuthenticationRepository.DeleteUser(ctx, admin.AuthUID); err != nil {
        return err
    }

    // Delete from database
    if err := i.adminRepository.Delete(ctx, param.TargetAdminID); err != nil {
        return err
    }
    // ...
}
```

## Interactor Dependency

Include authentication repository when IdP sync is needed:

```go
type adminAdminInteractor struct {
    transactable                  repository.Transactable
    adminRepository               repository.Admin
    adminAuthenticationRepository repository.AdminAuthentication  // Add this
}

func NewAdminAdminInteractor(
    transactable repository.Transactable,
    adminRepository repository.Admin,
    adminAuthenticationRepository repository.AdminAuthentication,
) AdminAdminInteractor {
    return &adminAdminInteractor{
        transactable:                  transactable,
        adminRepository:               adminRepository,
        adminAuthenticationRepository: adminAuthenticationRepository,
    }
}
```

## Transaction Considerations

### IdP Operations Within Transaction

IdP operations should be within the database transaction:
- If IdP operation fails, database changes are rolled back
- Ensures consistency between DB and IdP state

```go
if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
    // 1. Database operations
    if err := i.adminRepository.Update(ctx, admin); err != nil {
        return err
    }

    // 2. IdP operations (within same transaction)
    if err := i.adminAuthenticationRepository.StoreClaims(ctx, ...); err != nil {
        return err  // Will rollback database changes
    }

    return nil
}); err != nil {
    return nil, err
}
```

### Deletion Order

For deletion, delete from IdP before database:

```go
// Within transaction
// 1. Delete from IdP first
if err := i.adminAuthenticationRepository.DeleteUser(ctx, admin.AuthUID); err != nil {
    return err
}

// 2. Delete from database
if err := i.adminRepository.Delete(ctx, id); err != nil {
    return err
}
```

**Rationale**: If DB deletion fails after IdP deletion, the user can't log in anyway. If IdP deletion fails, the user data remains consistent.

## Email Service Integration

### Interface Definition

```go
type Emailer interface {
    SendAdminInvitation(ctx context.Context, invitation *model.AdminInvitation) error
    SendPasswordReset(ctx context.Context, user *model.User, token string) error
}
```

### Usage in Usecase

```go
// Within transaction - email failure will rollback
if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
    if err := i.invitationRepository.Create(ctx, invitation); err != nil {
        return err
    }

    // Send email within transaction
    if err := i.emailer.SendAdminInvitation(ctx, invitation); err != nil {
        return err  // Rollback if email fails
    }

    return nil
}); err != nil {
    return nil, err
}
```

## Best Practices

1. **Define in domain layer** - Repository interfaces in `domain/repository/`
2. **Implement in infrastructure** - Actual implementation in `infrastructure/`
3. **Sync within transactions** - Include IdP operations in DB transactions
4. **Handle partial failures** - Consider rollback scenarios
5. **Delete IdP first** - On user deletion, remove from IdP before DB
6. **Mock for testing** - Use mockgen for unit tests
7. **Wrap errors** - Use domain errors for consistent error handling

## Common Patterns

### Conditional Sync

Only sync when relevant field changes:

```go
// Only sync if role was actually updated
if param.Role.Valid {
    if err := i.adminAuthenticationRepository.StoreClaims(...); err != nil {
        return err
    }
}
```

### nullable.Type for Optional Fields

Use `nullable.Type[T]` for optional update fields:

```go
type AdminUpdateAdmin struct {
    AdminID       string          `validate:"required"`
    AdminRole     model.AdminRole `validate:"required"`
    TargetAdminID string          `validate:"required"`
    Role          nullable.Type[model.AdminRole]  // Optional update
    RequestTime   time.Time `validate:"required"`
}
```

```go
// In handler
param := input.NewAdminUpdateAdmin(...)
if req.Role != nil {
    param.Role = nullable.TypeFrom(marshaller.AdminRoleToModel(*req.Role))
}
```
