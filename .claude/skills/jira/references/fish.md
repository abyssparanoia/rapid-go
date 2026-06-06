# JIRA Command Reference (fish)

JIRA API command reference for fish shell users.

## Setting Credentials

### Option A: Add to .envrc (recommended — auto-loaded via direnv)

Add to `.envrc`:

```bash
export JIRA_EMAIL='your-email@example.com'
export JIRA_API_KEY='your-api-token'
```

### Option B: Add to fish conf.d

Add to `~/.config/fish/conf.d/private.fish`:

```fish
set -gx JIRA_EMAIL 'your-email@example.com'
set -gx JIRA_API_KEY 'your-api-token'
```

## Authentication Prefix

Claude Code runs commands in bash, so for fish users credentials are sourced from `.envrc`.
Because the `-u` flag may not handle special characters in tokens correctly, the `Authorization: Basic` header is built using explicit Base64 encoding.

```bash
# Source .envrc to load environment variables, then build the Base64 auth header
set -a; source .envrc; set +a; JIRA_AUTH=$(echo -n "$JIRA_EMAIL:$JIRA_API_KEY" | base64)
```

Prepend this prefix to every command below.

## Commands

### List Projects

```bash
set -a; source .envrc; set +a; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  "https://YOUR_ORG.atlassian.net/rest/api/3/project" | jq '[.[] | {key, name}]'
```

### Search Issues (JQL)

```bash
set -a; source .envrc; set +a; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X POST -H "Content-Type: application/json" \
  -d '{"jql": "project=SP ORDER BY created DESC", "maxResults": 10, "fields": ["summary", "status", "assignee", "priority"]}' \
  "https://YOUR_ORG.atlassian.net/rest/api/3/search/jql" | jq '.issues[] | {key, summary: .fields.summary, status: .fields.status.name, assignee: .fields.assignee.displayName}'
```

### Get Issue Details

```bash
set -a; source .envrc; set +a; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue/ISSUE_KEY" | jq '{key, summary: .fields.summary, status: .fields.status.name, description: .fields.description}'
```

### Create an Issue

```bash
set -a; source .envrc; set +a; curl -s -H "Authorization: Basic $JIRA_AUTH" \
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
set -a; source .envrc; set +a; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue/ISSUE_KEY/transitions" | jq '.transitions[] | {id, name}'
```

```bash
# Execute a transition
set -a; source .envrc; set +a; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X POST -H "Content-Type: application/json" \
  -d '{"transition": {"id": "TRANSITION_ID"}}' \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue/ISSUE_KEY/transitions"
```

### Add a Comment

```bash
set -a; source .envrc; set +a; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X POST -H "Content-Type: application/json" \
  -d '{"body": {"type": "doc", "version": 1, "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Comment text here"}]}]}}' \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue/ISSUE_KEY/comment" | jq '{id, author: .author.displayName}'
```

### Add / Remove Labels

```bash
# Add a label
set -a; source .envrc; set +a; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X PUT -H "Content-Type: application/json" \
  -d '{"update": {"labels": [{"add": "new-label"}]}}' \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue/ISSUE_KEY"
```

```bash
# Remove a label
set -a; source .envrc; set +a; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X PUT -H "Content-Type: application/json" \
  -d '{"update": {"labels": [{"remove": "old-label"}]}}' \
  "https://YOUR_ORG.atlassian.net/rest/api/3/issue/ISSUE_KEY"
```

### Sprint Management

```bash
# Get board ID
set -a; source .envrc; set +a; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  "https://YOUR_ORG.atlassian.net/rest/agile/1.0/board?projectKeyOrId=PROJECT_KEY" | jq '.values[] | {id, name}'
```

```bash
# Get active sprint
set -a; source .envrc; set +a; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  "https://YOUR_ORG.atlassian.net/rest/agile/1.0/board/BOARD_ID/sprint?state=active" | jq '.values[] | {id, name, state}'
```

```bash
# Move issues into a sprint
set -a; source .envrc; set +a; curl -s -H "Authorization: Basic $JIRA_AUTH" \
  -X POST -H "Content-Type: application/json" \
  -d '{"issues": ["ISSUE_KEY_1", "ISSUE_KEY_2"]}' \
  "https://YOUR_ORG.atlassian.net/rest/agile/1.0/sprint/SPRINT_ID/issue"
```
