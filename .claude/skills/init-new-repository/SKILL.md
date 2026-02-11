---
name: init-new-repository
description: "Initialize a new repository from rapid-go template by renaming all identifiers (Go modules, proto namespaces, Docker networks, .claude files) and selecting database backend (MySQL or PostgreSQL). Use when porting rapid-go to a new project or creating a new microservice based on rapid-go architecture. Automates: Go import path replacement, proto namespace changes, database code removal, Docker configuration updates, and .claude documentation updates."
---

# Initialize New Repository from rapid-go Template

This skill automates the process of creating a new repository from the rapid-go template by:
- Renaming all rapid-go specific identifiers to your new project name
- Selecting and configuring database backend (MySQL or PostgreSQL)
- Updating all configuration files and documentation
- Cleaning up generated code for regeneration

## Prerequisites

Before running this skill, ensure you have:

1. **Copied the rapid-go repository** to your new project location
2. **Decided on your project details**:
   - Go module path (e.g., `github.com/mycompany/awesome-api`)
   - Service name (e.g., `awesome` - used for proto namespace and Docker network)
   - Database choice (`mysql` or `postgresql`)
   - Project title (optional - defaults to service name in uppercase)

3. **Python 3.x installed** - The initialization script requires Python 3

## Workflow

### Step 1: Prepare Project Information

Gather the following information:

| Parameter | Example | Purpose |
|-----------|---------|---------|
| `--go-module` | `github.com/mycompany/awesome-api` | New Go module import path |
| `--service-name` | `awesome` | Service identifier (proto namespace, Docker network) |
| `--database` | `postgresql` | Database to use (`mysql` or `postgresql`) |
| `--project-title` | `Awesome API` | Human-readable project name (optional) |

**Derived values** (automatically calculated):
- Buf organization: Extracted from `--go-module` (e.g., `mycompany` from `github.com/mycompany/awesome-api`)
- Docker network: `{service-name}-network` (e.g., `awesome-network`)

### Step 2: Run Initialization Script

Navigate to your copied repository and run:

```bash
python3 .claude/skills/init-new-repository/scripts/init_repository.py \
  --go-module github.com/mycompany/awesome-api \
  --service-name awesome \
  --database postgresql
```

**Dry-run mode** (preview changes without applying):
```bash
python3 .claude/skills/init-new-repository/scripts/init_repository.py \
  --go-module github.com/mycompany/awesome-api \
  --service-name awesome \
  --database postgresql \
  --dry-run
```

#### What the script does:

**1. Renames proto directories:**
- `schema/proto/rapid/` → `schema/proto/{service-name}/`

**2. Performs text replacements across all files:**
- Go imports: `github.com/abyssparanoia/rapid-go` → your module path
- Proto packages: `package rapid.*` → `package {service-name}.*`
- Proto imports: `import "rapid/*"` → `import "{service-name}/*"`
- Docker networks: `rapid-go-network` → `{service-name}-network`
- Buf registry: `buf.build/abyssparanoia/rapid` → `buf.build/{org}/{service-name}`
- Project titles: `RAPID GO` → your project title
- .claude references: Updates all documentation

**3. Database selection:**
- Deletes unused database code:
  - If `mysql`: removes `db/postgresql/` and `internal/infrastructure/postgresql/`
  - If `postgresql`: removes `db/mysql/` and `internal/infrastructure/mysql/`
- Updates Go import statements:
  - `internal/infrastructure/dependency/dependency.go` - switches database import aliases (`mysql` ↔ `postgresql`)
  - `internal/infrastructure/grpc/internal/handler/public/handler.go` - switches database imports
  - `internal/infrastructure/cmd/internal/schema_migration_cmd/database_cmd/cmd.go` - toggles commented migration imports
- Updates configuration files:
  - `Makefile` - database-specific generation targets
  - `.envrc.tmpl` - environment variables
  - `docker-compose.yml` - removes unused database service
- Updates .claude documentation:
  - `.claude/rules/repository.md` - path pattern examples
  - `.claude/rules/dependency-injection.md` - import examples
  - `.claude/skills/add-domain-entity/references/repository-patterns.md` - import examples
  - `.claude/skills/add-domain-entity/references/marshaller-patterns.md` - import examples

**4. Cleans up generated code:**
- Deletes `internal/infrastructure/grpc/pb/rapid/` (will be regenerated)
- Deletes `schema/openapi/rapid/` (will be regenerated)

**5. Verifies database consistency:**
- Checks that all database-specific files are consistent with selected database
- Reports warnings if any inconsistencies are detected (script continues normally)

### Step 3: Regenerate Code

After initialization completes, regenerate all generated code:

```bash
# 1. Copy environment template and configure
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
# 1. Review all changes
git status
git diff

# 2. Ensure dependencies are correct
go mod tidy

# 3. Run linter
make lint.go

# 4. Run tests
make test

# 5. Verify server starts
make http.dev
```

### Step 5: Commit Changes

```bash
git add .
git commit -m "Initialize repository from rapid-go template"
```

## Parameters Reference

### --go-module (Required)

New Go module import path.

**Format**: `{host}/{organization}/{repository}`

**Examples:**
- `github.com/mycompany/awesome-api`
- `gitlab.com/myteam/payment-service`
- `bitbucket.org/myorg/user-service`

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

## Common Issues & Troubleshooting

### Issue: Script fails with "Permission denied"

**Solution**: Make the script executable:
```bash
chmod +x .claude/skills/init-new-repository/scripts/init_repository.py
```

### Issue: Binary files are corrupted

**Solution**: The script should skip binary files automatically. If issues persist, check the `binary_extensions` list in the script.

### Issue: Generated code compilation errors after initialization

**Solution**: Ensure you ran all regeneration steps:
```bash
make migrate.up
make generate.buf
make generate.mock
go mod tidy
```

### Issue: Database connection errors

**Solution**: Check `.envrc` configuration:
- MySQL: `DB_HOST="tcp(localhost:3306)"`
- PostgreSQL: `DB_HOST="localhost:5432"`

### Issue: Docker Compose network errors

**Solution**: The network name changed. Restart Docker Compose:
```bash
docker-compose down
docker-compose up -d
```

## What Files Are Modified

### Directory Structure Changes

| Original | New |
|----------|-----|
| `schema/proto/rapid/` | `schema/proto/{service-name}/` |
| `db/{unused-database}/` | Deleted |
| `internal/infrastructure/{unused-database}/` | Deleted |
| `internal/infrastructure/grpc/pb/rapid/` | Deleted (regenerated) |
| `schema/openapi/rapid/` | Deleted (regenerated) |

### File Content Changes

**Go files** (204+ files):
- All import paths updated

**Proto files** (19 files):
- Package declarations
- Import paths

**Configuration files**:
- `go.mod` - module declaration
- `buf.gen.yaml` - generated code path
- `schema/proto/buf.yaml` - registry name
- `db/{database}/sqlboiler.toml.tpl` - import paths
- `Makefile` - network names, generation targets
- `docker-compose.yml` - network name, removed unused DB service
- `.envrc.tmpl` - environment variables
- `.golangci.yml` - import path patterns
- Docker files - working directories

**.claude files** (30+ files):
- `CLAUDE.md` - project title
- All rule files - example import paths
- All skill files - example import paths
- Command files - repository references

For detailed replacement patterns, see [references/patterns.md](references/patterns.md).

## After Initialization

Your repository is now ready for development. Common next steps:

1. **Configure external services**:
   - Update `.envrc` with your AWS/GCP credentials
   - Configure authentication (Cognito or Firebase)
   - Set up storage buckets (S3 or GCS)

2. **Customize for your domain**:
   - Add new entities following the CRUD workflow
   - Modify proto definitions for your API
   - Implement business logic in domain services

3. **Set up CI/CD**:
   - Configure GitHub Actions or equivalent
   - Set up deployment pipelines
   - Configure environment-specific variables

4. **Update documentation**:
   - Modify README.md with project-specific info
   - Update .claude/ rules if needed
   - Add project-specific skills

## Related Skills

- `add-database-table` - Add new database tables
- `add-domain-entity` - Create domain models and repositories
- `add-api-endpoint` - Implement new API endpoints

## Related Rules

See `.claude/rules/` for detailed guidelines:
- `domain-model.md` - Domain model patterns
- `repository.md` - Repository implementation
- `usecase-interactor.md` - Business logic patterns
- `grpc-handler.md` - API handler patterns
