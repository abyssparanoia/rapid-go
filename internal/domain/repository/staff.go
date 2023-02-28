package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/volatiletech/null/v8"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Staff interface {
	Get(
		ctx context.Context,
		query GetStaffQuery,
	) (*model.Staff, error)
	Create(
		ctx context.Context,
		staff *model.Staff,
	) (*model.Staff, error)
}

type GetStaffQuery struct {
	BaseGetOptions
	ID      null.String
	AuthUID null.String
}
