---
name: skill-creator
description: Guide for creating and updating new skills. Triggers on "/skill-creator", "create a new skill", "create skill", "new skill", "build a skill", etc.
---

# Skill Creator

A skill for designing, creating, and updating skills. Builds modular, self-contained packages that extend Claude Code's capabilities.

## What Is a Skill

A skill is a modular package that extends Claude's capabilities. It provides domain-specific knowledge, workflows, and tool integrations — acting as an "onboarding guide" that transforms a general-purpose agent into a specialized one.

### What Skills Provide

1. **Specialized workflows** — Multi-step procedures for specific domains
2. **Tool integrations** — Instructions for working with specific file formats or APIs
3. **Domain knowledge** — Organization-specific knowledge, schemas, and business logic
4. **Bundled resources** — Scripts, references, and assets for complex, repetitive tasks

## Core Principles

### Brevity Is Key

The context window is a shared resource. A skill shares the window with the system prompt, conversation history, other skill metadata, and the user's request.

**Claude is already smart.** Only add context it doesn't yet have. For each piece of information, ask: "Does Claude really need this explanation?" and "Is this paragraph worth the token cost?"

Prefer concise examples over verbose explanations.

### Set the Right Level of Freedom

Adjust the level of specificity based on how fragile and variable the task is:

| Freedom | When to Use | Approach |
|---------|-------------|----------|
| High | Multiple valid approaches, context-dependent | Text-based instructions |
| Medium | Recommended pattern exists but variation is acceptable | Pseudocode or parameterized scripts |
| Low | Operations are fragile, error-prone, or consistency is critical | Concrete scripts, few parameters |

### Skill Structure

```
skill-name/
├── SKILL.md (required)
│   ├── YAML front matter (required)
│   │   ├── name: (required)
│   │   └── description: (required)
│   └── Markdown instructions (required)
└── Bundled resources (optional)
    ├── scripts/          - Executable code (Python/Bash, etc.)
    ├── references/       - Documents to load into context when needed
    └── assets/           - Files used in output (templates, icons, fonts, etc.)
```

#### SKILL.md (required)

- **Front matter** (YAML): `name` and `description` fields are required. `description` is the trigger mechanism for the skill, so clearly and comprehensively describe what it does and when to use it.
- **Body** (Markdown): Instructions for using the skill. Only loaded after the skill is triggered.

#### Bundled resources (optional)

- **scripts/**: Include when code would otherwise be rewritten repeatedly, or when deterministic reliability is required. Token-efficient and executable without loading into context.
- **references/**: Documents Claude should consult while working. Keeps SKILL.md slim; loaded only when needed. For files over 10k words, include grep search patterns in SKILL.md.
- **assets/**: Files used in final output (templates, images, fonts, etc.) — not loaded into context.

**Note:** Do not create unnecessary resource directories. Do not create supplementary documents like README.md or CHANGELOG.md.

### Progressive Disclosure

Manage context efficiently with a three-level loading system:

1. **Metadata (name + description)** — Always in context (~100 words)
2. **SKILL.md body** — Loaded when skill is triggered (<5k words, 500 lines or fewer)
3. **Bundled resources** — Loaded when Claude determines they are needed (unlimited)

Keep the SKILL.md body to 500 lines or fewer. If it exceeds that, split content into reference files and reference them clearly from SKILL.md.

See `references/workflows.md` and `references/output-patterns.md` for design pattern details.

## Skill Creation Process

### Step 1: Understand the Skill Through Concrete Examples

To create an effective skill, first clarify concrete use cases:
- "What functionality should this skill support?"
- "What are the intended use cases?"
- "What should trigger this skill?"

Avoid asking too many questions at once — start with the most important one.

### Step 2: Plan Reusable Content

For each concrete use case:
1. Think through how you would execute it from scratch
2. Identify scripts / references / assets that would be useful on repeated runs

### Step 3: Initialize the Skill

For new skills, run the `init_skill.py` script:

```bash
scripts/init_skill.py <skill-name> --path <output-directory>
```

This generates a template SKILL.md and stubs for scripts/, references/, and assets/.

### Step 4: Edit the Skill

- **Start with resources**: Implement the identified scripts / references / assets first
- **Test scripts**: Run any scripts you add to verify they work correctly
- **Remove unused resources**: Delete any stub files generated during initialization that are not needed

#### Editing SKILL.md

**Front matter:**
- `name`: Skill name (kebab-case)
- `description`: Primary trigger mechanism. Include what it does and when to use it. "When to use" belongs here, not in the body (the body is only read after triggering)

**Body:**
- Write in imperative or infinitive form
- Keep to 500 lines or fewer
- Include information that is helpful and non-obvious to other Claude instances

#### Style Guides

Follow `.claude/rules/claude/naming-convention.md` for skill naming conventions and `.claude/rules/claude/rule-authoring.md` for documentation style and structure.

### Step 5: Package the Skill

When development is complete, create a `.skill` file for distribution:

```bash
scripts/package_skill.py <path/to/skill-folder> [output-directory]
```

Validation and packaging run automatically. If there are validation errors, fix them and re-run.

### Step 6: Iterate

Use the skill on real tasks and improve any inefficiencies.

## meta-crud Integration

The existing `meta-crud` skill's operation 2 (creating a skill) is also a useful reference when creating new skills.

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| Directory already exists error in init_skill.py | A skill with that name already exists | Use a different name or update the existing skill |
| Validation failure | Invalid front matter, missing name/description | Follow the error message to fix |
| description too long | Exceeds 1024 characters | Move content to the body and summarize in description |
| Packaging failure | Unresolved validation errors | Use `quick_validate.py` to check individually |
