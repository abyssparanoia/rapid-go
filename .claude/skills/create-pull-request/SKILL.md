---
name: create-pull-request
description: Create pull requests following project conventions. Use when: (1) creating a PR with gh pr create, (2) writing commit messages, (3) need branch naming guidance. Covers branch naming patterns, PR title/body format, pre-commit checks, MCP-first automation, and draft PR support.
---

# PR Creation Quick Reference

## Branch Naming

| Type | Pattern | Example |
|------|---------|---------|
| Feature | `feature/{description}` | `feature/impl-admin-example-api` |
| Bugfix | `fix/{description}` | `fix/example-creation-error` |
| Chore | `chore/{description}` | `chore/update-dependencies` |

## PR Title

Use English. Be concise and descriptive.

```
Implement AdminService.GetExample
Add validation for Example create API
Fix null pointer in tenant lookup
```

If the branch name contains a ticket identifier (e.g. `feature/TICKET-1234-description`), prefix the title with `[TICKET-1234]`:

```
[TICKET-1234] Implement AdminService.GetExample
```

If no ticket identifier is found in the branch name, omit the prefix entirely.

## PR Body Template

```markdown
## Proposed Changes

- Brief description of what was changed

## Implementation

{Verification details - curl commands, CLI output, or file changes}
```

## Implementation Section Examples

### API Endpoint

```markdown
## Implementation

### POST /admin/v1/tenants/{tenant_id}/examples

curl -X POST "http://localhost:8080/admin/v1/tenants/xxx/examples" \
    -H "Authorization: Bearer {token}" \
    -H "Content-Type: application/json" \
    -d '{"name": "Test", "description": "Test desc"}'

Response:
{"example": {"id": "xxx", "name": "Test", "created_at": "2025-01-01T00:00:00Z"}}
```

### Non-API Changes

```markdown
## Implementation

- Updated `internal/domain/model/example.go` with new validation
- Added unit tests in `example_test.go`
- Regenerated mocks with `make generate.mock`
```

## Pre-PR Checklist

```bash
make lint.go     # Must pass
make test        # Must pass
```

If applicable:
- `make migrate.up` - after DB changes
- `make generate.buf` - after proto changes
- `make generate.mock` - after repository interface changes

## Automated PR Creation Workflow

When asked to create a PR, follow these steps in order.

### Step 0: Analyze current state (run in parallel)

```bash
git status                          # Check for uncommitted changes
git branch --show-current           # Current branch name
git log main..HEAD --oneline        # Commits to include in the PR
git diff main...HEAD --stat         # Summary of changed files
```

### Step 1: Branch and push

- If currently on `main` or `master`: offer to create a feature branch first.
- If there are unpushed commits: run `git push -u origin <branch>`.
- If already pushed and up to date: no action needed.

### Step 2: Detect the GitHub MCP tool

```
ToolSearch("select:mcp__github__create_pull_request")
```

- If the schema is returned: use the **MCP path** below.
- If the tool is not found: use the **gh CLI path** below.

### Step 3: Check for an existing PR

Before creating, check whether a PR already exists for the current branch:

- **MCP**: use `mcp__github__list_pull_requests` with a `head` filter.
- **CLI**: `gh pr list --head $(git branch --show-current)`

If a PR already exists, report its URL and ask the user whether to update it or proceed with a new one.

### Step 4: Generate PR metadata

**Title**: derive from commit messages or use a user-supplied title (English). Apply a ticket prefix if the branch name contains one (see PR Title section above).

**Body**: use the PR Body Template above. Fill in `## Proposed Changes` and `## Implementation` based on the diff.

**Draft**: if the user requested a draft PR (e.g. "create a draft PR", "make it a draft"), set draft mode.

### Step 5: Create the PR

**MCP path**:

```
mcp__github__create_pull_request(
  owner: "{org}",
  repo: "{repo}",
  title: "<generated title>",
  head: "<branch>",
  base: "main",
  body: "<generated body>",
  draft: <true|false>
)
```

**gh CLI path**:

```bash
gh pr create --title "<title>" --body "$(cat <<'EOF'
## Proposed Changes

- ...

## Implementation

{Details}
EOF
)" [--draft]
```

### Step 6: Report the result

Display:
- PR URL
- PR number
- Title
- Branch: `<head> → <base>`

## Error Handling

| Error | Action |
|-------|--------|
| GitHub MCP tool unavailable | Fall back to `gh pr create` |
| `gh` not authenticated | Suggest running `gh auth login` |
| No diff from base branch | Abort with a clear message |
| Uncommitted changes present | Ask the user to commit or stash first |
| PR already exists for branch | Report the URL and ask whether to update |

## Important Notes

- Do NOT add AI attribution (no "Generated with Claude Code" or "Co-Authored-By: Claude")
- Keep commit messages and PR descriptions clean
- Run `self-review` skill before creating PR to catch issues early
