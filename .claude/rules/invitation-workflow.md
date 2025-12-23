---
description: Invitation and approval workflow patterns for user onboarding
globs:
  - "internal/domain/model/*invitation*.go"
  - "internal/usecase/*invitation*.go"
---

# Invitation Workflow Guidelines

## Overview

Invitation workflows are used for user onboarding processes where:
- An authorized user creates an invitation
- The invitee receives notification (email)
- The invitee accepts/rejects the invitation
- Upon acceptance, a new user entity is created

## Domain Model Structure

### Invitation Entity

```go
type AdminInvitation struct {
    ID             string
    InvitationCode string                // Unique code for URL
    Status         AdminInvitationStatus
    Role           AdminRole             // Role to be assigned
    Email          string                // Invitee email
    DisplayName    string
    SentAt         time.Time
    ExpiresAt      time.Time             // Expiration timestamp
    AcceptedAt     null.Time
    RejectedAt     null.Time
    InvalidatedAt  null.Time
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

### Key Fields

| Field | Purpose |
|-------|---------|
| `InvitationCode` | URL-safe unique code for invitation link |
| `ExpiresAt` | Time-limited validity (e.g., 7 days) |
| `Status` | Current state of invitation |
| Timestamp fields | Audit trail for each state transition |

### Status Enum

```go
type AdminInvitationStatus string

const (
    AdminInvitationStatusUnknown     AdminInvitationStatus = "unknown"
    AdminInvitationStatusPending     AdminInvitationStatus = "pending"
    AdminInvitationStatusAccepted    AdminInvitationStatus = "accepted"
    AdminInvitationStatusRejected    AdminInvitationStatus = "rejected"
    AdminInvitationStatusInvalidated AdminInvitationStatus = "invalidated"
)
```

### State Diagram

```
                    ┌─────────────┐
                    │   Pending   │
                    └──────┬──────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
         ▼                 ▼                 ▼
┌─────────────┐   ┌─────────────┐   ┌─────────────────┐
│  Accepted   │   │  Rejected   │   │  Invalidated    │
└─────────────┘   └─────────────┘   └─────────────────┘
```

## Constructor Pattern

```go
const (
    AdminInvitationExpirationDays = 7
)

func NewAdminInvitation(
    role AdminRole,
    emailAddr string,
    displayName string,
    t time.Time,
) *AdminInvitation {
    return &AdminInvitation{
        ID:             id.New(),
        InvitationCode: uuid.UUIDBase64(),  // URL-safe code
        Status:         AdminInvitationStatusPending,
        Role:           role,
        Email:          email.NormalizeEmail(emailAddr),
        DisplayName:    displayName,
        SentAt:         t,
        ExpiresAt:      t.Add(time.Hour * 24 * AdminInvitationExpirationDays),
        AcceptedAt:     null.Time{},
        RejectedAt:     null.Time{},
        InvalidatedAt:  null.Time{},
        CreatedAt:      t,
        UpdatedAt:      t,
    }
}
```

### Key Points

- Generate URL-safe `InvitationCode` (not same as ID)
- Normalize email for consistent comparison
- Calculate `ExpiresAt` from creation time
- Use constants for expiration period

## State Transition Methods

### Accept

```go
func (m *AdminInvitation) Accept(t time.Time) (*AdminInvitation, error) {
    // Validate current state
    if m.Status != AdminInvitationStatusPending {
        return nil, errors.AdminInvitationAlreadyAcceptedErr.
            Errorf("status %s is not pending", m.Status).
            WithValue("admin_invitation", m)
    }

    // Check expiration
    if m.IsExpired(t) {
        return nil, errors.AdminInvitationExpiredErr.
            Errorf("invitation is expired").
            WithValue("admin_invitation", m)
    }

    // Apply state change
    m.Status = AdminInvitationStatusAccepted
    m.AcceptedAt = null.TimeFrom(t)
    m.UpdatedAt = t
    return m, nil
}
```

### Reject

```go
func (m *AdminInvitation) Reject(t time.Time) (*AdminInvitation, error) {
    switch m.Status {
    case AdminInvitationStatusPending:
        m.Status = AdminInvitationStatusRejected
        m.RejectedAt = null.TimeFrom(t)
        m.UpdatedAt = t
    case AdminInvitationStatusAccepted:
        return nil, errors.AdminInvitationAlreadyAcceptedErr.Errorf("already accepted")
    case AdminInvitationStatusRejected:
        return nil, errors.AdminInvitationAlreadyRejectedErr.Errorf("already rejected")
    case AdminInvitationStatusInvalidated:
        return nil, errors.AdminInvitationAlreadyInvalidatedErr.Errorf("already invalidated")
    }
    return m, nil
}
```

### Invalidate (Admin cancellation)

```go
func (m *AdminInvitation) Invalidate(t time.Time) (*AdminInvitation, error) {
    if m.Status != AdminInvitationStatusPending {
        return nil, errors.AdminInvitationAlreadyInvalidatedErr.
            Errorf("status %s is not pending", m.Status)
    }

    m.Status = AdminInvitationStatusInvalidated
    m.InvalidatedAt = null.TimeFrom(t)
    m.UpdatedAt = t
    return m, nil
}
```

## Helper Methods

```go
func (m *AdminInvitation) IsExpired(t time.Time) bool {
    return !m.ExpiresAt.After(t)
}

func (m *AdminInvitation) IsPending() bool {
    return m.Status == AdminInvitationStatusPending
}

func (m *AdminInvitation) IsInvalidated() bool {
    return m.Status == AdminInvitationStatusInvalidated
}
```

## Usecase: Create Invitation

```go
func (i *adminAdminInvitationInteractor) Create(
    ctx context.Context,
    param *input.AdminCreateAdminInvitation,
) (*model.AdminInvitation, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    // Authorization: Only root can create invitations
    if !param.AdminRole.IsRoot() {
        return nil, errors.AdminForbiddenErr.Errorf("only root admin can create invitation")
    }

    // Create domain entity
    invitation := model.NewAdminInvitation(
        param.Role,
        param.Email,
        param.DisplayName,
        param.RequestTime,
    )

    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        // Persist invitation
        if err := i.adminInvitationRepository.Create(ctx, invitation); err != nil {
            return err
        }

        // Send notification email
        if err := i.emailer.SendAdminInvitation(ctx, invitation); err != nil {
            return err
        }

        return nil
    }); err != nil {
        return nil, err
    }

    return i.adminInvitationRepository.Get(ctx, ...)
}
```

## Usecase: Accept Invitation

```go
func (i *adminAdminInvitationInteractor) Accept(
    ctx context.Context,
    param *input.AdminAcceptAdminInvitation,
) (*model.Admin, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    var admin *model.Admin

    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        // 1. Get invitation with lock
        invitation, err := i.adminInvitationRepository.Get(ctx, repository.GetAdminInvitationQuery{
            InvitationCode: null.StringFrom(param.InvitationCode),
            BaseGetOptions: repository.BaseGetOptions{
                OrFail:    true,
                ForUpdate: true,
            },
        })
        if err != nil {
            return err
        }

        // 2. Validate invitation state
        if invitation.IsInvalidated() {
            return errors.AdminInvitationAlreadyInvalidatedErr.New()
        }

        // 3. Verify email matches
        if invitation.Email != email.NormalizeEmail(param.Email) {
            return errors.AdminInvitationInvalidEmailErr.New()
        }

        // 4. Accept via domain method (validates state & expiration)
        invitation, err = invitation.Accept(param.RequestTime)
        if err != nil {
            return err
        }
        if err := i.adminInvitationRepository.Update(ctx, invitation); err != nil {
            return err
        }

        // 5. Create new user entity
        admin = model.NewAdmin(
            invitation.Role,
            param.AuthUID,
            invitation.Email,
            invitation.DisplayName,
            param.RequestTime,
        )
        if err := i.adminRepository.Create(ctx, admin); err != nil {
            return err
        }

        // 6. Sync with IdP (store claims)
        if err := i.adminAuthenticationRepository.StoreClaims(
            ctx,
            param.AuthUID,
            model.NewAdminClaims(
                param.AuthUID,
                admin.Email,
                null.StringFrom(admin.ID),
                nullable.TypeFrom(invitation.Role),
            ),
        ); err != nil {
            return err
        }

        // 7. Send acceptance confirmation email
        if err := i.emailer.SendAdminInvitationAccept(ctx, admin); err != nil {
            return err
        }

        return nil
    }); err != nil {
        return nil, err
    }

    return i.adminRepository.Get(ctx, ...)
}
```

## Usecase: Get by Invitation Code (Pre-auth)

For endpoints that need to be called before authentication:

```go
func (i *adminAdminInvitationInteractor) GetByInvitationCode(
    ctx context.Context,
    param *input.AdminGetAdminInvitationByInvitationCode,
) (*model.AdminInvitation, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    invitation, err := i.adminInvitationRepository.Get(ctx, repository.GetAdminInvitationQuery{
        InvitationCode: null.StringFrom(param.InvitationCode),
        BaseGetOptions: repository.BaseGetOptions{
            OrFail: true,
        },
    })
    if err != nil {
        return nil, err
    }

    // Validate usability (without requiring authentication)
    if invitation.IsInvalidated() {
        return nil, errors.AdminInvitationAlreadyInvalidatedErr.New()
    }
    if invitation.AcceptedAt.Valid {
        return nil, errors.AdminInvitationAlreadyAcceptedErr.New()
    }
    if invitation.IsExpired(param.RequestTime) {
        return nil, errors.AdminInvitationExpiredErr.New()
    }

    return invitation, nil
}
```

## Error Definitions

```go
var (
    AdminInvitationNotFoundErr           = &Error{code: "ADMIN_INVITATION_NOT_FOUND", message: "admin invitation not found"}
    AdminInvitationAlreadyAcceptedErr    = &Error{code: "ADMIN_INVITATION_ALREADY_ACCEPTED", message: "admin invitation already accepted"}
    AdminInvitationAlreadyRejectedErr    = &Error{code: "ADMIN_INVITATION_ALREADY_REJECTED", message: "admin invitation already rejected"}
    AdminInvitationAlreadyInvalidatedErr = &Error{code: "ADMIN_INVITATION_ALREADY_INVALIDATED", message: "admin invitation already invalidated"}
    AdminInvitationExpiredErr            = &Error{code: "ADMIN_INVITATION_EXPIRED", message: "admin invitation expired"}
    AdminInvitationInvalidEmailErr       = &Error{code: "ADMIN_INVITATION_INVALID_EMAIL", message: "email does not match invitation"}
)
```

## Interactor Separation

Separate invitation-related logic into dedicated interactor:

```go
// admin_admin.go - User CRUD operations
type AdminAdminInteractor interface {
    Get(ctx context.Context, param *input.AdminGetAdmin) (*model.Admin, error)
    List(ctx context.Context, param *input.AdminListAdmins) (*output.AdminListAdmins, error)
    Update(ctx context.Context, param *input.AdminUpdateAdmin) (*model.Admin, error)
    Delete(ctx context.Context, param *input.AdminDeleteAdmin) error
}

// admin_admin_invitation.go - Invitation workflow
type AdminAdminInvitationInteractor interface {
    GetByInvitationCode(ctx context.Context, param *input.AdminGetAdminInvitationByInvitationCode) (*model.AdminInvitation, error)
    Create(ctx context.Context, param *input.AdminCreateAdminInvitation) (*model.AdminInvitation, error)
    Accept(ctx context.Context, param *input.AdminAcceptAdminInvitation) (*model.Admin, error)
}
```

## Email Integration

Send emails within transaction to ensure atomicity:

```go
// message/emailer.go interface
type Emailer interface {
    // Send when invitation is created
    SendAdminInvitation(ctx context.Context, invitation *model.AdminInvitation) error
    // Send when invitation is accepted
    SendAdminInvitationAccept(ctx context.Context, admin *model.Admin) error
}
```

### On Create - Send invitation email

```go
// In usecase - within transaction
if err := i.emailer.SendAdminInvitation(ctx, invitation); err != nil {
    return err  // Transaction will rollback
}
```

### On Accept - Send acceptance confirmation email

```go
// In usecase - within transaction, after StoreClaims
if err := i.emailer.SendAdminInvitationAccept(ctx, admin); err != nil {
    return err  // Transaction will rollback
}
```

## Best Practices

1. **Expiration is mandatory** - Always set `ExpiresAt` for security
2. **Validate in domain** - State transitions validate in domain methods
3. **Email normalization** - Normalize emails for consistent comparison
4. **Separate invitation code from ID** - Use URL-safe code for external links
5. **Lock before update** - Use `ForUpdate` when accepting/rejecting
6. **Sync with IdP** - Store claims after user creation
7. **Transaction scope** - Include email sending in transaction for rollback
8. **Send confirmation on accept** - Notify user after successful acceptance
