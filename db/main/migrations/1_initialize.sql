-- +goose Up
CREATE TABLE `tenants` (
  `id`                       VARCHAR(64)  NOT NULL COMMENT "id",
  `name`                     VARCHAR(256) NOT NULL COMMENT "name",
  `created_at`               TIMESTAMP    NOT NULL COMMENT "created date",
  `updated_at`               TIMESTAMP    NOT NULL COMMENT "update date",
  CONSTRAINT `tenants_pkey` PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "tenant";

CREATE TABLE `staff_roles` (
  `id`                       VARCHAR(32)    NOT NULL COMMENT "id",
  CONSTRAINT `staff_roles_pkey` PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "staff_role";

INSERT INTO `staff_roles`
    (`id`)
VALUES
    ('normal'),
    ('admin');


CREATE TABLE `staffs` (
  `id`                       VARCHAR(64)    NOT NULL COMMENT "id",
  `tenant_id`                VARCHAR(64)    NOT NULL COMMENT "tenant_id",
  `role`                     VARCHAR(32)    NOT NULL COMMENT "role",
  `auth_uid`                 VARCHAR(256)   NOT NULL COMMENT "auth_uid",
  `display_name`             VARCHAR(256)   NOT NULL COMMENT "display_name",
  `image_path`               VARCHAR(1024)  NOT NULL COMMENT "auth_uid",
  `email`                    VARCHAR(512)   NOT NULL COMMENT "email",
  `created_at`               TIMESTAMP      NOT NULL COMMENT "created date",
  `updated_at`               TIMESTAMP      NOT NULL COMMENT "update date",
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
