package service

import (
	"context"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_service
type User interface {
	Create(
		ctx context.Context,
		param UserCreateParam,
	) (*model.User, error)
}

type UserCreateParam struct {
	TenantID    string
	Email       string
	Password    string
	UserRole    model.UserRole
	DisplayName string
	ImagePath   string
	RequestTime time.Time
}
