package model

type TenantSortKey string

const (
	TenantSortKeyUnknown       TenantSortKey = "unknown"
	TenantSortKeyCreatedAtDesc TenantSortKey = "created_at_desc"
	TenantSortKeyCreatedAtAsc  TenantSortKey = "created_at_asc"
	TenantSortKeyNameAsc       TenantSortKey = "name_asc"
	TenantSortKeyNameDesc      TenantSortKey = "name_desc"
)

func NewTenantSortKey(s string) TenantSortKey {
	switch s {
	case TenantSortKeyCreatedAtDesc.String(),
		TenantSortKeyCreatedAtAsc.String(),
		TenantSortKeyNameAsc.String(),
		TenantSortKeyNameDesc.String():
		return TenantSortKey(s)
	default:
		return TenantSortKeyUnknown
	}
}

func (k TenantSortKey) String() string {
	return string(k)
}

func (k TenantSortKey) Valid() bool {
	return k != TenantSortKeyUnknown && k != ""
}
