---
name: jira
description: JIRA operations (issue search, creation, update, sprint management, label operations). Triggered by "/jira", "in JIRA", "SP-XXXX", "ticket", "Issue", "sprint", "search JIRA", "create ticket", etc.
---

# JIRA Skill

A skill that calls the JIRA REST API directly to search, create, and update issues, manage sprints, and manipulate labels.

## Prerequisites

The following environment variables must be set:

- `JIRA_EMAIL` — Your Atlassian account email address
- `JIRA_API_KEY` — An API token generated at [Atlassian Account Settings](https://id.atlassian.com/manage-profile/security/api-tokens)

Where you set these depends on your shell. See the references for details.

## Shell Detection and Authentication

When running commands, use the reference that matches the user's shell environment:

1. **Check the user's shell**: Look at the `Shell:` field in the system information
2. **Read the corresponding reference**:
   - fish → `references/fish.md`
   - zsh / bash → `references/zsh.md`
3. **Run curl commands using the authentication prefix described in that reference**

## Use Cases

### 1. List Projects

`GET /rest/api/3/project` filtered with jq.

### 2. Search Issues (JQL)

Send a JQL query to `POST /rest/api/3/search/jql`.

### 3. Get Issue Details

Fetch with `GET /rest/api/3/issue/{ISSUE_KEY}`.

### 4. Create an Issue

Send fields to `POST /rest/api/3/issue`.

### 5. Change Status (Transition)

1. `GET /rest/api/3/issue/{ISSUE_KEY}/transitions` to retrieve available transitions
2. `POST /rest/api/3/issue/{ISSUE_KEY}/transitions` to execute a transition

### 6. Add a Comment

Send an ADF-format body to `POST /rest/api/3/issue/{ISSUE_KEY}/comment`.

### 7. Add / Remove Labels

Send `update.labels` to `PUT /rest/api/3/issue/{ISSUE_KEY}`.

### 8. Sprint Management

1. `GET /rest/agile/1.0/board?projectKeyOrId={KEY}` to get the board ID
2. `GET /rest/agile/1.0/board/{BOARD_ID}/sprint?state=active` to get the active sprint
3. `POST /rest/agile/1.0/sprint/{SPRINT_ID}/issue` to move issues into a sprint

## Useful JQL Query Examples

| Purpose | JQL |
|---------|-----|
| My issues | `assignee = currentUser() AND status != Done` |
| Created this week | `project = SP AND created >= startOfWeek()` |
| High priority | `project = SP AND priority in (Highest, High) AND status != Done` |
| By label | `project = SP AND labels = "backend"` |
| In sprint | `project = SP AND sprint in openSprints()` |
| Incomplete tasks | `project = SP AND status NOT IN (Done, Closed)` |
| Recently updated | `project = SP AND updated >= -7d ORDER BY updated DESC` |

## Workflow

1. User requests a JIRA operation via `/jira <request>` or natural language
2. Detect the shell and load credentials using the command from the corresponding reference
3. Build and execute the appropriate curl command
4. Parse the result and report back to the user

## Comparison with MCP

| Aspect | Via MCP | Via Skill |
|--------|---------|-----------|
| Setup | Requires MCP server configuration | No setup needed (direct curl) |
| Token consumption | All tool schemas are loaded | Loaded only when needed |
| Customizability | Depends on the server | Full control over JQL and fields |
| Transparency | Black box | Commands are visible |
| How to invoke | Automatic | `/jira` or natural language |

## Notes

- Do not commit credentials to the repository
- Be mindful of rate limits (avoid bulk requests)
- Perform write operations (create/update) with care
