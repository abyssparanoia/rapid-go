package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
)

func StaffRoleToModel(staffRole admin_apiv1.StaffRole) model.StaffRole {
	switch staffRole {
	case admin_apiv1.StaffRole_STAFF_ROLE_NORMAL:
		return model.StaffRoleNormal
	case admin_apiv1.StaffRole_STAFF_ROLE_ADMIN:
		return model.StaffRoleAdmin
	case admin_apiv1.StaffRole_STAFF_ROLE_UNSPECIFIED:
		fallthrough
	default:
		return model.StaffRoleUnknown
	}
}

func StaffRoleToPB(staffRole model.StaffRole) admin_apiv1.StaffRole {
	switch staffRole {
	case model.StaffRoleNormal:
		return admin_apiv1.StaffRole_STAFF_ROLE_NORMAL
	case model.StaffRoleAdmin:
		return admin_apiv1.StaffRole_STAFF_ROLE_ADMIN
	case model.StaffRoleUnknown:
		fallthrough
	default:
		return admin_apiv1.StaffRole_STAFF_ROLE_UNSPECIFIED
	}
}
