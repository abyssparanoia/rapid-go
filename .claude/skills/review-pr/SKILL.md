---
name: review-pr
description: Self-review PR changes against project conventions before creating PR. Use when: (1) before creating a pull request, (2) after completing feature implementation, (3) when asked to review changes, (4) running '/review-pr' or '/self-review'. Checks domain models, repositories, handlers, tests against project rules.
---

# PR Self-Review Guide

Review PR changes against project conventions to catch issues before creating a PR.

## Review Process

```
1. Gather Changes      → git diff, list changed files
2. Categorize Files    → map to rule files
3. Apply Checklists    → check against rules (see references/)
4. Cross-Cutting Check → lint, tests, DI, code generation
5. Generate Summary    → report issues and status
```

## Step 1: Gather Changes

```bash
# List changed files
git diff --name-only origin/master...HEAD

# Show detailed diff
git diff origin/master...HEAD

# Show commit history
git log origin/master...HEAD --oneline
```

## Step 2: File Category Mapping

| File Pattern | Rule File | Focus |
|--------------|-----------|-------|
| `internal/domain/model/**` | `domain-model.md` | Entity, constructor, state methods |
| `internal/domain/service/**` | `domain-service.md` | Param/Result, no TX |
| `internal/domain/errors/**` | `domain-errors.md` | Error naming, codes |
| `internal/domain/repository/*.go` | `repository.md` | Interface, query structs |
| `internal/infrastructure/**/repository/**` | `repository.md` | Implementation, marshaller |
| `internal/usecase/**` | `usecase-interactor.md` | Input, TX, external sync |
| `internal/infrastructure/grpc/**/handler/**` | `grpc-handler.md` | Handler, marshaller |
| `internal/infrastructure/dependency/**` | `dependency-injection.md` | Registration |
| `schema/proto/**` | `proto-definition.md` | Naming, HTTP annotations |
| `db/**/migrations/**` | `migration.md` | Up/Down, constraints |
| `**/*_test.go` | `testing.md` | Table-driven, mocks |
| `*invitation*` | `invitation-workflow.md` | Status, expiration |
| `*authentication*`, `*cognito*` | `external-service-integration.md` | Claims sync |

## Step 3: Apply Checklists

Detailed checklists by category are in `references/checklists.md`.

Read the checklist file and apply relevant sections based on changed file categories.

## Step 4: Cross-Cutting Checks

### Run Verification Commands

```bash
/usr/bin/make lint.go   # Lint check
/usr/bin/make test      # Test check
```

### Code Generation Verification

- Migrations changed → Run `make migrate.up`
- Proto changed → Run `make generate.buf`
- Repository interfaces changed → Run `make generate.mock`

### Registration Verification

- New interactors registered in `dependency.go`
- New handlers added to gRPC server

## Step 5: Generate Summary

```markdown
## Self-Review Summary

### Files Reviewed
- `path/to/file.go` - Category: domain-model

### Issues Found

#### Critical (Must Fix)
1. **[Category]** Description
   - File: `path/to/file.go:123`
   - Rule: `rules/xxx.md`
   - Fix: Suggested fix

#### Warnings (Should Fix)
1. **[Category]** Description

### Checklist Status
- [x] Lint passes
- [x] Tests pass
- [ ] Mocks regenerated
- [x] DI registration complete

### Overall Assessment
Ready for PR / Needs fixes before PR
```

## Common Issues

See `references/common-issues.md` for frequently encountered problems.

## Related Skills

- **code-investigation** - Use before this skill to understand existing patterns
- **create-pull-request** - Use after this skill to create the PR
