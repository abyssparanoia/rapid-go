package cache

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/abyssparanoia/rapid-go/internal/domain/cache"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/internal/dbmodel"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/internal/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/transactable"
)

type assetPath struct{}

func NewAssetPath() cache.AssetPath {
	return &assetPath{}
}

func (c *assetPath) Get(
	ctx context.Context,
	id string,
) (string, error) {
	mods := []qm.QueryMod{}
	mods = append(mods, dbmodel.AssetWhere.ID.EQ(id))
	dbAsset, err := dbmodel.Assets(
		mods...,
	).One(ctx, transactable.GetContextExecutor(ctx))
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return dbAsset.Path, nil
}

func (c *assetPath) Set(
	ctx context.Context,
	asset *model.Asset,
) error {
	dst := marshaller.AssetToDBModel(asset)
	if err := dst.Insert(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (c *assetPath) Clear(
	ctx context.Context,
	id string,
) error {
	dst := &dbmodel.Asset{ //nolint:exhaustruct
		ID: id,
	}
	if _, err := dst.Delete(ctx, transactable.GetContextExecutor(ctx)); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}
