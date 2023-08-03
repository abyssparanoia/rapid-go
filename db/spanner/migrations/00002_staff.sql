CREATE TABLE `StaffRoles` (
  `StaffRoleID`              STRING(32)     NOT NULL -- staff role id
) PRIMARY KEY(`StaffRoleID`);

CREATE TABLE `Staffs` (
  `StaffID`                  STRING(36)     NOT NULL, -- staff id
  `TenantID`                 STRING(36)     NOT NULL, -- tenant id
  `Role`                     STRING(32)     NOT NULL, -- role
  `AuthUID`                  STRING(256)    NOT NULL, -- auth uid
  `DisplayName`              STRING(256)    NOT NULL, -- role
  `ImagePath`                STRING(MAX)    NOT NULL, -- image path
  `Email`                    STRING(512)    NOT NULL, -- email
  `CreatedAt`                TIMESTAMP      NOT NULL, -- creation date
  `UpdatedAt`                TIMESTAMP      NOT NULL, -- updation date
  CONSTRAINT `Staffs_FK_TenantID` FOREIGN KEY (`TenantID`) REFERENCES `Tenants` (`TenantID`),
  CONSTRAINT `Staffs_FK_Role` FOREIGN KEY (`Role`) REFERENCES `StaffRoles` (`StaffRoleID`)
) PRIMARY KEY(`StaffID`);