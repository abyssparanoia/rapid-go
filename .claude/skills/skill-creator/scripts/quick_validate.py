#!/usr/bin/env python3
"""
Quick validation script for skills
"""

import sys
import re
from pathlib import Path

try:
    import yaml
except ImportError:
    yaml = None


def parse_frontmatter(content):
    """Parse YAML frontmatter without pyyaml dependency."""
    match = re.match(r"^---\n(.*?)\n---", content, re.DOTALL)
    if not match:
        return None
    if yaml:
        return yaml.safe_load(match.group(1))
    # Simple fallback parser for key: value pairs
    result = {}
    for line in match.group(1).split("\n"):
        line = line.strip()
        if ":" in line and not line.startswith("-"):
            key, value = line.split(":", 1)
            result[key.strip()] = value.strip()
    return result


def validate_skill(skill_path):
    """Basic validation of a skill."""
    skill_path = Path(skill_path)

    skill_md = skill_path / "SKILL.md"
    if not skill_md.exists():
        return False, "SKILL.md not found"

    content = skill_md.read_text()
    if not content.startswith("---"):
        return False, "No YAML frontmatter found"

    frontmatter = parse_frontmatter(content)
    if frontmatter is None:
        return False, "Invalid frontmatter format"
    if not isinstance(frontmatter, dict):
        return False, "Frontmatter must be a YAML dictionary"

    if "name" not in frontmatter:
        return False, "Missing 'name' in frontmatter"
    if "description" not in frontmatter:
        return False, "Missing 'description' in frontmatter"

    name = str(frontmatter.get("name", "")).strip()
    if name:
        if not re.match(r"^[a-z0-9-]+$", name):
            return False, f"Name '{name}' should be kebab-case"
        if name.startswith("-") or name.endswith("-") or "--" in name:
            return False, f"Name '{name}' has invalid hyphen placement"
        if len(name) > 64:
            return False, f"Name too long ({len(name)} chars, max 64)"

    description = str(frontmatter.get("description", "")).strip()
    if description:
        if "<" in description or ">" in description:
            return False, "Description cannot contain angle brackets"
        if len(description) > 1024:
            return False, f"Description too long ({len(description)} chars, max 1024)"

    return True, "Skill is valid!"


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python quick_validate.py <skill_directory>")
        sys.exit(1)

    valid, message = validate_skill(sys.argv[1])
    print(message)
    sys.exit(0 if valid else 1)
