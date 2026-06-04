# Sync Algorithm Reference

Detailed specification for the classification and conflict-resolution algorithm used in `sync-claude-config`.

## Exclusion List

Never include these in the sync candidate set regardless of which side they appear on:

```
.claude/settings.local.json
.claude/worktrees/           (entire directory)
.claude/skills/init-new-repository/   (entire skill directory — intentionally template-bearing)
```

## In-Scope Paths

Enumerate recursively under these roots relative to the repo root:

```
.claude/CLAUDE.md                 (single file)
.claude/rules/                    (all .md files)
.claude/skills/                   (all files, all subdirs, EXCEPT init-new-repository/)
.claude/agents/                   (all .md files)
.claude/commands/                 (all .md files)
```

Include all file types (not just `.md`) under `skills/` — skills may contain Python scripts,
JSON, YAML, etc. that should also sync.

## Classification Algorithm

For each path `p` in the union of both sides:

```
exists_upstream = file exists at <UPSTREAM>/.claude/<p>
exists_local    = file exists at <LOCAL>/.claude/<p>

if exists_upstream and not exists_local:
    → PULL  (new file from upstream, localize tokens, write into LOCAL)

if exists_local and not exists_upstream:
    → PUSH  (new file from local, canonicalize tokens, write into UPSTREAM)

if exists_upstream and exists_local:
    canonical_upstream = read(<UPSTREAM>/.claude/<p>)
    canonical_local    = apply_reverse_map(read(<LOCAL>/.claude/<p>))

    if canonical_upstream == canonical_local:
        → SKIP (identical)
    else:
        diff = compute_diff(canonical_upstream, canonical_local)

        if is_db_variant_only(diff):
            → SKIP — DB-specific

        → CONFLICT: resolve per rules below
```

## DB-Variant Detection

A diff is "DB-variant only" if every changed line (additions and removals) satisfies at
least one of:

- The line contains the literal string `mysql` or `postgresql` (case-insensitive)
- The line is a quoting-style difference introduced by MySQL (backtick) vs PostgreSQL
  (double-quote) identifier quoting in SQL/Go ORDER BY strings

If ANY changed line does NOT satisfy this condition, the diff is NOT DB-variant only.

## Conflict Resolution

When a file exists on both sides and differs after normalization:

1. Run timestamp check:
   ```bash
   git -C <UPSTREAM> log -1 --format="%ct %h %s" -- .claude/<relpath>
   git -C <LOCAL>    log -1 --format="%ct %h %s" -- .claude/<relpath>
   ```
   Note: timestamp is unix epoch (seconds). Newer = more recently committed = likely the
   authoritative version. If a file has never been committed on one side (`git log` returns
   empty), treat that side as older.

2. Present to the user:
   ```
   CONFLICT: .claude/rules/example-rule.md
   ─────────────────────────────────────────
   upstream last commit: 2025-05-28 (abc1234) "update grpc handler rule"
   local    last commit: 2025-06-03 (def5678) "add new validation pattern"

   Normalized diff (upstream ← → local):
   --- upstream
   +++ local
   @@ -12,4 +12,8 @@
    existing line
   +new validation pattern added in local
   ...

   Which version should be canonical?
   [U] upstream  [L] local  [S] skip this file
   ```

3. Based on user choice:
   - **U (upstream)**: Add file to PULL plan (overwrite local with localized upstream version)
   - **L (local)**: Add file to PUSH plan (overwrite upstream with canonicalized local version)
   - **S (skip)**: Do not include in either PR

## Token Replacement Implementation

When applying the reverse map (LOCAL → canonical) or forward map (canonical → LOCAL):

1. Substitute literal strings in the file content (not regex — plain string replacement)
2. Apply in the order specified in `token-mapping.md` (longest/most-specific first)
3. Both forward and reverse maps are applied to the **entire file content** as a string
4. Binary files are skipped (apply only to UTF-8 text files)

To detect if a file is binary: if it contains a null byte (`\x00`), treat as binary and skip.

## Summary Table Format

Present this table before any changes are applied:

```
══════════════════════════════════════════════════════════
PULL → LOCAL  (these files will be added/updated locally)
══════════════════════════════════════════════════════════
  [NEW]      .claude/skills/new-skill/SKILL.md
  [NEW]      .claude/rules/new-convention.md
  [CONFLICT] .claude/rules/grpc-handler.md        (upstream newer: 2025-06-01)

══════════════════════════════════════════════════════════
PUSH → UPSTREAM  (these files will be added/updated in rapid-go)
══════════════════════════════════════════════════════════
  [NEW]      .claude/rules/local-only-rule.md
  [CONFLICT] .claude/rules/testing.md             (local newer: 2025-06-03)

══════════════════════════════════════════════════════════
SKIPPED
══════════════════════════════════════════════════════════
  [IDENTICAL]   14 files (no change after normalization)
  [DB-SPECIFIC]  2 files (db-variant differences only)
  [EXCLUDED]     1 file  (.claude/skills/init-new-repository/)
══════════════════════════════════════════════════════════

Proceed? (Y to continue, N to cancel)
```

Conflicts marked with `(local newer:...)` or `(upstream newer:...)` should already be
resolved before this table is shown (ask per-file in Step 2 before building the table).

## Edge Cases

**File renamed in one repo**: Will appear as NEW on one side + absent on the other. No
automatic rename detection — it will create a PR that adds the new name on one side. The
old name stays if present. This is acceptable for the current implementation.

**Empty file**: Treat as a valid file. An empty file differs from an absent file.

**Skill with both SKILL.md and supporting files (scripts, references/)**: All files under
the skill directory are in scope together. If a skill is new on one side, all its files are
pulled/pushed as a group.

**Circular sync**: If the same change is already present (identical after normalization),
it is classified as SKIP. Running the skill after PRs merge results in no changes (idempotent).

**No changes detected**: Report "No sync needed — all in-scope files are identical after
normalization" and exit without creating branches or PRs.
