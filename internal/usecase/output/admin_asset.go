package output

type AdminCreateAssetPresignedURL struct {
	AssetID      string
	PresignedURL string
}

func NewAdminCreateAssetPresignedURL(
	assetID string,
	presignedURL string,
) *AdminCreateAssetPresignedURL {
	return &AdminCreateAssetPresignedURL{
		AssetID:      assetID,
		PresignedURL: presignedURL,
	}
}
