# Meta-CRUD Templates

Format templates for creating Rules, Skills, and Agents.

## Rule Template (`.md`)

```markdown
---
description: {English description — the purpose of the domain this rule covers}
globs:
  - {glob-pattern-1}
  - {glob-pattern-2}
alwaysApply: false
---

# {Domain Name}

## Overview

{A one-paragraph summary of what this rule covers and why it exists.}

## {Primary Section}

{Domain-specific content: constraints, patterns, standards, etc.}

## {Additional Section (as needed)}

{Domain-specific content.}

## Related Files

- {Links to related rules, skills, and source files}
```

### Rule Frontmatter Fields

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `description` | string | Yes | English, concise statement of purpose |
| `globs` | string[] | Yes | Specific file patterns; avoid `**/*` |
| `alwaysApply` | boolean | Yes | `false` unless this is a project-wide standard |

---

## Skill Template (`SKILL.md`)

```markdown
---
name: {skill-name}
description: {English description. Triggered by "/{skill-name}", "trigger phrase", etc.}
---

# {Skill Name} Skill

{One-line description of what this skill does.}

## Prerequisites

- {Required tools, APIs, or configuration}

## Workflow

### Step 1: {Action Name}

{Concrete, machine-executable instructions.}
{Specify the exact tools or commands to use.}

### Step 2: {Action Name}

{Next step with clear inputs and outputs.}

## Error Handling

| Error | Cause | Resolution |
|-------|-------|------------|
| {Error description} | {Why it occurs} | {How to fix it} |

## Example Usage

{A concrete usage example with inputs and expected outputs.}
```

### Skill Frontmatter Fields

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `name` | string | Yes | Must match the directory name (kebab-case) |
| `description` | string | Yes | Include trigger keywords (English) |

---

## Agent Template (`.md`)

```markdown
---
name: {agent-name}
description: {English description — the agent's purpose and triggers}
tools:
  - {Tool1}
  - {Tool2}
model: sonnet
skills:
  - {skill-name-1}
  - {skill-name-2}
memory: project
---

# {Agent Name}

## Purpose

{What this agent does and when it should be invoked.}

## Workflow

### Step 1: {Phase Name}

{Instructions for this phase.}

### Step 2: {Phase Name}

{Instructions for this phase.}

## Constraints

- {What the agent must not do}
- {Maximum scope per execution}
- {Protected resources}

## Escalation

- {When to stop and ask a human}
- {Conditions that require approval}

## Team Integration

- **role**: {lead | specialist | reviewer}
- {How this agent collaborates with other agents}
```

### Agent Frontmatter Fields

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `name` | string | Yes | Agent identifier (kebab-case) |
| `description` | string | Yes | Description of purpose and triggers |
| `tools` | string[] | Yes | Minimum required tools |
| `model` | string | No | `sonnet` (default), `opus`, or `haiku` |
| `skills` | string[] | No | Skills to preload as references |
| `memory` | string | No | `user`, `project`, or `local` |
| `maxTurns` | number | No | Maximum number of execution turns |
