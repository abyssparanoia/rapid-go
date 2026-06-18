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

### Field Type Conventions

- **Optional primitives & time → `null/v8`, NOT `nullable.Type[T]`.** Use `null.Int64`,
  `null.Bool`, `null.Float64`, `null.Time`, `null.String` for optional primitive/time fields
  (struct fields, `Update` method params, repository query filters, input DTOs). `nullable.Type[T]`
  (`internal/pkg/nullable`) is **only** for custom types that `null/v8` cannot represent — domain
  enums (e.g. `model.EstimateStatus`), domain structs, and `civil.Date`. There is a `null/v8` type
  for every primitive, so `nullable.Type[int64/bool/float64/time.Time/string]` is always wrong.

  ```go
  // BAD - nullable.Type for primitives
  UsageDays    nullable.Type[int64]
  HasTransport nullable.Type[bool]
  ProjectionDM nullable.Type[float64]

  // GOOD - null/v8 for primitives; nullable.Type only for custom types
  UsageDays    null.Int64
  HasTransport null.Bool
  ProjectionDM null.Float64
  Status       nullable.Type[EstimateStatus] // custom enum → nullable.Type is correct
  ```

- **Date-only fields → `nullable.Type[civil.Date]`** (not `null.String` / `time.Time`). A calendar date with no time component must be modeled as a real date type, not a stringly-typed `"YYYY-MM-DD"`. Use `nullable.Type[civil.Date]` (`cloud.google.com/go/civil`) for an optional date — this is a custom type `null/v8` does not cover; the marshaller converts to/from the DB `custom_types.NullDate`. Money amounts stay `int64`/`null.Int64`; NUMERIC rate/decimal columns become `float64`/`null.Float64`.

  ```go
  // BAD - date kept as a string
  DesiredDate null.String // YYYY-MM-DD

  // GOOD - real date type
  DesiredDate nullable.Type[civil.Date]
  ```

- **Document why a field is nullable.** When a column is nullable for a non-obvious reason — especially a master-data FK kept alongside denormalized snapshots (so the row survives master deletion / supports ad-hoc entries) — add a short comment on the field stating the reason. Don't leave reviewers guessing why an FK or value is optional.

  ```go
  // GOOD - the reason the FK is optional is explicit
  // LeaseProductID references the source master product. Nullable because the
  // line item keeps its own snapshots (ProductName, prices), so it must survive
  // even if the master product is removed, and ad-hoc lines have no master.
  LeaseProductID null.String
  ```

- **Issue human-facing identifiers in the constructor**, not as an empty/nullable field filled later. A user-visible number/code (e.g. an estimate number) should be generated in `New{Entity}` from the creation time and stored as a plain `string`, not left `null.String{}` for a separate setter.

  ```go
  // BAD - left empty, set later
  EstimateNumber null.String

  // GOOD - issued at creation in the constructor
  EstimateNumber string // e.g. fmt.Sprintf("%s-%d", t.In(now.JST).Format("20060102"), t.Unix())
  ```

- **Optionally-loaded relations go in `ReadonlyReference`, not top-level `nullable.Type` fields.** A relation that is only populated when preloaded (including a 1:1 owned detail surfaced for a "get with detail" endpoint) belongs inside the `ReadonlyReference` struct, so its presence clearly signals "loaded vs not". Reserve top-level direct fields for data that is *always* loaded with the entity (see Owned Child Entities – Cascade Requirements).

## ReadonlyReference Pattern

`ReadonlyReference` is used to hold **read-only** related entities that are optionally loaded by the repository.

### When to Use ReadonlyReference

| Relationship Type | Where to Define | Load Behavior | Write Behavior |
|-------------------|-----------------|---------------|----------------|
| Reference data (lookup) | `ReadonlyReference` | Optional (via `Preload`) | Never written together |
| Owned child entities | Direct field | Always loaded | Written together with parent |

#### Field type inside `ReadonlyReference`

Use the field type that matches the relation's **nullability**, not always a pointer:

- **Nullable relation** (the FK is nullable, or only one of several mutually-exclusive relations is ever set) → `nullable.Type[T]`. The `.Valid` flag distinguishes "loaded, but absent" from "present".
- **Required relation** (the FK is `NOT NULL`, so it is always present whenever the reference is loaded) → `*T` pointer.

```go
// estimate: assignee FK is nullable; created_by FK is NOT NULL;
// exactly one type-specific detail is ever set.
ReadonlyReference *struct {
    Assignee       nullable.Type[Admin] // nullable FK
    CreatedBy      *Admin               // NOT NULL FK
    LeaseEstimate  nullable.Type[LeaseEstimate]  // one-of (nullable)
    AwningEstimate nullable.Type[AwningEstimate]
    FabricEstimate nullable.Type[FabricEstimate]
}
```

### Owned 1:1 Child Field Type Convention (New Entities)

For **new** entities with optional owned 1:1 children (where the child may or may not exist), prefer `nullable.Type[ChildStruct]` over `*ChildStruct` pointer:

```go
// Good - nullable.Type for optional owned 1:1 child (new entities)
type PaymentMethod struct {
    ID         string
    CreditCard nullable.Type[PaymentMethodCreditCard]  // Present when type=card
    BankAccount nullable.Type[PaymentMethodBankAccount] // Present when type=bank
}

// Acceptable for existing entities (pointer used historically)
type PaymentMethod struct {
    CreditCard  *PaymentMethodCreditCard   // Existing code — not required to migrate
    BankAccount *PaymentMethodBankAccount
}
```

**Why `nullable.Type[T]` for new entities:**
- Explicit `.Valid` flag prevents nil-dereference without nil checks
- Consistent with other optional fields in the codebase
- Clear semantics: `.Valid` = child exists, `.Value()` = access child data

**Note:** Existing entities using `*ChildStruct` do not need to be migrated. Apply the `nullable.Type[T]` convention only when creating new owned 1:1 child relationships.

### Owned Child Entities – Cascade Requirements

When an entity **cannot exist without its parent** (fully-owned / composed relationship), it must be a **direct field**, not in `ReadonlyReference`. The repository layer must enforce cascading behavior:

| Operation | Requirement |
|-----------|-------------|
| `Get` / `List` | Load owned children **regardless of `Preload` flag** |
| `Create` | Insert children in the **same transaction** as the parent |
| `Update` | Upsert / replace children in the **same transaction** as the parent |
| `Delete` | Delete children in the same transaction (or via `ON DELETE CASCADE`) |

For 0..1 owned relationships, use `nullable.Type[ChildModel]` (not a pointer):

```go
// 1..N owned - direct field, always loaded, cascade writes
type Order struct {
    ID    string
    Items OrderItems  // Not in ReadonlyReference
}

// 0..1 owned - nullable.Type, not *StaffProfile
type Staff struct {
    ID      string
    Profile nullable.Type[StaffProfile]
}
```

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

### ReadonlyReference と Proto Partial Pattern

ドメインモデルの`ReadonlyReference`はprotoの`{Entity}Partial`にマッピング：

| Domain Model | Proto Response |
|--------------|----------------|
| `TenantID string` | Protoには露出しない（内部用） |
| `ReadonlyReference.Tenant` | `TenantPartial tenant`（必須フィールド） |

**重要**:
- Partialパターンでは`ReadonlyReference.Tenant`が常にロードされている必要があります
- Repository queryで`Preload: true`を必ず指定してください
- gRPC marshallerで`TenantPartialToPB(m.ReadonlyReference.Tenant)`を使用して変換

```go
// Usecase - Preload required for response
staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
    ID: null.StringFrom(param.StaffID),
    BaseGetOptions: repository.BaseGetOptions{
        OrFail:  true,
        Preload: true,  // 必須 - Tenantをロード
    },
})

// gRPC Marshaller - Convert to Partial
func StaffToPB(m *model.Staff) *admin_apiv1.Staff {
    var tenant *admin_apiv1.TenantPartial
    if m.ReadonlyReference != nil && m.ReadonlyReference.Tenant != nil {
        tenant = TenantPartialToPB(m.ReadonlyReference.Tenant)
    }
    return &admin_apiv1.Staff{
        Id:     m.ID,
        Tenant: tenant,  // TenantPartial (not TenantId string)
        // ...
    }
}
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

All List operations have a corresponding SortKey type for ordering results. This is a **unified specification** across the codebase.

### Pattern

```go
type ExampleSortKey string

const (
    ExampleSortKeyUnknown       ExampleSortKey = "unknown"
    ExampleSortKeyCreatedAtDesc ExampleSortKey = "created_at_desc"
    ExampleSortKeyCreatedAtAsc  ExampleSortKey = "created_at_asc"
    ExampleSortKeyNameAsc       ExampleSortKey = "name_asc"
    ExampleSortKeyNameDesc      ExampleSortKey = "name_desc"
)

func NewExampleSortKey(s string) ExampleSortKey {
    switch s {
    case ExampleSortKeyCreatedAtDesc.String(),
        ExampleSortKeyCreatedAtAsc.String(),
        ExampleSortKeyNameAsc.String(),
        ExampleSortKeyNameDesc.String():
        return ExampleSortKey(s)
    default:
        return ExampleSortKeyUnknown
    }
}

func (k ExampleSortKey) String() string {
    return string(k)
}

func (k ExampleSortKey) Valid() bool {
    return k != ExampleSortKeyUnknown && k != ""
}
```

### Naming Convention

- Type: `{Entity}SortKey`
- Constants: `{Entity}SortKey{Field}{Direction}` (e.g., `StaffSortKeyCreatedAtDesc`)
- Always include `Unknown` constant
- Common fields: `CreatedAt`, `UpdatedAt`, entity-specific fields

### Default Value

The default sort key is **always** `CreatedAtDesc`. This default is applied in the input layer constructor, not in the domain model.

## Method Naming Convention

All mutating methods on domain models follow a strict naming scheme based on **semantics**, not syntax.

| Prefix / verb | Meaning | Signature | Examples |
|---|---|---|---|
| `Set{X}` | Fill a value that was previously empty (nullable/derived field, never overwrites meaningful state) | `*Entity` | `Invoice.SetStripe`, `InvoiceItem.SetStripe`, `InvoiceTax.SetStripe`, `Payment.SetStripe`, `Refund.SetStripe`, `Organization.SetStripe`, `Staff.SetImageURL` |
| `Update` / `Update{Field}` | Overwrite an existing value with no domain ceremony | `*Entity` or `(*Entity, error)` if it validates | `Tenant.Update`, `Admin.Update`, `Staff.Update`, `Invoice.UpdatePaymentMethodID`, `PaymentMethod.UpdateIsDefault` |
| Specific verb (`Finalize`, `MarkPaid`, `Void`, `Pay`, `Delete`…) | A domain state transition with business meaning | `(*Entity, error)` | `Invoice.Finalize`, `Payment.MarkSucceeded`, `PaymentMethod.Delete` |

**Rule of thumb**: if the field was guaranteed empty before, it's `Set`; if you are replacing a possibly-meaningful value with a generic one, it's `Update`; if the change carries business meaning or guards a state machine, use a specific verb.

### Return the Receiver

**All mutating methods must return the receiver** (`*Entity`), and `(*Entity, error)` when they also validate a precondition. This enables chaining and ensures a consistent shape across the codebase.

```go
// Good - returns receiver
func (m *Admin) Update(displayName null.String, t time.Time) *Admin {
    if displayName.Valid { m.DisplayName = displayName.String }
    m.UpdatedAt = t
    return m
}

// Good - returns (receiver, error) when validating
func (m *PaymentMethod) UpdateCreditCardDetails(...) (*PaymentMethod, error) {
    if m.PaymentMethodType != PaymentMethodTypeCard {
        return nil, errors.PaymentMethodTypeMismatchErr.New()...
    }
    // mutate
    return m, nil
}
```

Void mutators (`func (m *Entity) Mutate(...)`) are **forbidden** on domain models.

### Validate Inside the Method

If a precondition can be checked from the object's own state, the check belongs **inside the model method** (returning a domain error), never in the usecase/service immediately before the call. Duplicating it externally risks other callers missing it.

### Predicates over Inline Field Inspection

Usecase/service code must not branch on a model's raw exported fields. Expose an intent-revealing predicate on the model instead:

```go
// Bad - usecase inspects raw fields
if inv.Status != model.InvoiceStatusDraft { return nil }

// Good - model exposes a predicate
if !inv.IsDraft() { return nil }
```

Common predicates: `IsDraft() bool`, `HasPaymentMethod(id string) bool`, `IsHostedURLAvailable() bool`, `IsStripe() bool`, `IsDownloadable() bool`.

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
