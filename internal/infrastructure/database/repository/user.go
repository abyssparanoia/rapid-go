package repository

import (
	"context"
	"database/sql"

	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/repository"
	"github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/database/internal/dbmodel"
	"github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/database/internal/marshaller"
	"github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/database/transactable"
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type user struct{}

func NewUser() repository.User {
	return &user{}
}

func (r *user) Get(
	ctx context.Context,
	query repository.GetUserQuery,
	orFail bool,
	preload bool,
) (*model.User, error) {
	mods := []qm.QueryMod{}
	if query.ID.Valid {
		mods = append(mods, dbmodel.UserWhere.ID.EQ(query.ID.String))
	}
	if query.AuthUID.Valid {
		mods = append(mods, dbmodel.UserWhere.AuthUID.EQ(query.AuthUID.String))
	}
	if preload {
		mods = append(mods,
			qm.Load(dbmodel.UserRels.Tenant),
		)
	}
	dbUser, err := dbmodel.Users(
		mods...,
	).One(ctx, transactable.GetContextExecutor(ctx))
	if err != nil {
		if err == sql.ErrNoRows && !orFail {
			return nil, nil
		} else if err == sql.ErrNoRows {
			return nil, errors.NotFoundErr.Errorf("user is not found")
		}
		return nil, errors.InternalErr.Wrap(err)
	}
	return marshaller.OutputUserToModel(dbUser), nil
}

func (r *user) Create(
	ctx context.Context,
	user *model.User,
) (*model.User, error) {
	dst := marshaller.NewUserFromModel(user)
	if err := dst.Insert(ctx, transactable.GetContextExecutor(ctx), boil.Infer()); err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	return marshaller.OutputUserToModel(dst), nil
}
