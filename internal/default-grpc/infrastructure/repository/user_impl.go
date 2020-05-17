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
)

type user struct {
}

func (r *user) Get(ctx context.Context, userID string) (*model.User, error) {

	dbUser, err := defaultdb.Users(
		defaultdb.UserWhere.ID.EQ(userID),
	).One(ctx, gluesqlboiler.GetContextExecutor(ctx))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpcerror.UserNotFoundErr.New()
			// return nil, status.New(codes.NotFound, "user not found").Err()
		}
		return nil, err
	}

	user := entity.User{User: *dbUser}
	return user.OutputModel(), nil
}

// NewUser ... get user repository
func NewUser() repository.User {
	return &user{}
}
