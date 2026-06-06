---
name: meta-review
description: Analyze PR review comments to identify improvement opportunities for rules and skills. Use when: "/meta-review", "improve from review feedback", "analyze review comments", "extract patterns from PR reviews", or any request to turn review feedback into system improvements.
---

# Meta-Review Skill

Analyze PR review comments to identify improvement opportunities for rules, skills, and agents.

## Prerequisites

- GitHub MCP tools or `gh` CLI must be available
- Access to the repository's PR history is required

## Workflow

### Step 1: Fetch PR review comments

Prefer MCP; fall back to gh CLI when unavailable:

**MCP approach:**
```
ToolSearch("select:mcp__github__pull_request_read")
```
Use `mcp__github__pull_request_read` to retrieve PR details and comments.

**gh CLI fallback:**
```bash
# Fetch recently merged PRs
gh pr list --state merged --limit 10 --json number,title,url

# Fetch review comments for a specific PR
gh api repos/{owner}/{repo}/pulls/{pr_number}/comments --jq '.[] | {body, path, line, created_at}'

# Fetch issue comments (general PR comments)
gh api repos/{owner}/{repo}/issues/{pr_number}/comments --jq '.[] | {body, created_at}'
```

If a specific PR number is provided, use that. Otherwise, analyze the most recent 5 merged PRs.

### Step 2: Classify comments

Classify each review comment into one of the following categories:

| Category | Code | Description | Action |
|----------|------|-------------|--------|
| Missing rule | `[missing-rule]` | Feedback about a pattern not covered by existing rules | Create a new rule |
| Pattern violation | `[pattern-violation]` | Cases where a violation of an existing rule went undetected | Improve rule visibility (adjust globs) |
| Skill gap | `[skill-gap]` | Workflow problems or missing procedural guidance | Update or create a skill |
| False positive | `[false-positive]` | Rule is too strict or inaccurate | Relax or correct the rule |
| Irrelevant | `[irrelevant]` | Not related to the meta-system (business logic, typos, etc.) | Skip |

**Classification criteria:**
- If feedback is related to a coding pattern → search existing rules with `Grep`
- If feedback is related to a workflow → search existing skills
- If no match → classify as `[missing-rule]` or `[skill-gap]`

### Step 3: Prioritize

Count occurrences across PRs:

| Priority | Threshold | Action |
|----------|-----------|--------|
| High | 3 or more similar comments across PRs | Recommend immediate action |
| Medium | 2 similar comments | Address in the next cycle |
| Low | Only 1 occurrence | Monitor only, no action required |

Group similar comments by theme (e.g., "error handling", "naming conventions", "test patterns").

### Step 4: Generate an improvement report

Output a structured report:

```markdown
## Meta-Review Report

**PRs analyzed**: #{pr1}, #{pr2}, ...
**Total comments**: {N}
**Actionable**: {N} ({missing-rule}: {n}, {pattern-violation}: {n}, {skill-gap}: {n}, {false-positive}: {n})
**Skipped**: {N} (irrelevant)

### High Priority

| # | Category | Theme | Occurrences | Recommended Action |
|---|----------|-------|-------------|-------------------|
| 1 | [missing-rule] | {theme} | {count} | Create rule: {name} |

### Medium Priority

| # | Category | Theme | Occurrences | Recommended Action |
|---|----------|-------|-------------|-------------------|

### Low Priority (monitoring)

| # | Category | Theme | Occurrences | Recommended Action |
|---|----------|-------|-------------|-------------------|

### Raw Comments

<details>
<summary>Comment details</summary>

- PR #{number}: "{comment excerpt}" → [{category}]
- ...
</details>
```

### Step 5: Suggest next actions

Based on the report:
- **High priority items**: Recommend immediate creation/update using the `meta-crud` skill
- **Medium priority**: Add to the monitoring list
- **Low priority**: Record for future reference

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| PR comments not found | No merged PRs or reviews are empty | Expand the search scope or verify repository access |
| GitHub API rate limit | Too many API calls | Wait and retry, or reduce the number of PRs |
| MCP tools unavailable | MCP server not configured | Fall back to `gh` CLI |
| Cannot classify comment | Ambiguous feedback | Default to `[irrelevant]` and record for manual review |

## Execution Example

User: "Analyze review feedback from recent PRs"

```
1. gh pr list --state merged --limit 5 → PRs #42, #41, #39, #38, #35
2. Fetch comments for each PR
3. Classify:
   - "error handling is missing" x3 → [missing-rule], high priority
   - "no tests" x2 → [skill-gap], medium priority
   - "typo" x1 → [irrelevant], skip
4. Output report and present recommended actions
5. Suggestions: create error handling rule, update testing skill
```
