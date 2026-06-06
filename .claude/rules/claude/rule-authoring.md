---
description: Rules for authoring and maintaining .claude/rules/*.md files - quality standards for rule creation
globs:
  - .claude/rules/**/*.md
alwaysApply: false
---

# Rule Authoring Standards

## Overview

Rules are the smallest unit of the meta-hierarchy. They define declarative constraints, patterns, and quality standards that Claude automatically applies whenever a matching file is active.

## Frontmatter Requirements

All `.md` files must include the following YAML frontmatter:

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `description` | Yes | string | Describe the rule's purpose in English |
| `globs` | Yes | string[] | File patterns that activate this rule |
| `alwaysApply` | Yes | boolean | Normally `false`. Use `true` only for project-wide standards |

## File Naming

- **Format**: `kebab-case.md`
- **Convention**: `{service}/{topic-noun}.md` — categorized by subdirectory, filename ends with a noun
- **Subdirectories**: Group related rules under `rules/{service}/`
- **Examples**: `claude/rule-authoring.md`, `github/pr-creation.md`, `infrastructure/database/migration.md`

## Content Structure

```markdown
# {Domain Name}

## Overview
One-paragraph summary of what this rule covers.

## {Domain-specific sections}
(Varies based on rule content)

## Related Files
- Links to related rules, skills, and source files
```

- **H1**: Domain name (matches the purpose, not the filename)
- **First section**: Always `## Overview` with a concise summary
- **Body**: Domain-specific sections with actionable content
- **Last section**: `## Related Files` linking to related resources

## Quality Checklist

Check the following before creating or updating a rule:

1. **Self-contained**: The rule can be understood without reading other files
2. **Precise globs**: Matches only the intended files; not too broad
3. **No duplication**: Does not overlap with existing rules (verify with `Glob .claude/rules/**/*.md`)
4. **Actionable**: Provides concrete guidance Claude can follow, not vague advice
5. **Up to date**: Referenced file paths and patterns reflect the current state

## Anti-patterns

- **Globs too broad**: `**/*` or `*.ts` — keep them domain-specific
- **Stale code snippets**: Reference file paths instead of copying code
- **Mixed domains**: One rule covers one coherent domain
- **Redundant with CLAUDE.md**: Project-wide standards belong in CLAUDE.md, not in conditional rules
