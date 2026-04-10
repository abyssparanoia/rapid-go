---
name: fix-review-comments
description: Fetch unresolved GitHub PR review comments and automatically fix the code. Use when: (1) running '/fix-review-comments', (2) asked to "fix review comments", "address PR comments", "respond to review", or similar. Can target the PR for the current branch, or a specific PR number (e.g. '/fix-review-comments 123'). Skips already-resolved threads and discussion-only comments that require no code change.
---

# Fix Review Comments

Fetch unresolved PR review threads and fix the code automatically.

## Step 1: Detect PR

```bash
# From current branch (no argument)
gh pr view --json number,url | jq -r '"#\(.number) \(.url)"'

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

## Step 5: Verify

```bash
/usr/bin/make lint.go
```

Run tests if logic changed:
```bash
/usr/bin/make test
```

## Step 6: Report

```
## Fix Review Comments Summary

PR: #NNN

### Fixed (N)
1. `path/to/file.go:line` — short description of fix
   Comment: "@author: ..."

### Skipped (N)
1. `path/to/file.go:line` — reason (already resolved / discussion / no code change needed)
   Comment: "@author: ..."

### Verification
- [x] lint.go passes
- [x] tests pass
```

## Notes

- Do **not** auto-commit — leave that to the user
- Do **not** resolve threads via the API — the user confirms fixes via GitHub UI
- If a thread's `isOutdated == true` and the problem no longer exists in current code, note it in the Skipped section with reason "already fixed (outdated thread)"
