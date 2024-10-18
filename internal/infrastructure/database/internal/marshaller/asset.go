package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/internal/dbmodel"
)

func AssetToModel(e *dbmodel.Asset) *model.Asset {
	m := &model.Asset{
		ID:          e.ID,
		ContentType: e.ContentType,
		Type:        model.NewAssetType(e.Type),
		Path:        e.Path,
		ExpiresAt:   e.ExpiresAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}

	return m
}

func AssetsToModel(slice dbmodel.AssetSlice) model.Assets {
	dsts := make(model.Assets, len(slice))
	for idx, e := range slice {
		dsts[idx] = AssetToModel(e)
	}
	return dsts
}

func AssetToDBModel(m *model.Asset) *dbmodel.Asset {
	return &dbmodel.Asset{
		ID:          m.ID,
		ContentType: m.ContentType,
		Type:        m.Type.String(),
		Path:        m.Path,
		ExpiresAt:   m.ExpiresAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,

		R: nil,
		L: struct{}{},
	}
}
