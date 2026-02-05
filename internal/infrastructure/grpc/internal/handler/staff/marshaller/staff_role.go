package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	staff_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/staff_api/v1"
)

func StaffRoleToModel(staffRole staff_apiv1.StaffRole) model.StaffRole {
	switch staffRole {
	case staff_apiv1.StaffRole_STAFF_ROLE_NORMAL:
		return model.StaffRoleNormal
	case staff_apiv1.StaffRole_STAFF_ROLE_ADMIN:
		return model.StaffRoleAdmin
	case staff_apiv1.StaffRole_STAFF_ROLE_UNSPECIFIED:
		fallthrough
	default:
		return model.StaffRoleUnknown
	}
}

func StaffRoleToPB(staffRole model.StaffRole) staff_apiv1.StaffRole {
	switch staffRole {
	case model.StaffRoleNormal:
		return staff_apiv1.StaffRole_STAFF_ROLE_NORMAL
	case model.StaffRoleAdmin:
		return staff_apiv1.StaffRole_STAFF_ROLE_ADMIN
	case model.StaffRoleUnknown:
		fallthrough
	default:
		return staff_apiv1.StaffRole_STAFF_ROLE_UNSPECIFIED
	}
}
