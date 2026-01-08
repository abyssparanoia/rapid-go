---
description: Analyze PR review comments and propose updates to Claude Rules
---

# Sync Claude Rules from PR Review Comments

## Context

- Current branch: !`git branch --show-current`
- Repository: abyssparanoia/rapid-go
- Arguments: $ARGUMENTS

## Task

Analyze PR review comments and determine whether `.claude/rules/` files should be updated.

### Step 1: Parse Arguments and Find Target Comments

**If a PR comment URL is provided in arguments:**

Parse the URL to extract PR number and comment ID:
- Format: `https://github.com/abyssparanoia/rapid-go/pull/{pr_number}#discussion_r{comment_id}`
- Or: `https://github.com/abyssparanoia/rapid-go/pull/{pr_number}#pullrequestreview-{review_id}`

Then:
1. Use `mcp__github__get_pull_request_comments` with the PR number
2. Filter to find the specific comment matching the comment_id from the URL
3. Only analyze that specific comment

**If no argument is provided:**

1. Use GitHub MCP `mcp__github__list_pull_requests` to find the PR for current branch:
   - owner: "abyssparanoia"
   - repo: "rapid-go"
   - head: "abyssparanoia:{current_branch_name}"
   - state: "open"
   - If not found, also try with state: "all"
2. Fetch ALL review comments for that PR

### Step 2: Fetch Review Comments

Use `mcp__github__get_pull_request_comments` to get inline review comments.
Also use `mcp__github__get_pull_request_reviews` to get review body text (unless targeting a specific comment).

### Step 3: Load Claude Rules

Read relevant rule files from `.claude/rules/`:

| File | Domain |
|------|--------|
| `domain-model.md` | Domain model design patterns |
| `domain-service.md` | Domain service patterns |
| `domain-errors.md` | Error definition conventions |
| `repository.md` | Repository and marshaller patterns |
| `usecase-interactor.md` | Interactor implementation, external service sync |
| `grpc-handler.md` | gRPC handler patterns |
| `testing.md` | Testing conventions |
| `proto-definition.md` | Protocol Buffers style |
| `migration.md` | Database migration patterns |
| `dependency-injection.md` | DI configuration |
| `invitation-workflow.md` | Invitation/approval flow patterns |
| `external-service-integration.md` | Cognito/IdP integration patterns |

### Step 4: Analyze and Propose

Analyze review comments with these criteria:

**Should update rules when:**
- "Please unify this pattern across the codebase" → Add to relevant rule
- "Avoid this approach, instead use..." → Add as anti-pattern
- "Always do X in this case" → Add as best practice
- Design pattern or naming convention feedback
- Coding standard feedback

**Should NOT update rules when:**
- Simple typo or spelling fixes
- Business logic specific to this feature (not generalizable)
- Temporary fixes or workarounds
- Context-specific implementation details

### Step 5: Output Format

Output results in this format:

```markdown
## PR Review Analysis

### PR Information
- PR: #{number} - {title}
- Branch: {branch_name}
- Review comments: {count}

### Proposed Rule Updates

#### Proposal 1: {rule_file_name}

**Related review comment:**
> {comment_body}

**Proposed changes:**
{proposed_changes}

**Rationale:**
{justification}

---

### Comments Not Requiring Rule Updates

| Comment Summary | Reason |
|-----------------|--------|
| {summary} | {reason} |
```

### Step 6: User Confirmation and Update

Present proposals to user and only update rule files after approval.
When updating, match the existing language (English) of the rule files.
