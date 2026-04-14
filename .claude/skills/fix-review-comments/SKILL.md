---
name: fix-review-comments
description: Fetch unresolved GitHub PR review comments and automatically fix the code, commit, and push. Use when: (1) running '/fix-review-comments', (2) asked to "fix review comments", "address PR comments", "respond to review", or similar. Can target the PR for the current branch, or a specific PR number (e.g. '/fix-review-comments 123'). Skips already-resolved threads and discussion-only comments that require no code change. Also updates Claude rules and review-diff anti-patterns when a comment reveals a missing convention.
---

# Fix Review Comments

Fetch unresolved PR review threads, fix code, update rules if needed, then commit and push.

## Step 1: Detect PR

```bash
# From current branch (no argument)
gh pr view --json number,url,headRefName | jq -r '"#\(.number) \(.url) \(.headRefName)"'

# Or use the number passed as argument
```

Extract `owner` and `repo` from remote:
```bash
gh repo view --json owner,name | jq -r '"\(.owner.login) \(.name)"'
```

## Step 2: Fetch Unresolved Review Threads (GraphQL)

```bash
gh api graphql -f query='
query($owner: String!, $repo: String!, $number: Int!) {
  repository(owner: $owner, name: $repo) {
    pullRequest(number: $number) {
      reviewThreads(first: 100) {
        nodes {
          id
          isResolved
          isOutdated
          path
          line
          startLine
          comments(first: 20) {
            nodes {
              body
              author { login }
              diffHunk
              createdAt
            }
          }
        }
      }
    }
  }
}' -f owner=OWNER -f repo=REPO -F number=NUMBER
```

Filter to keep only threads where `isResolved == false`.

## Step 3: Classify Comments

For each unresolved thread, read the **first comment's `body`** (the original review comment). Classify:

| Type | Action |
|------|--------|
| Requires code change | Fix |
| Question / acknowledgment / "LGTM" / praise | Skip |
| Discussion / "let's discuss" / "not blocking" | Skip |
| Nit / style suggestion (non-blocking) | Fix unless explicitly marked optional |

**Use judgment**: a comment like "this looks a bit confusing, what do you think?" is discussion. A comment like "this should use `nullable.Type[T]` instead of a pointer" requires a fix.

## Step 4: Read Files and Apply Fixes

For each thread requiring a fix:

1. Read the file at `thread.path`
2. Use `thread.line` (or `startLine`..`line` for multi-line) to locate the relevant code
3. Use `diffHunk` from the first comment for extra context about the surrounding code
4. Apply the fix based on the comment instruction
5. Note: if `isOutdated == true`, the code may have already changed — verify the issue still exists before fixing

Fix all actionable threads before running any verification.

## Step 5: Update Rules (if applicable)

After fixing code, evaluate whether any comment reveals a **missing or incomplete coding convention** that should be captured in project rules. Ask:

- Could this mistake recur if the rule isn't documented?
- Is it a pattern-level issue (not just a one-off typo)?

If yes, update the appropriate file(s):

| Comment type | Target file |
|---|---|
| AI-generated code pattern mistake | `.claude/skills/review-diff/references/ai-antipatterns.md` — add new numbered pattern, renumber subsequent |
| Domain model / entity convention | `.claude/rules/domain-model.md` |
| Repository / marshaller pattern | `.claude/rules/repository.md` |
| Usecase / interactor convention | `.claude/rules/usecase-interactor.md` |
| gRPC handler pattern | `.claude/rules/grpc-handler.md` |
| Testing convention | `.claude/rules/testing.md` |
| Proto definition style | `.claude/rules/proto-definition.md` |
| Other layer-specific rule | Corresponding `.claude/rules/*.md` |

When adding to `ai-antipatterns.md`:
1. Assign the next sequential number
2. Include BAD/GOOD code examples from the actual review comment
3. Renumber all subsequent patterns
4. Update `.claude/skills/review-diff/SKILL.md` priority list if the pattern is critical

When adding to `.claude/rules/*.md`:
1. Add the rule in the appropriate section
2. Follow the existing format (tables, code blocks, etc.)

**Do NOT update rules for**: one-off typos, project-specific business logic, or comments that are already covered by existing rules.

## Step 6: Verify

```bash
/usr/bin/make lint.go
```

Run tests if logic changed:
```bash
/usr/bin/make test
```

## Step 7: Commit and Push

Stage, commit, and push all changes (code fixes + rule updates):

```bash
git add -A
git commit -m "fix: address PR review comments

- <summary of code fixes>
- <summary of rule updates, if any>"
git push
```

If code fixes and rule updates are logically separate, use two commits:
1. Code fixes first: `fix: address PR #NNN review comments`
2. Rule updates second: `docs: update rules from PR #NNN review feedback`

## Step 8: Report

```
## Fix Review Comments Summary

PR: #NNN

### Fixed (N)
1. `path/to/file.go:line` — short description of fix
   Comment: "@author: ..."

### Rules Updated (N)
1. `.claude/rules/xxx.md` — added rule about ...
   Triggered by: "@author: ..." comment on `path/to/file.go`
2. `.claude/skills/review-diff/references/ai-antipatterns.md` — added #NN: ...

### Skipped (N)
1. `path/to/file.go:line` — reason (already resolved / discussion / no code change needed)
   Comment: "@author: ..."

### Verification
- [x] lint.go passes
- [x] tests pass
- [x] committed and pushed
```

## Notes

- Do **not** resolve threads via the API — the user confirms fixes via GitHub UI
- If a thread's `isOutdated == true` and the problem no longer exists in current code, note it in the Skipped section with reason "already fixed (outdated thread)"
