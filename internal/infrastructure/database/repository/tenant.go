package repository

import (
	"context"
	"database/sql"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/internal/dbmodel"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/internal/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/transactable"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type tenant struct{}

func NewTenant() repository.Tenant {
	return &tenant{}
}

func (r *tenant) Get(
	ctx context.Context,
	query repository.GetTenantQuery,
) (*model.Tenant, error) {
	mods := []qm.QueryMod{}
	if query.ID.Valid {
		mods = append(mods, dbmodel.TenantWhere.ID.EQ(query.ID.String))
	}
	if query.ForUpdate {
		mods = append(mods, qm.For("UPDATE"))
	}
	dbTenant, err := dbmodel.Tenants(
		mods...,
	).One(ctx, transactable.GetContextExecutor(ctx))
	if err != nil {
		if err == sql.ErrNoRows && !query.OrFail {
			return nil, nil
		} else if err == sql.ErrNoRows {
			return nil, errors.TenantNotFoundErr.Errorf("tenant is not found")
		}
		return nil, errors.InternalErr.Wrap(err)
	}
	return marshaller.TenantToModel(dbTenant), nil
}

func (r *tenant) List(
	ctx context.Context,
	query repository.ListTenantsQuery,
) (model.Tenants, error) {
	mods := []qm.QueryMod{}
	if query.Page.Valid && query.Limit.Valid {
		mods = append(mods,
			qm.Limit(int(query.Limit.Uint64)),
			qm.Offset(int(query.Limit.Uint64*(query.Page.Uint64-1))),
		)
	}
	dbTenants, err := dbmodel.Tenants(
		mods...,
	).All(ctx, transactable.GetContextExecutor(ctx))
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	return marshaller.TenantsToModel(dbTenants), nil
}

func (r *tenant) Count(
	ctx context.Context,
	query repository.ListTenantsQuery,
) (uint64, error) {
	mods := []qm.QueryMod{}
	ttl, err := dbmodel.Tenants(
		mods...,
	).Count(ctx, transactable.GetContextExecutor(ctx))
	if err != nil {
		return 0, errors.InternalErr.Wrap(err)
	}
	return uint64(ttl), nil
}

func (r *tenant) Create(
	ctx context.Context,
	tenant *model.Tenant,
) (*model.Tenant, error) {
	dst := marshaller.TenantsToDBModel(tenant)
	if err := dst.Insert(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	return marshaller.TenantToModel(dst), nil
}

func (r *tenant) Update(
	ctx context.Context,
	tenant *model.Tenant,
) (*model.Tenant, error) {
	dst := marshaller.TenantsToDBModel(tenant)
	if _, err := dst.Update(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	return marshaller.TenantToModel(dst), nil
}

func (r *tenant) Delete(
	ctx context.Context,
	id string,
) error {
	dst := &dbmodel.Tenant{
		ID: id,
	}
	if _, err := dst.Delete(ctx, transactable.GetContextExecutor(ctx)); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}
