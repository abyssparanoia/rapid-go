package output

type StaffCreateAssetPresignedURL struct {
	AssetID      string
	PresignedURL string
}

func NewStaffCreateAssetPresignedURL(
	assetID string,
	presignedURL string,
) *StaffCreateAssetPresignedURL {
	return &StaffCreateAssetPresignedURL{
		AssetID:      assetID,
		PresignedURL: presignedURL,
	}
}
