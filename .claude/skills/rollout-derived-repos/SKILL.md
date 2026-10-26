---
name: rollout-derived-repos
description: "Fan out rapid-go's current state to every derived repository in one command: Go toolchain + dependency upgrades (as last done in rapid-go) and/or a pull-only .claude/ config sync, one agent per repo in parallel, one PR per workstream per repo, CI verified. Use when: (1) running '/rollout-derived-repos', (2) rapid-go just merged a Go/dependency bump or a .claude/ change that every derived project should receive, (3) asked to 'ばら撒く' / propagate / roll out template changes to all derived repos. Requires a local registry file .claude/derived-repos.local.yaml (gitignored). For a single repo in either direction use sync-claude-config instead."
disable-model-invocation: true
---

# rollout-derived-repos

One command, run **from rapid-go**, that propagates rapid-go's HEAD state to every registered derived repository:

| Workstream | What it does per repo | Branch |
|---|---|---|
| `deps` | Go version, direct dependencies and tool versions brought to rapid-go's current state (plus whatever small code/config fixes the bump needs) | `chore/update-go-<major.minor>-and-deps` |
| `claude` | Pull-only `.claude/` sync from rapid-go (token-mapped, project-specific content preserved) | `chore/sync-claude-config` (suffix `-YYYY-MM` if taken) |

Each repo gets its own agent, running in parallel, working in **git worktrees** so the user's checkouts are never touched. The orchestrator (you, in rapid-go) only surveys, briefs, collects, verifies and reports. The agents hold no decision authority: every rule they follow lives in `references/rollout-brief.md`.

## Usage

```
/rollout-derived-repos                 # both workstreams, all registered repos
/rollout-derived-repos deps            # dependencies only
/rollout-derived-repos claude          # .claude/ sync only
/rollout-derived-repos deps efso-core sunabe-server   # restrict to repos whose slug or path basename matches
```

## Prerequisites

1. **Registry** — `.claude/derived-repos.local.yaml` in rapid-go (gitignored; rapid-go is public and derived repos are not). Format: `references/derived-repos.example.yaml`. If it is missing, stop and ask the user to create it from the example — never guess repositories.
2. **GitHub access** — `gh auth status` must succeed for every org in the registry (agents create branches, PRs and issues with `gh`).
3. **rapid-go state** — run on a clean, up-to-date default branch of rapid-go. Its `go.mod`, `.golangci.yml`, `.github/workflows/ci.yml`, `docker/*.Dockerfile` and `.claude/` ARE the target state; nothing else defines it.
4. **Local clones (optional)** — when a registry entry has `path` and it exists, worktrees are created from that clone (after `git fetch origin`). Otherwise the agent clones with `gh repo clone` into its scratch directory.

## Step 0 — Parse arguments & registry

- Mode: `deps`, `claude`, or `all` (default). Optional repo filters match the registry `slug` (owner/repo) or the `path` basename.
- Read the registry with the Read tool. Skip entries with `enabled: false`. Keep each entry's `notes` — they are passed verbatim to that repo's agent.

## Step 1 — Capture the target state from rapid-go HEAD

Mechanically (grep, never by memory) collect and write a **Target State** block:

```bash
grep -m1 '^go ' go.mod                                           # GO_VERSION
grep -E 'golangci-lint/v2 ' go.mod                               # GOLANGCI_VERSION (tool block)
grep -E 'version: v[0-9]' .github/workflows/ci.yml               # must equal GOLANGCI_VERSION
grep -E 'golang.org/x/tools ' go.mod                             # XTOOLS_VERSION
grep -rhoE 'golang:[0-9.]+-alpine' docker/                       # GOLANG_IMAGE_TAG
git log --oneline -5 -- go.mod                                    # REFERENCE_PRS (the rapid-go bump PRs agents should read)
git log --oneline -5 -- .claude                                   # latest .claude changes (context for workstream claude)
```

Render the brief: copy `references/rollout-brief.md` to the scratchpad directory and replace the `{{...}}` placeholders in its "Target State" section with these values. Agents read the rendered copy, never the template.

## Step 2 — Survey the registered repos (read-only, mechanical)

For each repo collect, with `git -C <path>` / `gh ... -R <slug>` (never `cd` inside a loop — it changes the session's working directory):

- default branch (`gh repo view <slug> --json defaultBranchRef`), current checkout branch and dirtiness (informational only — the agent will not touch it)
- `go` directive, golangci-lint version in go.mod and in CI, x/tools version
- presence of `wrench`, `yo`, `memeduck`, `jwx` major(s), `null` major, `renovate.json`
- `.claude/` shape: counts of rules/skills/agents/commands, `settings.json` present?
- open PRs/issues that overlap (Renovate/Dependabot dependency PRs, previous `chore/sync-claude-config*` branches)

Treat the survey as a hint only: the working copy may be behind `origin/<default>` (in the 2026-10 rollout two repos had already merged part of the work). The brief tells agents to verify against `origin/<default>`.

Persist the survey table and the rendered brief path now (e.g. in your task checkpoint, if your setup keeps one) — a rollout runs for an hour or more and the session may be compacted before the reports arrive.

## Step 3 — Fan out: one agent per repo, in parallel

Launch all agents in ONE message (background; you'll be notified as each finishes). Per agent:

```
Agent(
  description: "<repo-basename> rollout (<mode>)",
  subagent_type: "general-purpose",
  prompt: """
Read the shared brief FIRST and follow it exactly: <absolute path of the rendered brief>

Mode: <deps|claude|all>
Repository: <local path or 'clone from GitHub'> (GitHub <slug>, default branch <default>).
Survey facts: <the one-line survey for this repo: go version, golangci versions, spanner tools, memeduck, jwx/null majors, docker tags, .claude shape, overlapping PRs>.
Registry notes: <notes verbatim, or 'none'>.
Checkout state: <branch/dirty summary> — do NOT touch it; worktrees only.

Do Workstream A then Workstream B (as the mode allows), each in its own worktree as the brief describes. Create the PRs, watch CI until finished, fix what is within the rules, remove the worktrees when done, and return the compact final report described in the brief.
"""
)
```

Notes:
- The agents need judgment (lint fallout, API changes, `.claude` conflict merges), so use a full-capability agent; a `model: opus` override is appropriate when many repos are large.
- Concurrency: on a 64 GB / 10-core machine ten agents at once is fine (builds and golangci-lint runs overlap). On smaller machines launch in waves of ~5.
- Never let two agents share a repo. Never run deps and claude for the same repo as two agents (they would race on `git worktree add`/branch names); one agent does both sequentially.
- While waiting, do nothing else in those repos. Interim "agent is waiting on its own background work" notifications are normal (CI watches).

## Step 4 — Collect, verify, report

As each report arrives, append its essentials to the same persisted notes (PR URLs, CI state, versions, code changes, held-back items + issue URLs, decisions requested, "worth pushing upstream" items). Do not trust the reports blindly — after the last one, verify mechanically:

```bash
# every PR: state + check conclusions
gh pr view <n> -R <slug> --json state,title,statusCheckRollup \
  -q '.state + " " + .title + " " + ([.statusCheckRollup[] | (.conclusion // .state)] | group_by(.) | map("\(.[0])=\(length)") | join(","))'
# no leftover worktrees
git -C <path> worktree list | grep -c scratchpad
```

A PR counts as green only when every check is `SUCCESS` or `SKIPPED` (gated jobs) and the branch is pushed. If a PR is red, send the repo's agent a follow-up message (SendMessage) with the failing job link rather than fixing it yourself.

Final message to the user (Japanese, same shape every time):

1. Table: repo | deps PR | claude PR (links), one line per repo, plus "CI 全 pass" or the exception.
2. Common changes applied (Go, golangci-lint, notable code fixes), one short paragraph.
3. Held-back items with their issue links, grouped by reason (e.g. Spanner tooling, deprecated APIs, module path moves).
4. Decisions for the user (anything agents flagged), as a bullet list.
5. Candidates to pull back into rapid-go (local improvements agents found) — these feed the next rapid-go PR, or `/sync-claude-config` run from that repo.

Close out your persisted notes, and carry the pull-back candidates and any new holds into rapid-go (see "Maintaining the brief") rather than leaving them in the chat.

## Relationship to other skills

- `sync-claude-config` — single repo, bidirectional, interactive confirmation. This skill reuses its token mapping and algorithm (`references/token-mapping.md`, `references/sync-algorithm.md`) in pull-only, non-interactive form, across many repos. Local improvements discovered during a rollout are pushed upstream by running `sync-claude-config` from that repo, not by this skill.
- `create-pull-request` — PR/commit conventions the agents must follow in every repo (no AI attribution).
- `init-new-repository` — defines the forward token map used when localizing `.claude/` files; its script already skips the sync tooling so a derived repo can keep syncing.

## Maintaining the brief

`references/rollout-brief.md` is the single source of truth for what agents may and may not do. When a rollout teaches something new (a dependency that must be held, a config schema rename, a recurring CI quirk), update the brief's **Holds** or **Known fallout** section in rapid-go so the next rollout inherits it — do not leave the knowledge in a chat log.
