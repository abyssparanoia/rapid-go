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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type authentication struct {
	cli          *cognitoidentityprovider.CognitoIdentityProvider
	userPoolID   string
	clientID     string
	publicKeySet jwk.Set
}

func NewAuthentication(
	ctx context.Context,
	cognitoCli *cognitoidentityprovider.CognitoIdentityProvider,
	userPoolID string,
	clientID string,
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
		clientID:     clientID,
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

	jwtClaims := &dto.AWSCognitoClaims{}
	if err := json.Unmarshal(jsonString, jwtClaims); err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}

	return marshaller.UserAttributesToModel(dto.NewUserAttributesFromClaims(jwtClaims)), nil
}

func (r *authentication) GetUserByEmail(
	ctx context.Context,
	email string,
) (*repository.AuthenticationGetUserByEmailResult, error) {
	req := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: aws.String(r.userPoolID),
		Filter:     aws.String(fmt.Sprintf("email = \"%s\"", email)),
		Limit:      aws.Int64(1),
	}
	res, err := r.cli.ListUsers(req)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	var user *cognitoidentityprovider.UserType
	for _, cognitoUser := range res.Users {
		for _, attr := range cognitoUser.Attributes {
			if attr.Name == aws.String("email") && attr.Value == aws.String(email) {
				user = cognitoUser
			}
		}
	}
	if user == nil {
		return &repository.AuthenticationGetUserByEmailResult{
			Exist: false,
		}, nil
	}
	return &repository.AuthenticationGetUserByEmailResult{
		AuthUID: *user.Username,
		Claims:  marshaller.UserAttributesToModel(dto.NewUserAttributesFromCognitoUser(user)),
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
	req := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:             aws.String(r.userPoolID),
		Username:               aws.String(authUID),
		UserAttributes:         attrs,
		DesiredDeliveryMediums: aws.StringSlice([]string{deliveryMediumTypeEmail}),
	}
	_, err := r.cli.AdminCreateUser(req)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}

	if param.Password.Valid {
		req := &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(r.userPoolID),
			Username:   aws.String(authUID),
			Password:   aws.String(param.Password.String),
			Permanent:  aws.Bool(true),
		}
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
		SetUserAttributes(marshaller.ClaimsToCustomUserAttributes(claims).ToSlice())
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

func (r *authentication) CreateIDToken(
	ctx context.Context,
	authUID string,
) (string, error) {
	req := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeAdminUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(authUID),
			"PASSWORD": aws.String("kkkkk412"),
		},
		ClientId:   aws.String(r.clientID),
		UserPoolId: aws.String(r.userPoolID),
	}
	res, err := r.cli.AdminInitiateAuth(req)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	if res == nil || res.AuthenticationResult == nil || res.AuthenticationResult.IdToken == nil {
		return "", errors.InternalErr.Errorf("failed to auth")
	}

	return *res.AuthenticationResult.IdToken, nil
}
