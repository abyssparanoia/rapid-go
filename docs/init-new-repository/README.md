# Initialize New Repository from rapid-go

## Overview

This guide explains how to initialize a new repository from the rapid-go template by renaming all project-specific identifiers and configuring your database backend.

The initialization process automates:
- Go module import path replacement
- Protocol Buffers namespace changes
- Docker network configuration updates
- Database backend selection (MySQL or PostgreSQL)
- .claude configuration updates

## Prerequisites

- Python 3.x installed
- Copied rapid-go repository to your new project location
- Project details prepared:
  - New Go module path (e.g., `github.com/mycompany/awesome-api`)
  - Service name for proto/Docker (e.g., `awesome`)
  - Database choice (`mysql` or `postgresql`)

## Quick Start

```bash
# Navigate to your copied repository
cd /path/to/your-new-project

# Run initialization script (dry-run first to preview changes)
python3 .claude/skills/init-new-repository/scripts/init_repository.py \
  --go-module github.com/mycompany/awesome-api \
  --service-name awesome \
  --database postgresql \
  --dry-run

# If dry-run looks good, run without --dry-run
python3 .claude/skills/init-new-repository/scripts/init_repository.py \
  --go-module github.com/mycompany/awesome-api \
  --service-name awesome \
  --database postgresql

# Regenerate code
make migrate.up
make generate.buf
make generate.mock

# Verify
make lint.go
make test
```

## Parameters

### --go-module (Required)

New Go module import path.

**Format**: `{host}/{organization}/{repository}`

**Examples:**
- `github.com/mycompany/awesome-api`
- `gitlab.com/myteam/payment-service`

**Used for:**
- Replacing all Go import statements
- Extracting Buf organization name
- Updating go.mod

### --service-name (Required)

Service identifier used throughout the codebase.

**Format**: Lowercase, alphanumeric with hyphens allowed

**Examples:**
- `awesome`
- `payment-service`
- `user-api`

**Used for:**
- Proto package namespace (`package awesome.admin_api.v1`)
- Proto directory structure (`schema/proto/awesome/`)
- Docker network name (`awesome-network`)
- Buf registry name (`buf.build/{org}/awesome`)
- Generated code paths

### --database (Required)

Database backend to use.

**Choices**: `mysql` or `postgresql`

**Effects:**
- **Deletes** all code for the non-selected database
- **Activates** configuration for the selected database
- **Updates** Makefile, .envrc.tmpl, docker-compose.yml

**Note**: Spanner code is preserved regardless of selection (used for read replicas)

### --project-title (Optional)

Human-readable project title.

**Default**: Service name in uppercase with hyphens replaced by spaces

**Examples:**
- If `--service-name awesome` → default is `AWESOME`
- Specify `--project-title "Awesome API"` for custom title

**Used for:**
- .claude/CLAUDE.md header
- Documentation references

## Step-by-Step Guide

### Step 1: Copy rapid-go Repository

```bash
# Clone or copy rapid-go to your new project location
cp -r /path/to/rapid-go /path/to/your-new-project
cd /path/to/your-new-project

# Initialize git (if not already a git repository)
git init
```

### Step 2: Run Initialization Script

**Always use dry-run first** to preview changes:

```bash
python3 .claude/skills/init-new-repository/scripts/init_repository.py \
  --go-module github.com/mycompany/awesome-api \
  --service-name awesome \
  --database postgresql \
  --dry-run
```

Review the output to confirm:
- Proto directory rename: `schema/proto/rapid` → `schema/proto/awesome`
- Number of files to be modified
- Database deletion confirmation

If everything looks correct, run without `--dry-run`:

```bash
python3 .claude/skills/init-new-repository/scripts/init_repository.py \
  --go-module github.com/mycompany/awesome-api \
  --service-name awesome \
  --database postgresql
```

### Step 3: Regenerate Code

After initialization, regenerate all generated code:

```bash
# 1. Copy and configure environment variables
cp .envrc.tmpl .envrc
# Edit .envrc with your database credentials

# 2. Run database migrations and generate SQLBoiler models
make migrate.up

# 3. Regenerate protocol buffer code
make generate.buf

# 4. Regenerate mock files
make generate.mock
```

### Step 4: Verify Changes

```bash
# 1. Check that dependencies are correct
go mod tidy

# 2. Run linter
make lint.go

# 3. Run tests
make test

# 4. Verify server starts
make http.dev
```

## Database Selection

The initialization script provides complete database backend selection:

### MySQL Selection

```bash
--database mysql
```

**Actions:**
- Deletes `db/postgresql/` directory (migrations, constants, sqlboiler config)
- Deletes `internal/infrastructure/postgresql/` directory (client, repositories, marshallers)
- Activates MySQL configuration in:
  - `internal/infrastructure/dependency/dependency.go`
  - `internal/infrastructure/cmd/internal/schema_migration_cmd/database_cmd/cmd.go`
  - `Makefile`
  - `.envrc.tmpl`
  - `docker-compose.yml`

### PostgreSQL Selection

```bash
--database postgresql
```

**Actions:**
- Deletes `db/mysql/` directory
- Deletes `internal/infrastructure/mysql/` directory
- Activates PostgreSQL configuration in all config files

### Spanner

Spanner code (`db/spanner/`, `internal/infrastructure/spanner/`) is **always preserved** regardless of database selection, as it's used for read replicas.

## What Gets Modified

### Directory Structure Changes

| Original | New (example: service-name=awesome) |
|----------|-------------------------------------|
| `schema/proto/rapid/` | `schema/proto/awesome/` |
| `db/{unused-database}/` | Deleted |
| `internal/infrastructure/{unused-database}/` | Deleted |
| `internal/infrastructure/grpc/pb/rapid/` | Deleted (regenerated as `pb/awesome/`) |
| `schema/openapi/rapid/` | Deleted (regenerated as `openapi/awesome/`) |

### File Content Changes

**Go files** (204+ files):
- All import paths updated from `github.com/abyssparanoia/rapid-go` to your module path

**Proto files** (19 files):
- Package declarations: `package rapid.*` → `package {service-name}.*`
- Import paths: `import "rapid/*"` → `import "{service-name}/*"`

**Configuration files**:
- `go.mod` - module declaration
- `buf.gen.yaml` - generated code path
- `schema/proto/buf.yaml` - registry name
- `db/{database}/sqlboiler.toml.tpl` - import paths
- `Makefile` - network names, generation targets
- `docker-compose.yml` - network name, removed unused DB service
- `.envrc.tmpl` - environment variables
- `.golangci.yml` - import path patterns

**.claude files** (30+ files):
- `CLAUDE.md` - project title
- All rule files - example import paths
- All skill files - example import paths
- Command files - repository references

## Troubleshooting

### Script Permission Denied

**Solution**: Make the script executable:
```bash
chmod +x .claude/skills/init-new-repository/scripts/init_repository.py
```

### Generated Code Compilation Errors

**Solution**: Ensure you ran all regeneration steps:
```bash
make migrate.up
make generate.buf
make generate.mock
go mod tidy
```

### Database Connection Errors

**Solution**: Check `.envrc` configuration:
- MySQL: `DB_HOST="tcp(localhost:3306)"`
- PostgreSQL: `DB_HOST="localhost:5432"`

### Docker Compose Network Errors

**Solution**: The network name changed. Restart Docker Compose:
```bash
docker-compose down
docker-compose up -d
```

### Proto Generation Fails

**Solution**: Buf registry name changed. Update `schema/proto/buf.yaml`:
```yaml
name: buf.build/{your-org}/{service-name}
```

## Example: Full Initialization

```bash
# 1. Copy repository
cp -r rapid-go awesome-api
cd awesome-api

# 2. Run initialization (dry-run first)
python3 .claude/skills/init-new-repository/scripts/init_repository.py \
  --go-module github.com/mycompany/awesome-api \
  --service-name awesome \
  --database postgresql \
  --project-title "Awesome API" \
  --dry-run

# 3. Verify dry-run output, then run actual initialization
python3 .claude/skills/init-new-repository/scripts/init_repository.py \
  --go-module github.com/mycompany/awesome-api \
  --service-name awesome \
  --database postgresql \
  --project-title "Awesome API"

# 4. Configure environment
cp .envrc.tmpl .envrc
# Edit .envrc with your credentials

# 5. Regenerate code
make migrate.up
make generate.buf
make generate.mock

# 6. Verify
go mod tidy
make lint.go
make test

# 7. Commit changes
git add .
git commit -m "Initialize repository from rapid-go template"
```

## Result

After initialization, your repository will have:
- All imports pointing to your new module path
- Proto namespace matching your service name
- Docker network configured with your service name
- Only your selected database backend (MySQL or PostgreSQL)
- All generated code cleaned and ready for regeneration
- .claude configuration updated with your project details

## Next Steps

- Configure external services (AWS/GCP credentials in `.envrc`)
- Set up authentication (Cognito or Firebase)
- Add your domain-specific entities following the CRUD workflow
- Update README.md with project-specific information
- Set up CI/CD pipelines

## Related Documentation

- [Development Setup Guide](../development-setup/README.md)
- [Skill Documentation](../../.claude/skills/init-new-repository/SKILL.md)
- [Replacement Patterns Reference](../../.claude/skills/init-new-repository/references/patterns.md)
