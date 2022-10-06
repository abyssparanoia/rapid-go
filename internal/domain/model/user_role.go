package model

type UserRole string

const (
	UserRoleUnknown UserRole = "unknown"
	UserRoleNormal  UserRole = "normal"
	UserRoleAdmin   UserRole = "admin"
)

func NewUserRole(role string) UserRole {
	switch role {
	case UserRoleNormal.String(),
		UserRoleAdmin.String():
		return UserRole(role)
	default:
		return UserRoleUnknown
	}
}

func (m UserRole) String() string {
	return string(m)
}

func (m UserRole) IsAdmin() bool {
	return m == UserRoleAdmin
}

func (m UserRole) Valid() bool {
	return m != UserRoleUnknown && m != ""
}
