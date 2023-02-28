package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/volatiletech/null/v8"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type StaffAuthentication interface {
	VerifyIDToken(
		ctx context.Context,
		idToken string,
	) (*model.StaffClaims, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (*StaffAuthenticationGetUserByEmailResult, error)
	CreateUser(
		ctx context.Context,
		param StaffAuthenticationCreateUserParam,
	) (string, error)
	StoreClaims(
		ctx context.Context,
		authUID string,
		staffClaims *model.StaffClaims,
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

type StaffAuthenticationCreateUserParam struct {
	Email    string
	Password null.String
}

type StaffAuthenticationGetUserByEmailResult struct {
	AuthUID     string
	StaffClaims *model.StaffClaims
	Exist       bool
}
