---
name: add-database-table
description: REQUIRED Step 1 of CRUD workflow. Use when creating migration files, defining table schemas, or adding constant tables (enums). Run 'make migrate.create' and follow this guide.
---

# Add Database Table

This skill guides you through creating a new database table with migrations.

## Prerequisites

- Ensure the database is running locally
- Know the entity name and its fields

## Step 1: Create Migration File

```bash
make migrate.create
```

A new SQL file will be created in `db/postgresql/migrations/` (or `db/mysql/`, `db/spanner/` depending on target database).

## Step 2: Write Migration SQL

### Main Table

```sql
-- +goose Up
CREATE TABLE "examples" (
    "id"            VARCHAR(64)     PRIMARY KEY,
    "tenant_id"     VARCHAR(64)     NOT NULL,
    "name"          VARCHAR(256)    NOT NULL,
    "description"   TEXT            NOT NULL,
    "status"        VARCHAR(64)     NOT NULL,
    "order"         INTEGER         NOT NULL DEFAULT 0,
    "created_at"    TIMESTAMPTZ     NOT NULL,
    "updated_at"    TIMESTAMPTZ     NOT NULL,

    CONSTRAINT "examples_fkey_tenant_id"
        FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id"),
    CONSTRAINT "examples_fkey_status"
        FOREIGN KEY ("status") REFERENCES "example_statuses" ("id")
);

CREATE INDEX "examples_idx_tenant_id" ON "examples" ("tenant_id");
CREATE INDEX "examples_idx_status" ON "examples" ("status");
CREATE INDEX "examples_idx_created_at" ON "examples" ("created_at" DESC);

-- +goose Down
DROP TABLE IF EXISTS "examples";
```

### Column Type Reference

| Go Type | PostgreSQL Type |
|---------|-----------------|
| `string` (ID) | `VARCHAR(64)` |
| `string` (short) | `VARCHAR(256)` |
| `string` (long) | `TEXT` |
| `int` | `INTEGER` |
| `bool` | `BOOLEAN` |
| `time.Time` | `TIMESTAMPTZ` |
| `map/struct` | `JSONB` |

## Step 3: Create Constant Table (if needed)

If the entity has status/enum fields, create a constant table:

### Migration SQL

```sql
-- +goose Up
CREATE TABLE "example_statuses" (
    "id" VARCHAR(64) PRIMARY KEY
);

-- +goose Down
DROP TABLE IF EXISTS "example_statuses";
```

### YAML Definition

Create `db/postgresql/constants/example_statuses.yaml`:

```yaml
table: example_statuses
records:
  - id: draft
  - id: published
  - id: archived
```

## Step 4: Run Migration

```bash
make migrate.up
```

This will:
- Execute migrations
- Sync constant values
- Generate SQLBoiler code in `internal/infrastructure/{mysql|postgresql}/internal/dbmodel/`

## Step 5: Verify

Check that the following file was generated:
- `internal/infrastructure/{mysql|postgresql}/internal/dbmodel/{table}.go`

## Checklist

- [ ] Migration file created
- [ ] Table has appropriate indexes
- [ ] Foreign key constraints added
- [ ] Constant table created (if needed)
- [ ] YAML constants defined (if needed)
- [ ] `make migrate.up` executed successfully
- [ ] SQLBoiler model generated

## Next Steps

After creating the table, use the **add-domain-entity** skill to create:
- Domain model
- Repository interface
- Repository implementation
- Marshaller
