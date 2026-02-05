package model

type StaffSortKey string

const (
	StaffSortKeyUnknown         StaffSortKey = "unknown"
	StaffSortKeyCreatedAtDesc   StaffSortKey = "created_at_desc"
	StaffSortKeyCreatedAtAsc    StaffSortKey = "created_at_asc"
	StaffSortKeyDisplayNameAsc  StaffSortKey = "display_name_asc"
	StaffSortKeyDisplayNameDesc StaffSortKey = "display_name_desc"
)

func NewStaffSortKey(s string) StaffSortKey {
	switch s {
	case StaffSortKeyCreatedAtDesc.String(),
		StaffSortKeyCreatedAtAsc.String(),
		StaffSortKeyDisplayNameAsc.String(),
		StaffSortKeyDisplayNameDesc.String():
		return StaffSortKey(s)
	default:
		return StaffSortKeyUnknown
	}
}

func (k StaffSortKey) String() string {
	return string(k)
}

func (k StaffSortKey) Valid() bool {
	return k != StaffSortKeyUnknown && k != ""
}
