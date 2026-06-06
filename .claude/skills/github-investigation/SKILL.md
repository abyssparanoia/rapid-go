---
name: github-investigation
description: Investigate GitHub Actions logs, PR status, Issues, and other GitHub resources. Uses a 3-tier fallback: MCP → CLI → account switch. Use when: "/github-investigate", "check GHA logs", "investigate workflow", "look into PR status", or any GitHub resource investigation.
---

# GitHub Investigation Skill

A skill for investigating resources on GitHub (Actions logs, PRs, Issues, commits, etc.).
Uses a 3-tier fallback pattern that automatically resolves authentication issues.

## Fallback Strategy

```
MCP GitHub tools (authenticated)
  ↓ if the target API is not available in MCP
gh CLI
  ↓ on 404/auth errors
gh auth switch → retry with gh CLI
  ↓ still failing
Report the situation to the user
```

## Workflow

### Step 1: Attempt investigation with MCP tools

Select the appropriate MCP tool based on the target resource:

| Target | MCP Tool |
|--------|----------|
| PR details | `mcp__github__pull_request_read` (method: get) |
| PR status | `mcp__github__pull_request_read` (method: get_status) |
| PR diff | `mcp__github__pull_request_read` (method: get_diff) |
| PR comments | `mcp__github__pull_request_read` (method: get_comments) |
| PR review comments | `mcp__github__pull_request_read` (method: get_review_comments) |
| Issue | `mcp__github__issue_read` |
| Commit | `mcp__github__get_commit` |
| File contents | `mcp__github__get_file_contents` |
| PR search | `mcp__github__search_pull_requests` |
| Issue search | `mcp__github__search_issues` |

**Note**: MCP GitHub tools do not include tools for the Actions API (workflow runs, logs).
For Actions-related investigation, proceed directly to Step 2.

### Step 2: Investigate with gh CLI

Use gh CLI when there is no applicable MCP tool, or when MCP does not provide sufficient information.

**Key commands for Actions:**
```bash
# List workflow runs
gh run list --repo {owner}/{repo} --limit 5

# View a specific run
gh run view {run_id} --repo {owner}/{repo}

# View logs for failed steps
gh run view {run_id} --repo {owner}/{repo} --log-failed

# Re-run a workflow
gh run rerun {run_id} --repo {owner}/{repo} --failed
```

**Other investigation commands:**
```bash
# PR check status
gh pr checks {pr_number} --repo {owner}/{repo}

# Direct API call (any endpoint)
gh api repos/{owner}/{repo}/actions/runs/{run_id}/jobs --jq '.jobs[] | {name, status, conclusion}'
```

### Step 3: Switch accounts on authentication errors

If gh CLI returns `404 Not Found` or `401 Unauthorized`:

```bash
# 1. Check current authentication state
gh auth status

# 2. List available accounts and switch
gh auth switch

# 3. Retry access
gh run view {run_id} --repo {owner}/{repo}
```

**If still failing after switching:**
- Suggest running `gh auth login`
- Verify required scopes (`repo`, `workflow`)

### Step 4: Report results

Report investigation findings in a structured way:

- **On success**: Present the status of the investigated resource, any error details, and root cause analysis
- **On Actions failure**: Present the failed step name, error message, and suggested fix
- **On account switch**: Clearly indicate which account succeeded

## Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| url | GitHub URL (automatically parses owner/repo/resource) | `https://github.com/your-org/your-repo/actions/runs/123` |
| owner | Repository owner | `your-org` |
| repo | Repository name | `claude-blueprints` |
| target | Type of resource to investigate | `run`, `pr`, `issue`, `commit` |
| id | ID of the target resource | `21973188790`, `2`, `#15` |

## Invocation Examples

- "Check the GHA logs" + URL → Fetch Actions logs and analyze the failure cause
- "Check the status of PR #2" → Retrieve PR details and status via MCP
- "Re-run this workflow" → `gh run rerun`
- "Check recent workflow runs" → `gh run list`

## Error Handling

| Error | Action |
|-------|--------|
| MCP tool not supported | Fall back to gh CLI |
| gh CLI 404 | Switch to another account with `gh auth switch` |
| 404 on all accounts | Report that the user lacks access to the repository |
| gh not installed | Suggest installation steps |
| No Actions API permission | Report that the token needs the `workflow` scope |
