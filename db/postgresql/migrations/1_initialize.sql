-- +goose Up
CREATE TABLE content_types (
    "id" VARCHAR(256) PRIMARY KEY
);

CREATE TABLE asset_types (
    "id" VARCHAR(256) PRIMARY KEY
);

CREATE TABLE assets (
    "id"              VARCHAR(64)   PRIMARY KEY,
    "content_type"    VARCHAR(256)  NOT NULL,
    "type"            VARCHAR(256)  NOT NULL,
    "path"            TEXT          NOT NULL,
    "expires_at"      TIMESTAMPTZ   NOT NULL,
    "created_at"      TIMESTAMPTZ   NOT NULL,
    "updated_at"      TIMESTAMPTZ   NOT NULL,
    CONSTRAINT "assets_fkey_type" FOREIGN KEY ("type") REFERENCES "asset_types" ("id"),
    CONSTRAINT "assets_fkey_content_type" FOREIGN KEY ("content_type") REFERENCES "content_types" ("id")
);

CREATE INDEX "assets_idx_type" ON "assets" ("type");
CREATE INDEX "assets_idx_content_type" ON "assets" ("content_type");

CREATE TABLE tenants (
    "id"          VARCHAR(64)   PRIMARY KEY,
    "name"        VARCHAR(256)  NOT NULL,
    "created_at"  TIMESTAMPTZ   NOT NULL,
    "updated_at"  TIMESTAMPTZ   NOT NULL
);

CREATE TABLE tenant_tag_types (
    "id" VARCHAR(256) PRIMARY KEY
);

CREATE TABLE tenant_tags (
    "id"          VARCHAR(64)   PRIMARY KEY,
    "tenant_id"   VARCHAR(64)   NOT NULL,
    "type"        VARCHAR(256)  NOT NULL,
    "created_at"  TIMESTAMPTZ   NOT NULL,
    "updated_at"  TIMESTAMPTZ   NOT NULL,
    CONSTRAINT "tenant_tags_unique_tenant_id_type" UNIQUE ("tenant_id", "type"),
    CONSTRAINT "tenant_tags_fkey_tenant_id" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id"),
    CONSTRAINT "tenant_tags_fkey_type" FOREIGN KEY ("type") REFERENCES "tenant_tag_types" ("id")
);

CREATE INDEX "tenant_tags_idx_tenant_id" ON "tenant_tags" ("tenant_id");
CREATE INDEX "tenant_tags_idx_type" ON "tenant_tags" ("type");

CREATE TABLE staff_roles (
    "id" VARCHAR(32) PRIMARY KEY
);

CREATE TABLE staffs (
    "id"              VARCHAR(64)   PRIMARY KEY,
    "tenant_id"       VARCHAR(64)   NOT NULL,
    "role"            VARCHAR(32)   NOT NULL,
    "auth_uid"        VARCHAR(256)  NOT NULL,
    "display_name"    VARCHAR(256)  NOT NULL,
    "image_path"      VARCHAR(1024) NOT NULL,
    "email"           VARCHAR(512)  NOT NULL,
    "created_at"      TIMESTAMPTZ   NOT NULL,
    "updated_at"      TIMESTAMPTZ   NOT NULL,
    CONSTRAINT "staffs_unique_email" UNIQUE ("email"),
    CONSTRAINT "staffs_unique_auth_uid" UNIQUE ("auth_uid"),
    CONSTRAINT "staffs_fkey_tenant_id" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id"),
    CONSTRAINT "staffs_fkey_role" FOREIGN KEY ("role") REFERENCES "staff_roles" ("id")
);

CREATE INDEX "staffs_idx_tenant_id" ON "staffs" ("tenant_id");
CREATE INDEX "staffs_idx_role" ON "staffs" ("role");

-- +goose Down
DROP TABLE IF EXISTS staffs;
DROP TABLE IF EXISTS staff_roles;
DROP TABLE IF EXISTS tenant_tags;
DROP TABLE IF EXISTS tenant_tag_types;
DROP TABLE IF EXISTS tenants;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS asset_types;
DROP TABLE IF EXISTS content_types;
