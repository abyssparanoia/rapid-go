---
name: meta-crud
description: Create, update, delete, and list rules, skills, and agents. Triggered by "/meta-crud", "create a rule", "add a skill", "create rule", "add skill", "create an agent", etc.
---

# Meta-CRUD Skill

Creates, reads, updates, and deletes rules, skills, and agents within the `.claude/` meta-system.

## Prerequisites

Before performing any operation, read the relevant authoring rules:
- Rules: `.claude/rules/claude/rule-authoring.md`
- Skills: `.claude/rules/claude/skill-authoring.md`
- Agents: `.claude/rules/claude/agent-authoring.md`

Format templates are in `.claude/skills/meta-crud/references/templates.md`.

## Workflow

### Operation 1: Create a Rule

1. Read `.claude/rules/claude/rule-authoring.md` to review quality standards
2. Read `references/templates.md` to review the rule template
3. Check for duplicate globs:
   ```
   Grep pattern="{proposed-glob-pattern}" path=".claude/rules/" glob="*.md"
   ```
4. If duplicates exist → warn the user and suggest adjusting the globs
5. Use the template to generate a `.md` file, filling in:
   - `description` (in English)
   - `globs` (domain-specific)
   - `alwaysApply` (normally `false`)
   - Content sections aligned with the authoring standards
6. Write the file to `.claude/rules/{category}/{name}.md`
7. Validate: read the created file and confirm the frontmatter is valid YAML

### Operation 2: Create a Skill

1. Read `.claude/rules/claude/skill-authoring.md` to review quality standards
2. Read `references/templates.md` to review the skill template
3. Check whether the directory already exists:
   ```
   Glob pattern=".claude/skills/{skill-name}/SKILL.md"
   ```
4. If it exists → warn the user and suggest using the update operation instead
5. Create the directory: `.claude/skills/{skill-name}/`
6. Use the template to generate `SKILL.md`, verifying:
   - `name` matches the directory name
   - `description` includes trigger keywords (in English)
   - The workflow has numbered, machine-executable steps
   - An error-handling section is included
7. If reference data is needed, create the `references/` directory and files
8. Validate: read the created SKILL.md and confirm the frontmatter

### Operation 3: Create an Agent

1. Read `.claude/rules/claude/agent-authoring.md` to review quality standards
2. Read `references/templates.md` to review the agent template
3. Check whether `.claude/agents/` exists (create it if not)
4. Validate that referenced skills exist:
   ```
   Glob pattern=".claude/skills/{skill-name}/SKILL.md"  (for each skill)
   ```
5. If a skill is not found → warn the user
6. Use the template to generate `.claude/agents/{agent-name}.md`, verifying:
   - `tools` follows the principle of least privilege
   - `skills` references only existing skills
   - The constraints section is explicit
   - An escalation policy is defined
7. Validate: read the created file and confirm the frontmatter

### Operation 4: List Inventory

1. Scan all rules:
   ```
   Glob pattern=".claude/rules/**/*.md"
   ```
2. Scan all skills:
   ```
   Glob pattern=".claude/skills/*/SKILL.md"
   ```
3. Scan all agents:
   ```
   Glob pattern=".claude/agents/*.md"
   ```
4. Read the frontmatter of each file to extract metadata
5. Output as a table:
   | Type | Name | Description | Status |
   |------|------|-------------|--------|

### Operation 5: Update

1. Read the target file
2. Read the authoring rules for that file type
3. Apply the changes
4. Run through the authoring-rule quality checklist
5. Write the updated file
6. Validate the changes

### Operation 6: Delete

1. Read the target file to review its contents
2. Check for references from other files:
   ```
   Grep pattern="{filename}" path=".claude/"
   ```
3. If references exist → list them and warn the user
4. Ask the user to confirm the deletion (destructive operation)
5. Delete the file
6. Update the inventory table in `meta-hierarchy-definition.md` if needed

## Error Handling

| Error | Cause | Resolution |
|-------|-------|------------|
| Duplicate globs detected | New rule conflicts with an existing one | Narrow the globs to be more specific |
| Skill directory already exists | Attempted to create a duplicate | Use the update operation instead |
| Referenced skill not found | Agent references a skill that does not exist | Create the skill first |
| Frontmatter parse error | Invalid YAML format | Fix indentation and field formatting |

## Example Usage

User: "Create an infrastructure database rule covering Prisma patterns"

```
1. Read rule-authoring.md
2. Read templates.md
3. Grep existing rules for "src/infrastructure/database/**/*" → no duplicates found
4. Write .claude/rules/infrastructure/database-prisma.md with:
   - description: "Prisma ORM patterns and conventions"
   - globs: ["src/infrastructure/database/**/*.ts"]
   - alwaysApply: false
5. Validate frontmatter
```
