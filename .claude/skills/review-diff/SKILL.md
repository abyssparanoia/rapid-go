---
name: review-diff
description: Review and auto-fix code changes against project conventions by diffing the current branch against the default branch (main/master). Use when: (1) asked to review current changes, (2) running '/review-diff', (3) after completing a feature implementation and wanting to catch issues before creating a PR. Launches 5 specialized review agents in parallel (spec, convention, bug, security/performance, test quality), aggregates findings, and automatically fixes all issues found. Does NOT require an existing PR.
---

# Diff Review & Auto-Fix (Parallel Agent Edition)

Review current branch changes against main/master using **5 specialized parallel agents**, then **automatically fix all issues found**.

## Workflow

```
1. Detect Default Branch  → find main or master
2. Gather Changes         → git diff, changed files, classify
3. Parallel Agent Review  → launch 5 agents simultaneously
4. Aggregate Findings     → collect, deduplicate, sort
5. Auto-Fix Issues        → edit files to fix all found problems
6. Verify Fixes           → lint, tests, code generation
7. Report                 → summary with agent-level breakdown
```

## Step 1: Detect Default Branch

```bash
git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's|refs/remotes/origin/||'
```

If that fails, try `origin/main` first, then `origin/master`. Store the result as `$BASE`.

## Step 2: Gather Changes + Classify

```bash
BASE=origin/main  # or origin/master

# Changed files list
git diff --name-only $BASE...HEAD

# Check if test files are included
git diff --name-only $BASE...HEAD | grep '_test\.go$'

# Commit history for spec-reviewer context
git log $BASE...HEAD --oneline
```

Note the changed file list — this is passed to agents in their prompts.

Determine which agents to launch:
- **spec-reviewer**: Always (spec compliance)
- **convention-reviewer**: Always (convention & rules compliance)
- **bug-reviewer**: Always (bug & anti-pattern detection)
- **security-perf-reviewer**: Always (security & performance)
- **test-reviewer**: Only if `*_test.go` files exist in the diff (test quality)

## Step 3: Parallel Agent Review

Launch all applicable agents in a **single message** using the Agent tool. All agents run with `run_in_background: true` for parallel execution.

### CRITICAL: Use `subagent_type`, Do NOT Embed the Agent Definition

Each agent is registered as a `subagent_type` via its frontmatter in `.claude/agents/{name}.md`. When you specify `subagent_type: "{name}"`, Claude Code automatically applies:

- The `model` field (e.g., `sonnet`) — prevents cost overruns from running on the parent's model
- The `tools` allow-list (e.g., `[Read, Glob, Grep, Bash]`) — enforces read-only, preventing unintended edits
- The agent's full body as its system prompt

**Never** read the `.claude/agents/{name}.md` file and embed its contents in the `prompt`. Doing so:

- Causes the agent to launch as `general-purpose` (all tools, including Edit/Write/Bash unrestricted) — breaks the read-only guarantee
- Bypasses the `model: sonnet` declaration — may run on an expensive model
- Duplicates the definition between the agent file and SKILL.md, causing drift

### Prompt Content (review context only)

The `prompt` argument should contain **only** the review context. The agent already has its role from the frontmatter/body.

**Template**:
```
Target diff: HEAD vs {$BASE}

Changed files:
{changed_files}

Follow your role definition exactly. Report findings in the format you define, including the Semantic Category field on every finding.
```

For **spec-reviewer**, also append:
```
Commit history:
{git log output}
```

For **test-reviewer**, filter `{changed_files}` to `*_test.go` only.

### Agent Tool Call Pattern

Launch all applicable agents in one message:

```
Agent(
  description: "Spec compliance review",
  subagent_type: "spec-reviewer",
  run_in_background: true,
  prompt: "<review context only>"
)
Agent(
  description: "Convention & rules review",
  subagent_type: "convention-reviewer",
  run_in_background: true,
  prompt: "<review context only>"
)
Agent(
  description: "Bug & anti-pattern review",
  subagent_type: "bug-reviewer",
  run_in_background: true,
  prompt: "<review context only>"
)
Agent(
  description: "Security & performance review",
  subagent_type: "security-perf-reviewer",
  run_in_background: true,
  prompt: "<review context only>"
)
Agent(                                    # only if test files exist
  description: "Test quality review",
  subagent_type: "test-reviewer",
  run_in_background: true,
  prompt: "<review context only>"
)
```

Wait for all agents to complete. You will be notified as each finishes.

## Step 4: Aggregate Findings

After all agents complete:

1. **Collect** all findings from agent responses (each finding now carries `semantic_category`)
2. **Deduplicate** by `(file, line, semantic_category)` — see dedup rationale below
3. **Sort** by severity: `error` first, then `warning`, then `info`
4. **Verify coverage** — ensure every changed file was reviewed by at least one agent
5. **Separate** auto-fixable from manual-only findings

### Dedup Rationale: Why `semantic_category`, Not `rule_id`

Each agent defines its own rule_id prefix (`AP-*`, `SEC-*`, `CL-*`, etc.), so the same underlying issue can be reported under different rule_ids. Example: a missing `Preload: true` on a return-path `Get` could surface as `AP-21` (bug-reviewer), `PERF-...` (security-perf-reviewer, `query_preload_missing`), or `CL-repo-*` (convention-reviewer) — `rule_id`-based dedup would miss the overlap. `semantic_category` is defined by each agent in its own definition file; the orchestrator uses it as a single cross-agent dedup key.

### Tie-breaking When Multiple Agents Flag the Same `semantic_category`

When the dedup key collides, keep the finding with:
1. **Highest severity** (`error` > `warning` > `info`)
2. If tied, **most specific fix** (non-empty `Fix:` block wins over empty)
3. If still tied, **convention-reviewer > security-perf-reviewer > bug-reviewer > test-reviewer > spec-reviewer** (preferred owner order)

Note the winner's `rule_id` in the final report but merge all agents' notes into the finding's Description so no context is lost.

### Finding Categories

| Rule Prefix | Source Agent | Example |
|---|---|---|
| `SPEC-*` | spec-reviewer | SPEC-1 (missing requirement) |
| `CL-*` | convention-reviewer | CL-usecase-3 (missing Validate), CL-migration-2, CL-proto-5 |
| `AP-*` | bug-reviewer | AP-36 (direct struct init) |
| `BUG-*` | bug-reviewer | BUG-1 (nil dereference) |
| `SEC-*` | security-perf-reviewer | SEC-1 (missing authorization) |
| `PERF-*` | security-perf-reviewer | PERF-1 (N+1 query) |
| `TX-*` | security-perf-reviewer | TX-1 (ForUpdate missing) |
| `TEST-*` | test-reviewer | TEST-1 (insufficient coverage) |

### Agent Scope Matrix

| Concern | Owner agent |
|---------|------------|
| Test files (`*_test.go`) | **test-reviewer** only — bug-reviewer skips tests |
| TX boundary (RWTx, ForUpdate, nesting, long TX) | **security-perf-reviewer** only — bug-reviewer delegates |
| IdP sync (StoreClaims, DeleteUser, order) | **security-perf-reviewer** only |
| Input validation (`param.Validate()` missing) | **security-perf-reviewer** |
| Migration safety | **convention-reviewer** |
| Proto backward-compatibility | **convention-reviewer** |
| Preload efficiency | security-perf-reviewer (`query_preload_unnecessary` / `query_preload_missing`) |
| Logic bugs & non-TX anti-patterns | **bug-reviewer** |

## Step 5: Auto-Fix Issues

**Fix all auto-fixable issues immediately.** Edit files using available tools.

Fix in this priority order:
1. **Security** (SEC-*) — missing authorization, input validation gaps
2. **Correctness** (BUG-*, AP-19, AP-20, AP-21) — ForUpdate, Preload, IdP sync
3. **Convention** (CL-*) — method ordering, naming, pattern compliance
4. **Anti-patterns** (AP-*) — gomock.Any(), direct field assignment
5. **Test quality** (TEST-*, AP-1~6) — mock strictness, t.Parallel()
6. **Performance** (PERF-*) — N+1, unnecessary queries

When fixing test files that use `gomock.Any()` incorrectly, replace with exact expected values using domain model constructors.

**SPEC-* findings are NOT auto-fixed** — they require human judgment on specification interpretation.

## Step 6: Verify Fixes (with Retry Loop)

### 6a. Code Generation (if applicable)

Run generation first — lint/test depend on generated code being up to date:

```bash
# Migration files changed
/usr/bin/make migrate.up

# Proto files changed
/usr/bin/make generate.buf

# Repository interfaces changed
/usr/bin/make generate.mock
```

### 6b. Lint + Test with Retry

Run verification in a **bounded retry loop** (max 2 retries = 3 total attempts):

```
attempt = 1
max_attempts = 3
while attempt <= max_attempts:
  run `/usr/bin/make lint.go`
  if test files changed: run `/usr/bin/make test`
  if all pass:
    break
  else:
    parse failure output into synthetic findings:
      - file / line extracted from error location
      - rule_id = LINT-{n} or TEST-FAIL-{n}
      - semantic_category = `lint_failure` or `test_failure`
      - severity = error
      - description = failure message
    append to aggregated findings
    return to Step 5 (Auto-Fix) scoped to only these new findings
    attempt += 1
```

### 6c. Escalation After Exhaustion

If attempt > max_attempts:
- Do **not** discard the remaining failures
- List each unresolved failure in Step 7's "Issues Requiring Manual Action"
- Include the failure message verbatim and the file/line
- Mark overall status as ⚠️ `N manual actions needed`

### Anti-pattern: Silent Success

Never report "✅ Ready" if lint or test failed and was not re-fixed. The verify-retry loop exists specifically so the reviewer cannot silently lie about verification.

## Step 7: Report

```markdown
## Diff Review Summary

### Scope
- Base: origin/main (or master)
- Changed files: N files
- Agents launched: N

### Agent Results
| Agent | Files Reviewed | Findings | Auto-Fixed |
|-------|---------------|----------|------------|
| spec-reviewer | N | N (E:N W:N I:N) | 0 |
| convention-reviewer | N | N (E:N W:N I:N) | N |
| bug-reviewer | N | N (E:N W:N I:N) | N |
| security-perf-reviewer | N | N (E:N W:N I:N) | N |
| test-reviewer | N | N (E:N W:N I:N) | N |
| **Total** | | **N** | **N** |

### Auto-Fixed Issues (N)
1. **[Rule ID]** Short description
   - `path/to/file.go:123`
   - Was: `<bad code>`
   - Fixed: `<good code>`

### Issues Requiring Manual Action (N)
1. **[Rule ID]** Short description
   - `path/to/file.go`
   - Reason: Requires spec confirmation / `make generate.buf` needed

### Verification
- [x] lint.go passes
- [x] tests pass
- [ ] needs `make generate.buf`

### Result
✅ Ready / ⚠️ N manual actions needed
```
