---
name: sync-claude-config
description: "Bidirectionally sync .claude/ content (skills, rules, agents, commands, CLAUDE.md) between this project and the rapid-go template. Use when: (1) running /sync-claude-config, (2) you want project-local Claude-config improvements pushed upstream to rapid-go and upstream improvements pulled down, (3) coordinating changes across multiple rapid-go-derived projects. Requires rapid-go added via `claude --add-dir <rapid-go-path>` before running."
---

# sync-claude-config

Bidirectionally syncs `.claude/` content between this project and the rapid-go upstream template, creating a PR in each repo for changes that flow in each direction.

## Prerequisites

Both repositories must be accessible:
- **LOCAL** — the current project (where you invoke the skill)
- **UPSTREAM** — rapid-go

Before running this skill, add rapid-go as an additional directory:

```
/add-dir /path/to/rapid-go
```

If only one repository root is visible, stop and prompt the user to run the above command. Do NOT proceed without both repos accessible — the skill cannot diff them.

## Step 0 — Identify Repos & Tokens

Determine which root is LOCAL and which is UPSTREAM.

```bash
# Current project root
git rev-parse --show-toplevel

# Identify rapid-go root among the accessible directories by checking
# (any one of these conditions is sufficient):
#   go.mod module == github.com/abyssparanoia/rapid-go
#   schema/proto/rapid/ exists
#   .claude/skills/init-new-repository/ exists
```

If the skill is invoked FROM rapid-go (i.e., current project IS rapid-go) and a derived project was added via `--add-dir`, swap roles: LOCAL = derived project, UPSTREAM = rapid-go. The logic is always: UPSTREAM = rapid-go, LOCAL = the other.

Read the LOCAL project's token values for normalization (see `references/token-mapping.md` for full mapping):

| Token | How to read |
|-------|-------------|
| `{go-module}` | First line of `go.mod`: `module github.com/...` |
| `{service-name}` | Directory name under `schema/proto/` that is NOT `google` or `protoc-gen-openapiv2` |
| `{project-title}` | H1 line in `.claude/CLAUDE.md`: `# Project Title` |
| `{buf-org}` | `buf.yaml` or `buf.gen.yaml`: the org before `/{service-name}` in registry URL |
| `{docker-network}` | `docker-compose.yml`: the custom network name (e.g. `myapp-network`) |
| `{org}/{repo}` | Derive from `{go-module}` after `github.com/` |
| `{database}` | Which of `db/mysql/` or `db/postgresql/` exists |

## Step 1 — Enumerate Candidates

Build the set of in-scope file paths on both sides. In-scope paths:

```
.claude/CLAUDE.md
.claude/rules/**/*.md
.claude/skills/**    (all files recursively)
.claude/agents/**/*.md
.claude/commands/**/*.md
```

**Always exclude** (never sync these):

```
.claude/settings.local.json
.claude/worktrees/**
.claude/skills/init-new-repository/**   ← intentionally template-bearing; stays in rapid-go only
```

Compute the **union** of relative paths from both sides.

## Step 2 — Normalize & Classify

For each path in the union, classify it into one of four buckets. Detailed algorithm in `references/sync-algorithm.md`. Summary:

1. Read file content from both sides (use `null` when the file doesn't exist on one side).
2. Normalize the LOCAL version by applying the **reverse token map** (local → rapid-go canonical). Full map in `references/token-mapping.md`. Apply replacements longest-first to avoid partial matches.
3. Compare normalized-LOCAL vs UPSTREAM content:

| Exists in UPSTREAM | Exists in LOCAL | Match after normalization | Classification |
|--------------------|-----------------|--------------------------|----------------|
| Yes | No | — | **pull** (add to LOCAL) |
| No | Yes | — | **push** (add to UPSTREAM) |
| Yes | Yes | Yes | **skip** (identical, no action) |
| Yes | Yes | No | **conflict** (changed on both sides) |

4. For **conflict** files:
   - Run `git -C <UPSTREAM> log -1 --format="%ct %H %s" -- <relpath>` and `git -C <LOCAL> log -1 --format="%ct %H %s" -- <relpath>`. Compare unix timestamps.
   - If only DB-variant lines differ (e.g. `` `mysql` `` vs `` `postgresql` `` in examples), mark **skip — DB-specific** instead.
   - Otherwise present the normalized diff and timestamp hint. Ask which side is canonical before proceeding.

## Step 3 — Present Plan & Confirm

Print a consolidated table before touching anything:

```
PULL  →  LOCAL       (changes from rapid-go into this project)
─────────────────────────────────────────────────────────────
.claude/rules/new-rule.md                [new file]
.claude/skills/new-skill/SKILL.md        [new file]
.claude/rules/grpc-handler.md            [conflict → upstream newer: 2025-06-01]

PUSH  →  UPSTREAM    (changes from this project into rapid-go)
─────────────────────────────────────────────────────────────
.claude/rules/project-specific.md        [new file]
.claude/rules/testing.md                 [conflict → local newer: 2025-06-03]

SKIP (no effective change after normalization): 14 files
SKIP (DB-specific difference): 2 files
```

**Do not apply any changes until the user explicitly confirms.**

For each conflict, show the diff (normalized) and ask: "Apply [upstream | local] version?" before including it in the plan.

If there are no changes in either direction, report that and stop — do not create empty branches/PRs.

## Step 4 — Apply Changes

After confirmation, apply the classified changes:

**Pull → LOCAL** (rapid-go content into this project):
- Take the UPSTREAM file content.
- Apply the **forward token map** (rapid-go canonical → local tokens) to localize it.
- Write/overwrite the file at `<LOCAL>/<relpath>`.

**Push → UPSTREAM** (local content into rapid-go):
- Take the LOCAL file content.
- Apply the **reverse token map** (local → rapid-go canonical) to canonicalize it.
- Write/overwrite the file at `<UPSTREAM>/<relpath>`.

Use the Write and Edit tools for LOCAL; use absolute paths under the UPSTREAM root for UPSTREAM files.

## Step 5 — Create PRs (Both Directions)

Follow `.claude/skills/create-pull-request/SKILL.md` conventions. **No AI attribution**.

### PR in LOCAL repo (pull changes from upstream)

```bash
# From LOCAL repo root
git checkout -b chore/sync-claude-config
git add .claude/
git commit -m "chore: sync .claude/ config from rapid-go"
git push -u origin chore/sync-claude-config
gh pr create \
  --title "Sync .claude/ config from rapid-go" \
  --body "$(cat <<'EOF'
## Proposed Changes

- Pulled .claude/ content updates from rapid-go upstream template

## Implementation

Files synced from rapid-go:
<list pulled files>

EOF
)"
```

### PR in UPSTREAM repo (push changes from local)

```bash
# All git commands use -C <UPSTREAM> to target the upstream repo
git -C <UPSTREAM> checkout -b chore/sync-claude-config-from-<service-name>
git -C <UPSTREAM> add .claude/
git -C <UPSTREAM> commit -m "chore: sync .claude/ config from <service-name>"
git -C <UPSTREAM> push -u origin chore/sync-claude-config-from-<service-name>
gh pr create \
  -R abyssparanoia/rapid-go \
  --title "Sync .claude/ config from <service-name>" \
  --body "$(cat <<'EOF'
## Proposed Changes

- Pulled .claude/ content updates contributed by <service-name>

## Implementation

Files synced:
<list pushed files>

EOF
)"
```

Only open a PR on a side that actually received file changes. If one side has nothing to receive, skip that PR.

## Step 6 — Done

Report the PR URLs created and remind the user:

> After the rapid-go PR merges, other projects can run `/sync-claude-config` (with rapid-go added via `--add-dir`) to pull in the newly landed changes.
