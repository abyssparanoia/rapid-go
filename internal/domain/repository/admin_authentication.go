package repository

import (
	"context"

	"github.com/aarondl/null/v9"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type AdminAuthentication interface {
	VerifyIDToken(ctx context.Context, idToken string) (*model.AdminClaims, error)
	GetUserByEmail(ctx context.Context, email string) (*AdminAuthenticationGetUserByEmailResult, error)
	CreateUser(ctx context.Context, param AdminAuthenticationCreateUserParam) (string, error)
	StoreClaims(ctx context.Context, authUID string, adminClaims *model.AdminClaims) error
	CreateCustomToken(ctx context.Context, authUID string) (string, error)
	CreateIDToken(ctx context.Context, authUID string, password string) (string, error)
}

type AdminAuthenticationCreateUserParam struct {
	Email    string
	Password null.String
}

type AdminAuthenticationGetUserByEmailResult struct {
	AuthUID     string
	AdminClaims *model.AdminClaims
	Exist       bool
}
