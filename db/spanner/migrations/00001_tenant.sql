CREATE TABLE `Tenants` (
  `TenantID`                 STRING(36)   NOT NULL, -- tenant id
  `Name`                     STRING(256)  NOT NULL, -- name
  `CreatedAt`                TIMESTAMP    NOT NULL, -- creation date
  `UpdatedAt`                TIMESTAMP    NOT NULL, -- updation date
) PRIMARY KEY(`TenantID`)