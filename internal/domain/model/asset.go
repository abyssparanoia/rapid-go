package model

import (
	"fmt"
	"mime"
	"strings"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/uuid"
)

type Asset struct {
	Key         string
	ContentType string
	AssetType   AssetType
	Path        string
}

func NewAsset(
	assetType AssetType,
	contentType string,
) (*Asset, error) {
	ext, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	return &Asset{
		Key:         uuid.UUIDBase64(),
		ContentType: contentType,
		AssetType:   assetType,
		Path:        fmt.Sprintf("%s/%s%s", assetType.String(), uuid.UUIDBase64(), ext[0]),
	}, nil
}

type AssetType string

const (
	AssetTypeUnknown   AssetType = "unknown"
	AssetTypeUserImage AssetType = "private/user_images"
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

func (m AssetType) IsPrivate() bool {
	return strings.HasPrefix(m.String(), "private")
}

func (m AssetType) IsPublic() bool {
	return strings.HasPrefix(m.String(), "public")
}

func ValidateAssetPath(
	assetType AssetType,
	path string,
) error {
	if !assetType.Valid() {
		return errors.AssetInvalidErr.Errorf("invalid asset type: %s", assetType.String())
	}
	if !strings.HasPrefix(path, assetType.String()) {
		return errors.AssetInvalidErr.Errorf("invalid asset path: %s", path)
	}
	return nil
}
