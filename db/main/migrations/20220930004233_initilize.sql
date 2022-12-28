-- +goose Up
CREATE TABLE `tenants` (
  `id`                       VARCHAR(64)  NOT NULL COMMENT "id",
  `name`                     VARCHAR(256) NOT NULL COMMENT "name",
  `created_at`               TIMESTAMP    NOT NULL COMMENT "created date",
  `updated_at`               TIMESTAMP    NOT NULL COMMENT "update date",
  CONSTRAINT `tenants_pkey` PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "tenant";

CREATE TABLE `user_roles` (
  `id`                       VARCHAR(32)    NOT NULL COMMENT "id",
  CONSTRAINT `user_roles_pkey` PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "user_role";

INSERT INTO `user_roles`
    (`id`)
VALUES
    ('normal'),
    ('admin');


CREATE TABLE `users` (
  `id`                       VARCHAR(64)    NOT NULL COMMENT "id",
  `tenant_id`                VARCHAR(64)    NOT NULL COMMENT "tenant_id",
  `role`                     VARCHAR(32)    NOT NULL COMMENT "role",
  `auth_uid`                 VARCHAR(256)   NOT NULL COMMENT "auth_uid",
  `display_name`             VARCHAR(256)   NOT NULL COMMENT "display_name",
  `image_path`               VARCHAR(1024)  NOT NULL COMMENT "auth_uid",
  `email`                    VARCHAR(512)   NOT NULL COMMENT "email",
  `created_at`               TIMESTAMP      NOT NULL COMMENT "created date",
  `updated_at`               TIMESTAMP      NOT NULL COMMENT "update date",
  CONSTRAINT `users_pkey` PRIMARY KEY (`id`),
  UNIQUE `users_unique_email` (`email`),
  UNIQUE `users_unique_auth_uid` (`auth_uid`),
  CONSTRAINT `users_fkey_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`),
  CONSTRAINT `users_fkey_role` FOREIGN KEY (`role`) REFERENCES `user_roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT "user";

-- +goose Down
DROP TABLE users;
DROP TABLE user_roles;
DROP TABLE tenants;
