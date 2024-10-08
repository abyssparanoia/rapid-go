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

type staffStaffAuthentication struct {
	cli          *auth.Client
	clientAPIKey string
}

func NewStaffAuthentication(
	firebaseAuthCli *auth.Client,
	firebaseClientAPIKey string,
) repository.StaffAuthentication {
	return &staffStaffAuthentication{
		cli:          firebaseAuthCli,
		clientAPIKey: firebaseClientAPIKey,
	}
}

func (r *staffStaffAuthentication) VerifyIDToken(
	ctx context.Context,
	idToken string,
) (*model.StaffClaims, error) {
	t, err := r.cli.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, errors.InvalidIDTokenErr.Wrap(err)
	}
	return marshaller.StaffClaimsToModel(t.UID, t.Claims), nil
}

func (r *staffStaffAuthentication) GetUserByEmail(
	ctx context.Context,
	email string,
) (*repository.StaffAuthenticationGetUserByEmailResult, error) {
	user, err := r.cli.GetUserByEmail(ctx, email)
	if auth.IsUserNotFound(err) {
		return &repository.StaffAuthenticationGetUserByEmailResult{
			Exist: false,
		}, nil
	}
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	return &repository.StaffAuthenticationGetUserByEmailResult{
		AuthUID:     user.UID,
		StaffClaims: marshaller.StaffClaimsToModel(user.UID, user.CustomClaims),
		Exist:       true,
	}, nil
}

func (r *staffStaffAuthentication) CreateUser(
	ctx context.Context,
	param repository.StaffAuthenticationCreateUserParam,
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

func (r *staffStaffAuthentication) StoreClaims(
	ctx context.Context,
	authUID string,
	claims *model.StaffClaims,
) error {
	if err := r.cli.SetCustomUserClaims(ctx, authUID, marshaller.StaffClaimsToMap(claims)); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (r *staffStaffAuthentication) CreateCustomToken(
	ctx context.Context,
	authUID string,
) (string, error) {
	customToken, err := r.cli.CustomToken(ctx, authUID)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return customToken, nil
}

func (r *staffStaffAuthentication) CreateIDToken(
	ctx context.Context,
	authUID string,
	password string,
) (string, error) {
	customToken, err := r.cli.CustomToken(ctx, authUID)
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
