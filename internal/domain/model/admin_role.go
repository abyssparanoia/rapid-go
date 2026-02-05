package model

type AdminRole string

const (
	AdminRoleUnknown AdminRole = "unknown"
	AdminRoleRoot    AdminRole = "root"
	AdminRoleNormal  AdminRole = "normal"
)

func NewAdminRole(role string) AdminRole {
	switch role {
	case AdminRoleRoot.String(),
		AdminRoleNormal.String():
		return AdminRole(role)
	default:
		return AdminRoleUnknown
	}
}

func (m AdminRole) String() string {
	return string(m)
}

func (m AdminRole) Valid() bool {
	return m == AdminRoleRoot || m == AdminRoleNormal
}

func (m AdminRole) IsRoot() bool {
	return m == AdminRoleRoot
}

func (m AdminRole) IsNormal() bool {
	return m == AdminRoleNormal
}
