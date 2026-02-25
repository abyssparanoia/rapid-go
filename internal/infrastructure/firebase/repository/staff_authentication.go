package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/firebase/internal/dto"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/firebase/internal/marshaller"
)

type staffStaffAuthentication struct {
	cli          *auth.Client
	clientAPIKey string
	emulatorHost string
}

func NewStaffAuthentication(
	firebaseAuthCli *auth.Client,
	firebaseClientAPIKey string,
	emulatorHost string,
) repository.StaffAuthentication {
	return &staffStaffAuthentication{
		cli:          firebaseAuthCli,
		clientAPIKey: firebaseClientAPIKey,
		emulatorHost: emulatorHost,
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
	// Extract email from standard claims
	email := ""
	if emailVal, ok := t.Claims["email"]; ok {
		email = emailVal.(string) //nolint:errcheck
	}
	return marshaller.StaffClaimsToModel(t.UID, email, t.Claims), nil
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
		StaffClaims: marshaller.StaffClaimsToModel(user.UID, user.Email, user.CustomClaims),
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

	var reqBody io.Reader
	var apiURL string
	var contentType string
	var jsonBody []byte
	var customToken string
	var req *http.Request

	if r.emulatorHost != "" {
		// Emulator: Use signInWithPassword (avoids CustomToken emulatedSigner issue)
		apiURL = "http://" + r.emulatorHost + "/identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + r.clientAPIKey
		jsonBody, err = json.Marshal(map[string]interface{}{
			"email":             email,
			"password":          password,
			"returnSecureToken": true,
		})
		if err != nil {
			return "", errors.InternalErr.Wrap(err)
		}
		reqBody = bytes.NewReader(jsonBody)
		contentType = "application/json"
	} else {
		// Production: Use CustomToken flow with real service account signing
		customToken, err = r.cli.CustomToken(ctx, result.AuthUID)
		if err != nil {
			return "", errors.InternalErr.Wrap(err)
		}

		apiURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken"
		values := url.Values{}
		values.Add("token", customToken)
		values.Add("returnSecureToken", "true")
		values.Add("key", r.clientAPIKey)
		reqBody = strings.NewReader(values.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, apiURL, reqBody)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req) //nolint:gosec // URL is constructed from trusted environment variables
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", errors.InternalErr.Errorf("firebase auth request failed: status=%d body=%s", resp.StatusCode, string(b))
	}

	var res dto.VerifyCustomTokenResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return "", errors.InternalErr.Wrap(err)
	}

	// Check for empty idToken
	if res.IDToken == "" {
		return "", errors.InternalErr.Errorf("empty id_token in response: body=%s", string(b))
	}

	return res.IDToken, nil
}
