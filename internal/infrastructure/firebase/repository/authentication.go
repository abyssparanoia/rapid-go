package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/firebase/internal/dto"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/firebase/internal/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/pkg/errors"
)

type authentication struct {
	cli          *auth.Client
	clientApiKey string
}

func NewAuthentication(
	firebaseAuthCli *auth.Client,
	firebaseClientApiKey string,
) repository.Authentication {
	return &authentication{
		cli:          firebaseAuthCli,
		clientApiKey: firebaseClientApiKey,
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

func (r *authentication) CreateIDToken(
	ctx context.Context,
	authUID string,
	password string,
) (string, error) {
	customToken, err := r.cli.CustomToken(ctx, authUID)
	if err != nil {
		return "", err
	}

	values := url.Values{}
	values.Add("token", customToken)
	values.Add("returnSecureToken", "true")
	values.Add("key", r.clientApiKey)

	resp, err := http.Post(
		"https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken",
		"application/x-www-form-urlencoded",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res dto.VerifyCustomTokenResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return "", err
	}

	return res.IDToken, nil
}
