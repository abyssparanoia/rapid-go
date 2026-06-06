---
description: Meta-hierarchy definition - relationships between Rules, Skills, and Agents in the Claude Code knowledge system
globs:
  - .claude/rules/claude/**/*
  - .claude/skills/meta-*/**/*
  - .claude/agents/**/*
alwaysApply: false
---

# Meta-Hierarchy: Rules → Skills → Agents

## Overview

Claude Code knowledge is managed in a three-tier hierarchy. Each tier has its own distinct role, scope, and authoring standards. This definition serves as the common standard across the organization.

## Hierarchy Structure

```
Agents (fewest, largest scope)         ← Orchestrates multiple Skills for complex tasks
  └── Skills (middle tier)             ← Combines Rules into concrete workflows
        └── Rules (most, smallest scope)   ← Declarative constraints, patterns, quality standards
```

| Tier | Location | Format | Trigger | Scope |
|------|----------|--------|---------|-------|
| Rules | `.claude/rules/**/*.md` | YAML frontmatter + Markdown | Automatic (globs match) | Single domain |
| Skills | `.claude/skills/*/SKILL.md` | YAML frontmatter + Markdown | Semantic matching / command | Single workflow |
| Agents | `.claude/agents/*.md` | YAML frontmatter + Markdown | Explicit invocation | Multi-skill tasks |

## Naming Conventions

See `.claude/rules/claude/naming-convention.md` for the full specification.

The pattern used is `{what}-{object?}-{action}`.

## What to Create

| Situation | Create |
|-----------|--------|
| New coding standards or constraints for a file pattern | **Rule** |
| Reusable multi-step workflow that uses tools | **Skill** |
| Autonomous multi-skill task with explicit boundaries | **Agent** |
| One-time instruction for the entire project | **CLAUDE.md entry** (not a Rule) |

### Decision Flow

1. Is it a declarative constraint triggered by a file pattern? → **Rule**
2. Is it a step-by-step procedure triggered by user intent? → **Skill**
3. Does it require autonomous orchestration of multiple Skills? → **Agent**
4. None of the above → Add to `CLAUDE.md` or domain documentation

## Feedback Loop

```
PR review comments
       │
       ▼
  meta-review skill        ← Classifies and prioritizes feedback
       │
       ▼
  meta-improvement agent   ← Decides what to create or update
       │
       ▼
  meta-crud skill          ← Executes changes to Rules/Skills/Agents
       │
       ▼
  Pull Request             ← Human reviews before merging
       │
       ▼
  Improved Rules/Skills    ← Better guidance for future sessions
```

## Related Files

- `.claude/rules/claude/naming-convention.md` — Naming conventions for all tiers
- `.claude/rules/claude/rule-authoring.md` — Quality standards for Rules
- `.claude/rules/claude/skill-authoring.md` — Quality standards for Skills
- `.claude/rules/claude/agent-authoring.md` — Quality standards for Agents
