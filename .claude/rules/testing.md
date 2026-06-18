---
description: Testing conventions and patterns
globs:
  - "**/*_test.go"
---

# Testing Guidelines

## Test File Organization

- Test files are co-located with source: `example.go` -> `example_test.go`
- Use same package name (not `_test` suffix)

## Table-Driven Tests

```go
func TestExampleInteractor_Create(t *testing.T) {
    type args struct {
        ctx   context.Context
        param *input.AdminCreateExample
    }
    type usecase struct {
        exampleRepository func(ctrl *gomock.Controller) repository.Example
        transactable      func(ctrl *gomock.Controller) repository.Transactable
    }
    type want struct {
        result *model.Example
        err    error
    }

    tests := map[string]func(t *testing.T) (args, usecase, want){
        "success": func(t *testing.T) (args, usecase, want) {
            // Setup test case
            return args{...}, usecase{...}, want{...}
        },
        "validation error": func(t *testing.T) (args, usecase, want) {
            return args{...}, usecase{...}, want{...}
        },
    }

    for name, setup := range tests {
        t.Run(name, func(t *testing.T) {
            args, uc, want := setup(t)
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            interactor := usecase.NewAdminExampleInteractor(
                uc.transactable(ctrl),
                uc.exampleRepository(ctrl),
            )

            got, err := interactor.Create(args.ctx, args.param)

            // Assertions
            if want.err != nil {
                assert.ErrorIs(t, err, want.err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, want.result, got)
        })
    }
}
```

## Usecase interactor tests (MUST)

Tests for `internal/usecase/*_impl.go` interactors are held to these hard requirements. The canonical
references are `admin_tenant_impl_test.go` and `admin_staff_impl_test.go` — match their shape.

1. **Table-driven only.** Every interactor method's test MUST be a single
   `Test{Receiver}_{Method}` using `tests := map[string]testcaseFunc{...}` + a
   `for name, tc := range tests { t.Run(name, ...) }` loop. Flat sequential tests (a `Test*` function
   that calls the interactor directly with no `map[string]testcaseFunc`) are **prohibited** — every case,
   including error cases, lives in the one table. (See ai-antipatterns **#3**.)
2. **Test data from the factory only.** All fixtures MUST come from `factory.NewFactory()`. When you need
   a new entity, or a **computed/derived value** (a `*Calculation`, box-faces set, etc.), **extend the
   factory** — never inline `&model.X{...}` literals and never add a package-level
   `func newThing(...)`/`func thingCalculation(...)` helper in the test file. (See ai-antipatterns
   **#6** and **#42**.) Use `factory.CloneValue` for an independent copy.
3. **No package-level funcs in test files.** Only `func Test*` is allowed at package level; any helper is
   a closure inside the `Test*` function (per the helper rule below). (See ai-antipatterns **#31**.)
4. **Case coverage.** Each method covers `invalid argument` (when the input has a validatable field),
   `not found` (for Get/Update/Delete that fetch first), and `success`.
5. **Parallel.** `t.Parallel()` at both the outer test and every inner `t.Run`. When a case needs a
   deterministic ID, call `id.Mock()` inside that case's closure (matches the canonical references).
6. **Exact mock matching.** Only `context` uses `gomock.Any()`; the `requestTime` argument of
   `BatchSet*URLs` may also use `gomock.Any()` per the documented pattern. Everything else is matched
   exactly. (See ai-antipatterns **#1**.)

```go
// BAD - flat sequential usecase test, inline literal fixture, package-level helper
func awningCalculation(t time.Time) *model.AwningEstimateCalculation { return &model.AwningEstimateCalculation{...} }

func TestAdminAwningEstimateInteractor_Calculate(t *testing.T) {
    t.Parallel()
    testdata := factory.NewFactory()
    calc := awningCalculation(testdata.RequestTime) // ad-hoc helper + inline literal
    // ... single case, no map[string]testcaseFunc
}

// GOOD - table-driven, factory fixture, all cases in one Test func
func TestAdminAwningEstimateInteractor_Calculate(t *testing.T) {
    t.Parallel()
    type want struct { calculation *model.AwningEstimateCalculation; err error }
    type testcase struct { param *input.AdminCalculateAwningEstimate; usecase AdminAwningEstimateInteractor; want want }
    type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase
    tests := map[string]testcaseFunc{
        "invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase { /* ... */ },
        "success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
            testdata := factory.NewFactory()
            calc := testdata.AwningEstimateCalculation // factory fixture, not a helper
            // ...
        },
    }
    for name, tc := range tests {
        t.Run(name, func(t *testing.T) { t.Parallel(); /* ... */ })
    }
}
```

## Test Case Structure

```go
type testcaseFunc func(t *testing.T) (args, usecase, want)

tests := map[string]testcaseFunc{
    "success case description": func(t *testing.T) (args, usecase, want) { ... },
    "error case description": func(t *testing.T) (args, usecase, want) { ... },
}
```

Key patterns:
- Use `map[string]testcaseFunc` (not slice)
- Function returns `(args, usecase/service, want)` tuple
- Descriptive test names as map keys

### 同一メソッドのテストケースは 1 関数のテーブルに集約する

1 メソッド (例: `repository.List`) のテストは複数の `Test*` 関数に分けず、`Test{Receiver}_{Method}` の table-driven テスト 1 つに集約する。エラーケース (zero-value rejection 等) も別関数化せず、`wantErr bool` のような field を使ってテーブルに混ぜる。

```go
// GOOD - エラーケースもテーブルに統合
tests := map[string]struct {
    query   ListXxxQuery
    wantLen int
    wantErr bool
}{
    "success": { query: ListXxxQuery{Limit: 10}, wantLen: 3 },
    "rejects zero Limit": { query: ListXxxQuery{Limit: 0}, wantErr: true },
}
for name, tc := range tests {
    t.Run(name, func(t *testing.T) {
        got, err := repo.List(ctx, tc.query)
        if tc.wantErr {
            require.Error(t, err)
            return
        }
        require.NoError(t, err)
        require.Len(t, got, tc.wantLen)
    })
}

// BAD - エラーケースだけ別関数に分離
func TestRepo_List(t *testing.T) { /* 正常系 table */ }
func TestRepo_List_RejectsZeroLimit(t *testing.T) { /* zero limit only */ }
```

### テスト helper はパッケージレベルではなくテスト関数内の closure として定義する

このプロジェクトでは package-level 関数を避ける方針のため、テスト固有の helper は `Test*` 関数内で closure として定義する。複数の `Test*` 関数で共有が必要になった時点で初めてパッケージレベルに昇格させる。

```go
// GOOD - Test 関数内で closure として定義
func TestRepo_List(t *testing.T) {
    assertNoUnreliable := func(t *testing.T, logs model.HWBotLocations) {
        t.Helper()
        for _, l := range logs { /* ... */ }
    }
    // ...
}

// BAD - package-level の helper
func assertNoUnreliable(t *testing.T, logs model.HWBotLocations) {
    t.Helper()
    // ...
}

func TestRepo_List(t *testing.T) {
    // ... 1 箇所からしか呼ばれない
    assertNoUnreliable(t, got)
}
```

## Mock Setup

```go
type usecase struct {
    exampleRepository func(ctrl *gomock.Controller) repository.Example
    transactable      func(ctrl *gomock.Controller) repository.Transactable
}

// In test case setup
func(t *testing.T) (args, usecase, want) {
    return args{
        ctx:   context.Background(),
        param: &input.AdminCreateExample{...},
    }, usecase{
        exampleRepository: func(ctrl *gomock.Controller) repository.Example {
            m := mock_repository.NewMockExample(ctrl)
            m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
            m.EXPECT().Get(gomock.Any(), gomock.Any()).Return(expectedExample, nil)
            return m
        },
        transactable: func(ctrl *gomock.Controller) repository.Transactable {
            m := mock_repository.NewMockTransactable(ctrl)
            m.EXPECT().RWTx(gomock.Any(), gomock.Any()).DoAndReturn(
                func(ctx context.Context, fn func(context.Context) error) error {
                    return fn(ctx)
                },
            )
            return m
        },
    }, want{
        result: expectedExample,
        err:    nil,
    }
}
```

## Strict Mock Expectations Pattern (gomock.Any() Rules)

**CRITICAL RULE**: `gomock.Any()` is ONLY allowed for `context.Context` parameters (typically the first parameter). All other parameters MUST use exact object matching.

### Reference Pattern

See `internal/usecase/admin_tenant_impl_test.go` as the canonical correct example for this pattern.

### Correct Patterns

```go
// Repository Get - Exact query struct
mockTenantRepo.EXPECT().
    Get(gomock.Any(),  // ONLY context uses gomock.Any()
        repository.GetTenantQuery{  // EXACT struct matching
            ID: null.StringFrom(tenant.ID),
            BaseGetOptions: repository.BaseGetOptions{
                OrFail:  true,
                Preload: true,
            },
        }).
    Return(tenant, nil)

// Repository List - Exact query struct with all fields
mockStaffRepo.EXPECT().
    List(
        gomock.Any(),  // ONLY context
        repository.ListStaffQuery{
            TenantID: null.StringFrom(tenant.ID),
            BaseListOptions: repository.BaseListOptions{
                Page:    null.Uint64From(2),
                Limit:   null.Uint64From(30),
                Preload: true,
            },
            SortKey: nullable.TypeFrom(model.StaffSortKeyCreatedAtDesc),
        }).
    Return(model.Staffs{staff}, nil)

// Domain Service - Exact param struct
mockStaffService.EXPECT().
    Create(gomock.Any(),
        service.StaffCreateParam{
            TenantID:    tenant.ID,
            Email:       staff.Email,
            Password:    "random1234",
            StaffRole:   staff.Role,
            DisplayName: staff.DisplayName,
            ImagePath:   staff.ImagePath,
            RequestTime: requestTime,
        }).
    Return(staff, nil)

// Asset Service - Exact model slice + gomock.Any() for requestTime
mockAssetService.EXPECT().
    BatchSetStaffURLs(gomock.Any(), model.Staffs{staff}, gomock.Any()).  // EXACT model slice, Any for time
    Return(nil)

// Repository Create/Update - Exact domain object
admin := model.NewAdmin(
    model.AdminRoleRoot,
    authUID,
    email,
    displayName,
    requestTime,
)
mockAdminRepo.EXPECT().
    Create(gomock.Any(), admin).  // Exact object, not gomock.Any()
    Return(nil)

// Authentication Repository - Exact claims object
claims := model.NewAdminClaims(
    authUID,
    email,
    null.StringFrom(admin.ID),
    nullable.TypeFrom(admin.Role),
)
mockAdminAuthRepo.EXPECT().
    StoreClaims(gomock.Any(), authUID, claims).  // Exact claims object
    Return(nil)
```

### Anti-Patterns (PROHIBITED)

```go
// Bad - Using gomock.Any() for non-context parameters
mockStaffRepo.EXPECT().
    Get(gomock.Any(), gomock.Any()).  // WRONG - second param should be exact
    Return(staff, nil)

mockStaffService.EXPECT().
    Create(gomock.Any(), gomock.Any()).  // WRONG - second param should be exact
    Return(staff, nil)

mockAssetService.EXPECT().
    BatchSetStaffURLs(gomock.Any(), gomock.Any(), gomock.Any()).  // WRONG - second param should be exact
    Return(nil)

// Bad - Using DoAndReturn to avoid exact matching
mockAdminRepo.EXPECT().
    Create(gomock.Any(), gomock.Any()).DoAndReturn(  // WRONG - use exact object instead
        func(ctx context.Context, admin *model.Admin) error {
            return nil
        },
    )
```

### Why This Pattern

1. **Explicit expectations**: Tests clearly show what data is expected
2. **Catches bugs early**: Wrong field values cause test failures
3. **Self-documenting**: Test expectations serve as documentation
4. **Prevents false positives**: Tests verify actual behavior, not just that methods were called

### Using Model Constructors

Use domain model constructors to create exact expected objects:

```go
// Use NewAdmin constructor
admin := model.NewAdmin(
    model.AdminRoleRoot,
    authUID,
    email,
    displayName,
    requestTime,
)
mockAdminRepo.EXPECT().
    Create(gomock.Any(), admin).
    Return(nil)

// Use NewAdminClaims constructor
claims := model.NewAdminClaims(
    authUID,
    email,
    null.StringFrom(admin.ID),
    nullable.TypeFrom(admin.Role),
)
mockAdminAuthRepo.EXPECT().
    StoreClaims(gomock.Any(), authUID, claims).
    Return(nil)

// Use NewStaffClaims constructor
claims := model.NewStaffClaims(
    authUID,
    email,
    null.StringFrom(tenant.ID),
    null.StringFrom(staff.ID),
    nullable.TypeFrom(model.StaffRoleAdmin),
)
```

### Test Data Fixtures

Use `factory.NewFactory()` and `id.Mock()` for deterministic test data:

```go
testdata := factory.NewFactory()
staff := testdata.Staff
mockID := id.Mock()
staff.ID = mockID

// Use exact IDs in expectations
mockStaffRepo.EXPECT().
    Get(gomock.Any(),
        repository.GetStaffQuery{
            ID: null.StringFrom(mockID),
            BaseGetOptions: repository.BaseGetOptions{
                OrFail:  true,
                Preload: true,
            },
        }).
    Return(staff, nil)
```

## Transactable Mock Pattern

```go
// For RWTx - execute the function
m.EXPECT().RWTx(gomock.Any(), gomock.Any()).DoAndReturn(
    func(ctx context.Context, fn func(context.Context) error) error {
        return fn(ctx)
    },
)

// For ROTx - same pattern
m.EXPECT().ROTx(gomock.Any(), gomock.Any()).DoAndReturn(
    func(ctx context.Context, fn func(context.Context) error) error {
        return fn(ctx)
    },
)
```

## Assertion Patterns

```go
// Error expected
if want.err != nil {
    assert.ErrorIs(t, err, want.err)
    return
}

// No error expected
assert.NoError(t, err)

// Compare results
assert.Equal(t, want.result, got)

// Partial comparison
assert.Equal(t, want.result.ID, got.ID)
assert.Equal(t, want.result.Name, got.Name)
```

## Mock Generation

Mocks are generated with `go.uber.org/mock/mockgen`:

```bash
make generate.mock
```

Generated files are in `mock/` subdirectory:
- `internal/domain/repository/mock/example.go`
- `internal/usecase/mock/admin_example.go`

## Test Helpers

### Creating Test Fixtures

```go
func newTestExample(t *testing.T) *model.Example {
    t.Helper()
    return &model.Example{
        ID:        "test-id",
        TenantID:  "tenant-id",
        Name:      "Test Example",
        Status:    model.ExampleStatusDraft,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}
```

### Context with Values

```go
func newTestContext(t *testing.T) context.Context {
    t.Helper()
    return context.Background()
}
```

### Avoid file-level globals for test fixtures — declare them inline in the test function

固定値 (時刻、テナント ID、車両 ID 等) はテスト関数のローカル変数として宣言する。`var` でファイル先頭に置くと、複数テストが暗黙的に共有してしまい、片方のテストの修正が別テストに波及するため。

ヘルパー関数が値を必要とするときは、グローバル変数で共有するのではなく **引数として渡す**:

```go
// GOOD - テスト関数のローカル変数として宣言、ヘルパには引数で渡す
func setupTenantStaffContext(
    tenantID string,
    role model.StaffRole,
    deviceGroupID string,
    permType model.DeviceGroupPermissionType,
    requestTime time.Time, // ← caller から渡す
) context.Context {
    ctx := context.Background()
    ctx = request_interceptor.SetRequestTime(ctx, requestTime)
    // ...
}

func TestTenantHandler_ListXxx(t *testing.T) {
    t.Parallel()
    requestTime := time.Date(2026, 4, 28, 10, 0, 0, 0, time.UTC) // テスト関数内で宣言
    ctx := setupTenantStaffContext(tenantID, role, "", "", requestTime)
    // EXPECT 側でも requestTime を直接使う
}

// BAD - グローバル変数として共有
var fixedHandlerRequestTime = time.Date(2026, 4, 28, 10, 0, 0, 0, time.UTC)

func setupTenantStaffContext(...) context.Context {
    ctx = request_interceptor.SetRequestTime(ctx, fixedHandlerRequestTime) // 暗黙参照
    // ...
}
```

## Running Tests

```bash
# Run all tests
make test

# Run specific package
go test ./internal/usecase/...

# Run with verbose output
go test -v ./internal/usecase/...

# Run specific test
go test -run TestExampleInteractor_Create ./internal/usecase/...
```

## Repository Tests

This project requires tests for repository implementations. Repository tests verify the actual database queries work correctly.

### Repository Test Location

Location: `internal/infrastructure/{mysql,postgresql,spanner}/repository/{entity}_test.go`

### Repository Test Pattern

```go
func TestExampleRepository_Get(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()

    repo := repository.NewExample()

    tests := map[string]func(t *testing.T){
        "success": func(t *testing.T) {
            // Insert test data
            // Call repository method
            // Assert results
        },
        "not found with OrFail": func(t *testing.T) {
            // Test OrFail behavior
        },
        "not found without OrFail returns nil": func(t *testing.T) {
            // Test that OrFail=false returns nil instead of error
        },
    }

    for name, tc := range tests {
        t.Run(name, tc)
    }
}
```

### Required Test Coverage for Repositories

When adding a new repository, include tests for:

- **Get**: Various query conditions, OrFail behavior, ForUpdate behavior
- **List**: Pagination, filters, sorting
- **Create**: Successful persistence, constraint violations
- **Update**: Field changes, optimistic locking (if applicable)
- **Delete**: Successful removal, cascade behavior
- **BatchGet**: Multiple ID retrieval

### Why Repository Tests Matter

- Verify SQLBoiler query builders produce correct SQL
- Catch marshaller bugs between DB and domain models
- Ensure query conditions (WHERE clauses) work as expected
- Validate pagination and sorting logic
