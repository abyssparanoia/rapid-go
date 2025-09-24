package repository

import (
	"context"
	"database/sql"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/postgresql/internal/dbmodel"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/postgresql/internal/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/postgresql/transactable"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type staff struct{}

func NewStaff() repository.Staff {
	return &staff{}
}

func (r *staff) Get(
	ctx context.Context,
	query repository.GetStaffQuery,
) (*model.Staff, error) {
	mods := []qm.QueryMod{}
	if query.ID.Valid {
		mods = append(mods, dbmodel.StaffWhere.ID.EQ(query.ID.String))
	}
	if query.AuthUID.Valid {
		mods = append(mods, dbmodel.StaffWhere.AuthUID.EQ(query.AuthUID.String))
	}
	if query.Preload {
		mods = append(mods, r.buildPreload(mods)...)
	}
	mods = addForUpdateFromBaseGetOptions(mods, query.BaseGetOptions)
	dbStaff, err := dbmodel.Staffs(
		mods...,
	).One(ctx, transactable.GetContextExecutor(ctx))
	if err != nil {
		if err == sql.ErrNoRows && !query.OrFail {
			return nil, nil
		} else if err == sql.ErrNoRows {
			return nil, errors.StaffNotFoundErr.New().
				WithDetail("staff is not found").
				WithValue("query", query)
		}
		return nil, errors.InternalErr.Wrap(err)
	}
	return marshaller.StaffToModel(dbStaff), nil
}

func (r *staff) List(
	ctx context.Context,
	query repository.ListStaffQuery,
) (model.Staffs, error) {
	mods := []qm.QueryMod{}
	mods = append(mods, r.buildListQuery(query)...)
	if query.Page.Valid && query.Limit.Valid {
		mods = append(mods,
			qm.Limit(int(query.Limit.Uint64)),
			qm.Offset(int(query.Limit.Uint64*(query.Page.Uint64-1))),
		)
	}
	if query.Preload {
		mods = append(mods, r.buildPreload(mods)...)
	}
	mods = addForUpdateFromBaseListOptions(mods, query.BaseListOptions)
	dbStaffs, err := dbmodel.Staffs(
		mods...,
	).All(ctx, transactable.GetContextExecutor(ctx))
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	return marshaller.StaffsToModel(dbStaffs), nil
}

func (r *staff) Count(
	ctx context.Context,
	query repository.ListStaffQuery,
) (uint64, error) {
	mods := []qm.QueryMod{}
	mods = append(mods, r.buildListQuery(query)...)
	ttl, err := dbmodel.Staffs(
		mods...,
	).Count(ctx, transactable.GetContextExecutor(ctx))
	if err != nil {
		return 0, errors.InternalErr.Wrap(err)
	}
	return uint64(ttl), nil
}

func (r *staff) buildListQuery(query repository.ListStaffQuery) []qm.QueryMod {
	mods := []qm.QueryMod{}
	if query.TenantID.Valid {
		mods = append(mods, dbmodel.StaffWhere.TenantID.EQ(query.TenantID.String))
	}
	return mods
}

func (r *staff) buildPreload(mods []qm.QueryMod) []qm.QueryMod {
	return append(mods,
		qm.Load(dbmodel.StaffRels.Tenant),
	)
}

func (r *staff) Create(
	ctx context.Context,
	staff *model.Staff,
) error {
	dst := marshaller.StaffToDBModel(staff)
	if err := dst.Insert(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (r *staff) BatchCreate(
	ctx context.Context,
	staffs model.Staffs,
) error {
	dsts := marshaller.StaffsToDBModel(staffs)
	if _, err := dsts.InsertAll(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (r *staff) Update(
	ctx context.Context,
	staff *model.Staff,
) error {
	dst := marshaller.StaffToDBModel(staff)
	if _, err := dst.Update(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (r *staff) Delete(
	ctx context.Context,
	id string,
) error {
	dst := marshaller.StaffToDBModel(&model.Staff{ID: id}) //nolint:exhaustruct
	if _, err := dst.Delete(ctx, transactable.GetContextExecutor(ctx)); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}
