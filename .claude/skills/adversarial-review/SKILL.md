---
name: adversarial-review
description: "Adversarially verify a PR or the current branch diff — assume the change is wrong and try to prove it across 5 perspectives (設計 / 理解容易性 / 命名 / テスト / 機能・性能). Use when: (1) running '/adversarial-review' (current branch) or '/adversarial-review 123' (PR number), (2) asked for 敵対的検証 / adversarial verification / \"try to break this\", (3) invoked from review-diff via the adversarial-reviewer agent. Limited to concerns no other reviewer, linter, or CI job owns: traces entity lifecycles across all flows, checks state left behind on mid-flow failures, empirically verifies generated outputs, and lists mutations the tests would not catch. Read-only."
---

# Adversarial Review (敵対的検証)

Assume the change is wrong and try to prove it. Read-only.

## Scope

Report only what nothing else owns. **Do not check or report**:

| Delegated concern | Owner |
|-------------------|-------|
| Layer/package placement, DI, naming conventions, `.claude/rules/`, domain state-change methods | convention-reviewer |
| AI anti-patterns, redundant nil checks, nil deref, error swallowing, AP-28 | bug-reviewer |
| Test conventions, test quality, per-method case coverage | test-reviewer |
| PR/issue requirements, scope creep, declared 未対応 list | spec-reviewer |
| TX boundaries, authz, input validation, secrets in logs/responses, N+1, Preload | security-perf-reviewer |
| Unused code/params, complexity, nesting, function length, magic strings | golangci-lint (`unused`, `unparam`, `gocognit`, `cyclop`, `funlen`, `goconst`) via `make lint.go` |
| Spelling | `lint-typos` (CI) |
| Build / vet / unit tests / CI status | review-diff Step 6, CI |

Do **not** run `go build`, `go vet`, `go test ./...`, `make`, or `gh pr checks`, and do not read the PR body or issues in embedded mode.

## Modes

| Mode | Invoked by | Steps | Output |
|------|-----------|-------|--------|
| Embedded | `adversarial-reviewer` from review-diff | 2–4, **no command execution** beyond `git diff` / `grep` (no temp tests, `openssl`, `gh api`) | Step 5b |
| Standalone | `/adversarial-review [PR]` | 1–4, 5a, 6 | Step 5a; Step 6 optional |

## Step 1: Resolve Target (standalone)

`gh pr view {N} --json title,body,files` or `git diff --name-only origin/main...HEAD` (or `origin/master`). Use the PR body only to skip items it declares 未対応 (unless the change makes them reachable in production).

## Step 2: Read Narrowly

- Non-test files: `git diff <base>...HEAD -- <path>`. Do **not** Read whole files; open a range only around a hunk you are attacking.
- Test files (`*_test.go`): never diff or Read them in full. Grep only case names and expectations, then open a narrow range only to confirm a specific mutation gap:
  ```bash
  git diff <base>...HEAD -- <test> | grep -nE '^\+.*("[^"]+": func|NewMock|\.EXPECT\(\)|require\.|assert\.)'
  ```
- Skip generated code (`**/mock/**`, `**/pb/**`, `**/dbmodel/**`), `go.sum`, binaries/images.
- Lifecycle tracing outside the diff: `grep -n` the entity's state-change methods and repository calls, then read only those function bodies.

## Step 3: Attack

### 設計
- **責務設計**: variant-specific policy hidden in a generic method (`if provider == X` inside a shared transition); one component doing several jobs.
- **変更容易性**: how many places change when the next variant (provider/status/tenant type) is added.
- **一貫性**: the same state/concept means different things across variants.
- **テスト容易性**: changed paths that cannot be tested as written (side effects not behind an interface, constructors reading env or panicking).

### 理解容易性
- **簡潔性**: single-caller abstractions, behavior-switching flags.
- **可読性**: comments that are misleading or removed while the behavior they explained remains.
- **冗長性**: duplicated logic, re-fetching what is already loaded.
- **インライン func**: in non-test code, a local func literal helper (`build := func(...)`) or an immediately-invoked func inside a method — a hidden private helper (AP-25/31 workaround). Fix: inline flat, or move to `domain/model` / `domain/service`. Not findings: API-required closures (`RWTx`/`ROTx`, `sort.Slice`, `errgroup.Go`, `defer`), test closures, AP-28.

### 命名
- **責務・役割の明確性**: a name that does not hold in every branch (a `MarkDeleted` that no-ops for some states, a `Get` that writes).
- **英文法**: grammar in comments and error messages.

### テスト
- **E2E**: asserts only HTTP status where the API claims a DB/object state change; scripts relying on `source` order.
- **Mutation** (main test item): plausible bad edits that stay green — swallowed error path, swapped side-effect order, removed early return, signing/format option dropped (digest, detached signature, embedded cert). Phrase as "〇〇 に変えても検出されない" + the test to add.

### 機能・性能
- **Lifecycle across flows** (highest value): walk create → external events/webhooks → update → delete → re-create → transfer per variant. Look for silent no-ops behind 200, uniqueness checks that block forever, deleted entities still served elsewhere, features reachable in production whose prerequisite is missing.

  **Mandatory questions** — answer each explicitly for every operation the diff newly makes callable in production (new enum value accepted, new branch, new endpoint). A "no" is a finding at the severity shown:
  1. Using **only fields in the API response** (check the marshaller / output DTO), can the user perform the next step (download, open, accept, pay)? List every other entry point for that next step (grep RPCs/routes and their auth requirements); if each needs something the user was never given → **High**.
  2. After delete, can the same resource be created again (existence checks without status filter, unique constraints)? If it is blocked forever → **High**, unless that is the documented intent.
  3. After delete or transfer, does any other endpoint (webhook, worker, download, sync) still serve it as valid? → **High**.
  4. Does each existing state mean the same thing for the new variant as for the old ones? → **Medium**.
- **Refactor equivalence**: extracted/replaced logic identical in every branch.
- **Failure-time consistency**: state left behind when a step fails mid-flow or a concurrent duplicate loses (not *where* the TX boundary is — that is security-perf).
- **Empirical output check** — **standalone mode only**, and only when the diff adds/changes a generator of signed files, archives, or external formats: generate once with throwaway keys in `mktemp -d`, verify with one independent tool (`openssl cms -verify`, `unzip -l`, `jq`), and inspect fallback rendering when optional fields are absent. Delete the temp dir.
- **External platform spec**: deprecated or misused fields of the external service (IdP, payment, push, storage, wallet vendors, RFC) — from knowledge; no web fetch.
- **Startup & deploy**: new `notEmpty` env vars + panicking constructors vs deploy order; misconfig (key/cert mismatch, NotBefore, identifier mismatch, expiry while running) detected only by clients.
- **パフォーマンス**: per-request CPU-heavy work (signing, compression, parsing) that could be done once.
- **過剰なデータ通信**: whole objects fetched from storage when metadata suffices; response fields the client does not need.
- **脆弱性**: unescaped values in generated markup/URLs (by reading the code; no execution in embedded mode); ID/token charset & length. **New modules in `go.mod` — standalone mode only**: license, archived?, `gh api /repos/{o}/{r}/security-advisories` for the pinned version.

Pre-existing weaknesses the change makes reachable: tag `(既存、この変更で実害化)`.

### Refute before reporting

For each candidate, look for a guard elsewhere (DB constraint, another usecase, middleware) with one targeted grep. Drop it if refuted.

**Embedded mode exception** — time first: do not spend extra calls refuting. Report unconfirmed findings/mutation gaps prefixed `未確認:` with the exact condition to check; review-diff Step 4 resolves them. The rules below apply to standalone mode.

**No hedged findings.** If a finding depends on "unless X exists elsewhere" (another endpoint, handler, job, or guard outside the diff), grep for X once — e.g. other RPCs/routes returning the field, other callers of the method — then either confirm it and keep the full severity, or drop it. Never report with the hedge and never lower severity because of an unverified hedge.

For each **mutation gap**, grep the test files for the symbol/field the mutation would break (e.g. the response field, the status constant, the signing option, the mocked method name) and read the matching assertion lines. If any assertion would fail under the mutation, drop it. Beware of setup lines next to assertions (e.g. `require.Empty(out)` immediately followed by `out = build()` is still an assertion on the preceding step).

Also check **absent expectations** — gomock fails on any unexpected call, so a missing `EXPECT()` is itself an assertion. For order/skip mutations (swapped side effects, error swallowed and execution continuing), look at the failure cases: if a case makes an earlier step fail and creates the later dependency's mock with **no** `EXPECT()` for that call (e.g. "create failed" with `NewMockAsset(ctrl)` but no `PutObject`), the mutation is detected — drop it. `gomock.InOrder` is not required for that. Severity: **High** = wrong behavior / stuck data / security hole in a normal flow; **Medium** = inconsistency under failure or race, design that breaks on the next variant, test gap letting a High through; **Low** = hardening, deprecated field, edge-case display, naming.

## Step 4: Clean Up

Delete temp dirs and any temporary `_test.go`. `git status` must show no changes from you.

## Step 5a: Report (standalone, Japanese)

```markdown
敵対的検証を行いました。<総評 1 行>

## 問題なしと確認した点
- **[<観点>] <対象>**: <試したこと>(最大 5 件、各 1 行)

## 指摘
### [High][<観点>/<細目>] <タイトル>
- 場所 / 問題 / 再現シナリオ / 対処案 (a)(b)

## テストで検出できない変更(mutation)
- <変更> に変えても検出されない → <追加すべきテスト>
```

## Step 5b: Report (embedded)

Keep it short: no coverage table, Verified OK at most 5 one-line items, each finding ≤ 6 lines.

```markdown
## Adversarial Review Findings

### Verified OK
- [<perspective>] <target>: <evidence>

### Findings
#### [error|warning|info] file/path.go:L42 — Title
- **Rule**: ADV-{design|read|name|test|func}-{n} / **Perspective**: 設計/一貫性 / **Semantic Category**: {key} / **Auto-fixable**: yes|no
- **Scenario**: input/state → wrong result
- **Fix**: options (a)/(b)

### Mutation Gaps
- <layer>: <mutation> not detected → <test to add> (ADV-test-{n}, Auto-fixable: yes|no)

### Summary
Findings: N (error: N, warning: N, info: N), Mutation gaps: N
```

High → `error`, Medium → `warning`, Low → `info`. `Auto-fixable: yes` only for mechanical fixes (escaping, rename, adding a well-defined test case).

`semantic_category` — only these keys; a finding that needs another agent's key is out of scope:

| Perspective | Keys |
|-------------|------|
| 設計 | `adv_responsibility` `adv_changeability` `adv_consistency` `adv_testability` |
| 理解容易性 | `adv_simplicity` `adv_readability` `adv_redundancy` `adv_inline_func` |
| 命名 | `adv_naming_role` `adv_grammar` |
| テスト | `adv_e2e_gap` `adv_mutation_gap` |
| 機能・性能 | `adv_lifecycle` `adv_refactor_equivalence` `adv_output_verification` `adv_external_spec` `adv_failure_consistency` `adv_deploy_safety` `adv_performance` `adv_overfetch` `adv_escaping` `adv_dependency` |

## Step 6: Post as PR Review (standalone, optional)

Only for a PR target and only after the user confirms. Inline comments for lines in the diff (`**[Severity][観点] タイトル**`), everything else in the review body. Post once:

```bash
gh api repos/{owner}/{repo}/pulls/{N}/reviews --method POST --input <tmp>/review.json
# {"event":"COMMENT","body":"...","comments":[{"path":"...","line":123,"side":"RIGHT","body":"..."}]}
```
