---
description: Multi-agent orchestration and context budget policy (repository-agnostic, always active)
---

# Multi-Agent Orchestration Policy

> **Design principle: think globally, read locally.**
> The main agent owns goals, dependencies, and sequencing. Bulk reading — repository
> exploration, external research, large diffs, multiple big files — happens in isolated
> subagents that return compressed summaries, never raw dumps.

This policy is repository-agnostic and portable: it assumes no specific language,
framework, directory layout, or PR workflow. Repository-specific rules (CLAUDE.md,
`.claude/rules/`, skills) always take precedence for *what* to build; this policy only
governs *how* work is distributed.

## Task Triage (do this first, briefly)

At the start of each task, classify context risk. Do not announce the classification;
just act on it.

| Risk | Signals | Strategy |
|------|---------|----------|
| **Low** | 1–2 small files, clear local fix/rename/typo, no repo-wide exploration needed | Main agent handles directly. **No subagents.** |
| **Medium** | Several files, unfamiliar area, architecture/pattern lookup needed, external docs referenced, moderate diff | Delegate exploration/research; main agent decides and edits (or one focused-editor) |
| **High** | Broad repo exploration, cross-repository work, large PR/diff, many or huge files, long logs, architectural redesign, repo-wide review | Fully orchestrate: parallel research → central decision → split editors → isolated review |

Never multi-agent a task that is cheaper to just do. Delegation has overhead; use it to
protect the main context, not for ceremony.

## Generic Subagents & Model Routing

| Agent | Model | Use for |
|-------|-------|---------|
| `repo-explorer` | Sonnet | Finding relevant files/symbols/patterns in the repository; returns a summary map, not file bodies |
| `researcher` | Sonnet | External docs, API/library specs, referenced PRs/issues; returns facts/constraints, not raw pages |
| `deep-reasoner` | Opus | Architecture, domain modeling, hard design decisions, tradeoff/gap analysis — fed with summaries from explorers, not raw repo access from scratch |
| `focused-editor` | Sonnet | Editing one file (or one small independent file group) against decisions the main agent already made |
| `diff-reviewer` | Opus | Final review of the current diff: correctness, consistency, requirement coverage |

Model routing intent: the **main agent (Fable/top-tier)** orchestrates, decomposes,
integrates, makes final decisions, does small clear edits, and runs git/workflow
operations. **Sonnet** does volume work (reading, searching, mechanical edits,
summarizing). **Opus** is reserved for genuinely hard reasoning and final review —
never for exploration or mechanical changes.

If a repository defines more specific agents (e.g. specialized reviewers), prefer those
for their domains; these generic agents are the fallback and the context-isolation layer.

## Main Agent Context Budget Rules

- Never read the repository broadly "to understand it first". Search, then read only the
  narrowed targets. Delegate multi-file investigation to `repo-explorer`.
- When a subagent returned an adequate summary, do **not** re-read the same files in full.
- If a change spans 3+ large files, split editing across `focused-editor` agents
  (one per file/group) instead of holding all files in the main context.
- Never repeatedly load full PR diffs or huge `git diff` output. Escalate gradually:
  `git status` → `git diff --stat` → per-file diff only where needed.
- For huge test/build/lint logs, extract only the failing sections (grep/tail), or have a
  subagent run and summarize them.
- Final review of complex changes goes to `diff-reviewer`, not to the main agent
  re-reading every changed file.

## Subagent Contract

Every delegated task prompt must state: the goal, the scope boundary (what NOT to
explore), what is already known (prior summaries — don't re-investigate), and the
expected return shape. Subagents must:

- stay inside their assigned scope; search before reading; read only relevant
  sections/symbols/line ranges of large files;
- return **conclusions, evidence, and key paths** (roughly 5–15 bullets), never large raw
  file contents or tool output.

## Delegation & Parallelization

- Select only the roles the task needs — never invoke the full pipeline by default.
- Run independent investigations **in parallel** (e.g. repo exploration + external
  research in one batch).
- Do **not** parallelize: edits to the same file, duplicate investigation of the same
  question, or implementation before a pending design decision.
- Canonical flow for complex work:
  `parallel research → central decision (main or deep-reasoner) → independent edits → centralized review`

## Editing & Final Review

- Give each `focused-editor` the settled requirements/design decisions and its file scope;
  it returns changed sections, decisions reflected, and remaining concerns.
- Verify change scope via `git status` / `git diff --stat` first. For complex changes,
  hand the diff to `diff-reviewer`.
- Route Critical/Major review findings back to the responsible `focused-editor` (or fix
  directly if trivial). The main agent must not re-read the whole changeset to fix issues.

## Examples

- *"Fix this typo in README"* → Low risk. Main edits directly. No subagents.
- *"Fix this bug"* (location unknown) → Medium. `repo-explorer` locates cause/patterns →
  main decides fix → main or one `focused-editor` applies → review only if warranted.
- *"Port this PR to another repo"* → High. Parallel: source analysis + target-repo
  exploration (Sonnet) → `deep-reasoner` maps the adaptation (Opus) → `focused-editor`
  per file → `diff-reviewer`.
- *"Update the domain modeling docs"* → High. Parallel: existing-model exploration +
  external/domain research → `deep-reasoner` settles the model → one editor per document →
  `diff-reviewer`.
- *"Review this design"* → Medium/High. `repo-explorer` gathers current state →
  `deep-reasoner` analyzes → main synthesizes the verdict.
