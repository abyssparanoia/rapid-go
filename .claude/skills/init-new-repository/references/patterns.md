# Replacement Patterns Reference

This document provides detailed information about all the text replacement patterns used during repository initialization.

## Table of Contents

1. [Go Module Paths](#go-module-paths)
2. [Proto Definitions](#proto-definitions)
3. [Docker Configuration](#docker-configuration)
4. [Configuration Files](#configuration-files)
5. [Database-Specific Patterns](#database-specific-patterns)
6. [Documentation Updates](#documentation-updates)

---

## Go Module Paths

### Pattern: Import Statements

**Old:** `github.com/abyssparanoia/rapid-go`
**New:** User-specified Go module path

**Affected files:** 204+ `.go` files

**Examples:**

```go
// Before
import (
    "github.com/abyssparanoia/rapid-go/internal/domain/model"
    "github.com/abyssparanoia/rapid-go/internal/domain/repository"
)

// After (with --go-module github.com/mycompany/awesome-api)
import (
    "github.com/mycompany/awesome-api/internal/domain/model"
    "github.com/mycompany/awesome-api/internal/domain/repository"
)
```

### Pattern: go.mod Module Declaration

**File:** `go.mod` (line 1)

```go
// Before
module github.com/abyssparanoia/rapid-go

// After
module github.com/mycompany/awesome-api
```

### Pattern: SQLBoiler Custom Types

**Files:**
- `db/mysql/sqlboiler.toml.tpl` (lines 44, 55)
- `db/postgresql/sqlboiler.toml.tpl` (lines 44, 55)

```toml
# Before
import-override.date = "github.com/abyssparanoia/rapid-go/db/mysql/custom_types.Date"

# After
import-override.date = "github.com/mycompany/awesome-api/db/mysql/custom_types.Date"
```

### Pattern: Golangci-lint Configuration

**File:** `.golangci.yml` (lines 75-80, 136-138)

```yaml
# Before
exhaustruct:
  include:
    - github.com/abyssparanoia/rapid-go/internal/domain/model.*
    - github.com/abyssparanoia/rapid-go/internal/usecase/input.*

# After
exhaustruct:
  include:
    - github.com/mycompany/awesome-api/internal/domain/model.*
    - github.com/mycompany/awesome-api/internal/usecase/input.*
```

---

## Proto Definitions

### Pattern: Package Declarations

**Old:** `package rapid.{api_type}.v1`
**New:** `package {service-name}.{api_type}.v1`

**Affected files:** All `.proto` files in `schema/proto/rapid/`

**Examples:**

```protobuf
// Before
package rapid.admin_api.v1;

// After (with --service-name awesome)
package awesome.admin_api.v1;
```

### Pattern: Import Paths

**Old:** `import "rapid/{path}"`
**New:** `import "{service-name}/{path}"`

**Examples:**

```protobuf
// Before
import "rapid/admin_api/v1/model_tenant.proto";

// After (with --service-name awesome)
import "awesome/admin_api/v1/model_tenant.proto";
```

### Pattern: Generated Code Imports

**Old:** `pb/rapid/{path}`
**New:** `pb/{service-name}/{path}`

**Affected:** All Go files importing generated proto code

**Examples:**

```go
// Before
import admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"

// After (combined with Go module change)
import admin_apiv1 "github.com/mycompany/awesome-api/internal/infrastructure/grpc/pb/awesome/admin_api/v1"
```

### Pattern: Directory Structure

**Old path:** `schema/proto/rapid/`
**New path:** `schema/proto/{service-name}/`

**Subdirectories affected:**
- `admin_api/v1/`
- `staff_api/v1/`
- `public_api/v1/`
- `debug_api/v1/`

### Pattern: Buf Registry Name

**File:** `schema/proto/buf.yaml` (line 2)

```yaml
# Before
name: buf.build/abyssparanoia/rapid

# After (with --go-module github.com/mycompany/awesome-api, --service-name awesome)
name: buf.build/mycompany/awesome
```

**Organization extraction logic:**
- From `github.com/{org}/{repo}` → extract `{org}`
- Example: `github.com/mycompany/awesome-api` → `mycompany`

### Pattern: Buf Code Generation Config

**File:** `buf.gen.yaml` (line 5)

```yaml
# Before
default: github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb

# After (with --go-module github.com/mycompany/awesome-api)
default: github.com/mycompany/awesome-api/internal/infrastructure/grpc/pb
```

### Pattern: OpenAPI Output Paths

**File:** `Makefile` (lines 25-26)

```makefile
# Before
./schema/openapi/rapid/admin_api/v1
./schema/openapi/rapid/public_api/v1

# After (with --service-name awesome)
./schema/openapi/awesome/admin_api/v1
./schema/openapi/awesome/public_api/v1
```

---

## Docker Configuration

### Pattern: Docker Network Name

**Old:** `rapid-go-network`
**New:** `{service-name}-network`

**Affected files:**
- `docker-compose.yml` (lines 14, 27, 37, 56, 67, 77, 95, 98)
- `Makefile` (lines 113, 134)

**Examples:**

```yaml
# docker-compose.yml - Before
networks:
  - rapid-go-network

networks:
  rapid-go-network:

# After (with --service-name awesome)
networks:
  - awesome-network

networks:
  awesome-network:
```

```makefile
# Makefile - Before
--network rapid-go_rapid-go-network

# After (with --service-name awesome)
--network awesome_awesome-network
```

**Note:** Docker Compose prefixes network names with the project directory name.
If your project is in `awesome/`, the full network name becomes `awesome_awesome-network`.

### Pattern: Docker Workdir

**Files:**
- `docker/development.Dockerfile` (line 3)
- `docker/production.Dockerfile` (lines 3, 16)

```dockerfile
# Before
WORKDIR /go/src/github.com/abyssparanoia/rapid-go/

# After (with --go-module github.com/mycompany/awesome-api)
WORKDIR /go/src/github.com/mycompany/awesome-api/
```

---

## Configuration Files

### Pattern: Project Title

**Old:** `RAPID GO`
**New:** User-specified project title (or service name uppercase)

**Affected files:**
- `.claude/CLAUDE.md` (line 1)
- Various documentation files

**Examples:**

```markdown
# Before
# RAPID GO

# After (with --service-name awesome, no --project-title specified)
# AWESOME

# After (with --project-title "Awesome API")
# Awesome API
```

### Pattern: Environment Variables Template

**File:** `.envrc.tmpl` (lines 5-15)

**Behavior:** Toggle comments based on database selection

```bash
# Before (MySQL active)
# for mysql
export DB_HOST="tcp(localhost:3306)"
export DB_USER="root"
export DB_PASSWORD="password"
export DB_DATABASE="maindb"

# for postgresql
# export DB_HOST="localhost:5432"
# export DB_USER="postgres"
# export DB_PASSWORD="postgres"
# export DB_DATABASE="maindb"

# After (PostgreSQL selected with --database postgresql)
# for mysql
# export DB_HOST="tcp(localhost:3306)"
# export DB_USER="root"
# export DB_PASSWORD="password"
# export DB_DATABASE="maindb"

# for postgresql
export DB_HOST="localhost:5432"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_DATABASE="maindb"
```

---

## Database-Specific Patterns

### MySQL Selection (`--database mysql`)

**Deleted directories:**
- `db/postgresql/`
- `internal/infrastructure/postgresql/`

**Modified files:**

**`internal/infrastructure/dependency/dependency.go`:**
```go
// Imports remain as-is (already using mysql)
database "github.com/{org}/{repo}/internal/infrastructure/mysql"
database_cache "github.com/{org}/{repo}/internal/infrastructure/mysql/cache"
database_repository "github.com/{org}/{repo}/internal/infrastructure/mysql/repository"
database_transactable "github.com/{org}/{repo}/internal/infrastructure/mysql/transactable"
```

**`internal/infrastructure/cmd/internal/schema_migration_cmd/database_cmd/cmd.go`:**
```go
// Already uncommented
migration "github.com/{org}/{repo}/internal/infrastructure/mysql/migration"
// migration "github.com/{org}/{repo}/internal/infrastructure/postgresql/migration"
```

**`Makefile`:**
```makefile
# Uncommented
make generate.mermaid.mysql
make generate.sqlboiler.mysql

# Commented out
# make generate.mermaid.postgresql
# make generate.sqlboiler.postgresql
```

**`docker-compose.yml`:**
- PostgreSQL service block removed
- MySQL service retained

### PostgreSQL Selection (`--database postgresql`)

**Deleted directories:**
- `db/mysql/`
- `internal/infrastructure/mysql/`

**Modified files:**

**`internal/infrastructure/dependency/dependency.go`:**
```go
// Imports updated from mysql to postgresql
database "github.com/{org}/{repo}/internal/infrastructure/postgresql"
database_cache "github.com/{org}/{repo}/internal/infrastructure/postgresql/cache"
database_repository "github.com/{org}/{repo}/internal/infrastructure/postgresql/repository"
database_transactable "github.com/{org}/{repo}/internal/infrastructure/postgresql/transactable"
```

**`internal/infrastructure/cmd/internal/schema_migration_cmd/database_cmd/cmd.go`:**
```go
// Commented out
// migration "github.com/{org}/{repo}/internal/infrastructure/mysql/migration"
// Uncommented
migration "github.com/{org}/{repo}/internal/infrastructure/postgresql/migration"
```

**`Makefile`:**
```makefile
# Commented out
# make generate.mermaid.mysql
# make generate.sqlboiler.mysql

# Uncommented
make generate.mermaid.postgresql
make generate.sqlboiler.postgresql
```

**`docker-compose.yml`:**
- MySQL service block removed
- PostgreSQL service retained

### Spanner (Always Retained)

**Directories preserved:**
- `db/spanner/`
- `internal/infrastructure/spanner/`

Spanner is used for read replicas and is independent of the primary database selection.

---

## Documentation Updates

### Pattern: .claude/ Files

**Affected files:** 30+ files in `.claude/` directory

**Types of updates:**

1. **Import path examples:**
```go
// Before
import "github.com/abyssparanoia/rapid-go/internal/domain/model"

// After
import "github.com/mycompany/awesome-api/internal/domain/model"
```

2. **Proto path examples:**
```go
// Before
admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"

// After
admin_apiv1 "github.com/mycompany/awesome-api/internal/infrastructure/grpc/pb/awesome/admin_api/v1"
```

3. **Repository references:**
```markdown
# Before
Repository: abyssparanoia/rapid-go

# After
Repository: mycompany/awesome-api
```

4. **Skill descriptions:**
```yaml
# Before
description: Efficient codebase investigation for rapid-go

# After
description: Efficient codebase investigation for awesome
```

5. **Database-specific paths (PostgreSQL selection only):**
```go
// Before
import "github.com/mycompany/awesome-api/internal/infrastructure/mysql/repository"

// After (when postgresql is selected)
import "github.com/mycompany/awesome-api/internal/infrastructure/postgresql/repository"
```

**Database-specific files updated (PostgreSQL selection):**
- `.claude/rules/repository.md` - path pattern examples
- `.claude/rules/dependency-injection.md` - import examples
- `.claude/skills/add-domain-entity/references/repository-patterns.md` - import examples
- `.claude/skills/add-domain-entity/references/marshaller-patterns.md` - import examples

**Note:** Files under `.claude/skills/init-new-repository/` maintain both MySQL and PostgreSQL examples as templates.

**Updated files include:**
- `CLAUDE.md`
- All files in `rules/`
- All files in `skills/*/`
- All files in `commands/`

---

## Generated Code Cleanup

The script deletes the following generated code directories (they will be regenerated):

| Directory | Purpose | Regenerate with |
|-----------|---------|-----------------|
| `internal/infrastructure/grpc/pb/rapid/` | Generated proto code | `make generate.buf` |
| `schema/openapi/rapid/` | Generated OpenAPI specs | `make generate.buf` |
| `internal/infrastructure/mysql/internal/dbmodel/` | SQLBoiler models (MySQL) | `make migrate.up` |
| `internal/infrastructure/postgresql/internal/dbmodel/` | SQLBoiler models (PostgreSQL) | `make migrate.up` |

**Why delete?**
- Directory names have changed (rapid → {service-name})
- Import paths have changed
- Fresh regeneration ensures consistency

---

## Special Cases & Edge Cases

### Case 1: Docker Compose Network Prefix

Docker Compose automatically prefixes network names with the project directory name.

**Example:**
- If project is in directory `awesome/`
- Network defined as `awesome-network`
- Actual network name becomes `awesome_awesome-network`

The script handles this by replacing both patterns:
- `rapid-go-network` → `{service-name}-network` (in docker-compose.yml)
- `rapid-go_rapid-go-network` → `{service-name}_{service-name}-network` (in Makefile)

### Case 2: Proto Package vs Directory

Proto package names use dots, but directories use slashes:

**Package:** `package awesome.admin_api.v1`
**Directory:** `schema/proto/awesome/admin_api/v1/`
**Import:** `import "awesome/admin_api/v1/api.proto"`

The script correctly handles both patterns.

### Case 3: Binary Files

The script automatically skips binary files based on extension:
- Images: `.png`, `.jpg`, `.jpeg`, `.gif`
- Archives: `.zip`, `.tar`, `.gz`
- Executables: `.exe`, `.bin`, `.so`, `.dylib`
- Skills: `.skill`

If a file is incorrectly processed as text, it will show a `UnicodeDecodeError` and be skipped.

### Case 4: Hidden Files

The script processes specific hidden files but skips others:

**Processed:**
- `.envrc`
- `.envrc.tmpl`
- `.golangci.yml`
- `.gitignore`
- `.dockerignore`

**Skipped:**
- `.git/` (entire directory)
- Other hidden files/directories

### Case 5: Buf Organization Extraction

The script extracts organization from the second path component:

| Go Module | Extracted Org |
|-----------|---------------|
| `github.com/mycompany/awesome-api` | `mycompany` |
| `gitlab.com/myteam/payment-service` | `myteam` |
| `example.com/payment-service` | `payment-service` |

For single-component paths, the entire path is used as org.

---

## Troubleshooting Patterns

### Issue: Incomplete Replacement

**Symptom:** Some files still reference `rapid-go`

**Cause:** File was in excluded directory or is binary

**Solution:**
1. Check if file is in `.git/`, `node_modules/`, `vendor/`, or `data/`
2. Check file extension - is it binary?
3. Manually review and update if needed

### Issue: Proto Compilation Errors

**Symptom:** `make generate.buf` fails after initialization

**Cause:** Proto directory not renamed or imports not updated

**Solution:**
1. Verify `schema/proto/{service-name}/` exists
2. Check proto files have correct package declarations
3. Check `buf.gen.yaml` has correct import path

### Issue: Database Migration Fails

**Symptom:** `make migrate.up` fails

**Cause:** Wrong database selected or .envrc not configured

**Solution:**
1. Verify correct database directory exists (`db/mysql/` or `db/postgresql/`)
2. Check `.envrc` has correct `DB_HOST` format:
   - MySQL: `tcp(localhost:3306)`
   - PostgreSQL: `localhost:5432`
3. Verify Docker Compose has correct service running

### Issue: Import Cycle

**Symptom:** Compilation fails with import cycle

**Cause:** Stale generated code

**Solution:**
```bash
# Clean all generated code
rm -rf internal/infrastructure/grpc/pb/{service-name}/
rm -rf internal/infrastructure/{database}/internal/dbmodel/
rm -rf schema/openapi/{service-name}/

# Regenerate
make migrate.up
make generate.buf
make generate.mock
```

---

## Database Consistency Verification

After initialization, the script performs automatic verification to ensure all database-specific files are consistent with the selected database.

### What is Checked

1. **dependency.go imports**
   - Verifies database import aliases match selected database
   - Example (PostgreSQL): Should contain `internal/infrastructure/postgresql`, not `mysql`

2. **database_cmd/cmd.go migration imports**
   - Verifies correct migration import is active (uncommented)
   - Verifies incorrect migration import is commented out

3. **public/handler.go imports**
   - Verifies database import matches selected database

### Warning Output Example

If inconsistencies are detected, warnings are displayed but the script continues normally:

```
========================================================================
Step 5: Verifying database consistency
========================================================================

⚠️  Warnings detected (please verify manually):
  • dependency.go still contains 'mysql' imports
  • database_cmd/cmd.go: postgresql migration import is not active

========================================================================
Initialization Complete!
========================================================================
```

### Manual Verification

If warnings are shown, manually check the affected files:

```bash
# Check dependency.go imports
grep -n "infrastructure/mysql\|infrastructure/postgresql" \
  internal/infrastructure/dependency/dependency.go

# Check migration imports
grep -n "migration.*migration\"" \
  internal/infrastructure/cmd/internal/schema_migration_cmd/database_cmd/cmd.go

# Check public handler imports
grep -n "infrastructure/mysql\|infrastructure/postgresql" \
  internal/infrastructure/grpc/internal/handler/public/handler.go
```

---

## Replacement Pattern Summary

| Pattern | Old Value | New Value | Files Affected |
|---------|-----------|-----------|----------------|
| Go module | `github.com/abyssparanoia/rapid-go` | User-specified | 204+ |
| Proto namespace | `rapid` | `{service-name}` | 19+ |
| Proto directory | `schema/proto/rapid/` | `schema/proto/{service-name}/` | 1 |
| Docker network | `rapid-go-network` | `{service-name}-network` | 2 |
| Buf registry | `buf.build/abyssparanoia/rapid` | `buf.build/{org}/{service-name}` | 1 |
| Project title | `RAPID GO` | User-specified or uppercase service name | 5+ |
| Database code | Both MySQL & PostgreSQL | Selected database only | 100+ |

**Total estimated file changes:** 250+ files
