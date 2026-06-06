---
description: Standards for authoring and maintaining .claude/agents/*.md files - quality standards for agent creation
globs:
  - .claude/agents/**/*.md
alwaysApply: false
---

# Agent Authoring Standards

## Overview

Agents are the highest-scope, lowest-count tier of the meta-hierarchy. They autonomously execute complex multi-step tasks by orchestrating multiple Skills, with explicit boundaries and escalation policies.

## File Placement

```
.claude/agents/{agent-name}.md
```

- **Format**: `kebab-case.md`
- **Location**: always under `.claude/agents/`

## Frontmatter Requirements

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `name` | Yes | string | Agent identifier |
| `description` | Yes | string | Description of purpose and triggers |
| `tools` | Yes | string[] | Minimum required tools (principle of least privilege) |
| `model` | No | string | `sonnet` (default) / `opus` / `haiku` |
| `skills` | No | string[] | Skills to preload as references |
| `memory` | No | string | Memory scope: `user` / `project` / `local` |
| `maxTurns` | No | number | Upper limit on execution turns |

### Tool Selection Principle

Grant only the tools the agent actually needs:

```yaml
# Good - minimal tools for read-only analysis
tools: [Read, Grep, Glob]

# Bad - unnecessary write permissions
tools: [Read, Grep, Glob, Write, Edit, Bash]
```

## Content Structure

```markdown
# {Agent Name}

## Purpose
{What this agent does and when it should be invoked}

## Workflow
### Step 1: {Phase name}
{Instructions}
...

## Constraints
- {What the agent must not do}
- {Escalation conditions}

## Escalation
- {When to check with a human}

## Team Integration
- role: {lead | specialist | reviewer}
- {How this agent coordinates with other agents}
```

## Design Principles

1. **Single responsibility**: One agent handles one coherent task domain
2. **Explicit boundaries**: Clearly state what the agent must not do
3. **Deterministic workflow**: Steps are reproducible with consistent results
4. **Graceful degradation**: Define behavior on tool failure or missing data
5. **Human in the loop**: Destructive or irreversible actions require human approval

## Agent Teams Compatibility

For agents that participate in team coordination:

- Declare `role` in the Team Integration section
- Support the shared task list format (TaskCreate/TaskUpdate/TaskList)
- Document communication patterns between agents
- Specify which agents may delegate tasks to this agent

## Quality Checklist

1. **Minimal tools**: Only tools actually needed are listed
2. **Skills exist**: All Skills referenced in the `skills` field are implemented
3. **Boundaries defined**: A clear list of what the agent must not do
4. **Escalation policy**: Conditions under which the agent stops and checks with a human
5. **Testable workflow**: Can be invoked manually and verified
