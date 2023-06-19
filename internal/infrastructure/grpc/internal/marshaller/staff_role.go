package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	modelv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/model/v1"
)

func StaffRoleToModel(staffRole modelv1.StaffRole) model.StaffRole {
	switch staffRole {
	case modelv1.StaffRole_STAFF_ROLE_NORMAL:
		return model.StaffRoleNormal
	case modelv1.StaffRole_STAFF_ROLE_ADMIN:
		return model.StaffRoleAdmin
	case modelv1.StaffRole_STAFF_ROLE_UNSPECIFIED:
		fallthrough
	default:
		return model.StaffRoleUnknown
	}
}

func StaffRoleToPB(staffRole model.StaffRole) modelv1.StaffRole {
	switch staffRole {
	case model.StaffRoleNormal:
		return modelv1.StaffRole_STAFF_ROLE_NORMAL
	case model.StaffRoleAdmin:
		return modelv1.StaffRole_STAFF_ROLE_ADMIN
	case model.StaffRoleUnknown:
		fallthrough
	default:
		return modelv1.StaffRole_STAFF_ROLE_UNSPECIFIED
	}
}
