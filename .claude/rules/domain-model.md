---
description: Domain model and entity conventions for the domain layer
globs:
  - "internal/domain/model/**/*.go"
---

# Domain Model Guidelines

## Core Principles

- Domain models are **pure business logic** with no external dependencies
- No database packages, no infrastructure concerns
- Use standard library and domain-specific packages only

## Entity Structure

```go
type Example struct {
    // Primary fields
    ID          string
    TenantID    string
    Name        string
    Status      ExampleStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time

    // Relations (read-only, populated by repository)
    ReadonlyReference *struct {
        Tenant *Tenant
        // Add related entities here
    }
}
```

### Key Points

- Use `string` for IDs (not `uuid.UUID`)
- Use custom types for status/enum fields
- Relations are in `ReadonlyReference` struct pointer (nil when not loaded)
- Never modify `ReadonlyReference` in domain logic

## ReadonlyReference Pattern

`ReadonlyReference` is used to hold **read-only** related entities that are optionally loaded by the repository.

### When to Use ReadonlyReference

| Relationship Type | Where to Define | Load Behavior | Write Behavior |
|-------------------|-----------------|---------------|----------------|
| Reference data (lookup) | `ReadonlyReference` | Optional (via `Preload`) | Never written together |
| Owned child entities | Direct field | Always loaded | Written together with parent |

```go
// ReadonlyReference - for reference/lookup relations
type Example struct {
    ID        string
    TenantID  string
    CreatedAt time.Time

    // Tenant is reference data, not owned by Example
    ReadonlyReference *struct {
        Tenant *Tenant
    }
}

// Direct field - for owned/composed entities
type Order struct {
    ID        string
    Items     OrderItems  // Always loaded, written together
    CreatedAt time.Time
}
```

### Rules for ReadonlyReference

1. **ReadonlyReference within ReadonlyReference is always nil**
   - When loading related entities, do NOT recursively load their ReadonlyReferences
   - This prevents circular dependencies and excessive data loading

```go
// In marshaller - correct pattern
func ExampleToModel(e *dbmodel.Example) *model.Example {
    m := &model.Example{
        ID:                e.ID,
        TenantID:          e.TenantID,
        ReadonlyReference: nil,
    }

    if e.R != nil && e.R.Tenant != nil {
        // Tenant's ReadonlyReference is NOT populated (remains nil)
        m.ReadonlyReference = &struct{ Tenant *model.Tenant }{
            Tenant: TenantToModel(e.R.Tenant),  // Tenant.ReadonlyReference = nil
        }
    }
    return m
}
```

2. **Never modify ReadonlyReference in domain logic**
   - ReadonlyReference is populated only by repository/marshaller
   - Domain methods should not touch ReadonlyReference

3. **Constructor always sets ReadonlyReference to nil**
   - New entities don't have loaded relations

```go
func NewExample(...) *Example {
    return &Example{
        ID:                id.New(),
        // ...
        ReadonlyReference: nil,  // Always nil in constructor
    }
}
```

4. **Use Preload option to load relations**
   - Repository query options control whether relations are loaded

```go
// Without preload - ReadonlyReference is nil
example, _ := repo.Get(ctx, GetExampleQuery{ID: id})
// example.ReadonlyReference == nil

// With preload - ReadonlyReference is populated
example, _ := repo.Get(ctx, GetExampleQuery{
    ID: id,
    BaseGetOptions: BaseGetOptions{Preload: true},
})
// example.ReadonlyReference.Tenant != nil
```

## Type Aliases

```go
type Staffs []*Staff
type Tenants []*Tenant
type Assets []*Asset
```

Always define:
- `{Entity}s` (plural) - for slices
- `{Entity}MapByID` - for batch operations (optional)

## Constructor Pattern

### With ReadonlyReference (Staff)

```go
func NewStaff(
    tenantID string,
    role StaffRole,
    authUID string,
    displayName string,
    imagePath string,
    email string,
    t time.Time,
) *Staff {
    return &Staff{
        ID:          id.New(),
        TenantID:    tenantID,
        Role:        role,
        AuthUID:     authUID,
        DisplayName: displayName,
        ImagePath:   imagePath,
        Email:       email,
        CreatedAt:   t,
        UpdatedAt:   t,

        ReadonlyReference: nil,

        ImageURL: null.String{},
    }
}
```

### Without ReadonlyReference (Tenant)

```go
func NewTenant(
    name string,
    t time.Time,
) *Tenant {
    return &Tenant{
        ID:        id.New(),
        Name:      name,
        Tags:      make(TenantTags, 0),
        CreatedAt: t,
        UpdatedAt: t,
    }
}
```

- Constructor generates new ID using `internal/pkg/id.New()`
- Both `CreatedAt` and `UpdatedAt` use the same time parameter
- ReadonlyReference is always nil in constructor (when present)

## Update Methods

```go
func (m *Tenant) Update(
    name null.String,
    t time.Time,
) {
    if name.Valid {
        m.Name = name.String
    }

    m.UpdatedAt = t
}
```

- Use `null.String`, `null.Int64`, etc. for optional update fields
- Always update `UpdatedAt`
- Return type can be void or `*Entity` for method chaining

## Status/Enum Types

```go
type AssetType string

const (
    AssetTypeUnknown   AssetType = "unknown"
    AssetTypeUserImage AssetType = "private/user_images"
)

func NewAssetType(str string) AssetType {
    switch str {
    case AssetTypeUserImage.String():
        return AssetType(str)
    default:
        return AssetTypeUnknown
    }
}

func (m AssetType) String() string {
    return string(m)
}

func (m AssetType) Valid() bool {
    return m != "" && m != AssetTypeUnknown
}

// Helper methods for type-specific checks
func (m AssetType) IsPrivate() bool {
    return strings.HasPrefix(m.String(), "private")
}

func (m AssetType) IsPublic() bool {
    return strings.HasPrefix(m.String(), "public")
}
```

- Always include `Unknown` as first constant
- Implement `String()` and `Valid()` methods
- Add `New{Type}(str string)` constructor for parsing
- Add helper methods for type-specific checks when useful

## Sort Key Types

```go
type ExampleSortKey string

const (
    ExampleSortKeyUnknown       ExampleSortKey = "unknown"
    ExampleSortKeyCreatedAtDesc ExampleSortKey = "created_at_desc"
    ExampleSortKeyNameAsc       ExampleSortKey = "name_asc"
)

func (k ExampleSortKey) Valid() bool {
    return k != ExampleSortKeyUnknown && k != ""
}
```

## State Change Methods (Domain Logic First)

**IMPORTANT**: All state changes MUST be performed through domain model methods, not direct field assignment in usecase layer.

### Why Domain Methods

- Encapsulates business rules within the domain
- Ensures consistent validation and state transitions
- Makes business logic testable and reusable
- Prevents scattered logic across usecase layer

### Pattern: Single Field Update

```go
// Good - Domain method for role update
func (m *Admin) UpdateRole(role AdminRole, t time.Time) *Admin {
    m.Role = role
    m.UpdatedAt = t
    return m
}

// Usage in usecase
admin.UpdateRole(param.Role.Value(), param.RequestTime)
```

### Pattern: State Transition with Validation

```go
// Good - Domain method with validation
func (m *Invitation) Accept(t time.Time) (*Invitation, error) {
    if m.Status != InvitationStatusPending {
        return nil, errors.InvitationAlreadyAcceptedErr.
            Errorf("status %s is not pending", m.Status)
    }
    if m.IsExpired(t) {
        return nil, errors.InvitationExpiredErr.Errorf("invitation is expired")
    }

    m.Status = InvitationStatusAccepted
    m.AcceptedAt = null.TimeFrom(t)
    m.UpdatedAt = t
    return m, nil
}
```

### Anti-Pattern: Direct Field Assignment

```go
// Bad - Direct assignment in usecase
func (i *interactor) Update(...) {
    admin.Role = param.Role  // Don't do this!
    admin.UpdatedAt = param.RequestTime
}

// Good - Use domain method
func (i *interactor) Update(...) {
    admin.UpdateRole(param.Role, param.RequestTime)
}
```

## Role/Type Helper Methods

For enum types representing roles or permissions, add helper methods:

```go
type AdminRole string

const (
    AdminRoleUnknown AdminRole = "unknown"
    AdminRoleRoot    AdminRole = "root"
    AdminRoleNormal  AdminRole = "normal"
)

// Role check helpers
func (r AdminRole) IsRoot() bool {
    return r == AdminRoleRoot
}

func (r AdminRole) IsNormal() bool {
    return r == AdminRoleNormal
}

func (r AdminRole) Valid() bool {
    return r == AdminRoleRoot || r == AdminRoleNormal
}
```

### Usage in Authorization

```go
// In usecase - clear authorization check
if !param.AdminRole.IsRoot() {
    return nil, errors.AdminForbiddenErr.Errorf("only root admin can perform this action")
}
```

## Helper Methods on Slices

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
```
