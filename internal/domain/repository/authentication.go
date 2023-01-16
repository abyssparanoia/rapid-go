package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/volatiletech/null/v8"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Authentication interface {
	VerifyIDToken(
		ctx context.Context,
		idToken string,
	) (*model.Claims, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (*AuthenticationGetUserByEmailResult, error)
	CreateUser(
		ctx context.Context,
		param AuthenticationCreateUserParam,
	) (string, error)
	StoreClaims(
		ctx context.Context,
		authUID string,
		claims *model.Claims,
	) error
	CreateCustomToken(
		ctx context.Context,
		authUID string,
	) (string, error)
	CreateIDToken(
		ctx context.Context,
		authUID string,
		password string,
	) (string, error)
}

type AuthenticationCreateUserParam struct {
	Email    string
	Password null.String
}

type AuthenticationGetUserByEmailResult struct {
	AuthUID string
	Claims  *model.Claims
	Exist   bool
}
