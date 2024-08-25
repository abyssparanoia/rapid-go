CREATE TABLE `goose_db_version` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `version_id` bigint NOT NULL,
  `is_applied` tinyint(1) NOT NULL,
  `tstamp` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
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
  `image_path` varchar(1024) NOT NULL COMMENT 'auth_uid',
  `email` varchar(512) NOT NULL COMMENT 'email',
  `created_at` timestamp NOT NULL COMMENT 'created date',
  `updated_at` timestamp NOT NULL COMMENT 'update date',
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
  `created_at` timestamp NOT NULL COMMENT 'created date',
  `updated_at` timestamp NOT NULL COMMENT 'update date',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='tenant';

