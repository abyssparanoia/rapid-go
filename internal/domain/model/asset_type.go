package model

type AssetType string

const (
	AssetTypeUnknown   AssetType = "unknown"
	AssetTypeUserImage AssetType = "user_images"
)

func NewAssetType(str string) AssetType {
	switch str {
	case AssetTypeUserImage.String():
		return AssetType(str)
	default:
		return AssetTypeUnknown
	}
}

func (m AssetType) String() string {
	return string(m)
}

func (m AssetType) Valid() bool {
	return m != "" && m != AssetTypeUnknown
}
