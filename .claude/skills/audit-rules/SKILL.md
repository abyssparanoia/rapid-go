---
name: audit-rules
description: Audit the entire codebase for compliance with .claude/rules/ conventions. Use when: (1) running '/audit-rules', (2) asked to check whether the whole codebase follows project conventions, (3) running a periodic convention sweep before a major release. Partitions source files by rule category, launches convention-reviewer and test-reviewer agents in parallel (audit mode), aggregates findings, auto-fixes all fixable violations, verifies with lint/test, and opens a PR. Accepts an optional category argument to scope to one partition (e.g. '/audit-rules domain').
---

# Whole-Codebase Rule-Compliance Audit

Audit every source file against `.claude/rules/` using **parallel reviewer agents in Audit Mode**, auto-fix all violations, then open a PR with the fixes.

## Workflow

```
1. Determine Scope      → parse optional arg, enumerate files with git ls-files
2. Partition Files      → split by rule category (domain, usecase, grpc, db-repo,
                          proto, migration, infra-other, tests)
3. Parallel Agent Audit → launch convention-reviewer per partition + test-reviewer
4. Aggregate Findings   → collect, deduplicate, sort
5. Auto-Fix Issues      → edit files to fix all auto-fixable violations
6. Verify Fixes         → codegen + lint/test retry loop
7. Create PR            → branch, commit, gh pr create
8. Report               → per-partition results table
```

## Step 1: Determine Scope

Check for an optional argument. If provided (e.g. `/audit-rules domain`), run only that
partition. Otherwise run all partitions.

Enumerate auditable files (exclude generated code):

```bash
# All tracked Go source files, excluding generated
git ls-files 'internal/**/*.go' \
  | grep -v '/mock/' \
  | grep -v '/dbmodel/' \
  | grep -v '\.pb\.go$' \
  | grep -v '_grpc\.pb\.go$'

# Proto files
git ls-files 'schema/proto/**/*.proto'

# Migration + constant files
git ls-files 'db/**/migrations/**/*.sql' 'db/**/constants/**/*.yaml'
```

## Step 2: Partition Files

Assign each file to exactly one partition using the table below. Files matching multiple
patterns go to the **first** matching partition.

Every non-test partition's `convention-reviewer` applies the **full** `.claude/rules/*.md` set, the
matching `checklists.md` section, AND the **whole `ai-antipatterns.md` catalog** (#7–#50 are non-test
conventions). The "Rule files" column below lists the *primary* rules for each partition; it is not an
exhaustive allow-list — every rule and anti-pattern applies wherever its layer matches.

| Partition | glob patterns | Reviewer | Rule files |
|-----------|--------------|----------|------------|
| `domain` | `internal/domain/**` (non-test) | convention-reviewer | domain-model.md, domain-errors.md, domain-service.md, repository.md, ai-antipatterns.md (non-test) |
| `usecase` | `internal/usecase/**` (non-test) | convention-reviewer | usecase-interactor.md, domain-service.md, ai-antipatterns.md (non-test, incl. #50) |
| `grpc` | `internal/infrastructure/grpc/internal/handler/**` (non-test) | convention-reviewer | grpc-handler.md, ai-antipatterns.md (non-test) |
| `db-repo` | `internal/infrastructure/{mysql,postgresql,spanner}/**` (non-test), `internal/infrastructure/dependency/**` | convention-reviewer | repository.md, dependency-injection.md, external-service-integration.md, ai-antipatterns.md (non-test) |
| `proto` | `schema/proto/**/*.proto` | convention-reviewer | proto-definition.md |
| `migration` | `db/**/migrations/**/*.sql`, `db/**/constants/**/*.yaml` | convention-reviewer | migration.md |
| `infra-other` | `internal/infrastructure/{cognito,firebase,gcs,s3,redis,http,aws,cmd}/**` (non-test), `cmd/**` (non-test) | convention-reviewer | external-service-integration.md, webhook-implementation.md, job-system.md, worker-pattern.md, cli-command-pattern.md, ai-antipatterns.md (non-test) |
| `tests` | `**/*_test.go` | **test-reviewer** | testing.md, ai-antipatterns.md (#1-#6) |

If a single partition has more than ~60 files, split it into sub-batches and run agents
sequentially for that partition (to avoid context overflow). Each sub-batch agent still uses
AUDIT MODE with its own file list.

## Step 3: Parallel Agent Audit

### CRITICAL: Use `subagent_type`, Do NOT Embed the Agent Definition

Specify `subagent_type: "convention-reviewer"` (or `"test-reviewer"`). Never read the agent
`.md` and paste its body into `prompt`. Doing so breaks the `model: sonnet` and read-only
`tools` enforcement, and causes drift.

### Audit Mode Prompt Template

Pass **only** context — the agent's role comes from its own definition:

```
AUDIT MODE

Partition: {partition_name}

Files to audit (review complete file content — this is NOT a diff):
{file list, one path per line}

Review every listed file in full against your role definition. Do NOT run git diff or
detect $BASE. Apply all applicable rule checks to the complete file content.
Report findings in your defined Output Format, including Semantic Category on every finding.
```

### Agent Launch Pattern

Launch all partitions in **one message**, all with `run_in_background: true`:

```
Agent(
  description: "Audit domain layer",
  subagent_type: "convention-reviewer",
  run_in_background: true,
  prompt: "AUDIT MODE\n\nPartition: domain\n\nFiles to audit:\n{P1 file list}"
)
Agent(
  description: "Audit usecase layer",
  subagent_type: "convention-reviewer",
  run_in_background: true,
  prompt: "AUDIT MODE\n\nPartition: usecase\n\nFiles to audit:\n{P2 file list}"
)
Agent(
  description: "Audit gRPC handlers",
  subagent_type: "convention-reviewer",
  run_in_background: true,
  prompt: "AUDIT MODE\n\nPartition: grpc\n\nFiles to audit:\n{P3 file list}"
)
Agent(
  description: "Audit DB repositories + dependency",
  subagent_type: "convention-reviewer",
  run_in_background: true,
  prompt: "AUDIT MODE\n\nPartition: db-repo\n\nFiles to audit:\n{P4 file list}"
)
Agent(
  description: "Audit proto definitions",
  subagent_type: "convention-reviewer",
  run_in_background: true,
  prompt: "AUDIT MODE\n\nPartition: proto\n\nFiles to audit:\n{P5 file list}"
)
Agent(
  description: "Audit migrations + constants",
  subagent_type: "convention-reviewer",
  run_in_background: true,
  prompt: "AUDIT MODE\n\nPartition: migration\n\nFiles to audit:\n{P6 file list}"
)
Agent(
  description: "Audit infra-other (cognito/firebase/gcs/cmd/...)",
  subagent_type: "convention-reviewer",
  run_in_background: true,
  prompt: "AUDIT MODE\n\nPartition: infra-other\n\nFiles to audit:\n{P7 file list}"
)
Agent(
  description: "Audit test files",
  subagent_type: "test-reviewer",
  run_in_background: true,
  prompt: "AUDIT MODE\n\nFiles to audit (all *_test.go):\n{P-test file list}"
)
```

Wait for all agents to complete before proceeding.

## Step 4: Aggregate Findings

1. **Collect** all findings from agent results
2. **Deduplicate** by `(file, line, semantic_category)` — same rationale as review-diff:
   agents may surface the same violation under different `rule_id` prefixes
3. **Sort** by severity: `error` → `warning` → `info`
4. **Separate** auto-fixable from manual-only findings

### Tie-breaking (same dedup key from multiple agents)

1. Highest severity wins
2. Most specific Fix block wins (non-empty over empty)
3. Owner preference: convention-reviewer > test-reviewer

### Finding Categories Reference

| Rule Prefix | Owner Agent |
|-------------|-------------|
| `CL-*` | convention-reviewer |
| `TEST-*`, `AP-1`–`AP-6` | test-reviewer |

## Step 5: Auto-Fix Issues

**Fix all auto-fixable findings immediately.** Edit files using Edit/Write tools.

Priority order (same as review-diff):

1. **Correctness** (CL-* that break runtime behavior: missing ForUpdate, Preload, state transition via wrong method)
2. **Convention** (CL-*: method ordering, naming, pattern compliance, null.v8/now.Now(), ReadonlyReference, and the **return pattern (#50)** — a resource returned without `Preload: true` + `BatchSet{Entity}URLs` is an auto-fixable convention finding, even for master/lookup lists and singleton gets; add the defensive no-op `BatchSet` to `service.Asset` if it is missing)
3. **Anti-patterns** (AP-1–AP-6 in tests: gomock.Any() misuse, missing t.Parallel(), table-driven pattern)
4. **Test structure** (TEST-*/AP-3/AP-6/AP-42: converting a flat sequential test to table-driven `map[string]testcaseFunc`, and moving inline-literal / package-level-helper fixtures into `factory.NewFactory()`). This is **mechanical and IS auto-fixable** — do it, do not defer it.
5. **Test quality** (TEST-*: weak assertions, and adding obvious `invalid argument`/`not found`/`success` cases modeled on an existing case)

Do **not** auto-fix findings that require spec interpretation or cross-file architectural
decisions — list these as manual actions in the report. Note that "not table-driven" and
"ad-hoc fixture instead of factory" are **NOT** in that deferred bucket: they are mechanical
conversions (priority 4 above) and must be fixed, not deferred as "coverage gaps."

## Step 6: Verify Fixes (with Retry Loop)

### 6a. Code Generation (if applicable)

Run first — lint/test depend on generated code:

```bash
# Only if proto files were modified
/usr/bin/make generate.buf

# Only if migration files were modified
/usr/bin/make migrate.up

# Only if repository interfaces changed
/usr/bin/make generate.mock
```

### 6b. Lint + Test with Retry

Bounded retry loop (max 2 retries = 3 total attempts):

```
attempt = 1
max_attempts = 3
while attempt <= max_attempts:
  run `/usr/bin/make lint.go`
  if test files were touched: run `/usr/bin/make test`
  if all pass:
    break
  else:
    parse failures → synthetic findings (rule_id LINT-{n}/TEST-FAIL-{n},
    semantic_category lint_failure/test_failure, severity error)
    append to findings
    return to Step 5 for only these new findings
    attempt += 1
```

### 6c. Escalation After Exhaustion

If still failing after 3 attempts, list unresolved failures verbatim in the PR body and
report as "⚠️ N manual actions needed". Never report "Ready" if verification is failing.

## Step 7: Create PR

### Branch Naming

```
chore/audit-rules-compliance
```

If scoped to one partition:

```
chore/audit-rules-{partition}    # e.g. chore/audit-rules-domain
```

### Commit

Stage only the auto-fixed source files (NOT `.serena/`, NOT `.claude/` unless you're fixing
a rule file, NOT unrelated tracked changes):

```bash
git checkout -b chore/audit-rules-compliance
git add {only the fixed source files listed in auto-fixed findings}
git commit -m "chore: fix rule-compliance violations found by audit-rules"
git push -u origin chore/audit-rules-compliance
```

### PR Body (heredoc)

```bash
gh pr create --title "chore: fix rule-compliance violations (audit-rules)" --body "$(cat <<'EOF'
## Proposed Changes

- Auto-fixed rule-compliance violations detected by `/audit-rules`
- Partitions audited: domain, usecase, grpc, db-repo, proto, migration, infra-other, tests

## Implementation

### Audit Results

| Partition | Files | Findings | Auto-Fixed |
|-----------|-------|----------|------------|
| domain | N | N (E:N W:N) | N |
| usecase | N | N (E:N W:N) | N |
| grpc | N | N (E:N W:N) | N |
| db-repo | N | N (E:N W:N) | N |
| proto | N | N (E:N W:N) | N |
| migration | N | N (E:N W:N) | N |
| infra-other | N | N (E:N W:N) | N |
| tests | N | N (E:N W:N) | N |
| **Total** | | **N** | **N** |

### Issues Requiring Manual Action

{List any findings that could not be auto-fixed}

### Verification

- [x] `make lint.go` passes
- [x] `make test` passes
EOF
)"
```

**Do NOT add AI attribution** ("Generated with Claude Code", "Co-Authored-By: Claude").

## Step 8: Report

After the PR is created, output a summary:

```markdown
## Audit-Rules Summary

### Scope
- Partitions audited: N (or: single partition: {name})
- Total files reviewed: N

### Results by Partition
| Partition | Files | Findings | Auto-Fixed | Manual |
|-----------|-------|----------|------------|--------|
| domain | N | N (E:N W:N I:N) | N | N |
| ... | | | | |
| **Total** | | **N** | **N** | **N** |

### Auto-Fixed Issues (N)
1. **[CL-xxx]** Short description — `path/file.go:42`

### Issues Requiring Manual Action (N)
1. **[CL-xxx]** Short description — `path/file.go:42` — Reason

### Verification
- [x] lint.go passes
- [x] tests pass

### PR
{PR URL}

### Result
✅ Ready / ⚠️ N manual actions needed
```
