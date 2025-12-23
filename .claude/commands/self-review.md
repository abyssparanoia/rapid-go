---
description: Self-review PR changes against project rules before creating a PR
---

# Self-Review Current Changes

## Context

- Current branch: !`git branch --show-current`
- Base branch: !`git symbolic-ref refs/remotes/origin/HEAD --short 2>/dev/null | sed 's|origin/||' || echo "main"`
- Changed files: !`git diff --name-only origin/main...HEAD 2>/dev/null || git diff --name-only HEAD~1`
- Uncommitted changes: !`git status --porcelain`

## Task

Perform a comprehensive self-review of all changes on the current branch against project rules and conventions.

### Step 1: Gather Changes

First, understand what has changed:

```bash
# List all changed files (compared to main)
git diff --name-only origin/main...HEAD

# Get detailed diff
git diff origin/main...HEAD

# Show commit history on this branch
git log origin/main...HEAD --oneline
```

Also check for any uncommitted changes that should be reviewed.

### Step 2: Categorize Files

Map each changed file to its relevant rule file:

| File Pattern | Rule File |
|--------------|-----------|
| `internal/domain/model/**` | `domain-model.md` |
| `internal/domain/service/**` | `domain-service.md` |
| `internal/domain/errors/**` | `domain-errors.md` |
| `internal/domain/repository/*.go` | `repository.md` |
| `internal/infrastructure/{mysql,postgresql,spanner}/repository/**` | `repository.md` |
| `internal/infrastructure/{mysql,postgresql,spanner}/internal/marshaller/**` | `repository.md` |
| `internal/usecase/**` | `usecase-interactor.md` |
| `internal/infrastructure/grpc/internal/handler/**` | `grpc-handler.md` |
| `internal/infrastructure/dependency/**` | `dependency-injection.md` |
| `schema/proto/**` | `proto-definition.md` |
| `db/{mysql,postgresql,spanner}/**` | `migration.md` |
| `**/*_test.go` | `testing.md` |
| `*invitation*` | `invitation-workflow.md` |
| `*authentication*`, `*cognito*`, `*firebase*` | `external-service-integration.md` |

### Step 3: Review Each Category

For each category of changed files:
1. Read the relevant rule file from `.claude/rules/`
2. Read the changed files
3. Check against the checklist from the rule file

Key checks per category:

**Domain Model** - State changes via domain methods, `ReadonlyReference`, enum `Unknown` value
**Repository** - `nullable.Type[T]` for optional enums, `GetContextExecutor(ctx)`, `OrFail` handling
**Usecase** - `param.Validate()` first, TX boundaries, IdP sync within TX
**Handler** - Optional field nil checks, marshaller field mapping completeness
**Proto** - `UNSPECIFIED = 0`, optional keyword, HTTP annotations
**Tests** - Table-driven with map, `(args, usecase, want)` tuple pattern

### Step 4: Cross-Cutting Checks

Verify:
- [ ] If migrations changed: `make migrate.up` was run (SQLBoiler regenerated)
- [ ] If proto changed: `make generate.buf` was run
- [ ] If repository interfaces changed: `make generate.mock` was run
- [ ] New interactors registered in `dependency.go`
- [ ] New handlers added to gRPC server

Run:

```bash
make lint.go
make test
```

### Step 5: Output Summary

Provide a review summary:

```markdown
## Self-Review Summary

### Files Reviewed
- `path/to/file.go` - Category: {category}

### Issues Found

#### Critical (Must Fix)
1. **[Category]** Description
   - File: `path/to/file.go:line`
   - Rule: See `rules/xxx.md`
   - Fix: Suggested fix

#### Warnings (Should Fix)
1. **[Category]** Description

#### Suggestions (Nice to Have)
1. Description

### Checklist Status
- [ ] Lint passes
- [ ] Tests pass
- [ ] Mocks regenerated (if needed)
- [ ] DI registration complete

### Overall Assessment
Ready for PR / Needs fixes before PR
```

## Reference

For detailed checklists, see `.claude/skills/review-pr.md`

