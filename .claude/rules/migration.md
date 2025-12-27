---
description: Database migration file conventions and patterns
globs:
  - "db/mysql/migrations/**/*.sql"
  - "db/mysql/constants/**/*.yaml"
  - "db/postgresql/migrations/**/*.sql"
  - "db/postgresql/constants/**/*.yaml"
  - "db/spanner/migrations/**/*.sql"
---

# Migration Guidelines

## Database Structure

This project supports multiple database backends:

```
db/
├── mysql/
│   ├── migrations/     # MySQL migration files
│   └── constants/      # Constant table definitions (YAML)
├── postgresql/
│   ├── migrations/     # PostgreSQL migration files
│   └── constants/      # Constant table definitions (YAML)
└── spanner/
    └── migrations/     # Spanner migration files
```

## Creating Migrations

```bash
make migrate.create
```

This creates a new timestamped SQL file in the appropriate `db/{database}/migrations/` directory.

## File Naming

Migration files follow the pattern: `{timestamp}_{description}.sql`

Example: `20240115120000_create_examples_table.sql`

## Migration Structure

```sql
-- +goose Up
-- SQL statements for applying the migration

-- +goose Down
-- SQL statements for reverting the migration
```

## Table Creation Pattern

```sql
-- +goose Up
CREATE TABLE "examples" (
    "id"            VARCHAR(64)     PRIMARY KEY,
    "tenant_id"     VARCHAR(64)     NOT NULL,
    "name"          VARCHAR(256)    NOT NULL,
    "description"   TEXT            NOT NULL,
    "status"        VARCHAR(64)     NOT NULL,
    "order"         INTEGER         NOT NULL DEFAULT 0,
    "is_active"     BOOLEAN         NOT NULL DEFAULT TRUE,
    "metadata"      JSONB,
    "created_at"    TIMESTAMPTZ     NOT NULL,
    "updated_at"    TIMESTAMPTZ     NOT NULL,

    -- Foreign key constraints
    CONSTRAINT "examples_fkey_tenant_id"
        FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id"),
    CONSTRAINT "examples_fkey_status"
        FOREIGN KEY ("status") REFERENCES "example_statuses" ("id")
);

-- Indexes for common query patterns
CREATE INDEX "examples_idx_tenant_id" ON "examples" ("tenant_id");
CREATE INDEX "examples_idx_status" ON "examples" ("status");
CREATE INDEX "examples_idx_created_at" ON "examples" ("created_at" DESC);

-- +goose Down
DROP TABLE IF EXISTS "examples";
```

## Column Type Guidelines

| Go Type | PostgreSQL Type | Notes |
|---------|-----------------|-------|
| `string` (ID) | `VARCHAR(64)` | Primary keys, foreign keys |
| `string` (short) | `VARCHAR(256)` | Names, titles |
| `string` (long) | `TEXT` | Descriptions, content |
| `int` | `INTEGER` | Counts, order |
| `int64` | `BIGINT` | Large numbers |
| `bool` | `BOOLEAN` | Flags |
| `time.Time` | `TIMESTAMPTZ` | Timestamps (always with timezone) |
| `map/struct` | `JSONB` | Flexible data |
| `[]string` | `TEXT[]` | String arrays |

## Constraint Naming

| Type | Pattern | Example |
|------|---------|---------|
| Foreign Key | `{table}_fkey_{column}` | `examples_fkey_tenant_id` |
| Unique | `{table}_uq_{columns}` | `examples_uq_tenant_id_name` |
| Check | `{table}_chk_{description}` | `examples_chk_order_positive` |

## Index Naming

| Type | Pattern | Example |
|------|---------|---------|
| Single column | `{table}_idx_{column}` | `examples_idx_tenant_id` |
| Composite | `{table}_idx_{col1}_{col2}` | `examples_idx_tenant_id_status` |
| Unique | `{table}_uidx_{column}` | `examples_uidx_email` |

## Constant Tables (Enum Equivalent)

### Table Definition

```sql
-- +goose Up
CREATE TABLE "example_statuses" (
    "id" VARCHAR(64) PRIMARY KEY
);

-- +goose Down
DROP TABLE IF EXISTS "example_statuses";
```

### YAML Constant Definition

Location: `db/{mysql,postgresql}/constants/{table_name}.yaml`

```yaml
# db/{mysql,postgresql}/constants/example_statuses.yaml
table: example_statuses
records:
  - id: draft
  - id: published
  - id: archived
```

## Adding Columns

```sql
-- +goose Up
ALTER TABLE "examples" ADD COLUMN "new_field" VARCHAR(256);

-- +goose Down
ALTER TABLE "examples" DROP COLUMN "new_field";
```

## Adding NOT NULL Column with Default

```sql
-- +goose Up
ALTER TABLE "examples" ADD COLUMN "priority" INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE "examples" DROP COLUMN "priority";
```

## Running Migrations

```bash
# Apply all pending migrations and regenerate SQLBoiler
make migrate.up

# Check migration status
make migrate.status
```

## Best Practices

1. **Always include Down migration** - Enable rollback capability
2. **One logical change per migration** - Don't mix unrelated changes
3. **Test Down migration** - Ensure it properly reverts changes
4. **Add indexes for foreign keys** - Improve JOIN performance
5. **Use TIMESTAMPTZ** - Always store timestamps with timezone
6. **Consider data migration** - Handle existing data when adding constraints

## Common Patterns

### Soft Delete

```sql
ALTER TABLE "examples" ADD COLUMN "deleted_at" TIMESTAMPTZ;
CREATE INDEX "examples_idx_deleted_at" ON "examples" ("deleted_at") WHERE "deleted_at" IS NULL;
```

### Unique Constraint with Soft Delete

```sql
CREATE UNIQUE INDEX "examples_uidx_tenant_id_name"
    ON "examples" ("tenant_id", "name")
    WHERE "deleted_at" IS NULL;
```

### Composite Primary Key (Junction Table)

```sql
CREATE TABLE "example_tags" (
    "example_id" VARCHAR(64) NOT NULL,
    "tag_id"     VARCHAR(64) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    PRIMARY KEY ("example_id", "tag_id"),
    CONSTRAINT "example_tags_fkey_example_id"
        FOREIGN KEY ("example_id") REFERENCES "examples" ("id") ON DELETE CASCADE,
    CONSTRAINT "example_tags_fkey_tag_id"
        FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON DELETE CASCADE
);
```
