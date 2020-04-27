package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/abyssparanoia/rapid-go/internal/dbmodels/defaultdb"
	"github.com/abyssparanoia/rapid-go/internal/default/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/default/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/default/infrastructure/entity"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluesqlboiler"
	"github.com/abyssparanoia/rapid-go/internal/pkg/httperror"
	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type user struct {
}

func (r *user) Get(ctx context.Context, userID string) (*model.User, error) {

	dbUser, err := defaultdb.Users(
		defaultdb.UserWhere.ID.EQ(userID),
	).One(ctx, gluesqlboiler.GetContextExecutor(ctx))

	if err != nil {
		if err == sql.ErrNoRows {
			msg := fmt.Sprintf("user %s not found", userID)
			err = errors.Wrap(err, msg)
			log.Errorf(ctx, msg, zap.Error(err))
			return nil, httperror.NotFoundError(err)
		}
		return nil, err
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
		return nil, err
	}

	return dbUser.OutputModel(), nil
}

// NewUser ... get user repository
func NewUser() repository.User {
	return &user{}
}
