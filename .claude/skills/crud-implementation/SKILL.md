---
name: crud-implementation
description: REQUIRED when implementing a new entity with CRUD operations. Read this FIRST to understand the full workflow, then follow referenced skills in order. Covers database, domain, and API layers.
---

# CRUD Implementation Guide

This guide provides the complete workflow for implementing a new entity with CRUD operations.

## Implementation Flow

```
┌─────────────────────┐
│  add-database-table │  ← Step 1: Migration & Constants
└──────────┬──────────┘
           │
┌──────────▼──────────┐
│  add-domain-entity  │  ← Step 2: Model, Repository, Marshaller
└──────────┬──────────┘
           │
┌──────────▼──────────┐
│  add-api-endpoint   │  ← Step 3: Usecase, Proto, Handler, DI
└─────────────────────┘
```

## Step-by-Step

### Step 1: Database Layer

Use the **add-database-table** skill:

1. Create migration file: `make migrate.create`
2. Write table DDL with indexes and foreign keys
3. Create constant table if entity has status/enum fields
4. Define YAML constants in `db/postgresql/constants/`
5. Run `make migrate.up` to apply and generate SQLBoiler

### Step 2: Domain Layer

Use the **add-domain-entity** skill:

1. Create domain model in `internal/domain/model/`
2. Add domain error in `internal/domain/errors/`
3. Define repository interface in `internal/domain/repository/`
4. Create marshaller in `internal/infrastructure/{mysql|postgresql|spanner}/internal/marshaller/`
5. Implement repository in `internal/infrastructure/{mysql|postgresql|spanner}/repository/`
6. Generate mocks: `make generate.mock`

### Step 3: API Layer

Use the **add-api-endpoint** skill:

1. Define input/output in `internal/usecase/input/` and `output/`
2. Create interactor interface in `internal/usecase/`
3. Implement interactor in `internal/usecase/`
4. Define proto messages in `schema/proto/`
5. Generate proto: `make generate.buf`
6. Implement gRPC handler in `internal/infrastructure/grpc/internal/handler/`
7. Register in `internal/infrastructure/dependency/dependency.go`

## Quick Commands Reference

| Step | Command |
|------|---------|
| Create migration | `make migrate.create` |
| Apply migration | `make migrate.up` |
| Generate proto | `make generate.buf` |
| Generate mocks | `make generate.mock` |
| Run tests | `make test` |
| Lint code | `make lint.go` |

## Example: Adding "Example" Entity

```bash
# 1. Create migration
make migrate.create
# Edit: db/postgresql/migrations/YYYYMMDD_create_examples.sql

# 2. Apply migration (generates SQLBoiler)
make migrate.up

# 3. Create domain components
# - internal/domain/model/example.go
# - internal/domain/errors/errors.go (add ExampleNotFoundErr)
# - internal/domain/repository/example.go
# - internal/infrastructure/{mysql|postgresql|spanner}/internal/marshaller/example.go
# - internal/infrastructure/{mysql|postgresql|spanner}/repository/example.go

# 4. Generate mocks
make generate.mock

# 5. Create usecase components
# - internal/usecase/input/admin_example.go
# - internal/usecase/output/admin_example.go
# - internal/usecase/admin_example.go
# - internal/usecase/admin_example_impl.go

# 6. Create API components
# - schema/proto/rapid/admin_api/v1/example.proto
# Generate: make generate.buf
# - internal/infrastructure/grpc/internal/handler/admin/example.go

# 7. Register dependencies
# - internal/infrastructure/dependency/dependency.go

# 8. Test
make test
```

## See Also

- **add-database-table** - Detailed migration and constant table creation
- **add-domain-entity** - Detailed domain layer implementation
- **add-api-endpoint** - Detailed API layer implementation
