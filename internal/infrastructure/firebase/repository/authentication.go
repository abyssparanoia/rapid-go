package repository

import (
	"context"

	"firebase.google.com/go/auth"
	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/repository"
	"github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/firebase/internal/marshaller"
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/errors"
)

type authentication struct {
	cli *auth.Client
}

func NewAuthentication(
	firebaseAuthCli *auth.Client,
) repository.Authentication {
	return &authentication{
		cli: firebaseAuthCli,
	}
}

func (r *authentication) VerifyIDToken(
	ctx context.Context,
	idToken string,
) (*model.Claims, error) {
	t, err := r.cli.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, errors.UnauthorizedErr.Wrap(err)
	}
	return marshaller.ClaimsToModel(t.UID, t.Claims), nil
}

func (r *authentication) GetUserByEmail(
	ctx context.Context,
	email string,
) (*repository.AuthenticationGetUserByEmailResult, error) {
	user, err := r.cli.GetUserByEmail(ctx, email)
	if auth.IsUserNotFound(err) {
		return &repository.AuthenticationGetUserByEmailResult{
			Exist: false,
		}, nil
	}
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	return &repository.AuthenticationGetUserByEmailResult{
		AuthUID: user.UID,
		Claims:  marshaller.ClaimsToModel(user.UID, user.CustomClaims),
		Exist:   true,
	}, nil
}

func (r *authentication) CreateUser(
	ctx context.Context,
	param repository.AuthenticationCreateUserParam,
) (string, error) {
	dto := &auth.UserToCreate{}
	dto = dto.Email(param.Email)
	if param.Password.Valid {
		dto = dto.Password(param.Password.String)
	}
	res, err := r.cli.CreateUser(ctx, dto)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return res.UID, nil
}

func (r *authentication) StoreClaims(
	ctx context.Context,
	authUID string,
	claims *model.Claims,
) error {
	if err := r.cli.SetCustomUserClaims(ctx, authUID, marshaller.ClaimsToMap(claims)); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (r *authentication) CreateCustomToken(
	ctx context.Context,
	authUID string,
) (string, error) {
	customToken, err := r.cli.CustomToken(ctx, authUID)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return customToken, nil
}
