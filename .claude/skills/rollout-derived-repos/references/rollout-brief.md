# Rollout brief: Go/deps upgrade + `.claude/` sync for a rapid-go-derived repository

You are working on ONE repository (given in your prompt) and produce up to two PRs in it:
(A) Go + dependency upgrade to rapid-go's current state, (B) pull-only `.claude/` sync from rapid-go.
Do (A) first, then (B), as the `Mode` in your prompt allows. Everything below is a SETTLED decision —
do not re-decide; if something cannot be done within these rules, hold it, open an issue, and escalate in your final report.

## Target State (rendered by the orchestrator from rapid-go HEAD — never edit by hand)

| Item | Value |
|---|---|
| rapid-go root | `{{RAPID_GO_ROOT}}` |
| Go version (`go` directive) | `{{GO_VERSION}}` |
| golang docker tag | `{{GOLANG_IMAGE_TAG}}` (repos that pin a patch tag use `golang:{{GO_VERSION}}-alpine`) |
| golangci-lint (go.mod tool AND CI action `version:`) | `{{GOLANGCI_VERSION}}` |
| golang.org/x/tools | `{{XTOOLS_VERSION}}` |
| Reference PRs in rapid-go (read them for exact style, PR body wording and the fixes they needed) | {{REFERENCE_PRS}} |
| Latest rapid-go `.claude/` changes (context for B) | {{CLAUDE_CHANGES}} |

rapid-go's `go.mod`, `.golangci.yml`, `.github/workflows/ci.yml`, `docker/*.Dockerfile` and `.claude/` are the target; read them when in doubt.

## Ground rules (both workstreams)

1. **Never touch the user's checkout** (it may be on a feature branch with uncommitted files). Work ONLY in git worktrees under your scratchpad directory:
   ```bash
   REPO=<local clone path>            # or: gh repo clone <slug> <scratch>/<repo> and use that as $REPO
   SLUG=<owner/repo>
   git -C "$REPO" fetch origin
   DEFAULT=$(gh repo view "$SLUG" --json defaultBranchRef -q .defaultBranchRef.name)
   WT=<scratchpad>/wt/<repo-basename>-<deps|claude>
   git -C "$REPO" worktree add -b <branch> "$WT" "origin/$DEFAULT"
   ```
   If `<branch>` already exists locally or on origin (from an earlier merged rollout), use `<branch>-YYYY-MM` instead — never delete or force-push an existing branch. Do all edits, builds, commits and pushes inside `$WT`. After the PR's CI is green: `git -C "$REPO" worktree remove "$WT"` (the branch stays pushed).
2. Use `/usr/bin/make` (plain `make` may be shadowed). Never `cd` in loops spanning repos; prefer `git -C` / `make -C` / absolute paths. If a repo has git submodules (`.gitmodules`), run `git -C "$WT" submodule update --init` before building.
3. Conventions: follow the repo's own `.claude/skills/create-pull-request/SKILL.md` if present, else rapid-go's. **No AI attribution anywhere** (no "Generated with Claude", no Co-Authored-By, no trailing attribution tokens). Commit messages in the repo's prevailing style (`git log --oneline -15`). PR title in English, PR body in Japanese (same shape as the reference PRs).
4. Do not merge anything. Do not force-push. Do not close other people's PRs (Renovate/Dependabot PRs your PR supersedes: list them in the PR body). Do not touch other repos.
5. After `gh pr create`, watch CI: `gh pr checks <n> --watch --interval 20` (if it prints "no checks reported", sleep 30 and retry). Fix failures that are within these rules with follow-up commits. If a job fails for reasons outside your control (missing secrets, flaky infra, failure that also exists on the default branch), say so explicitly with the job link.
6. Read `.claude/` files with the Read tool (never cat/sed) so the repo's path-scoped rules load; mechanical bulk replacements (sed) are fine for `.go` / `go.mod` / YAML edits you have already decided.
7. Final report (compact, ≤ 25 lines): repo; PR URLs with CI state; versions changed (Go, golangci-lint, notable deps); code changes made; items held back + issue URLs; "local improvements worth pushing upstream"; anything needing the user's decision.

## Workstream A — Go + dependencies (branch `chore/update-go-<major.minor>-and-deps`)

Target = the Target State table. Verify the repo's actual state from `origin/<default>` first (the survey in your prompt may be stale); if part of the work is already merged, do only what remains and say so.

1. `go mod edit -go=<GO_VERSION>`; drop a `toolchain` line if it is now redundant. Docker images that build THIS app → the target golang tag (keep digest pinning if the repo pins by digest — look up the new digest the same way the repo's last Go bump PR did). Fix other hard-coded Go versions (CI `go-version:`, `.tool-versions`, `mise.toml`, Makefile, docs/CLAUDE.md "Go x.y+" text). `go-version-file: ./go.mod` needs no change.
2. Upgrade direct dependencies (require block AND `tool` block, or `tools/tools.go` in older repos) with explicit `go get <mod>@latest ...` (NOT `go get -u ./...`), then `go mod tidy`. List candidates with
   `go list -m -u -f '{{if and .Update (not .Indirect)}}{{.Path}} {{.Version}} -> {{.Update.Version}}{{end}}' all`
   plus each tool module individually (`go list -m -u <tool-module>`; `all` omits tool-only modules).
   - Avoid retracted versions (watch `go get` warnings; e.g. mimetype v1.4.15 → use v1.4.16) and pre-release/rc versions.
   - Respect `renovate.json` `allowedVersions` / disabled packages in the repo.
   - If a bump pulls an indirect module whose old version no longer builds (e.g. `otelzap`), bump that indirect module explicitly.
3. **golangci-lint** → target version in BOTH the go.mod `tool` block and the CI workflow (`golangci/golangci-lint-action` `version:`); the two must match. **x/tools** → target version. (Go export-data format changes require matching golangci-lint/x/tools: an old linter binary fails with "export data version N is greater than maximum supported version".)
   Then `go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint config verify` and fix schema errors in `.golangci.yml` with behaviour-equivalent settings (see Known fallout). Run lint on tracked packages (`run ./cmd/... ./internal/...` or the Makefile target). Resolve findings:
   - wrapcheck on the project's own error base (`github.com/abyssparanoia/goerr` or the repo's equivalent): add it to `wrapcheck.ignore-package-globs` (as rapid-go did).
   - Mechanical fixes (govet `inline` such as `reflect.Ptr` → `reflect.Pointer`, `WriteString(fmt.Sprintf)` → `fmt.Fprintf`, goconst constants): apply directly.
   - If a bump produces many findings needing real refactors (roughly > 30 non-mechanical code edits, or semantic changes), keep golangci-lint at the highest version that passes (never below the minimum the Go version needs), create an issue listing the findings by linter, and continue.
   - "parallel golangci-lint is running" means another lint process (possibly another repo's agent) holds the lock — wait and retry, do not kill it.
4. Code fixes for API changes in upgraded libraries are ALLOWED when small and mechanical (the user permits code changes in derived repos). If an upgrade needs a non-trivial rewrite or changes behaviour (deprecated APIs with replacements that alter semantics, module path moves such as `cel.dev/cel-go` or `github.com/aarondl/*`), hold that module at its current version, create an issue describing the breakage and the migration, and continue.
5. Generated code: proto generation via BSR remote plugins needs no regen. If local protoc plugins in the tool block were bumped, run the repo's buf generate target and commit the deterministic diff. If `go.uber.org/mock` was bumped, run the mock generation target and commit. Do NOT regenerate SQLBoiler/yo models (needs databases); say so in the PR body. If the new Go/gofumpt reformats files, commit the formatting-only diff separately.
6. Validate before pushing: `go build ./...`, `go vet ./...`, `go build -o /dev/null <every tool module>`, lint (0 issues), `make test` (or equivalent), proto lint if present, formatting target if present.
7. PR body: same structure as the reference PRs (Proposed Changes incl. the code changes and why they were needed; Direct dependencies table; Tools table; Implementation = what was validated; superseded Renovate/Dependabot PRs; held-back items with issue links).

### Holds (do not upgrade across these; update within the current major only)

| Module | Rule | Why / tracking |
|---|---|---|
| `github.com/aarondl/null`, `github.com/volatiletech/null` | stay on the current major | rapid-go pins v8 to match SQLBoiler (renovate.json) |
| `github.com/lestrrat-go/jwx` | stay on the current major(s); keep both if v2 and v3 coexist | v4 needs GOEXPERIMENT=jsonv2 |
| `github.com/cloudspannerecosystem/wrench`, `go.mercari.io/yo` | HOLD in every repo; create ONE issue "Upgrade wrench and yo (Spanner tooling)" if present | newer versions need memefish v0.9.0 which breaks `github.com/abyssparanoia/memeduck` (rapid-go #1989); yo v0.11 regenerates Spanner models and needs template validation against an emulator. If the repo has no Spanner at all, the issue should propose removing them instead. |
| `github.com/aws/aws-sdk-go-v2/feature/s3/manager` | stay ≤ v1.21.1 | v1.22+ is deprecated (staticcheck SA1019); migration to transfermanager rewrites upload code → issue |
| `github.com/ericlagergren/decimal` | keep the version SQLBoiler's `types.Decimal` works with | newer versions break `pq: encode: unknown type for types.Decimal` at runtime (sunabe-server #89) |

When a hold stops being necessary (upstream fixed), remove the row in rapid-go so the next rollout picks it up.

### Known fallout (what past rollouts needed; apply when you hit the same thing)

- golangci-lint v2.15: `exhaustive.check: map` → `mapliteral` (verify rejects `map`; `run` silently ignores it, disabling the check); `gomodguard`/`exhaustruct` deprecation warnings are warnings only — leave them.
- golangci-lint v2.14+: goconst also counts strings in composite literals → if that floods findings, `exclude-types: [Call, CompositeLit]` keeps the previous scope.
- Go 1.27 gofmt re-indents multi-value returns (large whitespace-only diffs) — separate commit.
- Go 1.27 coverage blocks start at the first statement, not the `{` — tooling that matches `cover:ignore` markers by line may need to walk back over comments.
- go-elasticsearch v8.18+: `NestedQuery.Query` became a value type.
- buf ≥ v1.75 requires a protobuf pseudo-version (`v1.36.13-0.2026…`) — accept it, note it in the PR body.

## Workstream B — `.claude/` sync from rapid-go (branch `chore/sync-claude-config`, suffix `-YYYY-MM` if taken)

Pull-only sync of rapid-go's `.claude/` into this repo, in a second worktree branched from `origin/<default>`. Follow `{{RAPID_GO_ROOT}}/.claude/skills/sync-claude-config/SKILL.md` and its `references/` (token-mapping.md, sync-algorithm.md) EXCEPT: no push to rapid-go, and no interactive confirmation (the user pre-approved; you decide conflicts with the rules below).

Decision rules:
- UPSTREAM = rapid-go `.claude/` at HEAD. Apply the forward token map (rapid-go canonical → this repo's tokens: go module, service name, buf org, docker network, project title, org/repo, database) to everything you copy. Afterwards grep the written files for leftover `rapid-go` / `abyssparanoia/rapid` / `RAPID GO` tokens — legitimate ONLY inside `sync-claude-config`, `rollout-derived-repos`, `init-new-repository` and `config-sync-classifier`, which must stay byte-identical to upstream (restore them if an earlier init localized them).
- Pull a file when upstream is newer or structurally improved and the local file is the generic template version (possibly stale). Keep the local file when it carries project-specific content (domain rules, project-only skills/agents/commands, custom permissions, a different job/worker architecture). When both apply, MERGE: upstream's structure/improvements plus the project-specific sections. Never delete project-specific local files; delete local files only when upstream removed them because they were migrated (e.g. `commands/x.md` whose content now lives in `skills/x/SKILL.md`), and only after the replacement is in place.
- `.claude/settings.json`: merge (add generic env flags upstream added), never drop local permissions/hooks; create it with the env flag only if absent. Never touch `settings.local.json`.
- `.claude/CLAUDE.md`: keep the project's own overview/tech stack/directory structure; bring the Rules table, Skills table, Settings section, Implementation Preflight and PR checklist to upstream's shape, with rows matching what exists in THIS repo after the sync.
- Rules' `paths:` frontmatter: keep only globs for layers/directories this repo has (drop mysql/spanner/cognito/s3 globs where absent; keep a repo's more precise local globs if they are a superset of upstream's intent).
- DB/language specifics: keep what matches this repo (MySQL vs PostgreSQL vs Spanner; if the repo writes `.claude` docs in Japanese, translate new upstream text rather than switching the file to English).
- If the local version of a file is clearly better/more specific than upstream, keep local and list it under "local improvements worth pushing upstream" (do not push).
- If nothing meaningful changes after normalization, do not open a PR — but DO still fix broken sync tooling (see first bullet) and stale table rows, and say what you did.
- Validate: every touched rules/agents/skills file has frontmatter that parses as strict YAML (quote `description:` values containing `: `); `paths:` globs reference existing directories; no references to files that don't exist in this repo.
- PR body: table of files [new | updated | merged | removed] with one-line reasons; conflicts resolved; what was intentionally kept local.
