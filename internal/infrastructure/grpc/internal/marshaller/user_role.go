package marshaller

import (
	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	modelv1 "github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/grpc/pb/mmg/model/v1"
)

func UserRoleToModel(userRole modelv1.UserRole) model.UserRole {
	switch userRole {
	case modelv1.UserRole_USER_ROLE_NORMAL:
		return model.UserRoleNormal
	case modelv1.UserRole_USER_ROLE_ADMIN:
		return model.UserRoleAdmin
	default:
		return model.UserRoleUnknown
	}
}

func UserRoleToPB(userRole model.UserRole) modelv1.UserRole {
	switch userRole {
	case model.UserRoleNormal:
		return modelv1.UserRole_USER_ROLE_NORMAL
	case model.UserRoleAdmin:
		return modelv1.UserRole_USER_ROLE_ADMIN
	default:
		return modelv1.UserRole_USER_ROLE_UNSPECIFIED
	}
}
