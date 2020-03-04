package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/default/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/default/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/default/infrastructure/entity"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluemysql"
	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
)

type user struct {
	cli *gluemysql.Client
}

func (r *user) Get(ctx context.Context, userID string) (*model.User, error) {

	dsts := []*entity.User{}

	db := r.cli.GetDB(ctx).
		Where("id = ?", userID).
		Limit(1).
		Find(&dsts)

	if err := gluemysql.HandleErrors(db); err != nil {
		log.Errorm(ctx, "db.Find", err)
		return nil, err
	}

	if len(dsts) == 0 {
		return nil, nil
	}

	return entity.OutputUsers(dsts)[0], nil
}

// NewUser ... get user repository
func NewUser(cli *gluemysql.Client) repository.User {
	return &user{cli}
}
