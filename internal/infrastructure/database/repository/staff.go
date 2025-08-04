package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/internal/dbmodel"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/internal/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/transactable"
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
		mods = append(mods,
			qm.Load(dbmodel.StaffRels.Tenant),
		)
	}
	if query.ForUpdate {
		mods = append(mods, qm.For("UPDATE"))
	}
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
