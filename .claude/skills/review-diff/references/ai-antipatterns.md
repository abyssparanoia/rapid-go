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

### 31. Package-level private functions

Package-level private functions (non-receiver functions) are prohibited in domain, usecase, and infrastructure layers. They pollute the package namespace and scatter logic.

```go
// BAD - package-level private function
func buildStaffQuery(id string) repository.GetStaffQuery {
    return repository.GetStaffQuery{...}
}

// BAD - package-level private helper
func toStaffModel(param *input.AdminCreateStaff, t time.Time) *model.Staff {
    return model.NewStaff(param.TenantID, param.Role, ...)
}

// BAD - package-level validation helper
func validateStaffInput(param *input.AdminCreateStaff) error {
    if param.Email == "" {
        return errors.RequestInvalidArgumentErr.New()
    }
    return nil
}

// GOOD - inline the logic
staff, err := i.repo.Get(ctx, repository.GetStaffQuery{...})

// GOOD - method on domain model
func (m *Staff) Validate() error { ... }

// GOOD - method on domain service struct
func (s *staffService) Create(ctx context.Context, param StaffCreateParam) (*model.Staff, error) { ... }
```

**Decision guide:**
| Situation | Solution |
|---|---|
| Short logic (few lines) | Inline in the calling method |
| Reusable business logic | Method on domain model or domain service |
| Data conversion | Method on marshaller struct, or inline |
| Validation | Method on input struct or domain model |

**Exception**: Package-level private functions are acceptable only when they are pure utility functions with no domain knowledge (e.g., generic type conversion helpers). Even then, prefer placing them in a shared `pkg/` utility package.

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

```go
// BAD - direct struct literal in usecase
func (i *adminStaffInteractor) Create(...) (*model.Staff, error) {
    staff := &model.Staff{
        ID:          id.New(),
        TenantID:    param.TenantID,
        Role:        param.Role,
        DisplayName: param.DisplayName,
        CreatedAt:   param.RequestTime,
        UpdatedAt:   param.RequestTime,
    }
    return staff, nil
}

// GOOD - use domain model constructor
func (i *adminStaffInteractor) Create(...) (*model.Staff, error) {
    staff := model.NewStaff(
        param.TenantID,
        param.Role,
        param.AuthUID,
        param.DisplayName,
        param.ImagePath,
        param.Email,
        param.RequestTime,
    )
    return staff, nil
}
```

**Why**: Constructors encapsulate initialization logic (ID generation, default values, field constraints). Bypassing them risks missing required fields and duplicates logic across call sites.

**Note**: This is distinct from #7 (direct field assignment for updates). #7 is about mutation; this is about creation.

---

### 37. Unnecessary nil/valid checks on guaranteed values

```go
// BAD - staff is guaranteed non-nil by OrFail: true
staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
    ID:             null.StringFrom(param.StaffID),
    BaseGetOptions: repository.BaseGetOptions{OrFail: true},
})
if err != nil {
    return nil, err
}
if staff == nil {  // Unnecessary - OrFail guarantees non-nil on success
    return nil, errors.StaffNotFoundErr.New()
}

// GOOD - trust OrFail guarantee
staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
    ID:             null.StringFrom(param.StaffID),
    BaseGetOptions: repository.BaseGetOptions{OrFail: true},
})
if err != nil {
    return nil, err
}
// Use staff directly

// BAD - checking .Valid on a value just set with TypeFrom
role := nullable.TypeFrom(param.Role)
if role.Valid {  // Always true
    query.Role = role
}

// GOOD
query.Role = nullable.TypeFrom(param.Role)
```

**Fix**: Remove nil/valid checks when the preceding code guarantees the value. Trust `OrFail`, `TypeFrom`, `StringFrom`, and similar constructors.

---

### 38. Unnecessary intermediate variable declarations

```go
// BAD - unnecessary variable
result := someFunction(ctx, param)
return result

// GOOD
return someFunction(ctx, param)

// BAD - declaring then immediately returning
staffs, err := i.staffRepository.List(ctx, query)
if err != nil {
    return nil, err
}
result := &output.AdminListStaffs{
    Staffs:     staffs,
    TotalCount: totalCount,
}
return result, nil

// GOOD - return directly
staffs, err := i.staffRepository.List(ctx, query)
if err != nil {
    return nil, err
}
return &output.AdminListStaffs{
    Staffs:     staffs,
    TotalCount: totalCount,
}, nil
```

**Exception**: Variables are acceptable when (1) the value is used multiple times, (2) conditional mutation is needed (e.g., nullable timestamp fields in marshallers), or (3) readability significantly improves.

---

### 39. Field-by-field struct construction in conversion functions

```go
// BAD - empty struct then field-by-field assignment
func StaffToPB(m *model.Staff) *pb.Staff {
    result := &pb.Staff{}
    result.Id = m.ID
    result.DisplayName = m.DisplayName
    result.Email = m.Email
    result.Role = StaffRoleToPB(m.Role)
    result.CreatedAt = timestamppb.New(m.CreatedAt)
    result.UpdatedAt = timestamppb.New(m.UpdatedAt)
    return result
}

// GOOD - struct literal with all fields
func StaffToPB(m *model.Staff) *pb.Staff {
    return &pb.Staff{
        Id:          m.ID,
        DisplayName: m.DisplayName,
        Email:       m.Email,
        Role:        StaffRoleToPB(m.Role),
        CreatedAt:   timestamppb.New(m.CreatedAt),
        UpdatedAt:   timestamppb.New(m.UpdatedAt),
    }
}
```

**Why**: Field-by-field assignment makes it easy to miss fields silently. Struct literals cause compile errors when fields are added, catching omissions early.

**Exception**: When conditional field assignment is needed (e.g., `ReadonlyReference` or nullable timestamps), use `m := &XXX{...}` with conditional blocks, then `return m`. But all non-conditional fields must still be in the initial struct literal.
