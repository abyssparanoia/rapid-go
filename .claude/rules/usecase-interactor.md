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

## No Private Methods

Interactor implementations must **not** define private methods. Only public interface methods are implemented.

Private helpers scatter business logic inside the usecase layer and reduce testability. Move logic to the right place instead:

| Situation | Solution |
|-----------|----------|
| Short logic (a few lines) | Inline in the public method |
| Logic spans multiple entities | Extract to `domain/service/` |
| Single entity state transition | Add method to `domain/model/` |
| Repository query building | Inline (do not create a helper) |

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

### All input-derived validation belongs in the input layer — not in the interactor body

入力に対する検証は **input パッケージ内で完結** させ、interactor 側には漏らさない。役割は constructor と `Validate()` で分担する:

| Layer | 責務 |
|-------|------|
| `New{...}` constructor | **format パース** (string → time.Time など)。失敗時は `(nil, error)` を返す |
| `Validate()` | **タグベース検証** (`required`, `max`, `gte/lte`, `gtefield/ltefield`) + **タグでは表現できない検証** (相対距離・複雑な相互排他など) |
| Interactor | `param.Validate()` を呼ぶだけ。パース済み値 (`param.StartDate` 等) を直接使う |

**フィールド間比較は validator のタグ (`gtefield` / `ltefield` など) を優先**: time.Time 同士の比較も対応している。手動 if 文に書き下す前にタグで表現できないか検討する。「最大 N 日間」のような相対距離ガードはタグでは表現できないので Validate() 内で手動チェックする。

**パース済みフィールドのみ struct に持つ**: 入力日付文字列を constructor で `time.Time` にパースしたら、パース済みフィールド (`StartDate time.Time` 等) のみ struct に保持し、元の string フィールドは残さない。パース後に元文字列を参照しない限り string を併存させる理由はなく、`StartDate string` (元文字列) と `StartDay time.Time` (パース済み) を両持ちすると Date / Day の命名も曖昧になる。**例外**: パース後も元の入力文字列そのものを後続処理で使う必要がある場合に限り、string フィールドも併せて保持してよい。

`Validate()` で parse もしてしまうと「parse 失敗時のエラーが pre-condition チェックに混ざる」「Validate 前後でパース済みフィールドが zero value か否か変わる中間状態になる」ため、parse は constructor に寄せる。

#### Pattern: parse in constructor, range-check in Validate

```go
type TenantListVehicleLocationHistories struct {
    TenantID    string    `validate:"required"`
    RequestTime time.Time `validate:"required"`

    // StartDate / EndDate は入力日付文字列を constructor が JST 0:00 にパースした値
    // (両端含む期間)。パース失敗時は constructor がエラーを返すため、ここに到達する
    // 時点で必ず有効値。元の入力文字列はパース後参照しないため struct には保持しない。
    StartDate time.Time
    // gtefield=StartDate で start <= end をタグレベルで担保 (時系の time.Time 比較に対応)
    EndDate time.Time `validate:"gtefield=StartDate"`
}

// constructor は string → time.Time のパースに責任を持ち、失敗時はエラーを返す。
// 引数 startDate/endDate (string) はパース後は使わないため struct に保持しない。
func NewTenantListVehicleLocationHistories(
    tenantID, startDate, endDate string,
    requestTime time.Time,
) (*TenantListVehicleLocationHistories, error) {
    parsedStart, err := time.ParseInLocation(time.DateOnly, startDate, now.JST)
    if err != nil {
        return nil, errors.RequestInvalidArgumentErr.Wrap(err)
    }
    parsedEnd, err := time.ParseInLocation(time.DateOnly, endDate, now.JST)
    if err != nil {
        return nil, errors.RequestInvalidArgumentErr.Wrap(err)
    }
    return &TenantListVehicleLocationHistories{
        TenantID:    tenantID,
        RequestTime: requestTime,
        StartDate:   parsedStart,
        EndDate:     parsedEnd,
    }, nil
}

// Validate() は範囲・整合性チェックのみ。parse は前提として済んでいる。
// `EndDate >= StartDate` はタグ (gtefield) で担保済みなので、ここでは「相対距離が 7 日以内」のみチェックする。
func (p *TenantListVehicleLocationHistories) Validate() error {
    if err := validation.Validate(p); err != nil {
        return errors.RequestInvalidArgumentErr.Wrap(err)
    }
    maxEnd := p.StartDate.AddDate(0, 0, 6) // 7 日 = 両端含む
    if p.EndDate.After(maxEnd) {
        return errors.RequestInvalidArgumentErr.New().
            WithDetail("period must be within 7 days")
    }
    return nil
}

// Interactor は Validate() を呼ぶだけ。Validate 後は param.StartDate / EndDate を直接使う。
func (i *interactor) ListLocationHistories(ctx context.Context, param *input.TenantListVehicleLocationHistories) (...) {
    if err := param.Validate(); err != nil {
        return nil, err
    }
    rangeEnd := param.EndDate.AddDate(0, 0, 1)
    // ...
}

// Handler 側は constructor のエラーも透過する。
func (h *Handler) ListXxx(ctx context.Context, req *pb.Req) (*pb.Resp, error) {
    param, err := input.NewTenantListVehicleLocationHistories(...)
    if err != nil {
        return nil, err
    }
    got, err := h.interactor.ListLocationHistories(ctx, param)
    // ...
}
```

#### Anti-pattern: parsing / range-check in interactor body

```go
// BAD - 未パースの string を struct に残し、parse と範囲チェックが interactor 側に漏れている
func (i *interactor) ListLocationHistories(...) (..., error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }
    // param.StartDateStr / EndDateStr は未パースの string フィールド。
    // constructor でパースすべきものを interactor でやってしまっている。
    startDay, err := time.ParseInLocation(time.DateOnly, param.StartDateStr, now.JST)
    if err != nil { return nil, ... }
    endDay, err := time.ParseInLocation(time.DateOnly, param.EndDateStr, now.JST)
    if err != nil { return nil, ... }
    if endDay.Before(startDay) { return nil, ... }
    if endDay.After(maxEnd) { return nil, ... }
    // ... business logic
}
```

理由:
- 入力検証ロジックが分散すると、別の caller (CLI / worker / 直接テスト) から呼ばれた際に同じ検証が走らずバグになる
- Validate() の単一責任が崩れる
- テストで「validation エラー」を検証するときに interactor 全体をモックする必要が出る

#### Anti-pattern: parse を Validate() 側にまとめる

```go
// BAD - 未パースの string を struct に残し、Validate() で parse もやってしまう
func (p *TenantListVehicleLocationHistories) Validate() error {
    // ...
    parsedStart, err := time.ParseInLocation(time.DateOnly, p.StartDateStr, now.JST)
    if err != nil { return ... }
    // ...
    p.StartDate = parsedStart
    return nil
}
```

問題点:
- Validate 前後で `p.StartDate` の状態が変わる (zero value or parsed)。caller は Validate を呼んだか覚えておく必要がある
- parse エラーと「期間が 7 日超」のような pre-condition チェックエラーが同じ Validate() 内に混在する
- constructor 直後に `param.StartDate` を参照したら zero value、というバグを生みやすい

→ parse は constructor、Validate は範囲・整合性チェックという責務分離を守る。

#### Date format constants

時刻フォーマット文字列は **必ず標準ライブラリ定数 (`time.DateOnly`, `time.RFC3339`, `time.DateTime` など) を使う**。`"2006-01-02"` のようなマジックリテラルは禁止。

```go
// GOOD - constructor 内で引数 string をパース
parsed, err := time.ParseInLocation(time.DateOnly, startDate, now.JST)

// BAD - マジックリテラル
parsed, err := time.ParseInLocation("2006-01-02", startDate, now.JST)
```

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

**No exceptions for masters or singletons.** This applies to **every** resource-returning method —
including master/lookup `List`s (e.g. `LeaseProduct`, `Fabric`) and singleton `Get`s (e.g.
`LeasePricingSettings`) that have no asset field and no `ReadonlyReference` today. Each returned type
still needs its own `BatchSet{Entity}URLs` (a defensive no-op loop) on `service.Asset`, the query
still sets `Preload: true`, and the interactor still calls it. A singleton whose `Get` returns
`*model.X` is wrapped in its slice type (e.g. `model.LeasePricingSettingsList{settings}`) so the call
keeps the slice convention.

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

For optional update fields, use `nullable.Type[T]` instead of pointers — but **only for custom
types** (domain enums, `civil.Date`). Optional **primitive/time** input and `Update`-param fields
use `null/v8` (`null.Int64`, `null.Bool`, `null.Float64`, `null.Time`, `null.String`), never
`nullable.Type[int64/bool/float64/time.Time/string]`. In handlers, build these from optional proto
fields with the matching `null/v8` constructor (`null.StringFromPtr`, `null.Int64FromPtr`,
`null.BoolFromPtr`, `null.Float64FromPtr`) rather than ad-hoc `nullable` wrappers. See the type
table in `repository.md`.

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

## Asset Validation in Update Methods

When update methods accept optional asset fields (e.g., `ImageAssetID`), follow these critical patterns to avoid bugs:

### Critical Rules

1. **Validate assets INSIDE transaction** - Not before
2. **Use `var imagePath null.String`** - NOT `var imagePath string`
3. **Only set when provided** - `imagePath = null.StringFrom(path)` only in the `if` block

### Correct Pattern

```go
func (i *staffMeInteractor) Update(
    ctx context.Context,
    param *input.StaffUpdateMe,
) (*model.Staff, error) {
    if err := param.Validate(); err != nil {
        return nil, err
    }

    if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
        staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
            ID: null.StringFrom(param.StaffID),
            BaseGetOptions: repository.BaseGetOptions{
                OrFail:    true,
                ForUpdate: true,
            },
        })
        if err != nil {
            return err
        }

        // Validate asset if provided (INSIDE transaction)
        var imagePath null.String  // IMPORTANT: null.String, not string
        if param.ImageAssetID.Valid {
            authContext := model.NewStaffAssetAuthContext(param.StaffID)
            path, err := i.assetService.GetWithValidate(ctx, model.AssetTypeUserImage, param.ImageAssetID.String, authContext)
            if err != nil {
                return err
            }
            imagePath = null.StringFrom(path)  // Only set when asset provided
        }

        // Update with null.String (not converted)
        staff.Update(param.DisplayName, nullable.Type[model.StaffRole]{}, imagePath, param.RequestTime)

        return i.staffRepository.Update(ctx, staff)
    }); err != nil {
        return nil, err
    }

    // Return with relations
    return i.staffRepository.Get(ctx, repository.GetStaffQuery{
        ID: null.StringFrom(param.StaffID),
        BaseGetOptions: repository.BaseGetOptions{
            OrFail:  true,
            Preload: true,
        },
    })
}
```

### Anti-Pattern: Using `string` Instead of `null.String`

```go
// Bad - Using string causes bug
var imagePath string  // WRONG - This causes the bug!
if param.ImageAssetID.Valid {
    var err error
    imagePath, err = i.assetService.GetWithValidate(...)
    if err != nil {
        return nil, err
    }
}

// When ImageAssetID is NOT provided:
// - imagePath remains "" (empty string)
// - null.StringFrom("") creates {Valid: true, String: ""}
// - This OVERWRITES existing ImagePath with empty string!
staff.Update(param.DisplayName, nullable.Type[model.StaffRole]{}, null.StringFrom(imagePath), param.RequestTime)
```

### Why This Matters

| Approach | When Asset NOT Provided | Result |
|----------|-------------------------|--------|
| ✅ `var imagePath null.String` | `{Valid: false}` (zero value) | Field NOT updated (correct) |
| ❌ `var imagePath string` | `""` (empty string) | Field cleared to empty (BUG) |

**Key Insight**: Domain model's `Update()` method only updates fields where `null.String.Valid == true`. Using `null.String` zero value (`{Valid: false}`) correctly skips the update.

## Domain Method Usage (Domain Logic First)

**IMPORTANT**: All domain model operations (creation and state changes) MUST use domain model constructors and methods. Never directly initialize structs or assign fields in the usecase layer.

### Creation: Use Constructors

```go
// GOOD - use domain constructor
staff := model.NewStaff(param.TenantID, param.Role, param.AuthUID, param.DisplayName, param.ImagePath, param.Email, param.RequestTime)

// BAD - direct struct initialization in usecase
staff := &model.Staff{
    ID:          id.New(),
    TenantID:    param.TenantID,
    Role:        param.Role,
    CreatedAt:   param.RequestTime,
    UpdatedAt:   param.RequestTime,
}
```

### Updates: Use Domain Methods

```go
// GOOD
admin.UpdateRole(param.Role.Value(), param.RequestTime)

// BAD - direct field assignment
admin.Role = param.Role.Value()
admin.UpdatedAt = param.RequestTime
```

See `domain-model.md` for more details on domain method patterns.
