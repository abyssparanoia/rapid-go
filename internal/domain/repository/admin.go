package repository

import (
	"context"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
)

//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Admin interface {
	Get(
		ctx context.Context,
		query GetAdminQuery,
	) (*model.Admin, error)
	List(
		ctx context.Context,
		query ListAdminsQuery,
	) (model.Admins, error)
	Count(
		ctx context.Context,
		query ListAdminsQuery,
	) (uint64, error)
	Create(
		ctx context.Context,
		admin *model.Admin,
	) error
	Update(
		ctx context.Context,
		admin *model.Admin,
	) error
	Delete(
		ctx context.Context,
		id string,
	) error
}

type GetAdminQuery struct {
	BaseGetOptions
	ID      null.String
	AuthUID null.String
	Email   null.String
}

type ListAdminsQuery struct {
	BaseListOptions
	Role    nullable.Type[model.AdminRole]
	SortKey nullable.Type[model.AdminSortKey]
}
