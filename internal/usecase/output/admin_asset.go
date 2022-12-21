package output

type AdminCreateAssetPresignedURL struct {
	Path         string
	PresignedURL string
}

func NewAdminCreateAssetPresignedURL(
	path string,
	presignedURL string,
) *AdminCreateAssetPresignedURL {
	return &AdminCreateAssetPresignedURL{
		Path:         path,
		PresignedURL: presignedURL,
	}
}
