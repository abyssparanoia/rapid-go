---
name: sync-rules-from-pr
description: Analyze GitHub PR review comments and propose updates to .claude/rules/. Use when: (1) running '/sync-rules-from-pr', (2) asked to reflect reviewer feedback into project rules or conventions, (3) a PR review surfaced a pattern worth codifying. Accepts no argument (PR for current branch), a PR number, or a single review comment URL.
argument-hint: "[pr-number | comment-url]"
---

# Sync Claude Rules from PR Review Comments

Turn human review feedback into durable `.claude/rules/` updates. Proposals only — never edit a rule before the user approves.

## Step 1: Resolve the Target

| Argument | Target |
|----------|--------|
| (none) | PR of the current branch: `gh pr view --json number` |
| `123` | PR #123 |
| `https://github.com/{owner}/{repo}/pull/123#discussion_r456` | Only inline comment `456` of PR #123 |
| `https://github.com/{owner}/{repo}/pull/123#pullrequestreview-789` | Only review `789` of PR #123 |

```bash
gh repo view --json owner,name | jq -r '"\(.owner.login) \(.name)"'
```

## Step 2: Fetch Comments

```bash
# Inline review comments (path, line, body, id)
gh api "repos/{owner}/{repo}/pulls/{number}/comments" --paginate \
  | jq '.[] | {id, path, line, body, user: .user.login}'

# Review summaries (body text of each review)
gh api "repos/{owner}/{repo}/pulls/{number}/reviews" --paginate \
  | jq '.[] | select(.body != "") | {id, state, body, user: .user.login}'
```

When a comment/review id was given, keep only that entry. Skip bot comments (`user` ending in `[bot]`).

## Step 3: Map Each Comment to a Rule

Match the comment's `path` against the `paths:` frontmatter of each `.claude/rules/*.md` (`grep -A8 '^paths:' .claude/rules/*.md`). Review-level comments without a path: infer the layer from the comment text. The rule table in `.claude/CLAUDE.md` is the fallback map.

Read only the matched rule files, not all of them.

## Step 4: Decide Whether a Rule Update Is Warranted

**Propose an update when the comment is generalizable:**
- "Unify this pattern across the codebase" → add to the rule
- "Don't do X, do Y instead" → add as an anti-pattern (and consider `.claude/skills/review-diff/references/ai-antipatterns.md`)
- "Always do X in this situation" → add as a best practice
- Naming, ordering, layering, or transaction-boundary feedback

**Do not propose when:**
- Typo / wording fix
- Feature-specific business logic
- Temporary workaround
- Already covered by the matched rule (say so)

## Step 5: Present Proposals

```markdown
## PR Review Analysis — #{number} {title}

### Proposal 1: `.claude/rules/{file}.md` › {section}
**Source comment** ({user}, `{path}:{line}`):
> {comment body}

**Proposed change:**
{concrete text / code block to add or modify}

**Rationale:** {why this generalizes}

---

### Not Proposed
| Comment | Reason |
|---------|--------|
| {summary} | {reason} |
```

## Step 6: Apply After Approval

Only after the user approves each proposal:
1. Edit the rule file, matching its existing language (English) and section style
2. If an anti-pattern was added, append it to `ai-antipatterns.md` with the next number and update `checklists.md` if a checklist item applies
3. Suggest `/sync-claude-config` if this project is derived from rapid-go, so the change reaches upstream
