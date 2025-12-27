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
