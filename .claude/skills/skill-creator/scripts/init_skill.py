#!/usr/bin/env python3
"""
Skill Initializer - Creates a new skill from template

Usage:
    init_skill.py <skill-name> --path <path>

Examples:
    init_skill.py my-new-skill --path .claude/skills
    init_skill.py data-analyzer --path .claude/skills
"""

import sys
from pathlib import Path


SKILL_TEMPLATE = """---
name: {skill_name}
description: [TODO: Describe what this skill does and when to use it. Include trigger phrases like "/{skill_name}", "English trigger", etc.]
---

# {skill_title} Skill

[TODO: 1-2 sentence description of what this skill enables]

## Prerequisites

- [TODO: Required tools, APIs, configuration]

## Workflow

### Step 1: [TODO: Action Name]

[TODO: Specific, machine-executable instructions. Specify the exact tools or commands to use.]

### Step 2: [TODO: Action Name]

[TODO: Next step with clear inputs and outputs.]

## Error Handling

| Error | Cause | Solution |
|-------|-------|---------|
| [TODO] | [TODO] | [TODO] |

## Usage Examples

[TODO: Concrete usage examples with input and expected output]
"""

EXAMPLE_REFERENCE = """# {skill_title} Reference

[TODO: Write detailed reference documentation here.
Split long content here to keep SKILL.md slim.]
"""


def title_case_skill_name(skill_name):
    """Convert hyphenated skill name to Title Case for display."""
    return " ".join(word.capitalize() for word in skill_name.split("-"))


def init_skill(skill_name, path):
    """Initialize a new skill directory with template SKILL.md."""
    skill_dir = Path(path).resolve() / skill_name

    if skill_dir.exists():
        print(f"Error: Skill directory already exists: {skill_dir}")
        return None

    try:
        skill_dir.mkdir(parents=True, exist_ok=False)
        print(f"Created skill directory: {skill_dir}")
    except Exception as e:
        print(f"Error creating directory: {e}")
        return None

    skill_title = title_case_skill_name(skill_name)
    skill_content = SKILL_TEMPLATE.format(
        skill_name=skill_name, skill_title=skill_title
    )

    skill_md_path = skill_dir / "SKILL.md"
    try:
        skill_md_path.write_text(skill_content)
        print("Created SKILL.md")
    except Exception as e:
        print(f"Error creating SKILL.md: {e}")
        return None

    # Create references/ directory with example
    try:
        references_dir = skill_dir / "references"
        references_dir.mkdir(exist_ok=True)
        example_reference = references_dir / "reference.md"
        example_reference.write_text(
            EXAMPLE_REFERENCE.format(skill_title=skill_title)
        )
        print("Created references/reference.md")
    except Exception as e:
        print(f"Error creating references directory: {e}")
        return None

    print(f"\nSkill '{skill_name}' initialized at {skill_dir}")
    print("\nNext steps:")
    print("1. Fill in the [TODO] items in SKILL.md")
    print("2. Customize or delete the stubs in references/ as needed")
    print("3. Add scripts/ or assets/ directories if required")

    return skill_dir


def main():
    if len(sys.argv) < 4 or sys.argv[2] != "--path":
        print("Usage: init_skill.py <skill-name> --path <path>")
        print("\nSkill name requirements:")
        print("  - kebab-case (e.g., 'my-data-analyzer')")
        print("  - Lowercase letters, digits, and hyphens only")
        print("  - Max 64 characters")
        print("\nExamples:")
        print("  init_skill.py my-new-skill --path .claude/skills")
        print("  init_skill.py data-analyzer --path .claude/skills")
        sys.exit(1)

    skill_name = sys.argv[1]
    path = sys.argv[3]

    print(f"Initializing skill: {skill_name}")
    print(f"Location: {path}\n")

    result = init_skill(skill_name, path)
    sys.exit(0 if result else 1)


if __name__ == "__main__":
    main()
