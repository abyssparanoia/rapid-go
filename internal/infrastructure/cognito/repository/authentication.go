package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito/internal/dto"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito/internal/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/pkg/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/uuid"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type authentication struct {
	cli          *cognitoidentityprovider.CognitoIdentityProvider
	userPoolID   string
	publicKeySet jwk.Set
}

func NewAuthentication(
	ctx context.Context,
	cognitoCli *cognitoidentityprovider.CognitoIdentityProvider,
	userPoolID string,
	emulatorHost string,
) repository.Authentication {
	endpoint := cognitoCli.Endpoint
	if emulatorHost != "" {
		endpoint = emulatorHost
	}
	publicKeysURL := fmt.Sprintf("%s/%s/.well-known/jwks.json", endpoint, userPoolID)
	publicKeySet, err := jwk.Fetch(ctx, publicKeysURL)
	if err != nil {
		panic(err)
	}
	return &authentication{
		cli:          cognitoCli,
		userPoolID:   userPoolID,
		publicKeySet: publicKeySet,
	}
}

func (r *authentication) VerifyIDToken(
	ctx context.Context,
	idToken string,
) (*model.Claims, error) {
	CustomClaims := jwt.MapClaims{}

	jwtToken, err := jwt.ParseWithClaims(idToken, CustomClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.InternalErr.Errorf("unexpected signing method")
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.InternalErr.Errorf("kid header not found")
		}
		if _, ok := token.Claims.(*dto.AWSCognitoClaims); !ok {
			return nil, errors.InternalErr.Errorf("there is problem to get claims")
		}
		key, ok := r.publicKeySet.LookupKeyID(kid)
		if !ok {
			return nil, errors.InternalErr.Errorf("key %v not found", kid)
		}
		var tokenKey interface{}
		if err := key.Raw(&tokenKey); err != nil {
			return nil, errors.InternalErr.Errorf("failed to create token key")
		}

		return tokenKey, nil
	})
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}

	jsonString, err := json.Marshal(jwtToken.Claims)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}

	var jwtClaims *dto.AWSCognitoClaims

	// 定義したStructへ変換
	if err := json.Unmarshal(jsonString, jwtClaims); err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}

	return marshaller.ClaimsToModel(jwtClaims.Username), nil
}

func (r *authentication) GetUserByEmail(
	ctx context.Context,
	email string,
) (*repository.AuthenticationGetUserByEmailResult, error) {
	req := &cognitoidentityprovider.ListUsersInput{}
	req.SetUserPoolId(r.userPoolID).
		SetFilter(fmt.Sprintf("email=%s", email)).
		SetLimit(1)
	res, err := r.cli.ListUsers(req)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	if len(res.Users) == 0 {
		return &repository.AuthenticationGetUserByEmailResult{
			Exist: false,
		}, nil
	}
	user := res.Users[0]
	return &repository.AuthenticationGetUserByEmailResult{
		AuthUID: *user.Username,
		Claims:  marshaller.ClaimsToModel(*user.Username),
		Exist:   true,
	}, nil
}

func (r *authentication) CreateUser(
	ctx context.Context,
	param repository.AuthenticationCreateUserParam,
) (string, error) {
	authUID := uuid.UUIDBase64()
	emailAttr := &cognitoidentityprovider.AttributeType{}
	emailAttr.SetName(cognitoidentityprovider.UsernameAttributeTypeEmail).
		SetValue(param.Email)
	attrs := []*cognitoidentityprovider.AttributeType{emailAttr}
	deliveryMediumTypeEmail := cognitoidentityprovider.DeliveryMediumTypeEmail
	req := &cognitoidentityprovider.AdminCreateUserInput{}
	req.SetUserPoolId(r.userPoolID).
		SetUsername(authUID).
		SetUserAttributes(attrs).
		SetDesiredDeliveryMediums([]*string{&deliveryMediumTypeEmail})
	_, err := r.cli.AdminCreateUser(req)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}

	if param.Password.Valid {
		req := &cognitoidentityprovider.AdminSetUserPasswordInput{}
		req.SetUserPoolId(r.userPoolID).
			SetUsername(authUID).
			SetPassword(param.Password.String)
		_, err := r.cli.AdminSetUserPassword(req)
		if err != nil {
			return "", errors.InternalErr.Wrap(err)
		}
	}
	return authUID, nil
}

func (r *authentication) StoreClaims(
	ctx context.Context,
	authUID string,
	claims *model.Claims,
) error {
	req := &cognitoidentityprovider.AdminUpdateUserAttributesInput{}
	req.SetUserPoolId(r.userPoolID).
		SetUsername(authUID).
		SetUserAttributes(marshaller.ClaimsToUserAttributes(claims).ToSlice())
	_, err := r.cli.AdminUpdateUserAttributes(req)
	if err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (r *authentication) CreateCustomToken(
	ctx context.Context,
	authUID string,
) (string, error) {
	return "", errors.InternalErr.Errorf("can not create custom token in cognito")
}
