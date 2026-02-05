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

type admin struct{}

func NewAdmin() repository.Admin {
	return &admin{}
}

func (r *admin) Get(
	ctx context.Context,
	query repository.GetAdminQuery,
) (*model.Admin, error) {
	mods := []qm.QueryMod{}
	if query.ID.Valid {
		mods = append(mods, dbmodel.AdminWhere.ID.EQ(query.ID.String))
	}
	if query.AuthUID.Valid {
		mods = append(mods, dbmodel.AdminWhere.AuthUID.EQ(query.AuthUID.String))
	}
	if query.Email.Valid {
		mods = append(mods, dbmodel.AdminWhere.Email.EQ(query.Email.String))
	}
	mods = addForUpdateFromBaseGetOptions(mods, query.BaseGetOptions)
	dbAdmin, err := dbmodel.Admins(
		mods...,
	).One(ctx, transactable.GetContextExecutor(ctx))
	if err != nil {
		if err == sql.ErrNoRows && !query.OrFail {
			return nil, nil
		} else if err == sql.ErrNoRows {
			return nil, errors.AdminNotFoundErr.New().
				WithDetail("admin is not found").
				WithValue("query", query)
		}
		return nil, errors.InternalErr.Wrap(err)
	}
	return marshaller.AdminToModel(dbAdmin), nil
}

func (r *admin) List(
	ctx context.Context,
	query repository.ListAdminsQuery,
) (model.Admins, error) {
	mods := r.buildListQuery(query)

	// Sorting (BEFORE pagination)
	if query.SortKey.Valid && query.SortKey.Value().Valid() {
		switch query.SortKey.Value() {
		case model.AdminSortKeyCreatedAtDesc:
			mods = append(mods, qm.OrderBy("\"created_at\" DESC"))
		case model.AdminSortKeyCreatedAtAsc:
			mods = append(mods, qm.OrderBy("\"created_at\" ASC"))
		case model.AdminSortKeyDisplayNameAsc:
			mods = append(mods, qm.OrderBy("\"display_name\" ASC"))
		case model.AdminSortKeyDisplayNameDesc:
			mods = append(mods, qm.OrderBy("\"display_name\" DESC"))
		case model.AdminSortKeyUnknown:
			return nil, errors.InternalErr.Errorf("invalid sort key: %s", query.SortKey.Value())
		}
	}

	// Pagination (AFTER sorting)
	if query.Page.Valid && query.Limit.Valid {
		mods = append(mods,
			qm.Limit(int(query.Limit.Uint64)),
			qm.Offset(int(query.Limit.Uint64*(query.Page.Uint64-1))),
		)
	}
	mods = addForUpdateFromBaseListOptions(mods, query.BaseListOptions)
	dbAdmins, err := dbmodel.Admins(
		mods...,
	).All(ctx, transactable.GetContextExecutor(ctx))
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	return marshaller.AdminsToModel(dbAdmins), nil
}

func (r *admin) Count(
	ctx context.Context,
	query repository.ListAdminsQuery,
) (uint64, error) {
	mods := r.buildListQuery(query)
	ttl, err := dbmodel.Admins(
		mods...,
	).Count(ctx, transactable.GetContextExecutor(ctx))
	if err != nil {
		return 0, errors.InternalErr.Wrap(err)
	}
	return uint64(ttl), nil
}

func (r *admin) buildListQuery(query repository.ListAdminsQuery) []qm.QueryMod {
	mods := []qm.QueryMod{}
	if query.Role.Valid && query.Role.Value().Valid() {
		mods = append(mods, dbmodel.AdminWhere.Role.EQ(query.Role.Value().String()))
	}
	return mods
}

func (r *admin) Create(
	ctx context.Context,
	admin *model.Admin,
) error {
	dst := marshaller.AdminToDBModel(admin)
	if err := dst.Insert(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (r *admin) Update(
	ctx context.Context,
	admin *model.Admin,
) error {
	dst := marshaller.AdminToDBModel(admin)
	if _, err := dst.Update(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (r *admin) Delete(
	ctx context.Context,
	id string,
) error {
	dst := &dbmodel.Admin{ //nolint:exhaustruct
		ID: id,
	}
	if _, err := dst.Delete(ctx, transactable.GetContextExecutor(ctx)); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}
