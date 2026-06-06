---
description: Standards for authoring and maintaining .claude/skills/*/SKILL.md files - quality standards for skill creation
globs:
  - .claude/skills/**/SKILL.md
alwaysApply: false
---

# Skill Authoring Standards

## Overview

Skills are the middle tier of the meta-hierarchy, combining multiple rules into concrete workflows. They provide step-by-step instructions for Claude to follow when triggered by a user command or semantic matching.

## Directory Structure

```
.claude/skills/{skill-name}/
├── SKILL.md              # Main skill definition (required)
└── references/           # Supporting reference files (optional)
    ├── templates.md
    └── mappings.md
```

- **Directory name**: `kebab-case`, matching the `name` field in the frontmatter
- **SKILL.md**: Required entry point
- **references/**: Optional directory for reference data, templates, and mappings

## Frontmatter Requirements

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `name` | Yes | string | Must match the directory name |
| `description` | Yes | string | Include trigger keywords for semantic matching (both English and Japanese) |

### Trigger Keywords in Description

The `description` field should include keywords that help activate the skill:

```yaml
# Good - includes triggers
description: Automates GitHub PR creation. Activated by "/github-pr-create", "create a PR", "make a PR", etc.

# Bad - no trigger keywords
description: Creates a PR
```

## Content Structure

```markdown
# {Skill Name} Skill

{One-line description}

## Workflow

### Step 1: {Action name}
{Specific, machine-executable instructions}

### Step 2: {Action name}
...

## Error Handling

| Error | Cause | Resolution |
|-------|-------|------------|
| ... | ... | ... |

## Examples

{Concrete examples with input and expected output}
```

## Workflow Step Requirements

Each step must:

1. **Be machine-executable**: Claude can carry it out without human judgment
2. **Specify tools**: Name the exact tools or commands to use
3. **Have clear inputs/outputs**: State what goes in and what comes out
4. **Handle errors**: Account for expected failure modes

## Skill Composition

When a skill depends on another skill:

- Reference it explicitly: "Use the `meta-crud` skill to create the rule"
- Do not duplicate: Link to the other skill's workflow instead of copying it
- Declare dependencies in the `skills` field when used inside an agent

## Quality Checklist

1. **Executable steps**: All steps can be carried out mechanically
2. **Error handling**: Common failure modes are documented with resolutions
3. **Trigger coverage**: Description includes triggers in both English and Japanese
4. **Single responsibility**: One skill handles one coherent workflow domain
5. **References organized**: Large data lives in `references/` rather than inline
