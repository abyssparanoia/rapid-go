# SQL Patterns Reference

Advanced SQL patterns for database migrations.

## Table Creation Patterns

### Basic Table with Relations

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

### Junction Table (Many-to-Many)

```sql
-- +goose Up
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

CREATE INDEX "example_tags_idx_tag_id" ON "example_tags" ("tag_id");

-- +goose Down
DROP TABLE IF EXISTS "example_tags";
```

### Soft Delete Pattern

```sql
-- +goose Up
ALTER TABLE "examples" ADD COLUMN "deleted_at" TIMESTAMPTZ;

-- Partial index for efficient queries on non-deleted records
CREATE INDEX "examples_idx_deleted_at" ON "examples" ("deleted_at")
    WHERE "deleted_at" IS NULL;

-- +goose Down
DROP INDEX IF EXISTS "examples_idx_deleted_at";
ALTER TABLE "examples" DROP COLUMN "deleted_at";
```

### Unique Constraint with Soft Delete

```sql
-- +goose Up
-- Unique constraint only applies to non-deleted records
CREATE UNIQUE INDEX "examples_uidx_tenant_id_name"
    ON "examples" ("tenant_id", "name")
    WHERE "deleted_at" IS NULL;

-- +goose Down
DROP INDEX IF EXISTS "examples_uidx_tenant_id_name";
```

## Column Modification Patterns

### Add Column

```sql
-- +goose Up
ALTER TABLE "examples" ADD COLUMN "new_field" VARCHAR(256);

-- +goose Down
ALTER TABLE "examples" DROP COLUMN "new_field";
```

### Add NOT NULL Column with Default

```sql
-- +goose Up
ALTER TABLE "examples" ADD COLUMN "priority" INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE "examples" DROP COLUMN "priority";
```

### Add NOT NULL Column (Existing Data)

For tables with existing data, use a multi-step approach:

```sql
-- +goose Up
-- Step 1: Add nullable column
ALTER TABLE "examples" ADD COLUMN "category" VARCHAR(64);

-- Step 2: Backfill existing rows
UPDATE "examples" SET "category" = 'default' WHERE "category" IS NULL;

-- Step 3: Add NOT NULL constraint
ALTER TABLE "examples" ALTER COLUMN "category" SET NOT NULL;

-- +goose Down
ALTER TABLE "examples" DROP COLUMN "category";
```

### Rename Column

```sql
-- +goose Up
ALTER TABLE "examples" RENAME COLUMN "old_name" TO "new_name";

-- +goose Down
ALTER TABLE "examples" RENAME COLUMN "new_name" TO "old_name";
```

## Index Patterns

### Single Column Index

```sql
CREATE INDEX "examples_idx_tenant_id" ON "examples" ("tenant_id");
```

### Composite Index

```sql
CREATE INDEX "examples_idx_tenant_id_status" ON "examples" ("tenant_id", "status");
```

### Descending Index (for ORDER BY DESC)

```sql
CREATE INDEX "examples_idx_created_at" ON "examples" ("created_at" DESC);
```

### Unique Index

```sql
CREATE UNIQUE INDEX "examples_uidx_email" ON "examples" ("email");
```

### Partial Index

```sql
-- Index only active records
CREATE INDEX "examples_idx_active" ON "examples" ("tenant_id")
    WHERE "is_active" = TRUE;
```

### GIN Index for JSONB

```sql
CREATE INDEX "examples_idx_metadata" ON "examples" USING GIN ("metadata");
```

## Constraint Patterns

### Foreign Key with CASCADE

```sql
CONSTRAINT "example_items_fkey_example_id"
    FOREIGN KEY ("example_id") REFERENCES "examples" ("id") ON DELETE CASCADE
```

### Foreign Key with SET NULL

```sql
CONSTRAINT "examples_fkey_parent_id"
    FOREIGN KEY ("parent_id") REFERENCES "examples" ("id") ON DELETE SET NULL
```

### Check Constraint

```sql
CONSTRAINT "examples_chk_order_positive"
    CHECK ("order" >= 0)
```

### Unique Constraint (Multiple Columns)

```sql
CONSTRAINT "examples_uq_tenant_id_name"
    UNIQUE ("tenant_id", "name")
```

## Constant Table Patterns

### Simple Constant Table

```sql
-- +goose Up
CREATE TABLE "example_statuses" (
    "id" VARCHAR(64) PRIMARY KEY
);

-- +goose Down
DROP TABLE IF EXISTS "example_statuses";
```

### Constant Table with Display Name

```sql
-- +goose Up
CREATE TABLE "example_types" (
    "id"           VARCHAR(64)  PRIMARY KEY,
    "display_name" VARCHAR(256) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS "example_types";
```

## Migration Order

When creating multiple related tables, follow this order:

1. Constant tables (statuses, types, roles)
2. Parent tables (no foreign keys to other new tables)
3. Child tables (with foreign keys)
4. Junction tables (many-to-many relationships)

**Down migration order is reversed.**

```sql
-- +goose Up
-- 1. Constant table
CREATE TABLE "example_statuses" (...);

-- 2. Main table
CREATE TABLE "examples" (...);

-- 3. Child table
CREATE TABLE "example_items" (...);

-- +goose Down
-- Reverse order
DROP TABLE IF EXISTS "example_items";
DROP TABLE IF EXISTS "examples";
DROP TABLE IF EXISTS "example_statuses";
```

## Best Practices

1. **Always include Down migration** - Enable rollback capability
2. **One logical change per migration** - Easier to review and rollback
3. **Test Down migration locally** - Ensure it properly reverts changes
4. **Add indexes for foreign keys** - Improve JOIN performance
5. **Use TIMESTAMPTZ** - Always store timestamps with timezone
6. **Quote identifiers** - Use double quotes for table/column names
7. **Consider existing data** - Plan data migration for schema changes
