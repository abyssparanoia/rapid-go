-- +goose Up
CREATE TABLE `content_types` (
  `id`                          VARCHAR(256)    NOT NULL COMMENT "id",
  CONSTRAINT `content_types_pkey` PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "content_type";

CREATE TABLE `asset_types` (
  `id`                          VARCHAR(256)    NOT NULL COMMENT "id",
  CONSTRAINT `asset_types_pkey` PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "asset_type";

CREATE TABLE `assets` (
  `id`                       VARCHAR(64)    NOT NULL COMMENT "id",
  `content_type`             VARCHAR(256)  NOT NULL COMMENT "content_type",
  `type`                     VARCHAR(256)   NOT NULL COMMENT "type",
  `path`                     TEXT           NOT NULL COMMENT "path",
  `expires_at`               DATETIME       NOT NULL COMMENT "expires_at",
  `created_at`               DATETIME       NOT NULL COMMENT "created date",
  `updated_at`               DATETIME       NOT NULL COMMENT "update date",
  CONSTRAINT `assets_pkey` PRIMARY KEY (`id`),
  CONSTRAINT `assets_fkey_type` FOREIGN KEY (`type`) REFERENCES `asset_types` (`id`),
  CONSTRAINT `assets_fkey_content_type` FOREIGN KEY (`content_type`) REFERENCES `content_types` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "asset";

CREATE TABLE `admin_roles` (
  `id` VARCHAR(32) NOT NULL COMMENT "id",
  CONSTRAINT `admin_roles_pkey` PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "admin role";

CREATE TABLE `admins` (
  `id`           VARCHAR(64)   NOT NULL COMMENT "id",
  `role`         VARCHAR(32)   NOT NULL COMMENT "role",
  `auth_uid`     VARCHAR(256)  NOT NULL COMMENT "auth uid",
  `email`        VARCHAR(512)  NOT NULL COMMENT "email",
  `display_name` VARCHAR(256)  NOT NULL COMMENT "display name",
  `created_at`   DATETIME      NOT NULL COMMENT "created date",
  `updated_at`   DATETIME      NOT NULL COMMENT "update date",
  CONSTRAINT `admins_pkey` PRIMARY KEY (`id`),
  UNIQUE `admins_unique_auth_uid` (`auth_uid`),
  UNIQUE `admins_unique_email` (`email`),
  CONSTRAINT `admins_fkey_role` FOREIGN KEY (`role`) REFERENCES `admin_roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "admin";

CREATE TABLE `tenants` (
  `id`                       VARCHAR(64)  NOT NULL COMMENT "id",
  `name`                     VARCHAR(256) NOT NULL COMMENT "name",
  `created_at`               DATETIME     NOT NULL COMMENT "created date",
  `updated_at`               DATETIME     NOT NULL COMMENT "update date",
  CONSTRAINT `tenants_pkey` PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "tenant";

CREATE TABLE `tenant_tag_types` (
  `id`                       VARCHAR(256) NOT NULL COMMENT "id",
  CONSTRAINT `tenant_tag_types_pkey` PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "tenant_tag_type";

CREATE TABLE `tenant_tags` (
  `id`                       VARCHAR(64)  NOT NULL COMMENT "id",
  `tenant_id`                VARCHAR(64)  NOT NULL COMMENT "tenant_id",
  `type`                     VARCHAR(256) NOT NULL COMMENT "type",
  `created_at`               DATETIME     NOT NULL COMMENT "created date",
  `updated_at`               DATETIME     NOT NULL COMMENT "update date",
  CONSTRAINT `tenant_tags_pkey` PRIMARY KEY (`id`),
  UNIQUE `tenant_tags_unique_tenant_id_type` (`tenant_id`, `type`),
  CONSTRAINT `tenant_tags_fkey_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`),
  CONSTRAINT `tenant_tags_fkey_type` FOREIGN KEY (`type`) REFERENCES `tenant_tag_types` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "tenant_tag";

CREATE TABLE `staff_roles` (
  `id`                       VARCHAR(32)    NOT NULL COMMENT "id",
  CONSTRAINT `staff_roles_pkey` PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "staff_role";

CREATE TABLE `staffs` (
  `id`                       VARCHAR(64)    NOT NULL COMMENT "id",
  `tenant_id`                VARCHAR(64)    NOT NULL COMMENT "tenant_id",
  `role`                     VARCHAR(32)    NOT NULL COMMENT "role",
  `auth_uid`                 VARCHAR(256)   NOT NULL COMMENT "auth_uid",
  `display_name`             VARCHAR(256)   NOT NULL COMMENT "display_name",
  `image_path`               VARCHAR(1024)  NOT NULL COMMENT "image_path",
  `email`                    VARCHAR(512)   NOT NULL COMMENT "email",
  `created_at`               DATETIME       NOT NULL COMMENT "created date",
  `updated_at`               DATETIME       NOT NULL COMMENT "update date",
  CONSTRAINT `staffs_pkey` PRIMARY KEY (`id`),
  UNIQUE `staffs_unique_email` (`email`),
  UNIQUE `staffs_unique_auth_uid` (`auth_uid`),
  CONSTRAINT `staffs_fkey_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`),
  CONSTRAINT `staffs_fkey_role` FOREIGN KEY (`role`) REFERENCES `staff_roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "staff";

-- +goose Down
DROP TABLE staffs;
DROP TABLE staff_roles;
DROP TABLE tenant_tags;
DROP TABLE tenant_tag_types;
DROP TABLE tenants;
DROP TABLE admins;
DROP TABLE admin_roles;
DROP TABLE assets;
DROP TABLE asset_types;
DROP TABLE content_types;