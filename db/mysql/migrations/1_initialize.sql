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

CREATE TABLE `tenants` (
  `id`                       VARCHAR(64)  NOT NULL COMMENT "id",
  `name`                     VARCHAR(256) NOT NULL COMMENT "name",
  `created_at`               DATETIME     NOT NULL COMMENT "created date",
  `updated_at`               DATETIME     NOT NULL COMMENT "update date",
  CONSTRAINT `tenants_pkey` PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "tenant";

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
DROP TABLE tenants;
DROP TABLE assets;
DROP TABLE asset_types;
DROP TABLE content_types;