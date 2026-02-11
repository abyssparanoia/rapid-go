#!/usr/bin/env python3
"""
Repository initialization script for rapid-go template.

This script renames/replaces all rapid-go specific identifiers to create
a new repository from the rapid-go template.
"""

import argparse
import os
import re
import shutil
import sys
from pathlib import Path
from typing import List, Tuple


def parse_args():
    """Parse command line arguments."""
    parser = argparse.ArgumentParser(
        description="Initialize new repository from rapid-go template"
    )
    parser.add_argument(
        "--go-module",
        required=True,
        help="New Go module path (e.g., github.com/myorg/myproject)",
    )
    parser.add_argument(
        "--service-name",
        required=True,
        help="New service name (e.g., myservice)",
    )
    parser.add_argument(
        "--database",
        required=True,
        choices=["mysql", "postgresql"],
        help="Database to use (mysql or postgresql)",
    )
    parser.add_argument(
        "--project-title",
        help="Project title (default: SERVICE_NAME in uppercase)",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Show what would be changed without making changes",
    )

    return parser.parse_args()


def extract_buf_org(go_module: str) -> str:
    """
    Extract organization name from Go module path.

    Example: github.com/myorg/myproject -> myorg
    """
    parts = go_module.split("/")
    if len(parts) >= 2:
        return parts[1]
    return parts[0]


def get_repository_root() -> Path:
    """Get the repository root directory."""
    # Assume script is in .claude/skills/init-new-repository/scripts/
    script_dir = Path(__file__).parent
    repo_root = script_dir.parent.parent.parent.parent
    return repo_root.resolve()


def find_all_text_files(root_dir: Path, exclude_dirs: List[str]) -> List[Path]:
    """
    Find all text files (non-binary) in the repository.

    Args:
        root_dir: Root directory to search
        exclude_dirs: List of directory names to exclude

    Returns:
        List of Path objects for text files
    """
    text_files = []
    binary_extensions = {'.png', '.jpg', '.jpeg', '.gif', '.pdf', '.zip', '.tar', '.gz', '.exe', '.bin', '.so', '.dylib', '.skill'}

    for item in root_dir.rglob('*'):
        # Skip excluded directories
        if any(exclude in item.parts for exclude in exclude_dirs):
            continue

        # Skip non-files
        if not item.is_file():
            continue

        # Skip binary files by extension
        if item.suffix.lower() in binary_extensions:
            continue

        # Skip hidden files (except .envrc, .envrc.tmpl, .golangci.yml, .gitignore, etc.)
        if item.name.startswith('.') and item.name not in {
            '.envrc', '.envrc.tmpl', '.golangci.yml', '.gitignore', '.gitattributes',
            '.dockerignore', '.editorconfig'
        }:
            continue

        text_files.append(item)

    return text_files


def replace_in_file(file_path: Path, replacements: List[Tuple[str, str]], dry_run: bool = False) -> bool:
    """
    Replace all occurrences of old strings with new strings in a file.

    Args:
        file_path: Path to the file
        replacements: List of (old, new) tuples
        dry_run: If True, don't actually modify the file

    Returns:
        True if file was modified, False otherwise
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except UnicodeDecodeError:
        # Skip binary files
        return False
    except Exception as e:
        print(f"Warning: Could not read {file_path}: {e}")
        return False

    original_content = content

    # Apply all replacements
    for old, new in replacements:
        content = content.replace(old, new)

    if content != original_content:
        if not dry_run:
            try:
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(content)
            except Exception as e:
                print(f"Error: Could not write to {file_path}: {e}")
                return False
        return True

    return False


def rename_directory(old_path: Path, new_path: Path, dry_run: bool = False) -> bool:
    """
    Rename a directory.

    Args:
        old_path: Current path
        new_path: New path
        dry_run: If True, don't actually rename

    Returns:
        True if directory was renamed, False otherwise
    """
    if not old_path.exists():
        return False

    if new_path.exists():
        print(f"Warning: Target path already exists: {new_path}")
        return False

    if not dry_run:
        try:
            old_path.rename(new_path)
        except Exception as e:
            print(f"Error: Could not rename {old_path} to {new_path}: {e}")
            return False

    return True


def delete_directory(path: Path, dry_run: bool = False) -> bool:
    """
    Delete a directory and all its contents.

    Args:
        path: Path to delete
        dry_run: If True, don't actually delete

    Returns:
        True if directory was deleted, False otherwise
    """
    if not path.exists():
        return False

    if not dry_run:
        try:
            shutil.rmtree(path)
        except Exception as e:
            print(f"Error: Could not delete {path}: {e}")
            return False

    return True


def toggle_makefile_targets(makefile_path: Path, database: str, dry_run: bool = False):
    """
    Toggle commented Makefile targets for database selection.

    Args:
        makefile_path: Path to Makefile
        database: Selected database ('mysql' or 'postgresql')
        dry_run: If True, don't actually modify
    """
    if not makefile_path.exists():
        return

    try:
        with open(makefile_path, 'r', encoding='utf-8') as f:
            lines = f.readlines()
    except Exception as e:
        print(f"Warning: Could not read {makefile_path}: {e}")
        return

    modified = False
    for i, line in enumerate(lines):
        # Comment out unused database targets
        if database == 'mysql':
            if 'make generate.mermaid.mysql' in line and line.strip().startswith('#'):
                lines[i] = line.lstrip('#').lstrip()
                modified = True
            elif 'make generate.sqlboiler.mysql' in line and line.strip().startswith('#'):
                lines[i] = line.lstrip('#').lstrip()
                modified = True
            elif 'make generate.mermaid.postgresql' in line and not line.strip().startswith('#'):
                lines[i] = '# ' + line
                modified = True
            elif 'make generate.sqlboiler.postgresql' in line and not line.strip().startswith('#'):
                lines[i] = '# ' + line
                modified = True
        else:  # postgresql
            if 'make generate.mermaid.postgresql' in line and line.strip().startswith('#'):
                lines[i] = line.lstrip('#').lstrip()
                modified = True
            elif 'make generate.sqlboiler.postgresql' in line and line.strip().startswith('#'):
                lines[i] = line.lstrip('#').lstrip()
                modified = True
            elif 'make generate.mermaid.mysql' in line and not line.strip().startswith('#'):
                lines[i] = '# ' + line
                modified = True
            elif 'make generate.sqlboiler.mysql' in line and not line.strip().startswith('#'):
                lines[i] = '# ' + line
                modified = True

    if modified and not dry_run:
        try:
            with open(makefile_path, 'w', encoding='utf-8') as f:
                f.writelines(lines)
        except Exception as e:
            print(f"Error: Could not write to {makefile_path}: {e}")


def toggle_envrc_template(envrc_path: Path, database: str, dry_run: bool = False):
    """
    Toggle commented environment variables in .envrc.tmpl for database selection.

    Args:
        envrc_path: Path to .envrc.tmpl
        database: Selected database ('mysql' or 'postgresql')
        dry_run: If True, don't actually modify
    """
    if not envrc_path.exists():
        return

    try:
        with open(envrc_path, 'r', encoding='utf-8') as f:
            lines = f.readlines()
    except Exception as e:
        print(f"Warning: Could not read {envrc_path}: {e}")
        return

    modified = False
    in_mysql_section = False
    in_postgresql_section = False

    for i, line in enumerate(lines):
        # Detect sections
        if '# for mysql' in line.lower():
            in_mysql_section = True
            in_postgresql_section = False
            continue
        elif '# for postgresql' in line.lower():
            in_mysql_section = False
            in_postgresql_section = True
            continue
        elif line.strip() == '' or not line.startswith('export DB_') and not line.startswith('# export DB_'):
            in_mysql_section = False
            in_postgresql_section = False
            continue

        # Toggle comments based on selected database
        if in_mysql_section:
            if database == 'mysql':
                if line.strip().startswith('# export DB_'):
                    lines[i] = line.lstrip('# ')
                    modified = True
            else:  # postgresql
                if line.strip().startswith('export DB_'):
                    lines[i] = '# ' + line
                    modified = True
        elif in_postgresql_section:
            if database == 'postgresql':
                if line.strip().startswith('# export DB_'):
                    lines[i] = line.lstrip('# ')
                    modified = True
            else:  # mysql
                if line.strip().startswith('export DB_'):
                    lines[i] = '# ' + line
                    modified = True

    if modified and not dry_run:
        try:
            with open(envrc_path, 'w', encoding='utf-8') as f:
                f.writelines(lines)
        except Exception as e:
            print(f"Error: Could not write to {envrc_path}: {e}")


def remove_docker_compose_db_service(compose_path: Path, database: str, dry_run: bool = False):
    """
    Remove unused database service from docker-compose.yml.

    Args:
        compose_path: Path to docker-compose.yml
        database: Selected database ('mysql' or 'postgresql')
        dry_run: If True, don't actually modify
    """
    if not compose_path.exists():
        return

    try:
        with open(compose_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        print(f"Warning: Could not read {compose_path}: {e}")
        return

    # Determine which service to remove
    remove_service = 'postgresql' if database == 'mysql' else 'mysql'

    # Find and remove the service block
    # This is a simple approach that works for typical docker-compose.yml structure
    lines = content.split('\n')
    result_lines = []
    skip = False
    indent_level = 0

    for line in lines:
        # Check if this is the start of the service to remove
        if re.match(rf'^  {remove_service}:', line):
            skip = True
            indent_level = len(line) - len(line.lstrip())
            continue

        # If we're skipping, check if we've reached the next service or end of services
        if skip:
            current_indent = len(line) - len(line.lstrip())
            # If line is not empty and has same or less indentation, we've reached the next section
            if line.strip() and current_indent <= indent_level:
                skip = False

        if not skip:
            result_lines.append(line)

    new_content = '\n'.join(result_lines)

    if new_content != content and not dry_run:
        try:
            with open(compose_path, 'w', encoding='utf-8') as f:
                f.write(new_content)
        except Exception as e:
            print(f"Error: Could not write to {compose_path}: {e}")


def toggle_database_imports(repo_root: Path, database: str, dry_run: bool = False) -> bool:
    """
    Toggle database import paths in Go files.

    For PostgreSQL selection, replaces 'mysql' with 'postgresql' in database import aliases.
    For MySQL selection, no changes needed (mysql is the default).

    Args:
        repo_root: Repository root directory
        database: Selected database ('mysql' or 'postgresql')
        dry_run: If True, don't actually modify files

    Returns:
        True if any files were modified, False otherwise
    """
    if database == 'mysql':
        # MySQL is default, no changes needed for imports
        return False

    # PostgreSQL selected - need to switch imports from mysql to postgresql
    files_to_update = [
        repo_root / "internal" / "infrastructure" / "dependency" / "dependency.go",
        repo_root / "internal" / "infrastructure" / "grpc" / "internal" / "handler" / "public" / "handler.go",
    ]

    modified = False
    for file_path in files_to_update:
        if not file_path.exists():
            continue

        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
        except Exception as e:
            print(f"Warning: Could not read {file_path}: {e}")
            continue

        # Replace mysql with postgresql in import paths
        # Pattern: internal/infrastructure/mysql → internal/infrastructure/postgresql
        new_content = content.replace(
            'internal/infrastructure/mysql',
            'internal/infrastructure/postgresql'
        )

        if new_content != content:
            if not dry_run:
                try:
                    with open(file_path, 'w', encoding='utf-8') as f:
                        f.write(new_content)
                    modified = True
                except Exception as e:
                    print(f"Error: Could not write to {file_path}: {e}")
            else:
                modified = True

    return modified


def toggle_migration_imports(repo_root: Path, database: str, dry_run: bool = False) -> bool:
    """
    Toggle commented migration imports in database_cmd/cmd.go.

    For MySQL (default):
        migration "github.com/.../mysql/migration"
        // migration "github.com/.../postgresql/migration"

    For PostgreSQL:
        // migration "github.com/.../mysql/migration"
        migration "github.com/.../postgresql/migration"

    Args:
        repo_root: Repository root directory
        database: Selected database ('mysql' or 'postgresql')
        dry_run: If True, don't actually modify the file

    Returns:
        True if file was modified, False otherwise
    """
    cmd_path = repo_root / "internal" / "infrastructure" / "cmd" / "internal" / "schema_migration_cmd" / "database_cmd" / "cmd.go"

    if not cmd_path.exists():
        return False

    try:
        with open(cmd_path, 'r', encoding='utf-8') as f:
            lines = f.readlines()
    except Exception as e:
        print(f"Warning: Could not read {cmd_path}: {e}")
        return False

    modified = False
    for i, line in enumerate(lines):
        # Check for migration import lines
        if 'migration "' in line and '/migration"' in line:
            if database == 'mysql':
                # Uncomment mysql, comment postgresql
                if '/mysql/migration"' in line and line.strip().startswith('//'):
                    lines[i] = line.lstrip('# ').lstrip('/')
                    lines[i] = lines[i].lstrip()  # Remove any leading whitespace after removing comment
                    modified = True
                elif '/postgresql/migration"' in line and not line.strip().startswith('//'):
                    lines[i] = '// ' + line
                    modified = True
            else:  # postgresql
                # Comment mysql, uncomment postgresql
                if '/mysql/migration"' in line and not line.strip().startswith('//'):
                    lines[i] = '\t// ' + line.lstrip('\t')
                    modified = True
                elif '/postgresql/migration"' in line and line.strip().startswith('//'):
                    # Remove the comment marker while preserving indentation
                    stripped = line.lstrip()
                    if stripped.startswith('// '):
                        lines[i] = line.replace('// ', '', 1)
                    elif stripped.startswith('//'):
                        lines[i] = line.replace('//', '', 1)
                    modified = True

    if modified and not dry_run:
        try:
            with open(cmd_path, 'w', encoding='utf-8') as f:
                f.writelines(lines)
        except Exception as e:
            print(f"Error: Could not write to {cmd_path}: {e}")
            return False

    return modified


def verify_database_consistency(repo_root: Path, database: str) -> List[str]:
    """
    Verify all database-specific files are consistent.

    Checks:
    1. dependency.go imports match selected database
    2. database_cmd/cmd.go has correct import uncommented
    3. public/handler.go imports match selected database

    Args:
        repo_root: Repository root directory
        database: Selected database ('mysql' or 'postgresql')

    Returns:
        List of warning messages (empty if all checks pass)
    """
    warnings = []
    expected_db = database
    unexpected_db = 'postgresql' if database == 'mysql' else 'mysql'

    # Check dependency.go
    dep_file = repo_root / "internal" / "infrastructure" / "dependency" / "dependency.go"
    if dep_file.exists():
        try:
            content = dep_file.read_text(encoding='utf-8')
            if f'infrastructure/{unexpected_db}' in content:
                warnings.append(f"dependency.go still contains '{unexpected_db}' imports")
        except Exception:
            pass

    # Check database_cmd/cmd.go
    cmd_file = repo_root / "internal" / "infrastructure" / "cmd" / "internal" / "schema_migration_cmd" / "database_cmd" / "cmd.go"
    if cmd_file.exists():
        try:
            content = cmd_file.read_text(encoding='utf-8')
            lines = content.split('\n')

            # Check if expected import is uncommented
            expected_active = False
            unexpected_commented = False

            for line in lines:
                if f'/{expected_db}/migration"' in line:
                    if not line.strip().startswith('//'):
                        expected_active = True
                elif f'/{unexpected_db}/migration"' in line:
                    if line.strip().startswith('//'):
                        unexpected_commented = True

            if not expected_active:
                warnings.append(f"database_cmd/cmd.go: {expected_db} migration import is not active")
            if not unexpected_commented:
                warnings.append(f"database_cmd/cmd.go: {unexpected_db} migration import is not commented out")
        except Exception:
            pass

    # Check public/handler.go
    handler_file = repo_root / "internal" / "infrastructure" / "grpc" / "internal" / "handler" / "public" / "handler.go"
    if handler_file.exists():
        try:
            content = handler_file.read_text(encoding='utf-8')
            if f'infrastructure/{unexpected_db}' in content:
                warnings.append(f"public/handler.go still contains '{unexpected_db}' imports")
        except Exception:
            pass

    return warnings


def main():
    """Main execution function."""
    args = parse_args()

    # Calculate derived values
    go_module = args.go_module
    service_name = args.service_name
    database = args.database
    project_title = args.project_title or service_name.upper().replace('-', ' ')
    buf_org = extract_buf_org(go_module)
    docker_network = f"{service_name}-network"

    # Get repository root
    repo_root = get_repository_root()

    print("=" * 80)
    print("Repository Initialization")
    print("=" * 80)
    print(f"Repository root: {repo_root}")
    print(f"Go module:       {go_module}")
    print(f"Service name:    {service_name}")
    print(f"Database:        {database}")
    print(f"Project title:   {project_title}")
    print(f"Buf org:         {buf_org}")
    print(f"Docker network:  {docker_network}")

    if args.dry_run:
        print("\n*** DRY RUN MODE - No changes will be made ***\n")

    print("\n" + "=" * 80)
    print("Step 1: Renaming proto directories")
    print("=" * 80)

    # Rename proto directory: schema/proto/rapid/ -> schema/proto/{service_name}/
    proto_old_dir = repo_root / "schema" / "proto" / "rapid"
    proto_new_dir = repo_root / "schema" / "proto" / service_name

    if rename_directory(proto_old_dir, proto_new_dir, args.dry_run):
        print(f"✓ Renamed: {proto_old_dir} -> {proto_new_dir}")
    else:
        if proto_old_dir.exists():
            print(f"× Failed to rename: {proto_old_dir}")
        else:
            print(f"ℹ Already renamed or not found: {proto_old_dir}")

    print("\n" + "=" * 80)
    print("Step 2: Performing text replacements")
    print("=" * 80)

    # Define replacements
    replacements = [
        # Go module path
        ("github.com/abyssparanoia/rapid-go", go_module),
        # Proto namespace (after directory rename, update package declarations and imports)
        ("package rapid.", f"package {service_name}."),
        ("import \"rapid/", f"import \"{service_name}/"),
        ("pb/rapid/", f"pb/{service_name}/"),
        # Buf registry
        ("buf.build/abyssparanoia/rapid", f"buf.build/{buf_org}/{service_name}"),
        # Docker network
        ("rapid-go-network", docker_network),
        ("rapid-go_rapid-go-network", f"{service_name}_{docker_network}"),
        # Project title
        ("RAPID GO", project_title),
        ("# RAPID GO", f"# {project_title}"),
        # OpenAPI output paths
        ("./schema/openapi/rapid/", f"./schema/openapi/{service_name}/"),
        # Skill description
        ("codebase investigation for rapid-go", f"codebase investigation for {service_name}"),
        ("Repository: abyssparanoia/rapid-go", f"Repository: {go_module.replace('github.com/', '')}"),
    ]

    # Find all text files (exclude .git, node_modules, vendor, etc.)
    exclude_dirs = ['.git', 'node_modules', 'vendor', 'data', '__pycache__', '.idea', '.vscode']
    text_files = find_all_text_files(repo_root, exclude_dirs)

    print(f"Found {len(text_files)} text files to process")

    modified_count = 0
    for file_path in text_files:
        if replace_in_file(file_path, replacements, args.dry_run):
            modified_count += 1
            rel_path = file_path.relative_to(repo_root)
            print(f"  ✓ Modified: {rel_path}")

    print(f"\n✓ Modified {modified_count} files")

    print("\n" + "=" * 80)
    print(f"Step 3: Database selection - using {database}")
    print("=" * 80)

    # Determine which database to remove
    db_to_remove = 'postgresql' if database == 'mysql' else 'mysql'

    # Delete unused database directories
    db_dir_to_remove = repo_root / "db" / db_to_remove
    infra_dir_to_remove = repo_root / "internal" / "infrastructure" / db_to_remove

    if delete_directory(db_dir_to_remove, args.dry_run):
        print(f"✓ Deleted: {db_dir_to_remove}")
    else:
        print(f"ℹ Already deleted or not found: {db_dir_to_remove}")

    if delete_directory(infra_dir_to_remove, args.dry_run):
        print(f"✓ Deleted: {infra_dir_to_remove}")
    else:
        print(f"ℹ Already deleted or not found: {infra_dir_to_remove}")

    # Toggle Makefile database targets
    makefile_path = repo_root / "Makefile"
    toggle_makefile_targets(makefile_path, database, args.dry_run)
    print(f"✓ Updated Makefile for {database}")

    # Toggle .envrc.tmpl environment variables
    envrc_tmpl_path = repo_root / ".envrc.tmpl"
    toggle_envrc_template(envrc_tmpl_path, database, args.dry_run)
    print(f"✓ Updated .envrc.tmpl for {database}")

    # Remove unused database service from docker-compose.yml
    compose_path = repo_root / "docker-compose.yml"
    remove_docker_compose_db_service(compose_path, database, args.dry_run)
    print(f"✓ Updated docker-compose.yml (removed {db_to_remove} service)")

    # Toggle database imports in Go files
    if toggle_database_imports(repo_root, database, args.dry_run):
        print(f"✓ Updated database imports in dependency.go and public/handler.go for {database}")
    else:
        print(f"ℹ No database import changes needed (using {database} as default)")

    # Toggle migration imports in database_cmd/cmd.go
    if toggle_migration_imports(repo_root, database, args.dry_run):
        print(f"✓ Updated migration imports in database_cmd/cmd.go for {database}")

    print("\n" + "=" * 80)
    print("Step 4: Cleaning up generated code")
    print("=" * 80)

    # Delete generated proto directories
    generated_pb_dir = repo_root / "internal" / "infrastructure" / "grpc" / "pb" / "rapid"
    generated_openapi_dir = repo_root / "schema" / "openapi" / "rapid"

    if delete_directory(generated_pb_dir, args.dry_run):
        print(f"✓ Deleted: {generated_pb_dir}")
    else:
        print(f"ℹ Already deleted or not found: {generated_pb_dir}")

    if delete_directory(generated_openapi_dir, args.dry_run):
        print(f"✓ Deleted: {generated_openapi_dir}")
    else:
        print(f"ℹ Already deleted or not found: {generated_openapi_dir}")

    # Verify database consistency
    print("\n" + "=" * 80)
    print("Step 5: Verifying database consistency")
    print("=" * 80)

    if args.dry_run:
        print("ℹ Skipping verification in dry-run mode (files not modified yet)")
    else:
        warnings = verify_database_consistency(repo_root, database)
        if warnings:
            print("\n⚠️  Warnings detected (please verify manually):")
            for warning in warnings:
                print(f"  • {warning}")
        else:
            print("✓ All database-specific files are consistent")

    print("\n" + "=" * 80)
    print("Initialization Complete!")
    print("=" * 80)

    if not args.dry_run:
        print("\nNext steps:")
        print("1. Review the changes:")
        print("   git status")
        print("   git diff")
        print("")
        print("2. Copy .envrc.tmpl to .envrc and configure:")
        print("   cp .envrc.tmpl .envrc")
        print("")
        print("3. Regenerate code:")
        print("   make migrate.up      # Run migrations + generate SQLBoiler models")
        print("   make generate.buf    # Regenerate proto code")
        print("   make generate.mock   # Regenerate mocks")
        print("")
        print("4. Verify the changes compile:")
        print("   go mod tidy")
        print("   make lint.go")
        print("   make test")
        print("")
        print("5. Commit the changes:")
        print("   git add .")
        print(f"   git commit -m 'Initialize repository from rapid-go template'")
    else:
        print("\nDry run completed. Use without --dry-run to apply changes.")


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\nInterrupted by user")
        sys.exit(1)
    except Exception as e:
        print(f"\n\nError: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc()
        sys.exit(1)
