package model

import (
	"fmt"
	"mime"
	"strings"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/id"
	"github.com/abyssparanoia/rapid-go/internal/pkg/uuid"
)

type Asset struct {
	ID          string
	ContentType string
	Type        AssetType
	Path        string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Assets []*Asset

func NewAsset(
	assetType AssetType,
	contentType string,
	t time.Time,
) (*Asset, error) {
	ext, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	return &Asset{
		ID:          id.New(),
		ContentType: contentType,
		Type:        assetType,
		Path:        fmt.Sprintf("%s/%s%s", assetType.String(), uuid.UUIDBase64(), ext[0]),
		ExpiresAt:   t.Add(15 * time.Minute),
		CreatedAt:   t,
		UpdatedAt:   t,
	}, nil
}

func (m *Asset) Expiration() time.Duration {
	return m.ExpiresAt.Sub(m.CreatedAt)
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
