# Type Mappings Reference

Complete column type mappings between Go and PostgreSQL.

## Standard Types

| Go Type | PostgreSQL Type | SQLBoiler Type | Notes |
|---------|-----------------|----------------|-------|
| `string` | `VARCHAR(n)` | `string` | Length varies by use case |
| `string` | `TEXT` | `string` | Unlimited length |
| `int` | `INTEGER` | `int` | 32-bit integer |
| `int64` | `BIGINT` | `int64` | 64-bit integer |
| `float64` | `DOUBLE PRECISION` | `float64` | 64-bit float |
| `bool` | `BOOLEAN` | `bool` | True/false |
| `time.Time` | `TIMESTAMPTZ` | `time.Time` | Timestamp with timezone |
| `[]byte` | `BYTEA` | `[]byte` | Binary data |

## Nullable Types

| Go Type | PostgreSQL Type | SQLBoiler Type | Notes |
|---------|-----------------|----------------|-------|
| `null.String` | `VARCHAR(n)` | `null.String` | Nullable string |
| `null.Int` | `INTEGER` | `null.Int` | Nullable int |
| `null.Int64` | `BIGINT` | `null.Int64` | Nullable int64 |
| `null.Float64` | `DOUBLE PRECISION` | `null.Float64` | Nullable float |
| `null.Bool` | `BOOLEAN` | `null.Bool` | Nullable boolean |
| `null.Time` | `TIMESTAMPTZ` | `null.Time` | Nullable timestamp |
| `null.Bytes` | `BYTEA` | `null.Bytes` | Nullable binary |

## String Length Guidelines

| Use Case | Length | PostgreSQL |
|----------|--------|------------|
| ID (primary key) | 64 | `VARCHAR(64)` |
| ID (foreign key) | 64 | `VARCHAR(64)` |
| Auth UID | 256 | `VARCHAR(256)` |
| Status/enum | 64 | `VARCHAR(64)` |
| Role | 32 | `VARCHAR(32)` |
| Name/title | 256 | `VARCHAR(256)` |
| Email | 512 | `VARCHAR(512)` |
| URL/path | 1024 | `VARCHAR(1024)` |
| Description | unlimited | `TEXT` |
| Content | unlimited | `TEXT` |

## Complex Types

### JSON/JSONB

```sql
-- PostgreSQL
"metadata" JSONB

-- Go domain model
Metadata map[string]interface{}

-- Or with typed struct
Metadata *ExampleMetadata
```

**JSONB vs JSON:**
- Use `JSONB` (default) - binary format, indexable, faster queries
- Use `JSON` only when preserving key order matters

### Arrays

```sql
-- PostgreSQL
"tags" TEXT[]

-- Go domain model
Tags []string
```

### Enum via Constant Table

Instead of PostgreSQL ENUM, use a constant table:

```sql
-- Constant table
CREATE TABLE "example_statuses" (
    "id" VARCHAR(64) PRIMARY KEY
);

-- Main table with foreign key
CREATE TABLE "examples" (
    "status" VARCHAR(64) NOT NULL,
    CONSTRAINT "examples_fkey_status"
        FOREIGN KEY ("status") REFERENCES "example_statuses" ("id")
);
```

**Go domain model:**

```go
type ExampleStatus string

const (
    ExampleStatusUnknown   ExampleStatus = "unknown"
    ExampleStatusDraft     ExampleStatus = "draft"
    ExampleStatusPublished ExampleStatus = "published"
)
```

## Timestamp Guidelines

**Always use TIMESTAMPTZ:**

```sql
"created_at" TIMESTAMPTZ NOT NULL
"updated_at" TIMESTAMPTZ NOT NULL
"deleted_at" TIMESTAMPTZ          -- nullable for soft delete
"expires_at" TIMESTAMPTZ NOT NULL -- for time-limited records
```

**Go domain model:**

```go
CreatedAt time.Time      // Required
UpdatedAt time.Time      // Required
DeletedAt null.Time      // Optional (soft delete)
ExpiresAt time.Time      // Required
```

## Default Values

| Type | PostgreSQL Default | Example |
|------|-------------------|---------|
| Integer | `DEFAULT n` | `"order" INTEGER NOT NULL DEFAULT 0` |
| Boolean | `DEFAULT TRUE/FALSE` | `"is_active" BOOLEAN NOT NULL DEFAULT TRUE` |
| String | `DEFAULT 'value'` | `"status" VARCHAR(64) NOT NULL DEFAULT 'draft'` |
| Timestamp | `DEFAULT NOW()` | `"created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()` |

**Note:** This project sets timestamps in Go code, not via database defaults.

## MySQL Differences

If targeting MySQL instead of PostgreSQL:

| PostgreSQL | MySQL |
|------------|-------|
| `TIMESTAMPTZ` | `DATETIME(6)` |
| `TEXT` | `LONGTEXT` |
| `BOOLEAN` | `TINYINT(1)` |
| `JSONB` | `JSON` |
| `TEXT[]` | Not supported (use JSON) |
| `SERIAL` | `AUTO_INCREMENT` |

## Spanner Differences

If targeting Cloud Spanner:

| PostgreSQL | Spanner |
|------------|---------|
| `VARCHAR(n)` | `STRING(n)` |
| `TEXT` | `STRING(MAX)` |
| `INTEGER` | `INT64` |
| `BOOLEAN` | `BOOL` |
| `TIMESTAMPTZ` | `TIMESTAMP` |
| `JSONB` | `JSON` |
| Foreign keys | Not enforced (use application logic) |
