# JIRA Command Reference (zsh / bash)

JIRA API command reference for zsh / bash shell users.

## Setting Credentials

### Option A: Add to .envrc (recommended — auto-loaded via direnv)

Add to `.envrc`:

```bash
export JIRA_EMAIL='your-email@example.com'
export JIRA_API_KEY='your-api-token'
```

### Option B: Add to ~/.zsh_private

Add to `~/.zsh_private`:

```bash
export JIRA_EMAIL='your-email@example.com'
export JIRA_API_KEY='your-api-token'
```

## Authentication Prefix

If the environment variables are not set, fall back to loading from `.zsh_private`.
Because the `-u` flag may not handle special characters in tokens correctly, the `Authorization: Basic` header is built using explicit Base64 encoding.

```bash
# Prefer .envrc, fall back to ~/.zsh_private, then build the Base64 auth header
if [ -z "$JIRA_EMAIL" ] || [ -z "$JIRA_API_KEY" ]; then
  if [ -f ~/.zsh_private ]; then
    source ~/.zsh_private
  fi
fi
JIRA_AUTH=$(echo -n "$JIRA_EMAIL:$JIRA_API_KEY" | base64)
```

Prepend this prefix to every command below.

## Commands

### List Projects

```bash
source ~/.zsh_private 2>/dev/null; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  "https://YOUR_ORG.atlassian.net/rest/api/3/project" | jq '[.[] | {key, name}]'
```

### Search Issues (JQL)

```bash
source ~/.zsh_private 2>/dev/null; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X POST -H "Content-Type: application/json" \
  -d '{"jql": "project=SP ORDER BY created DESC", "maxResults": 10, "fields": ["summary", "status", "assignee", "priority"]}' \
  "https://YOUR_ORG.atlassian.net/rest/api/3/search/jql" | jq '.issues[] | {key, summary: .fields.summary, status: .fields.status.name, assignee: .fields.assignee.displayName}'
```

### Get Issue Details

```bash
source ~/.zsh_private 2>/dev/null; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue/ISSUE_KEY" | jq '{key, summary: .fields.summary, status: .fields.status.name, description: .fields.description}'
```

### Create an Issue

```bash
source ~/.zsh_private 2>/dev/null; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X POST -H "Content-Type: application/json" \
  -d '{
    "fields": {
      "project": {"key": "PROJECT_KEY"},
      "summary": "Issue title",
      "description": {"type": "doc", "version": 1, "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Description here"}]}]},
      "issuetype": {"name": "Task"}
    }
  }' \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue" | jq '{key, self}'
```

### Change Status (Transition)

```bash
# Get available transitions
source ~/.zsh_private 2>/dev/null; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue/ISSUE_KEY/transitions" | jq '.transitions[] | {id, name}'
```

```bash
# Execute a transition
source ~/.zsh_private 2>/dev/null; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X POST -H "Content-Type: application/json" \
  -d '{"transition": {"id": "TRANSITION_ID"}}' \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue/ISSUE_KEY/transitions"
```

### Add a Comment

```bash
source ~/.zsh_private 2>/dev/null; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X POST -H "Content-Type: application/json" \
  -d '{"body": {"type": "doc", "version": 1, "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Comment text here"}]}]}}' \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue/ISSUE_KEY/comment" | jq '{id, author: .author.displayName}'
```

### Add / Remove Labels

```bash
# Add a label
source ~/.zsh_private 2>/dev/null; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X PUT -H "Content-Type: application/json" \
  -d '{"update": {"labels": [{"add": "new-label"}]}}' \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue/ISSUE_KEY"
```

```bash
# Remove a label
source ~/.zsh_private 2>/dev/null; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X PUT -H "Content-Type: application/json" \
  -d '{"update": {"labels": [{"remove": "old-label"}]}}' \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue/ISSUE_KEY"
```

### Sprint Management

```bash
# Get board ID
source ~/.zsh_private 2>/dev/null; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  "https://YOUR_ORG.atlassian.net/rest/agile/1.0/board?projectKeyOrId=PROJECT_KEY" | jq '.values[] | {id, name}'
```

```bash
# Get active sprint
source ~/.zsh_private 2>/dev/null; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  "https://YOUR_ORG.atlassian.net/rest/agile/1.0/board/BOARD_ID/sprint?state=active" | jq '.values[] | {id, name, state}'
```

```bash
# Move issues into a sprint
source ~/.zsh_private 2>/dev/null; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X POST -H "Content-Type: application/json" \
  -d '{"issues": ["ISSUE_KEY_1", "ISSUE_KEY_2"]}' \
  "https://YOUR_ORG.atlassian.net/rest/agile/1.0/sprint/SPRINT_ID/issue"
```
