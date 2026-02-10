package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/firebase/internal/dto"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/firebase/internal/marshaller"
)

type adminAuthentication struct {
	cli          *auth.Client
	clientAPIKey string
}

func NewAdminAuthentication(
	firebaseAuthCli *auth.Client,
	firebaseClientAPIKey string,
) repository.AdminAuthentication {
	return &adminAuthentication{
		cli:          firebaseAuthCli,
		clientAPIKey: firebaseClientAPIKey,
	}
}

func (r *adminAuthentication) VerifyIDToken(
	ctx context.Context,
	idToken string,
) (*model.AdminClaims, error) {
	t, err := r.cli.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, errors.InvalidIDTokenErr.Wrap(err)
	}

	// Get user to retrieve email
	user, err := r.cli.GetUser(ctx, t.UID)
	if err != nil {
		return nil, errors.InvalidIDTokenErr.Wrap(err)
	}

	return marshaller.AdminClaimsToModel(t.UID, user.Email, t.Claims), nil
}

func (r *adminAuthentication) GetUserByEmail(
	ctx context.Context,
	email string,
) (*repository.AdminAuthenticationGetUserByEmailResult, error) {
	user, err := r.cli.GetUserByEmail(ctx, email)
	if auth.IsUserNotFound(err) {
		return &repository.AdminAuthenticationGetUserByEmailResult{
			Exist: false,
		}, nil
	}
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	return &repository.AdminAuthenticationGetUserByEmailResult{
		AuthUID:     user.UID,
		AdminClaims: marshaller.AdminClaimsToModel(user.UID, user.Email, user.CustomClaims),
		Exist:       true,
	}, nil
}

func (r *adminAuthentication) CreateUser(
	ctx context.Context,
	param repository.AdminAuthenticationCreateUserParam,
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

func (r *adminAuthentication) StoreClaims(
	ctx context.Context,
	authUID string,
	claims *model.AdminClaims,
) error {
	if err := r.cli.SetCustomUserClaims(ctx, authUID, marshaller.AdminClaimsToMap(claims)); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (r *adminAuthentication) CreateCustomToken(
	ctx context.Context,
	authUID string,
) (string, error) {
	customToken, err := r.cli.CustomToken(ctx, authUID)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return customToken, nil
}

func (r *adminAuthentication) CreateIDToken(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	// Get user by email to retrieve authUID
	result, err := r.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if !result.Exist {
		return "", errors.InvalidIDTokenErr.New().WithDetail("user not found")
	}

	customToken, err := r.cli.CustomToken(ctx, result.AuthUID)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}

	values := url.Values{}
	values.Add("token", customToken)
	values.Add("returnSecureToken", "true")
	values.Add("key", r.clientAPIKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken", strings.NewReader(values.Encode()))
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}

	var res dto.VerifyCustomTokenResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return "", errors.InternalErr.Wrap(err)
	}

	return res.IDToken, nil
}
