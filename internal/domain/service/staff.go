package service

import (
	"context"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_service
type Staff interface {
	Create(
		ctx context.Context,
		param StaffCreateParam,
	) (*model.Staff, error)
}

type StaffCreateParam struct {
	TenantID    string
	Email       string
	Password    string
	StaffRole   model.StaffRole
	DisplayName string
	ImagePath   string
	RequestTime time.Time
}
