# AI Coding Anti-Patterns

Patterns that AI models commonly introduce. Check every changed file against these.

---

## Tests

### 1. `gomock.Any()` for non-context parameters

**Rule**: `gomock.Any()` is ONLY allowed for `context.Context` (first param). All other params must use exact values.

```go
// BAD
mockRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(staff, nil)
mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(staff, nil)

// GOOD
mockRepo.EXPECT().
    Get(gomock.Any(), repository.GetStaffQuery{
        ID: null.StringFrom(staff.ID),
        BaseGetOptions: repository.BaseGetOptions{OrFail: true, Preload: true},
    }).
    Return(staff, nil)
```

**Fix**: Replace `gomock.Any()` for non-context params with exact expected structs/values using domain constructors.

**No exceptions.** Every non-context parameter must be exact match. Using `DoAndReturn` to bypass parameter matching is also prohibited (see #5).

---

### 2. Missing `t.Parallel()` in table-driven tests

```go
// BAD
for name, tc := range tests {
    t.Run(name, func(t *testing.T) {
        tc := tc(ctx, ctrl)
        // ...
    })
}

// GOOD
for name, tc := range tests {
    t.Run(name, func(t *testing.T) {
        t.Parallel()
        // ...
    })
}
```

---

### 3. Not using table-driven test pattern

```go
// BAD - separate test functions per case
func TestCreate_Success(t *testing.T) { ... }
func TestCreate_Error(t *testing.T) { ... }

// GOOD - map[string]testcaseFunc
tests := map[string]testcaseFunc{
    "success": func(ctx context.Context, ctrl *gomock.Controller) testcase { ... },
    "invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase { ... },
    "not found": func(ctx context.Context, ctrl *gomock.Controller) testcase { ... },
}
```

---

### 4. Missing required test cases

Every interactor method must have tests for:
- `"invalid argument"` – validation error (empty required fields)
- `"not found"` – entity doesn't exist (for Get/Update/Delete)
- `"success"` – happy path

---

### 5. `DoAndReturn` used to avoid exact matching

```go
// BAD - hides what values are actually passed
mockRepo.EXPECT().
    Create(gomock.Any(), gomock.Any()).
    DoAndReturn(func(ctx context.Context, m *model.Admin) error {
        return nil
    })

// GOOD - use exact object
admin := model.NewAdmin(role, authUID, email, displayName, requestTime)
mockRepo.EXPECT().Create(gomock.Any(), admin).Return(nil)
```

---

### 6. Direct model initialization instead of using factory in tests

Test data must be created via `factory.NewFactory()`, not by directly initializing domain model structs. Use `factory.CloneValue()` when a modified copy is needed.

```go
// BAD - direct struct initialization with hardcoded values
staff := &model.Staff{
    ID:        "test-id",
    TenantID:  "tenant-id",
    Role:      model.StaffRoleAdmin,
    Email:     "test@example.com",
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}

// GOOD - use factory
testdata := factory.NewFactory()
staff := testdata.Staff

// GOOD - modified copy via CloneValue
testdata := factory.NewFactory()
updatedStaff := &model.Staff{}
factory.CloneValue(testdata.Staff, updatedStaff)
updatedStaff.DisplayName = "Updated Name"
```

**Why**: The factory generates consistent, cross-referenced test data with faker. Direct initialization risks missing fields, inconsistent references (e.g., `Staff.TenantID` not matching `Tenant.ID`), and breaks when new required fields are added to the model.

**Exception**: Empty structs used solely as `CloneValue` targets (`&model.Staff{}`) are acceptable.

---

## Domain Layer

### 7. Direct field assignment instead of domain methods

```go
// BAD - in usecase
admin.Role = param.Role
admin.UpdatedAt = param.RequestTime

// GOOD - use domain method
admin.UpdateRole(param.Role, param.RequestTime)
```

**Fix**: Move the assignment into a domain method on the model, then call it from usecase.

---

### 8. State transition without validation

```go
// BAD - no current state check
func (m *Invitation) Accept(t time.Time) {
    m.Status = InvitationStatusAccepted
    m.AcceptedAt = null.TimeFrom(t)
    m.UpdatedAt = t
}

// GOOD - validate state before transition
func (m *Invitation) Accept(t time.Time) (*Invitation, error) {
    if m.Status != InvitationStatusPending {
        return nil, errors.InvitationAlreadyAcceptedErr.Errorf("status %s is not pending", m.Status)
    }
    if m.IsExpired(t) {
        return nil, errors.InvitationExpiredErr.New()
    }
    m.Status = InvitationStatusAccepted
    m.AcceptedAt = null.TimeFrom(t)
    m.UpdatedAt = t
    return m, nil
}
```

---

### 9. Enum type missing `Unknown` constant

```go
// BAD
type Status string
const (
    StatusActive Status = "active"
)

// GOOD - Unknown always first
type Status string
const (
    StatusUnknown Status = "unknown"
    StatusActive  Status = "active"
)
```

---

### 10. Missing `Valid()` and `String()` on enum types

Every custom type must have both methods.

---

### 11. `ReadonlyReference` set in constructor

```go
// BAD - constructor sets ReadonlyReference
func NewStaff(...) *Staff {
    return &Staff{
        ReadonlyReference: &struct{ Tenant *Tenant }{Tenant: tenant},
    }
}

// GOOD - always nil in constructor
func NewStaff(...) *Staff {
    return &Staff{
        ReadonlyReference: nil,
    }
}
```

---

### 12. Fully-owned (completely dependent) entity placed in ReadonlyReference

Entities that **cannot exist without the parent** (owned/composed) must be direct fields, not in `ReadonlyReference`. `ReadonlyReference` is only for lookup/reference data that exists independently.

| Relationship type | Placement | Load behavior | Write behavior |
|---|---|---|---|
| Owned child (1..N) | Direct field `Items OrderItems` | Always loaded (ignore Preload flag) | Same TX as parent |
| Owned child (0..1) | Direct field `Profile nullable.Type[Profile]` | Always loaded | Same TX as parent |
| Reference / lookup | `ReadonlyReference.Xxx` | Optional (Preload flag) | Never written together |

```go
// BAD - OrderItems cannot exist without Order, but placed in ReadonlyReference
type Order struct {
    ID string
    ReadonlyReference *struct{ Items OrderItems }
}

// GOOD - direct field, always loaded, written together with parent
type Order struct {
    ID    string
    Items OrderItems
}

// GOOD - 0..1 owned relation uses nullable.Type (not pointer)
type Staff struct {
    ID      string
    Profile nullable.Type[StaffProfile]
}
```

**Repository cascade requirements (check these when reviewing):**
- `Get`/`List`: load owned children regardless of `Preload` flag
- `Create`: insert children in the same transaction as parent
- `Update`: upsert/replace children in the same transaction as parent
- `Delete`: delete children in the same transaction as parent (or via `ON DELETE CASCADE`)

---

## Repository Layer

### 13. Using pointer type for optional enum/custom type fields

```go
// BAD
type ListQuery struct {
    Status  *model.Status
    SortKey *model.SortKey
}

// GOOD
type ListQuery struct {
    Status  nullable.Type[model.Status]
    SortKey nullable.Type[model.SortKey]
}
```

---

### 14. Missing `transactable.GetContextExecutor(ctx)`

```go
// BAD
dbmodel.Examples(mods...).One(ctx, db)

// GOOD
dbmodel.Examples(mods...).One(ctx, transactable.GetContextExecutor(ctx))
```

---

### 15. Sorting applied after pagination

```go
// BAD - wrong order
mods = append(mods, qm.Limit(limit), qm.Offset(offset))
mods = append(mods, qm.OrderBy("`created_at` DESC"))  // After pagination!

// GOOD - sort before paginate
mods = append(mods, qm.OrderBy("`created_at` DESC"))
mods = append(mods, qm.Limit(limit), qm.Offset(offset))
```

---

### 16. `SortKey` unknown case does not return error

```go
// BAD - silently skips
case model.ExampleSortKeyUnknown:
    // no sorting

// GOOD - return error
case model.ExampleSortKeyUnknown:
    return nil, errors.InternalErr.Errorf("invalid sort key: %s", query.SortKey.Value())
```

---

### 17. Related entity's `ReadonlyReference` populated in marshaller

```go
// BAD - recursive population
func ExampleToModel(e *dbmodel.Example) *model.Example {
    if e.R != nil && e.R.Tenant != nil {
        m.ReadonlyReference = &struct{ Tenant *model.Tenant }{
            Tenant: &model.Tenant{
                ReadonlyReference: &struct{...}{...},  // BAD - recursive!
            },
        }
    }
}

// GOOD - related entity ReadonlyReference is always nil
Tenant: TenantToModel(e.R.Tenant),  // TenantToModel returns Tenant with ReadonlyReference=nil
```

---

## Usecase Layer

### 18. Missing `param.Validate()` call

Every interactor method must start with:
```go
if err := param.Validate(); err != nil {
    return nil, err
}
```

---

### 19. Missing `ForUpdate: true` before update/delete

```go
// BAD - no lock
entity, err := i.repo.Get(ctx, GetQuery{ID: id, BaseGetOptions: BaseGetOptions{OrFail: true}})
entity.Update(...)
i.repo.Update(ctx, entity)

// GOOD - lock row before modification
entity, err := i.repo.Get(ctx, GetQuery{
    ID: id,
    BaseGetOptions: BaseGetOptions{OrFail: true, ForUpdate: true},
})
```

---

### 20. IdP sync (`StoreClaims` / `DeleteUser`) outside transaction

```go
// BAD
if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
    return i.repo.Update(ctx, entity)
}); err != nil {
    return nil, err
}
i.authRepo.StoreClaims(ctx, ...)  // OUTSIDE TRANSACTION!

// GOOD - inside RWTx
if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
    if err := i.repo.Update(ctx, entity); err != nil {
        return err
    }
    return i.authRepo.StoreClaims(ctx, ...)  // Inside TX
}); err != nil { ... }
```

---

### 21. Missing `Preload: true` on returned entities

```go
// BAD
return i.repo.Get(ctx, GetQuery{ID: id, BaseGetOptions: BaseGetOptions{OrFail: true}})

// GOOD - always preload for returned entities
return i.repo.Get(ctx, GetQuery{
    ID: id,
    BaseGetOptions: BaseGetOptions{OrFail: true, Preload: true},
})
```

---

### 22. Missing `BatchSetXxxURLs` call after fetching entities

Even if the entity currently has no asset fields, always call:
```go
if err := i.assetService.BatchSetStaffURLs(ctx, model.Staffs{staff}, param.RequestTime); err != nil {
    return nil, err
}
```

---

### 23. Using `null.StringFrom(*req.Field)` instead of `null.StringFromPtr(req.Field)`

```go
// BAD
var displayName null.String
if req.DisplayName != nil {
    displayName = null.StringFrom(*req.DisplayName)
}

// GOOD
null.StringFromPtr(req.DisplayName)
```

---

### 24. IdP deleted after DB deletion (wrong order)

On user deletion, IdP must be deleted FIRST:
```go
// GOOD order
i.authRepo.DeleteUser(ctx, admin.AuthUID)  // First
i.repo.Delete(ctx, id)                     // Then DB
```

---

### 25. Private method in usecase interactor

Interactor implementations must NOT define private methods. Only implement public interface methods.

```go
// BAD - private helper in interactor
func (i *adminStaffInteractor) Create(ctx context.Context, param *input.AdminCreateStaff) (*model.Staff, error) {
    staff, err := i.buildStaff(param)  // calls private helper
    ...
}
func (i *adminStaffInteractor) buildStaff(param *input.AdminCreateStaff) *model.Staff { ... }

// GOOD (A) - inline the logic
func (i *adminStaffInteractor) Create(ctx context.Context, param *input.AdminCreateStaff) (*model.Staff, error) {
    staff := model.NewStaff(param.TenantID, param.Role, param.AuthUID, ...)
    ...
}

// GOOD (B) - extract to domain service if logic spans entities
func (i *adminStaffInteractor) Create(ctx context.Context, param *input.AdminCreateStaff) (*model.Staff, error) {
    staff, err := i.staffService.Create(ctx, service.StaffCreateParam{...})
    ...
}
```

**Decision guide:**
| Situation | Solution |
|---|---|
| Short logic (few lines) | Inline in the public method |
| Logic spans multiple entities | Extract to `domain/service/` |
| Single entity state transition | Add method to `domain/model/` |
| Repository query building | Inline (do not create helper) |

---

## gRPC Handler Layer

### 26. Wrong method ordering in handler file

Required order: Get → List → Create → Custom(no ID) → Update → Custom(with ID) → Delete

```go
// BAD - wrong order
func (h *Handler) UpdateStaff(...)  // Update before List
func (h *Handler) ListStaffs(...)
func (h *Handler) GetStaff(...)

// GOOD
func (h *Handler) GetStaff(...)
func (h *Handler) ListStaffs(...)
func (h *Handler) UpdateStaff(...)
```

---

### 27. Optional proto enum field not using `nullable.TypeFromPtr` / manual nil check

```go
// BAD - using pointer variable before constructor
var status *model.Status
if req.Status != nil {
    s := marshaller.StatusToModel(*req.Status)
    status = &s
}
param := input.NewAdminList(status)

// GOOD
param := &input.AdminList{}
if req.Status != nil {
    status := marshaller.StatusToModel(*req.Status)
    param.Status = nullable.TypeFrom(status)
}
```

---

### 28. Nullable timestamp field without variable declaration pattern in marshaller

```go
// BAD - inline ternary / anonymous function
return &pb.Example{
    AcceptedAt: func() *timestamppb.Timestamp {
        if m.AcceptedAt.Valid { return timestamppb.New(m.AcceptedAt.Time) }
        return nil
    }(),
}

// GOOD - variable declaration first
var acceptedAt *timestamppb.Timestamp
if m.AcceptedAt.Valid {
    acceptedAt = timestamppb.New(m.AcceptedAt.Time)
}
return &pb.Example{AcceptedAt: acceptedAt}
```

---

## Proto Definition

### 29. `SortKey` enum defined after the field that uses it

```protobuf
// BAD
optional ListStaffsSortKey sort_key = 4;
enum ListStaffsSortKey { ... }  // After field!

// GOOD - enum before field
enum ListStaffsSortKey { ... }
optional ListStaffsSortKey sort_key = 4;
```

---

### 30. Enum not starting with `_UNSPECIFIED = 0`

```protobuf
// BAD
enum StaffRole {
  STAFF_ROLE_NORMAL = 0;
}

// GOOD
enum StaffRole {
  STAFF_ROLE_UNSPECIFIED = 0;
  STAFF_ROLE_NORMAL = 1;
}
```

---

## General AI Patterns (Across All Layers)

### 31. Package-level private functions in domain/usecase/infrastructure layers

Package-level private (non-receiver) functions are prohibited in domain, usecase, and infrastructure layers. They scatter business logic and reduce testability.

```go
// BAD - package-level private functions
func buildStaffQuery(id string) repository.GetStaffQuery { ... }
func toStaffModel(e *dbmodel.Staff) *model.Staff { ... }
func validateStaffInput(param *input.AdminCreateStaff) error { ... }

// GOOD (A) - inline short logic
staff, err := i.repo.Get(ctx, repository.GetStaffQuery{
    ID: null.StringFrom(id),
    BaseGetOptions: repository.BaseGetOptions{OrFail: true},
})

// GOOD (B) - use domain model method for single entity logic
func (m *Staff) Validate() error { ... }

// GOOD (C) - use domain service struct method for multi-entity logic
func (s *staffService) buildQuery(param StaffQueryParam) repository.GetStaffQuery { ... }
```

**Decision guide:**

| Situation | Solution |
|---|---|
| Short logic (few lines) | Inline in the public method |
| Logic spans multiple entities | Extract to `domain/service/` struct method |
| Single entity behavior | Add method to `domain/model/` |
| Repository query building | Inline (do not create helper) |
| Conversion function (marshaller) | Inline or use struct receiver method |

**Exception**: Pure utility functions with no domain knowledge belong in `pkg/` packages (e.g., `pkg/id`, `pkg/email`).

---

### 32. Adding comments or docstrings to unchanged code

AI often adds `// CreateStaff creates a staff` style comments to functions it touches. Remove these unless the logic is genuinely non-obvious.

---

### 33. Adding error handling for impossible cases

```go
// BAD - fmt.Sprintf cannot fail
str := fmt.Sprintf("%s", name)
if str == "" {
    return errors.InternalErr.New()
}
```

---

### 34. Backwards-compatibility shims for removed code

```go
// BAD - renaming to _unused instead of deleting
func _unusedHelper() { ... }
// Or: // Removed: OldMethod was removed
```

Delete unused code entirely.

---

### 35. Feature flags or conditional code for "future use"

```go
// BAD
if featureEnabled {
    // new behavior
} else {
    // old behavior (never used)
}
```

Just implement the new behavior directly.

---

### 36. Direct domain model struct initialization in usecase (bypassing constructor)

Domain model structs must be created using their constructor functions in the usecase layer, not direct struct initialization.

```go
// BAD - direct struct initialization bypasses constructor guarantees
func (i *adminStaffInteractor) Create(ctx context.Context, param *input.AdminCreateStaff) (*model.Staff, error) {
    staff := &model.Staff{
        ID:          id.New(),
        TenantID:    param.TenantID,
        Role:        param.Role,
        DisplayName: param.DisplayName,
        // Missing fields risk silent bugs when new fields are added
    }
    // ...
}

// GOOD - use domain constructor
func (i *adminStaffInteractor) Create(ctx context.Context, param *input.AdminCreateStaff) (*model.Staff, error) {
    staff := model.NewStaff(param.TenantID, param.Role, param.AuthUID, param.DisplayName, param.ImagePath, param.Email, param.RequestTime)
    // ...
}
```

**Why**: Constructors encapsulate ID generation (`id.New()`), default values, field constraints, and ensure all required fields are set. Direct initialization bypasses these guarantees and breaks silently when new fields are added to the model.

---

### 37. Unnecessary nil/valid checks on guaranteed values

Do not add nil checks or `.Valid` checks on values that are guaranteed by preceding code.

```go
// BAD - OrFail: true guarantees non-nil; nil check is impossible
staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
    ID: null.StringFrom(param.StaffID),
    BaseGetOptions: repository.BaseGetOptions{OrFail: true},
})
if err != nil {
    return nil, err
}
if staff == nil {  // IMPOSSIBLE - OrFail returns error if not found
    return nil, errors.StaffNotFoundErr.New()
}

// BAD - nullable.TypeFrom guarantees Valid == true
role := nullable.TypeFrom(model.StaffRoleAdmin)
if role.Valid {  // ALWAYS true
    // ...
}

// GOOD - trust the guarantee
staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
    ID: null.StringFrom(param.StaffID),
    BaseGetOptions: repository.BaseGetOptions{OrFail: true},
})
if err != nil {
    return nil, err
}
// Use staff directly - guaranteed non-nil
```

**Why**: Unnecessary checks add noise, suggest distrust of established contracts, and mislead future readers about whether nil is actually possible.

---

### 38. Unnecessary intermediate variable declarations

Do not assign to a variable if the value is immediately returned without modification.

```go
// BAD - unnecessary intermediate variable
func StaffToPb(m *model.Staff) *pb.Staff {
    result := &pb.Staff{
        Id:   m.ID,
        Name: m.DisplayName,
    }
    return result
}

// GOOD - return directly
func StaffToPb(m *model.Staff) *pb.Staff {
    return &pb.Staff{
        Id:   m.ID,
        Name: m.DisplayName,
    }
}
```

**Exceptions**: Intermediate variables are acceptable when:
- The value is used more than once
- Conditional mutation is needed before return (e.g., nullable timestamp prep)
- Readability clearly benefits from a named variable

---

### 39. Field-by-field struct construction in conversion functions

Conversion/marshaller functions must use struct literal initialization with all fields listed, not field-by-field assignment on an empty struct.

```go
// BAD - field-by-field assignment on empty struct
func ExampleToModel(e *dbmodel.Example) *model.Example {
    m := &model.Example{}
    m.ID = e.ID
    m.TenantID = e.TenantID
    m.Name = e.Name
    m.Status = model.ExampleStatus(e.Status)
    m.CreatedAt = e.CreatedAt
    m.UpdatedAt = e.UpdatedAt
    return m
}

// GOOD - struct literal with all fields
func ExampleToModel(e *dbmodel.Example) *model.Example {
    return &model.Example{
        ID:                e.ID,
        TenantID:          e.TenantID,
        Name:              e.Name,
        Status:            model.ExampleStatus(e.Status),
        CreatedAt:         e.CreatedAt,
        UpdatedAt:         e.UpdatedAt,
        ReadonlyReference: nil,
    }
}
```

**Exception**: When conditional field assignment is required (e.g., `ReadonlyReference` populated only when `e.R != nil`, nullable timestamps), prepare variables first, then use a single struct literal with those variables.

**Fix**: Replace with struct literal return. If `ReadonlyReference` needs conditional population, use the `var` + struct literal pattern from `repository.md`.

---

### 40. Manual nil-pointer conversion when `.Ptr()` exists

For `null.String`, `null.Int64`, `null.Time` (from `github.com/aarondl/null/v8`), use `.Ptr()` directly instead of a manual `if .Valid` block.

```go
// BAD - manual block
var paymentMethodID *string
if m.PaymentMethodID.Valid {
    paymentMethodID = &m.PaymentMethodID.String
}
result := &pb.Invoice{
    PaymentMethodId: paymentMethodID,
}

// GOOD - inline .Ptr()
result := &pb.Invoice{
    PaymentMethodId: m.PaymentMethodID.Ptr(),
}
```

**Why**: `null.String.Ptr()` returns `*string` (nil when `!Valid`). The manual block is dead weight and visually obscures the field mapping.

**Fix**: Replace `var x *T; if .Valid { x = &.Field }` with `.Ptr()` inline in the struct literal.

---

### 41. Redundant `var _ Interface = (*impl)(nil)` assertion

```go
// BAD - redundant assertion, AI-generated boilerplate
var _ StaffInvoiceInteractor = (*staffInvoiceInteractor)(nil)

// GOOD - constructor return type already enforces the contract
func NewStaffInvoiceInteractor(...) StaffInvoiceInteractor {
    return &staffInvoiceInteractor{...}
}
```

**Why**: When a constructor returns `InterfaceType`, the compiler already rejects any implementation that doesn't satisfy the interface. The `var _ Interface = (*impl)(nil)` is redundant and adds noise.

**Fix**: Delete the assertion line entirely.

---

### 42. Ad-hoc per-test fixture helper instead of `factory.NewFactory()`

```go
// BAD - ad-hoc per-test helper duplicates factory responsibilities
func newTestInvoice() *model.Invoice {
    return &model.Invoice{
        ID:             "test-invoice-id",
        OrganizationID: "test-org-id",
        Status:         model.InvoiceStatusDraft,
        // ...
    }
}

// In tests:
invoice := newTestInvoice()

// GOOD - use factory for deterministic, consistent test data
testdata := factory.NewFactory()
invoice := testdata.Invoice
organizationID := invoice.OrganizationID
```

**Why**: `factory.NewFactory()` provides deterministic go-faker fixtures with consistent related entity references. Ad-hoc helpers diverge from factory data, cause inconsistency between tests, and duplicate responsibilities.

**Fix**: Delete the helper function. Use `factory.NewFactory()` and access the entity via `testdata.Invoice` (or whichever entity field). If the factory doesn't have the entity yet, add it to `factory.go` following the existing pattern.

---

### 43. Conditional preload of fully-owned child entities

```go
// BAD - preloadStripe bool flag on owned children
func (r *invoice) buildPreload(preloadStripe bool) []qm.QueryMod {
    mods := []qm.QueryMod{
        qm.Load(dbmodel.InvoiceRels.InvoiceItems),
    }
    if preloadStripe {
        mods = append(mods, qm.Load(dbmodel.InvoiceRels.InvoiceStripe))
    }
    return mods
}

// GOOD - always load owned children unconditionally
func (r *invoice) buildPreload() []qm.QueryMod {
    return []qm.QueryMod{
        qm.Load(dbmodel.InvoiceRels.InvoiceItems),
        qm.Load(dbmodel.InvoiceRels.InvoiceStripe),
        qm.Load(fmt.Sprintf("%s.%s", dbmodel.InvoiceRels.InvoiceItems, dbmodel.InvoiceItemRels.InvoiceItemStripe)),
    }
}
```

**Why**: Fully-owned 1:1/1:N children (entities that cannot exist without their parent, composed relationship) must always be loaded — they are part of the parent's aggregate. Optional `preload` flags are only appropriate for *reference* relations (e.g., `ReadonlyReference`). Conditional loading of owned children leads to incomplete domain models and nil-dereference bugs.

---

### 44. pkg layer importing domain layer

```go
// BAD - internal/pkg/phone/phone.go imports domain/errors
import (
    "github.com/abyssparanoia/rapid-go/internal/domain/errors"
    "github.com/nyaruka/phonenumbers"
)

func ParseE164(rawNumber string, countryCode string) (*phonenumbers.PhoneNumber, error) {
    num, err := phonenumbers.Parse(rawNumber, countryCode)
    if err != nil {
        return nil, errors.RequestInvalidArgumentErr.Wrap(err)  // domain dependency in pkg
    }
    ...
}

// GOOD - internal/pkg/phone/phone.go uses fmt.Errorf only
import (
    "fmt"
    "github.com/nyaruka/phonenumbers"
)

func ParseE164(rawNumber string, countryCode string) (*phonenumbers.PhoneNumber, error) {
    num, err := phonenumbers.Parse(rawNumber, countryCode)
    if err != nil {
        return nil, fmt.Errorf("invalid phone number format: %w", err)
    }
    ...
}

// Caller in handler/usecase wraps with domain error
phoneNumber, err := phone.ParseE164(req.GetPhoneNumber(), "")
if err != nil {
    return nil, errors.RequestInvalidArgumentErr.Wrap(err)  // wrapping happens at call site
}
```

**Why**: `internal/pkg` is a shared utility layer (logging, IDs, phone parsing, etc.) that must have no dependencies on `domain`, `usecase`, or `infrastructure`. It is used by all layers including domain itself. If `pkg` imported `domain/errors`, a circular import would result when domain tried to use pkg utilities. Domain error wrapping is the responsibility of the caller (handler or usecase layer), not the utility function.

**Fix**: Remove the bool parameter. Always load all owned children in `buildPreload()`.

---

### 45. `Set*` method that overwrites non-empty state (should be `Update` or specific verb)

`Set` is reserved for **fill-empty** semantics: the field was guaranteed empty before this call (nullable/derived field set for the first time). Using `Set` for an overwrite that replaces a possibly-meaningful existing value is misleading.

```go
// BAD - PaymentMethodID may already hold a value; calling this "Set" implies it was empty
func (m *Invoice) SetPaymentMethodID(id string, t time.Time) {
    m.PaymentMethodID = null.StringFrom(id)
    m.UpdatedAt = t
}

// BAD - IsDefault is an existing bool being toggled; "Set" implies first-time assignment
func (m *PaymentMethod) SetDefault(isDefault bool, t time.Time) {
    m.IsDefault = isDefault
    m.UpdatedAt = t
}

// GOOD - "Update" signals that an existing value is being replaced
func (m *Invoice) UpdatePaymentMethodID(id string, t time.Time) *Invoice {
    m.PaymentMethodID = null.StringFrom(id)
    m.UpdatedAt = t
    return m
}

func (m *PaymentMethod) UpdateIsDefault(isDefault bool, t time.Time) *PaymentMethod {
    m.IsDefault = isDefault
    m.UpdatedAt = t
    return m
}
```

**Rule of thumb**: if the field was guaranteed empty before the call → `Set`; if you are replacing a possibly-meaningful existing value → `Update`; if the change has business meaning or a state guard → use a specific verb (`Finalize`, `MarkPaid`, `Void`, `Delete`…).

**Keep as `Set`**: `SetStripe` family (Invoice, InvoiceItem, InvoiceTax, Payment, Refund, Organization), `Staff.SetImageURL` — these genuinely fill a derived/nullable field that starts empty.

**Fix**: Rename `Set{X}` → `Update{X}` (or a domain-specific verb) when the field may already hold a value. Update all call sites.

---

### 46. Void mutator on a domain model (should return receiver)

All mutating methods on domain models must return the receiver (`*Entity`), and `(*Entity, error)` when they also validate a precondition. Void mutators (`func (m *Entity) Mutate(...)`) are forbidden.

```go
// BAD - void return; caller cannot chain and shape is inconsistent
func (m *Invoice) ApplyTotals(subtotal int64, tax int64) {
    m.Subtotal = subtotal
    m.TaxAmount = tax
    m.TotalAmount = subtotal + tax
}

func (m *Organization) Update(name null.String, t time.Time) {
    if name.Valid { m.Name = name.String }
    m.UpdatedAt = t
}

// GOOD - return receiver enables chaining and consistent shape
func (m *Invoice) ApplyTotals(subtotal int64, tax int64) *Invoice {
    m.Subtotal = subtotal
    m.TaxAmount = tax
    m.TotalAmount = subtotal + tax
    return m
}

func (m *Organization) Update(name null.String, t time.Time) *Organization {
    if name.Valid { m.Name = name.String }
    m.UpdatedAt = t
    return m
}

// When validating a precondition, return (*Entity, error)
func (m *PaymentMethod) UpdateCreditCardDetails(...) (*PaymentMethod, error) {
    if m.PaymentMethodType != PaymentMethodTypeCard {
        return nil, errors.PaymentMethodTypeMismatchErr.New()
    }
    // mutate
    return m, nil
}
```

**Why**: Consistent return shape makes the codebase predictable, enables chaining, and distinguishes pure mutations from validated state transitions at a glance.

**Fix**: Add `return m` (or `return m, nil`) to every void domain model mutator. Change signature from `void` to `*Entity` (or `(*Entity, error)` if validation is present). Update call sites if they were discarding an `error` return that is now a `(*Entity, error)`.

---

### 47. Usecase branching on raw domain model fields instead of calling a predicate

Usecase and service code must not branch on a domain model's exported fields or status constants. Expose an intent-revealing predicate on the model and call it instead.

```go
// BAD - usecase inspects raw fields
if inv.Status != model.InvoiceStatusDraft {
    return nil, errors.InvoiceNotDraftErr.New()
}

// BAD - composite guard inline in usecase
if !invoice.IsStripe() || !invoice.IsDownloadable() || !invoice.Stripe.Valid {
    return "", errors.InvoiceHostedURLNotAvailableErr.New()
}

// BAD - idempotency check on raw fields in usecase
if inv.PaymentMethodID.Valid && inv.PaymentMethodID.String == pm.ID {
    return nil
}

// GOOD - intent-revealing predicates on the model
if !inv.IsDraft() {
    return nil, errors.InvoiceNotDraftErr.New()
}

if !invoice.IsHostedURLAvailable() {
    return "", errors.InvoiceHostedURLNotAvailableErr.New()
}

if inv.HasPaymentMethod(pm.ID) {
    return nil
}
```

**Common predicates to add**: `IsDraft() bool`, `IsPending() bool`, `HasPaymentMethod(id string) bool`, `IsHostedURLAvailable() bool`, `IsStripe() bool`, `IsDownloadable() bool`.

**Why**: Predicates encapsulate domain knowledge inside the model. If another caller forgets the guard, behavior diverges. Raw field inspection leaks domain semantics into the application layer and makes refactors error-prone.

**Fix**: Add an intent-revealing predicate method to the domain model. Replace the inline field check in usecase with a call to the predicate.

---

### 48. Precondition duplicated in usecase immediately before a model method that already checks it

If a model method validates a precondition internally and returns a domain error, the usecase must not check the same condition right before calling the method. Duplicating the guard risks the two checks diverging over time and misleads readers about where the authoritative check lives.

```go
// BAD - usecase duplicates Invoice.Finalize's internal Status != Draft guard
invoice, err := i.invoiceRepository.Get(ctx, ...)
if err != nil { return err }
if invoice.Status != model.InvoiceStatusDraft {         // duplicate of Finalize's guard
    return errors.InvoiceNotDraftErr.New()
}
invoice, err = invoice.Finalize(pm, t)                  // also checks Status == Draft
if err != nil { return err }

// GOOD - rely on the model method's own guard; delete the external duplicate
invoice, err := i.invoiceRepository.Get(ctx, ...)
if err != nil { return err }
invoice, err = invoice.Finalize(pm, t)                  // single authoritative check
if err != nil { return err }
```

**Exception — idempotency short-circuit**: A webhook handler that wants a silent `return nil` (not an error) on re-delivery may check a predicate *before* calling the model method, because the model method would return an error that the webhook cannot propagate. In this case, the predicate (e.g., `!inv.IsDraft()`) serves as an idempotency guard, not a duplicate validation:

```go
// OK - webhook idempotency short-circuit (silent nil, not an error)
if !inv.IsDraft() {
    logger.L(ctx).Info("invoice already finalized (idempotent)")
    return nil  // silent return, not an error
}
inv, err = inv.Finalize(pm, t)  // model's own guard remains the safety net
```

**Fix**: Delete the usecase/service check. Let the model method's own guard produce the authoritative error. If a silent return is needed (webhook idempotency), keep the predicate call but document it as an idempotency guard, not a validation.
