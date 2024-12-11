package repository

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/abyssparanoia/memeduck"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/spanner/internal/dbmodel"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/spanner/internal/marshaller"
	"google.golang.org/grpc/codes"
)

type tenant struct{}

func NewTenant() repository.Tenant {
	return &tenant{}
}

func (r *tenant) Get(
	ctx context.Context,
	query repository.GetTenantQuery,
) (*model.Tenant, error) {
	conds := []memeduck.WhereCond{}
	params := map[string]interface{}{}
	if query.ID.Valid {
		conds = append(conds, memeduck.Eq(memeduck.Ident("TenantID"), memeduck.Param("TenantID")))
		params["TenantID"] = query.ID.String
	}
	sql, err := memeduck.Select(
		dbmodel.TenantTableName(),
		dbmodel.TenantColumns(),
	).
		Where(conds...).
		Limit(1).
		SQL()
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	rows, err := dbmodel.GetSpannerTransaction(ctx).QueryContext(ctx, sql, params)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	if ok, err := rows.Next(); err != nil && spanner.ErrCode(err) != codes.NotFound {
		return nil, errors.InternalErr.Wrap(err)
	} else if !ok {
		if !query.OrFail {
			return nil, nil
		} else {
			return nil, errors.TenantNotFoundErr.New().
				WithDetail("tenant is not found").
				WithValue("query", query)
		}
	}

	var dst dbmodel.Tenant
	if err := rows.ToStruct(&dst); err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}

	return marshaller.TenantToModel(&dst), nil
}

func (r *tenant) buildListQuery(_ repository.ListTenantsQuery) ([]memeduck.WhereCond, map[string]interface{}) {
	return []memeduck.WhereCond{}, map[string]interface{}{}
}

func (r *tenant) List(
	ctx context.Context,
	query repository.ListTenantsQuery,
) (model.Tenants, error) {
	conds, params := r.buildListQuery(query)
	stmt := memeduck.Select(
		dbmodel.TenantTableName(),
		dbmodel.TenantColumns(),
	).
		Where(conds...)
	if query.Page.Valid && query.Limit.Valid {
		stmt = stmt.LimitOffset(int(query.Limit.Uint64), int(query.Limit.Uint64*(query.Page.Uint64-1)))
	}
	sql, err := stmt.SQL()
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	rows, err := dbmodel.GetSpannerTransaction(ctx).QueryContext(ctx, sql, params)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	dsts := dbmodel.TenantSlice{}
	for {
		if ok, err := rows.Next(); err != nil && spanner.ErrCode(err) != codes.NotFound {
			return nil, errors.InternalErr.Wrap(err)
		} else if !ok {
			break
		}

		var dst dbmodel.Tenant
		if err := rows.ToStruct(&dst); err != nil {
			return nil, errors.InternalErr.Wrap(err)
		}
		dsts = append(dsts, &dst)
	}

	return marshaller.TenantsToModel(dsts), nil
}

func (r *tenant) Count(
	ctx context.Context,
	query repository.ListTenantsQuery,
) (uint64, error) {
	conds, params := r.buildListQuery(query)
	sql, err := memeduck.Select(
		dbmodel.TenantTableName(),
		[]string{"COUNT(*)"},
	).
		Where(conds...).
		SQL()
	if err != nil {
		return 0, errors.InternalErr.Wrap(err)
	}
	rows, err := dbmodel.GetSpannerTransaction(ctx).QueryContext(ctx, sql, params)
	if err != nil {
		return 0, errors.InternalErr.Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	if ok, err := rows.Next(); err != nil {
		return 0, errors.InternalErr.Wrap(err)
	} else if !ok {
		return 0, nil
	}

	var count uint64
	if err := rows.Scan(&count); err != nil {
		return 0, errors.InternalErr.Wrap(err)
	}
	return count, nil
}

func (r *tenant) Create(
	ctx context.Context,
	tenant *model.Tenant,
) error {
	if err := marshaller.TenantToDBModel(tenant).Insert(ctx); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (r *tenant) Update(
	ctx context.Context,
	tenant *model.Tenant,
) error {
	if err := marshaller.TenantToDBModel(tenant).Update(ctx); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (r *tenant) Delete(
	ctx context.Context,
	id string,
) error {
	dst := dbmodel.Tenant{TenantID: id} //nolint:exhaustruct
	if err := dst.Delete(ctx); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}
