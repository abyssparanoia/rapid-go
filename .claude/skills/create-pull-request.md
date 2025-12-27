---
name: create-pull-request
description: REQUIRED when creating pull requests or commits. Covers branch naming, PR title/body format, and team conventions. ALWAYS read this before running git commit or gh pr create commands.
---

# Pull Request Creation Guide

This guide covers the conventions for creating pull requests in this repository.

## Branch Naming

| Type | Pattern | Example |
|------|---------|---------|
| Feature | `feature/{description}` | `feature/impl-admin-example-api` |
| Chore/Maintenance | `chore/{description}` | `chore/update-dependencies` |
| Bugfix | `fix/{description}` | `fix/example-creation-error` |

## PR Title

- Use English for this repository
- Be descriptive and concise
- Examples:
  - `Implement AdminService.GetExample`
  - `Add validation for Example create API`
  - `Update dependencies`

## PR Body Template

```markdown
## Proposed Changes

- Brief description of what was changed
- Additional context if needed

## Implementation

{Technical details with verification examples}
```

## Implementation Section Examples

### For API Endpoints

Include curl commands with headers and expected JSON response:

```markdown
## Implementation

### Local Environment

#### POST /admin/v1/tenants/{tenant_id}/examples

```shell
curl -X POST --location "http://localhost:8080/admin/v1/tenants/xxx/examples" \
    -H "Authorization: Bearer {token}" \
    -H "Content-Type: application/json" \
    -d '{
          "name": "Test Example",
          "description": "Test description"
        }'
```

```json
{
  "example": {
    "id": "xxx",
    "name": "Test Example",
    "description": "Test description",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```
```

### For CLI Commands

Include command and output:

```markdown
## Implementation

```sh
go run ./cmd/app task admin create-example --name "Test"
```

```json
{"level":"info","message":"completed task","example_id":"xxx"}
```
```

### For Non-API Changes

Simply describe the implementation details:

```markdown
## Implementation

- Updated `internal/domain/model/example.go` with new validation logic
- Added unit tests in `example_test.go`
- Regenerated mocks with `make generate.mock`
```

## Pre-PR Checklist

Before creating a PR, ensure:

1. **Code Quality**
   ```bash
   make lint.go    # Passes without errors
   make test       # All tests pass
   ```

2. **Generated Code**
   - Run `make migrate.up` if database changes were made
   - Run `make generate.buf` if proto files were modified
   - Run `make generate.mock` if repository interfaces changed

3. **Branch is Up-to-Date**
   ```bash
   git fetch origin main
   git rebase origin/main
   ```

## Creating PR with gh CLI

```bash
# Create and push branch
git checkout -b feature/my-feature
git add .
git commit -m "Implement feature"
git push -u origin feature/my-feature

# Create PR
gh pr create --title "Implement feature" --body "$(cat <<'EOF'
## Proposed Changes

- Implemented new feature

## Implementation

{Details here}
EOF
)"
```

## Notes for AI Assistants

- **Do NOT add Claude watermarks** - No need to append "Generated with Claude Code", "Co-Authored-By: Claude", or similar AI attribution markers to commits or PR descriptions
- Keep commit messages and PR descriptions clean and human-like

## Review Process

1. Assign yourself to the PR
2. Request review from team members
3. Address review comments
4. Squash merge when approved

## Common Patterns from Existing PRs

### Feature Implementation PR

```markdown
## Proposed Changes

- Implemented `AdminTenantInteractor.Get`

## Implementation

```sh
curl --request GET \
  --url http://localhost:8080/admin/v1/tenants/{tenant_id} \
  --header 'Authorization: Bearer {token}'
```

```json
{
  "tenant": {
    "id": "xxx",
    "name": "example tenant"
  }
}
```
```

