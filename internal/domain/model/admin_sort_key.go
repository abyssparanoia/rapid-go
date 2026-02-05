package model

type AdminSortKey string

const (
	AdminSortKeyUnknown         AdminSortKey = "unknown"
	AdminSortKeyCreatedAtDesc   AdminSortKey = "created_at_desc"
	AdminSortKeyCreatedAtAsc    AdminSortKey = "created_at_asc"
	AdminSortKeyDisplayNameAsc  AdminSortKey = "display_name_asc"
	AdminSortKeyDisplayNameDesc AdminSortKey = "display_name_desc"
)

func NewAdminSortKey(s string) AdminSortKey {
	switch s {
	case AdminSortKeyCreatedAtDesc.String(),
		AdminSortKeyCreatedAtAsc.String(),
		AdminSortKeyDisplayNameAsc.String(),
		AdminSortKeyDisplayNameDesc.String():
		return AdminSortKey(s)
	default:
		return AdminSortKeyUnknown
	}
}

func (k AdminSortKey) String() string {
	return string(k)
}

func (k AdminSortKey) Valid() bool {
	return k != AdminSortKeyUnknown && k != ""
}
