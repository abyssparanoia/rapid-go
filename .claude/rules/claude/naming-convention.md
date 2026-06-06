---
description: Naming conventions for Rules, Skills, and Agents
globs: []
alwaysApply: false
---

# Naming Conventions

## Base Pattern

```
{what}-{object?}-{action}
  ↑ service/layer   ↑ optional   ↑ gerund or action noun
```

Use kebab-case. The final segment should be an **action noun** (gerund or nominalized verb).

## Rules

The directory handles `{what}/`, so the filename is `{object?}-{action}.md`.

```
.claude/rules/{what}/{object?}-{action}.md
```

Examples: `github/pr-creation.md`, `claude/rule-authoring.md`, `infrastructure/database/migration.md`

## Skills

Directory name: `{what}-{object?}-{action}/`
YAML name (slash cmd): For Action-type skills only, convert the final segment to an **imperative verb**. Knowledge and Service Wrapper types keep the noun form.

| Type | Dir example | YAML name example |
|------|-------|-------------|
| Action | `github-pr-creation/` | `github-pr-create` |
| Knowledge | `react-best-practices/` | `react-best-practices` |
| Service Wrapper | `jira/` | `jira` |

## Agents

Filename: `{what}-{object?}-{action}.md`. The YAML name matches.
Examples: `meta-improvement.md`, `fe-implementation.md`

## Allowed Exceptions

- **Knowledge skills**: A knowledge domain name that stands on its own (`react-best-practices`, `next-cache-components`)
- **Service Wrappers**: The service name alone (`jira`, `figma`)
- **Acronyms**: Compound nouns representing well-known operations (`meta-crud`)
