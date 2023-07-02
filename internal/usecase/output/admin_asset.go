package output

type AdminCreateAssetPresignedURL struct {
	AssetKey     string
	PresignedURL string
}

func NewAdminCreateAssetPresignedURL(
	assetKey string,
	presignedURL string,
) *AdminCreateAssetPresignedURL {
	return &AdminCreateAssetPresignedURL{
		AssetKey:     assetKey,
		PresignedURL: presignedURL,
	}
}
