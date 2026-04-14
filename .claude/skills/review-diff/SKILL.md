---
name: review-diff
description: Review and auto-fix code changes against project conventions by diffing the current branch against the default branch (main/master). Use when: (1) asked to review current changes, (2) running '/review-diff', (3) after completing a feature implementation and wanting to catch issues before creating a PR. Reviews all changed files against .claude/rules/, detects common AI coding mistakes (gomock.Any() misuse, missing ForUpdate, missing Preload, direct field assignment, wrong method ordering, etc.), and automatically fixes all issues found. Does NOT require an existing PR.
---

# Diff Review & Auto-Fix

Review current branch changes against main/master and **automatically fix all issues found**.

## Workflow

```
1. Detect Default Branch  → find main or master
2. Gather Changes         → git diff, changed files
3. AI Anti-Pattern Check  → detect common AI mistakes
4. Rule-Based Review      → apply rule checklists per file category
5. Auto-Fix Issues        → edit files to fix all found problems
6. Verify Fixes           → lint, tests, code generation
7. Report                 → summary of what was found and fixed
```

## Step 1: Detect Default Branch

```bash
# Try to detect
git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's|refs/remotes/origin/||'
```

If that fails, try `origin/main` first, then `origin/master`.

## Step 2: Gather Changes

```bash
BASE=origin/main  # or origin/master

# Changed files
git diff --name-only $BASE...HEAD

# Full diff for review
git diff $BASE...HEAD

# Commit history
git log $BASE...HEAD --oneline
```

## Step 3: AI Anti-Pattern Check

Read `references/ai-antipatterns.md`. Check **every changed file** against all patterns listed.

**Priority patterns — must detect in every review:**
- **#1** `gomock.Any()` for non-context parameters (no exceptions)
- **#6** Direct model initialization instead of using `factory.NewFactory()` in tests
- **#12** Fully-owned entity placed in `ReadonlyReference` instead of direct field
- **#25** Private method defined in usecase interactor (receiver method)
- **#31** Package-level private functions in domain/usecase/infrastructure layers
- **#36** Direct domain model struct initialization in usecase (bypassing constructor)
- **#37** Unnecessary nil/valid checks on guaranteed values
- **#39** Field-by-field struct construction in conversion functions

Record each finding: file path, line number, pattern violated, fix to apply.

## Step 4: Rule-Based Review

Map each changed file to its rule category, then read `references/checklists.md` and apply relevant checklists.

| File Pattern | Category |
|---|---|
| `internal/domain/model/**` | domain-model |
| `internal/domain/service/**` | domain-service |
| `internal/domain/errors/**` | domain-errors |
| `internal/domain/repository/*.go` | repository-interface |
| `internal/infrastructure/**/repository/**` | repository-impl |
| `internal/infrastructure/**/marshaller/**` | marshaller |
| `internal/usecase/**` (non-test) | usecase |
| `internal/usecase/input/**` | input |
| `internal/infrastructure/grpc/**/handler/**` | grpc-handler |
| `internal/infrastructure/dependency/**` | dependency |
| `schema/proto/**` | proto |
| `db/**/migrations/**` | migration |
| `**/*_test.go` | testing |
| `*invitation*` | invitation-workflow |
| `*authentication*`, `cognito/**`, `firebase/**` | external-service |

## Step 5: Auto-Fix Issues

**Fix all issues immediately—do not just report.** Edit files using available tools.

Fix in this priority order:
1. **Correctness** – missing `ForUpdate`, IdP sync outside TX, wrong type for nullable fields
2. **Rule violations** – method ordering, naming, missing required methods/fields
3. **AI anti-patterns** – `gomock.Any()` misuse, direct field assignment, unnecessary abstractions
4. **Test completeness** – missing test cases, wrong mock patterns, missing `t.Parallel()`

When fixing test files that use `gomock.Any()` incorrectly, replace with exact expected values using domain model constructors.

## Step 6: Verify Fixes

Always run after fixing:

```bash
/usr/bin/make lint.go
```

Run if test files were changed or fixed:

```bash
/usr/bin/make test
```

Run code generation if applicable:

```bash
# Migration files changed
/usr/bin/make migrate.up

# Proto files changed
/usr/bin/make generate.buf

# Repository interfaces changed
/usr/bin/make generate.mock
```

## Step 7: Report

```markdown
## Diff Review Summary

### Scope
- Base: origin/main (or master)
- Changed files: N files across X categories

### Auto-Fixed Issues (N)
1. **[Category]** Short description
   - `path/to/file.go:123`
   - Was: `<bad code>`
   - Fixed: `<good code>`

### Issues Requiring Manual Action (N)
1. **[Category]** Short description
   - `path/to/file.go`
   - Reason: Requires `make generate.buf` / `make generate.mock`

### Verification
- [x] lint.go passes
- [x] tests pass
- [ ] needs `make generate.buf`

### Result
✅ Ready / ⚠️ N manual actions needed
```
