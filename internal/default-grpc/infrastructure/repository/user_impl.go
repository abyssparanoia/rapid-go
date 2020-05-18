package repository

import (
	"context"
	"database/sql"

	"github.com/abyssparanoia/rapid-go/internal/dbmodels/defaultdb"
	"github.com/abyssparanoia/rapid-go/internal/default-grpc/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/default-grpc/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/default-grpc/infrastructure/entity"
	"github.com/abyssparanoia/rapid-go/internal/pkg/error/grpcerror"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluesqlboiler"
	"github.com/volatiletech/sqlboiler/boil"
)

type user struct {
}

func (r *user) Get(
	ctx context.Context,
	userID string,
	orFail bool,
) (*model.User, error) {

	dbUser, err := defaultdb.Users(
		defaultdb.UserWhere.ID.EQ(userID),
	).One(ctx, gluesqlboiler.GetContextExecutor(ctx))

	if err != nil {
		if err == sql.ErrNoRows && orFail {
			return nil, grpcerror.UserNotFoundErr.New()
		}
		return nil, grpcerror.DBInternalErr.Wrap(err)
	}

	user := entity.User{User: *dbUser}
	return user.OutputModel(), nil
}

func (r *user) Create(
	ctx context.Context,
	user *model.User,
) (*model.User, error) {
	dbUser := entity.NewUserFromModel(user)

	err := dbUser.Insert(ctx, gluesqlboiler.GetContextExecutor(ctx), boil.Infer())
	if err != nil {
		return nil, grpcerror.DBInternalErr.Wrap(err)
	}

	return dbUser.OutputModel(), nil
}

// NewUser ... get user repository
func NewUser() repository.User {
	return &user{}
}
