package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/postgresql/internal/dbmodel"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/postgresql/internal/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/postgresql/transactable"
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
	mods = append(mods, r.buildPreload(query.Preload)...)
	mods = addForUpdateFromBaseGetOptions(mods, query.BaseGetOptions)
	dbTenant, err := dbmodel.Tenants(
		mods...,
	).One(ctx, transactable.GetContextExecutor(ctx))
	if err != nil {
		if err == sql.ErrNoRows && !query.OrFail {
			return nil, nil
		} else if err == sql.ErrNoRows {
			return nil, errors.TenantNotFoundErr.New().
				WithDetail("tenant is not found").
				WithValue("query", query)
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

	// Sorting (BEFORE pagination)
	if query.SortKey.Valid && query.SortKey.Value().Valid() {
		switch query.SortKey.Value() {
		case model.TenantSortKeyCreatedAtDesc:
			mods = append(mods, qm.OrderBy("\"created_at\" DESC"))
		case model.TenantSortKeyCreatedAtAsc:
			mods = append(mods, qm.OrderBy("\"created_at\" ASC"))
		case model.TenantSortKeyNameAsc:
			mods = append(mods, qm.OrderBy("\"name\" ASC"))
		case model.TenantSortKeyNameDesc:
			mods = append(mods, qm.OrderBy("\"name\" DESC"))
		case model.TenantSortKeyUnknown:
			return nil, errors.InternalErr.Errorf("invalid sort key: %s", query.SortKey.Value())
		}
	}

	// Pagination
	if query.Page.Valid && query.Limit.Valid {
		mods = append(mods,
			qm.Limit(int(query.Limit.Uint64)),
			qm.Offset(int(query.Limit.Uint64*(query.Page.Uint64-1))),
		)
	}
	mods = append(mods, r.buildPreload(query.Preload)...)
	mods = addForUpdateFromBaseListOptions(mods, query.BaseListOptions)
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

func (r *tenant) buildPreload(_ bool) []qm.QueryMod {
	return []qm.QueryMod{
		qm.Load(dbmodel.TenantRels.TenantTags),
	}
}

func (r *tenant) Create(
	ctx context.Context,
	tenant *model.Tenant,
) error {
	dst := marshaller.TenantToDBModel(tenant)
	if err := dst.Insert(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
		return errors.InternalErr.Wrap(err)
	}

	if len(tenant.Tags) > 0 {
		tags := marshaller.TenantTagsToDBModel(tenant.Tags, tenant.ID)
		if _, err := tags.InsertAll(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
			return errors.InternalErr.Wrap(err)
		}
	}

	return nil
}

func (r *tenant) BatchCreate(
	ctx context.Context,
	tenants model.Tenants,
) error {
	dstTenants := marshaller.TenantsToDBModel(tenants)
	exec := transactable.GetContextExecutor(ctx)

	if _, err := dstTenants.InsertAll(ctx, exec, boil.Infer()); err != nil {
		return errors.InternalErr.Wrap(err)
	}

	tagSlice := dbmodel.TenantTagSlice{}
	for _, t := range tenants {
		if len(t.Tags) == 0 {
			continue
		}
		tagSlice = append(tagSlice, marshaller.TenantTagsToDBModel(t.Tags, t.ID)...)
	}

	if len(tagSlice) > 0 {
		if _, err := tagSlice.InsertAll(ctx, exec, boil.Infer()); err != nil {
			return errors.InternalErr.Wrap(err)
		}
	}

	return nil
}

func (r *tenant) Update(
	ctx context.Context,
	tenant *model.Tenant,
) error {
	dst := marshaller.TenantToDBModel(tenant)
	if _, err := dst.Update(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
		return errors.InternalErr.Wrap(err)
	}

	if _, err := dbmodel.TenantTags(
		dbmodel.TenantTagWhere.TenantID.EQ(tenant.ID),
	).DeleteAll(ctx, transactable.GetContextExecutor(ctx)); err != nil {
		return errors.InternalErr.Wrap(err)
	}

	if len(tenant.Tags) > 0 {
		tags := marshaller.TenantTagsToDBModel(tenant.Tags, tenant.ID)
		if _, err := tags.InsertAll(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
			return errors.InternalErr.Wrap(err)
		}
	}

	return nil
}

func (r *tenant) BatchUpdate(
	ctx context.Context,
	tenants model.Tenants,
) error {
	exec := transactable.GetContextExecutor(ctx)

	for _, tenant := range tenants {
		if _, err := marshaller.TenantToDBModel(tenant).Update(ctx, exec, boil.Infer()); err != nil {
			return errors.InternalErr.Wrap(err)
		}

		if _, err := dbmodel.TenantTags(
			dbmodel.TenantTagWhere.TenantID.EQ(tenant.ID),
		).DeleteAll(ctx, exec); err != nil {
			return errors.InternalErr.Wrap(err)
		}

		if len(tenant.Tags) > 0 {
			tags := marshaller.TenantTagsToDBModel(tenant.Tags, tenant.ID)
			if _, err := tags.InsertAll(ctx, exec, boil.Infer()); err != nil {
				return errors.InternalErr.Wrap(err)
			}
		}
	}

	return nil
}

func (r *tenant) Delete(
	ctx context.Context,
	id string,
) error {
	dst := &dbmodel.Tenant{ //nolint:exhaustruct
		ID: id,
	}
	if _, err := dst.Delete(ctx, transactable.GetContextExecutor(ctx)); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}
