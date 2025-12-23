---
name: review-pr
description: Self-review PR changes against project rules and conventions. Run this before creating a PR to catch issues early. Reviews domain models, repositories, handlers, tests, and more.
---

# PR Self-Review Guide

This skill guides you through a comprehensive self-review of PR changes based on project rules and conventions.

## When to Use This Skill

- Before creating a PR
- After completing a feature implementation
- When asked to review current changes

## Review Process Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    1. Gather PR Changes                      │
│         (git diff, list changed files, understand scope)     │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              2. Categorize Changed Files                     │
│     (domain, repository, usecase, handler, proto, etc.)     │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│           3. Apply Rule-Based Review per Category            │
│        (check each file against relevant rule file)         │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              4. Cross-Cutting Concerns Check                 │
│      (DI registration, mock generation, lint, tests)        │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│               5. Generate Review Summary                     │
│           (issues found, suggestions, checklist)            │
└─────────────────────────────────────────────────────────────┘
```

## Step 1: Gather PR Changes

First, understand what has changed:

```bash
# Get list of changed files
git diff --name-only origin/main...HEAD

# Get detailed diff
git diff origin/main...HEAD

# Show commit history
git log origin/main...HEAD --oneline
```

## Step 2: File Category Mapping

Map changed files to their review rules:

| File Pattern | Rule File | Key Checks |
|--------------|-----------|------------|
| `internal/domain/model/**` | `domain-model.md` | Entity structure, constructors, state change methods |
| `internal/domain/service/**` | `domain-service.md` | Param/Result pattern, no TX management |
| `internal/domain/errors/**` | `domain-errors.md` | Error naming, codes |
| `internal/domain/repository/*.go` | `repository.md` | Interface definition, query structs |
| `internal/infrastructure/{mysql,postgresql,spanner}/repository/**` | `repository.md` | Implementation pattern, marshaller |
| `internal/infrastructure/{mysql,postgresql,spanner}/internal/marshaller/**` | `repository.md` | ToModel/ToDBModel conversion |
| `internal/usecase/**` | `usecase-interactor.md` | Input validation, TX boundaries, external sync |
| `internal/infrastructure/grpc/internal/handler/**` | `grpc-handler.md` | Handler pattern, marshaller |
| `internal/infrastructure/dependency/**` | `dependency-injection.md` | Registration order |
| `schema/proto/**` | `proto-definition.md` | Naming, HTTP annotations |
| `db/{mysql,postgresql,spanner}/migrations/**` | `migration.md` | Up/Down, constraints |
| `**/*_test.go` | `testing.md` | Table-driven tests, mock setup |
| `*invitation*` | `invitation-workflow.md` | Status transitions, expiration |
| `*authentication*`, `*cognito*` | `external-service-integration.md` | Claims sync, TX order |

## Step 3: Category-Specific Review Checklists

### Domain Model Review (`internal/domain/model/**`)

- [ ] Entity has `ReadonlyReference` struct pointer for relations
- [ ] Constructor uses `id.New()` for ID generation
- [ ] Constructor sets both `CreatedAt` and `UpdatedAt` to same time
- [ ] Update methods use `null.String`, `null.Int64` for optional fields
- [ ] Update methods always update `UpdatedAt`
- [ ] Status/Enum types have `Unknown` as first constant
- [ ] Status types have `String()` and `Valid()` methods
- [ ] State changes are done via domain methods (not direct field assignment)
- [ ] Role types have helper methods like `IsRoot()`, `IsNormal()`
- [ ] Slice types have `IDs()` and `MapByID()` helpers
- [ ] Type aliases defined: `{Entity}MapByID`, `{Entity}s`

### Repository Interface Review (`internal/domain/repository/**`)

- [ ] Has `//go:generate` directive for mockgen
- [ ] Query structs use `null.String`, `null.Uint64` for optional string/numeric fields
- [ ] Query structs use `nullable.Type[T]` for optional enum/custom type fields
- [ ] `BaseGetOptions`, `BaseBatchGetOptions`, `BaseListOptions` properly embedded

### Repository Implementation Review (`internal/infrastructure/{mysql,postgresql,spanner}/repository/**`)

- [ ] Uses `transactable.GetContextExecutor(ctx)` for all queries
- [ ] `Get` method handles `OrFail` correctly (nil vs error on not found)
- [ ] `Get` method handles `ForUpdate` option
- [ ] `List` method applies pagination with `Page` and `Limit`
- [ ] `List` method validates sort key with `query.SortKey.Valid && query.SortKey.Ptr().Valid()`
- [ ] Preload helper defined if relations exist

### Marshaller Review (`internal/infrastructure/{mysql,postgresql,spanner}/internal/marshaller/**`)

- [ ] `ToModel` handles relations via `R != nil` check
- [ ] `ToDBModel` sets `R: nil, L: struct{}{}`
- [ ] Enum conversions handle all cases including Unknown/default
- [ ] Slice conversion function defined (`{Entity}sToModel`)

### Usecase Interactor Review (`internal/usecase/**`)

- [ ] Interface has `//go:generate` directive
- [ ] Implementation uses dependency injection via constructor
- [ ] All methods start with `param.Validate()` check
- [ ] Write operations wrapped in `transactable.RWTx`
- [ ] Get before update uses `ForUpdate: true` for locking
- [ ] State changes use domain methods (not direct field assignment)
- [ ] IdP sync (StoreClaims/DeleteUser) happens within transaction
- [ ] On delete: IdP deletion before database deletion
- [ ] Final return fetches entity with `Preload: true` for fresh data

### Input Struct Review (`internal/usecase/input/**`)

- [ ] Named as `{Actor}{Action}{Resource}`
- [ ] Has `RequestTime` field with `validate:"required"`
- [ ] Has `Validate()` method that uses `validation.Validate()`
- [ ] Optional update fields use `nullable.Type[T]` (not pointers)
- [ ] Validation includes business rule checks for optional fields

### gRPC Handler Review (`internal/infrastructure/grpc/internal/handler/**`)

- [ ] Gets claims via `request_interceptor.Get{Actor}Claims(ctx)`
- [ ] Gets request time via `request_interceptor.GetRequestTime(ctx)`
- [ ] Converts proto to input struct correctly
- [ ] Handles optional proto fields with `if req.Field != nil` pattern
- [ ] Returns error directly (interceptor handles conversion)
- [ ] Uses marshaller for domain-to-proto conversion

### Handler Marshaller Review (`internal/infrastructure/grpc/internal/handler/**/marshaller/**`)

- [ ] Each resource has its own file (not combined)
- [ ] `ToPb` handles nil input
- [ ] `ToPb` uses variable declaration pattern for optional/nullable fields
- [ ] Enum conversions have both `ToPb` and `ToModel` directions
- [ ] Slice conversion function defined (`{Entity}sToPb`)
- [ ] All proto fields are explicitly mapped (check for omissions)

### Proto Definition Review (`schema/proto/**`)

- [ ] Enum values start with `{ENUM_NAME}_UNSPECIFIED = 0`
- [ ] Field names use snake_case
- [ ] Request/Response named as `{Action}{Resource}Request/Response`
- [ ] HTTP annotations follow REST patterns
- [ ] Optional fields marked with `optional` keyword
- [ ] List requests have `page` and `limit` fields

### Migration Review (`db/{mysql,postgresql,spanner}/migrations/**`)

- [ ] Has both `+goose Up` and `+goose Down` sections
- [ ] Column types match Go types (see type mapping in rule)
- [ ] Foreign keys named as `{table}_fkey_{column}`
- [ ] Indexes named as `{table}_idx_{column}`
- [ ] Unique constraints named as `{table}_uq_{columns}`
- [ ] `TIMESTAMPTZ` used for all timestamps (not TIMESTAMP)
- [ ] Constant tables have corresponding YAML in `db/{mysql,postgresql}/constants/`

### Test Review (`**/*_test.go`)

- [ ] Uses table-driven tests with `map[string]testcaseFunc`
- [ ] Test function returns `(args, usecase/service, want)` tuple
- [ ] Mock setup uses closure pattern with `func(ctrl *gomock.Controller)`
- [ ] Transactable mock uses `DoAndReturn` to execute function
- [ ] Error assertions use `assert.ErrorIs(t, err, want.err)`

### Dependency Injection Review (`internal/infrastructure/dependency/**`)

- [ ] New repository registered in Dependency struct
- [ ] New interactor registered in Dependency struct
- [ ] Constructor call added in `Inject()` method
- [ ] Handler updated to include new interactor
- [ ] Injection order follows: clients → transactable → repos → services → interactors → handlers

### External Service Integration Review

- [ ] Claims model uses `null.String` and `nullable.Type[T]` for optional fields
- [ ] `StoreClaims` called after entity creation/update
- [ ] `DeleteUser` called before database deletion
- [ ] All IdP operations within transaction boundary

## Step 4: Cross-Cutting Concerns

### Pre-Commit Checks

```bash
# Lint check
make lint.go

# Test check
make test
```

### Code Generation Verification

- [ ] If migrations changed: `make migrate.up` was run
- [ ] If proto changed: `make generate.buf` was run
- [ ] If repository interfaces changed: `make generate.mock` was run

### Registration Verification

- [ ] New interactors registered in `dependency.go`
- [ ] New handlers added to gRPC server registration

## Step 5: Review Summary Template

After completing the review, provide a summary:

```markdown
## Self-Review Summary

### Files Reviewed
- `path/to/file1.go` - Category: domain-model
- `path/to/file2.go` - Category: usecase

### Issues Found

#### Critical (Must Fix)
1. **[Category]** Description of critical issue
   - File: `path/to/file.go:123`
   - Rule: See `rules/xxx.md`
   - Fix: Suggested fix

#### Warnings (Should Fix)
1. **[Category]** Description of warning
   - File: `path/to/file.go:45`
   - Suggestion: ...

#### Suggestions (Nice to Have)
1. Description of suggestion

### Checklist Status
- [x] Lint passes
- [x] Tests pass
- [ ] Mocks regenerated (needed for repository interface changes)
- [x] DI registration complete

### Overall Assessment
Ready for PR / Needs fixes before PR
```

## Common Issues to Watch For

### Domain Layer
- Direct field assignment instead of domain methods
- Missing `UpdatedAt` update in modification methods
- Enum without `Unknown` value

### Repository Layer
- Missing `nullable.Type[T]` for optional enum fields (using pointers instead)
- Not checking `OrFail` for nil return vs error
- Missing `transactable.GetContextExecutor(ctx)`

### Usecase Layer
- Validation not called at method start
- Missing transaction wrapper for write operations
- IdP sync outside transaction
- Direct field assignment instead of domain methods

### Handler Layer
- Missing nil check for optional request fields
- Field mapping omissions in marshallers
- Not using variable declaration pattern for nullable fields

### Tests
- Missing mock for new dependencies
- Not using table-driven test pattern
- Hardcoded values instead of test fixtures

## Related Skills

- **code-investigation** - Use before this skill to understand existing patterns
- **create-pull-request** - Use after this skill to create the PR
