package model

type StaffRole string

const (
	StaffRoleUnknown StaffRole = "unknown"
	StaffRoleNormal  StaffRole = "normal"
	StaffRoleAdmin   StaffRole = "admin"
)

func NewStaffRole(role string) StaffRole {
	switch role {
	case StaffRoleNormal.String(),
		StaffRoleAdmin.String():
		return StaffRole(role)
	default:
		return StaffRoleUnknown
	}
}

func (m StaffRole) String() string {
	return string(m)
}

func (m StaffRole) IsAdmin() bool {
	return m == StaffRoleAdmin
}

func (m StaffRole) Valid() bool {
	return m != StaffRoleUnknown && m != ""
}
