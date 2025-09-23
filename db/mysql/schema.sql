CREATE TABLE `asset_types` (
  `id` varchar(256) NOT NULL COMMENT 'id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='asset_type';

CREATE TABLE `assets` (
  `id` varchar(64) NOT NULL COMMENT 'id',
  `content_type` varchar(256) NOT NULL COMMENT 'content_type',
  `type` varchar(256) NOT NULL COMMENT 'type',
  `path` text NOT NULL COMMENT 'path',
  `expires_at` datetime NOT NULL COMMENT 'expires_at',
  `created_at` datetime NOT NULL COMMENT 'created date',
  `updated_at` datetime NOT NULL COMMENT 'update date',
  PRIMARY KEY (`id`),
  KEY `assets_fkey_type` (`type`),
  KEY `assets_fkey_content_type` (`content_type`),
  CONSTRAINT `assets_fkey_content_type` FOREIGN KEY (`content_type`) REFERENCES `content_types` (`id`),
  CONSTRAINT `assets_fkey_type` FOREIGN KEY (`type`) REFERENCES `asset_types` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='asset';

CREATE TABLE `content_types` (
  `id` varchar(256) NOT NULL COMMENT 'id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='content_type';

CREATE TABLE `goose_db_version` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `version_id` bigint NOT NULL,
  `is_applied` tinyint(1) NOT NULL,
  `tstamp` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb3;

CREATE TABLE `staff_roles` (
  `id` varchar(32) NOT NULL COMMENT 'id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='staff_role';

CREATE TABLE `staffs` (
  `id` varchar(64) NOT NULL COMMENT 'id',
  `tenant_id` varchar(64) NOT NULL COMMENT 'tenant_id',
  `role` varchar(32) NOT NULL COMMENT 'role',
  `auth_uid` varchar(256) NOT NULL COMMENT 'auth_uid',
  `display_name` varchar(256) NOT NULL COMMENT 'display_name',
  `image_path` varchar(1024) NOT NULL COMMENT 'image_path',
  `email` varchar(512) NOT NULL COMMENT 'email',
  `created_at` datetime NOT NULL COMMENT 'created date',
  `updated_at` datetime NOT NULL COMMENT 'update date',
  PRIMARY KEY (`id`),
  UNIQUE KEY `staffs_unique_email` (`email`),
  UNIQUE KEY `staffs_unique_auth_uid` (`auth_uid`),
  KEY `staffs_fkey_tenant_id` (`tenant_id`),
  KEY `staffs_fkey_role` (`role`),
  CONSTRAINT `staffs_fkey_role` FOREIGN KEY (`role`) REFERENCES `staff_roles` (`id`),
  CONSTRAINT `staffs_fkey_tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='staff';

CREATE TABLE `tenants` (
  `id` varchar(64) NOT NULL COMMENT 'id',
  `name` varchar(256) NOT NULL COMMENT 'name',
  `created_at` datetime NOT NULL COMMENT 'created date',
  `updated_at` datetime NOT NULL COMMENT 'update date',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='tenant';

