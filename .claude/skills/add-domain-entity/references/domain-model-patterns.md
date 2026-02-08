# Domain Model Patterns

Detailed code patterns for creating domain models.

## Entity Structure

```go
package model

import (
    "time"

    "github.com/abyssparanoia/rapid-go/internal/pkg/id"
    "github.com/aarondl/null/v9"
)

type Example struct {
    ID          string
    TenantID    string
    Name        string
    Description string
    Status      ExampleStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time

    // Relations (read-only, populated by repository)
    ReadonlyReference *struct {
        Tenant *Tenant
    }

    // Computed fields (set by service layer)
    ImageURL null.String
}
```

### Field Guidelines

| Field Type         | Go Type       | Notes                               |
| ------------------ | ------------- | ----------------------------------- |
| ID                 | `string`      | Use `id.New()` in constructor       |
| Foreign Key        | `string`      | Reference to parent entity          |
| Text               | `string`      | Short text, names                   |
| Status/Enum        | Custom type   | See Status Types section            |
| Timestamp          | `time.Time`   | Always required                     |
| Optional Timestamp | `null.Time`   | For nullable DB columns             |
| Computed           | `null.String` | Set by service layer, not persisted |

## Type Aliases

```go
type Examples []*Example
type ExampleMapByID map[string]*Example
```

## Constructor

```go
func NewExample(
    tenantID string,
    name string,
    description string,
    t time.Time,
) *Example {
    return &Example{
        ID:          id.New(),
        TenantID:    tenantID,
        Name:        name,
        Description: description,
        Status:      ExampleStatusDraft,  // Default status
        CreatedAt:   t,
        UpdatedAt:   t,

        ReadonlyReference: nil,  // Always nil in constructor

        ImageURL: null.String{},  // Empty for computed fields
    }
}
```

### Constructor Rules

1. Generate ID using `id.New()`
2. Set both `CreatedAt` and `UpdatedAt` to the same time parameter
3. Set `ReadonlyReference` to nil
4. Initialize computed fields as empty (`null.String{}`)
5. Set default status if applicable

## Update Method

```go
func (m *Example) Update(
    name null.String,
    description null.String,
    t time.Time,
) *Example {
    if name.Valid {
        m.Name = name.String
    }
    if description.Valid {
        m.Description = description.String
    }
    m.UpdatedAt = t
    return m
}
```

### Update Method Rules

1. Use `null.*` types for optional update fields
2. Check `.Valid` before applying changes
3. Always update `UpdatedAt`
4. Return `*Entity` for method chaining

## State Change Methods

For complex state transitions, create dedicated methods:

```go
func (m *Example) Publish(t time.Time) (*Example, error) {
    if m.Status != ExampleStatusDraft {
        return nil, errors.ExampleInvalidStatusErr.Errorf(
            "cannot publish: current status=%s", m.Status,
        )
    }

    m.Status = ExampleStatusPublished
    m.UpdatedAt = t
    return m, nil
}

func (m *Example) Archive(t time.Time) (*Example, error) {
    if m.Status == ExampleStatusArchived {
        return nil, errors.ExampleAlreadyArchivedErr.New()
    }

    m.Status = ExampleStatusArchived
    m.UpdatedAt = t
    return m, nil
}
```

## Helper Methods on Entity

```go
func (m *Example) Exist() bool {
    return m != nil
}

func (m *Example) SetImageURL(url string) {
    m.ImageURL = null.StringFrom(url)
}

func (m *Example) IsDraft() bool {
    return m.Status == ExampleStatusDraft
}

func (m *Example) IsPublished() bool {
    return m.Status == ExampleStatusPublished
}
```

## Helper Methods on Slice

```go
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

func (es Examples) FilterByStatus(status ExampleStatus) Examples {
    result := make(Examples, 0)
    for _, e := range es {
        if e.Status == status {
            result = append(result, e)
        }
    }
    return result
}
```

## Status/Enum Types

### Definition

```go
type ExampleStatus string

const (
    ExampleStatusUnknown   ExampleStatus = "unknown"
    ExampleStatusDraft     ExampleStatus = "draft"
    ExampleStatusPublished ExampleStatus = "published"
    ExampleStatusArchived  ExampleStatus = "archived"
)
```

### Required Methods

```go
func NewExampleStatus(str string) ExampleStatus {
    switch str {
    case ExampleStatusDraft.String():
        return ExampleStatusDraft
    case ExampleStatusPublished.String():
        return ExampleStatusPublished
    case ExampleStatusArchived.String():
        return ExampleStatusArchived
    default:
        return ExampleStatusUnknown
    }
}

func (s ExampleStatus) String() string {
    return string(s)
}

func (s ExampleStatus) Valid() bool {
    return s != ExampleStatusUnknown && s != ""
}
```

## Sort Key Types

```go
type ExampleSortKey string

const (
    ExampleSortKeyUnknown       ExampleSortKey = "unknown"
    ExampleSortKeyCreatedAtDesc ExampleSortKey = "created_at_desc"
    ExampleSortKeyCreatedAtAsc  ExampleSortKey = "created_at_asc"
    ExampleSortKeyNameAsc       ExampleSortKey = "name_asc"
    ExampleSortKeyNameDesc      ExampleSortKey = "name_desc"
)

func (k ExampleSortKey) Valid() bool {
    return k != ExampleSortKeyUnknown && k != ""
}
```

## Role Types with Authorization Helpers

```go
type ExampleRole string

const (
    ExampleRoleUnknown ExampleRole = "unknown"
    ExampleRoleAdmin   ExampleRole = "admin"
    ExampleRoleMember  ExampleRole = "member"
    ExampleRoleViewer  ExampleRole = "viewer"
)

func (r ExampleRole) IsAdmin() bool {
    return r == ExampleRoleAdmin
}

func (r ExampleRole) CanEdit() bool {
    return r == ExampleRoleAdmin || r == ExampleRoleMember
}

func (r ExampleRole) Valid() bool {
    return r == ExampleRoleAdmin || r == ExampleRoleMember || r == ExampleRoleViewer
}
```

## Entity with Invitation Pattern

For entities that follow an invitation workflow:

```go
type ExampleInvitation struct {
    ID             string
    InvitationCode string  // URL-safe unique code
    Status         ExampleInvitationStatus
    Email          string
    ExpiresAt      time.Time
    AcceptedAt     null.Time
    RejectedAt     null.Time
    CreatedAt      time.Time
    UpdatedAt      time.Time

    ReadonlyReference *struct {
        Inviter *User
    }
}

const ExampleInvitationExpirationDays = 7

func NewExampleInvitation(
    email string,
    t time.Time,
) *ExampleInvitation {
    return &ExampleInvitation{
        ID:             id.New(),
        InvitationCode: uuid.UUIDBase64(),  // URL-safe
        Status:         ExampleInvitationStatusPending,
        Email:          email.NormalizeEmail(email),
        ExpiresAt:      t.Add(time.Hour * 24 * ExampleInvitationExpirationDays),
        AcceptedAt:     null.Time{},
        RejectedAt:     null.Time{},
        CreatedAt:      t,
        UpdatedAt:      t,

        ReadonlyReference: nil,
    }
}

func (m *ExampleInvitation) IsExpired(t time.Time) bool {
    return !m.ExpiresAt.After(t)
}

func (m *ExampleInvitation) Accept(t time.Time) (*ExampleInvitation, error) {
    if m.Status != ExampleInvitationStatusPending {
        return nil, errors.ExampleInvitationAlreadyAcceptedErr.New()
    }
    if m.IsExpired(t) {
        return nil, errors.ExampleInvitationExpiredErr.New()
    }

    m.Status = ExampleInvitationStatusAccepted
    m.AcceptedAt = null.TimeFrom(t)
    m.UpdatedAt = t
    return m, nil
}
```

See `.claude/rules/invitation-workflow.md` for complete invitation patterns.
