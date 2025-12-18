package model

type TenantTagType string

const (
	TenantTagTypeUnknown       TenantTagType = "unknown"
	TenantTagTypeEntertainment TenantTagType = "entertainment"
	TenantTagTypeEducation     TenantTagType = "education"
	TenantTagTypeBusiness      TenantTagType = "business"
	TenantTagTypeOther         TenantTagType = "other"
)

func NewTenantTagType(tagType string) TenantTagType {
	switch tagType {
	case TenantTagTypeEntertainment.String(),
		TenantTagTypeEducation.String(),
		TenantTagTypeBusiness.String(),
		TenantTagTypeOther.String():
		return TenantTagType(tagType)
	default:
		return TenantTagTypeUnknown
	}
}

func (m TenantTagType) String() string {
	return string(m)
}

func (m TenantTagType) Valid() bool {
	return m != TenantTagTypeUnknown && m != ""
}
