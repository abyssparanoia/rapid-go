package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito/internal/dto"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito/internal/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/pkg/uuid"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type staffAuthentication struct {
	cli          *cognitoidentityprovider.Client
	userPoolID   string
	clientID     string
	publicKeySet jwk.Set
}

func NewStaffAuthentication(
	ctx context.Context,
	cognitoCli *cognitoidentityprovider.Client,
	userPoolID string,
	clientID string,
	emulatorHost string,
	region string,
) repository.StaffAuthentication {
	publicKeySet, err := cognito.NewPublicKeySet(ctx, cognitoCli, userPoolID, emulatorHost, region)
	if err != nil {
		panic(err)
	}
	return &staffAuthentication{
		cli:          cognitoCli,
		userPoolID:   userPoolID,
		clientID:     clientID,
		publicKeySet: publicKeySet,
	}
}

func (r *staffAuthentication) VerifyIDToken(
	ctx context.Context,
	idToken string,
) (*model.StaffClaims, error) {
	customClaims := jwt.MapClaims{}

	jwtToken, err := jwt.ParseWithClaims(idToken, customClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.InvalidIDTokenErr.New().WithDetail("unexpected signing method")
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.InvalidIDTokenErr.New().WithDetail("kid header not found")
		}
		key, ok := r.publicKeySet.LookupKeyID(kid)
		if !ok {
			return nil, errors.InvalidIDTokenErr.New().WithDetail(fmt.Sprintf("key %v not found", kid))
		}
		var tokenKey interface{}
		if err := key.Raw(&tokenKey); err != nil {
			return nil, errors.InvalidIDTokenErr.New().WithDetail("failed to create token key")
		}

		return tokenKey, nil
	})
	if err != nil {
		return nil, errors.InvalidIDTokenErr.Wrap(err)
	}

	jsonString, err := json.Marshal(jwtToken.Claims)
	if err != nil {
		return nil, errors.InvalidIDTokenErr.Wrap(err)
	}

	jwtClaims := &dto.AWSCognitoStaffClaims{}
	if err := json.Unmarshal(jsonString, jwtClaims); err != nil {
		return nil, errors.InvalidIDTokenErr.Wrap(err)
	}

	return marshaller.StaffUserAttributesToModel(dto.NewUserAttributesFromClaims(jwtClaims)), nil
}

func (r *staffAuthentication) GetUserByEmail(
	ctx context.Context,
	email string,
) (*repository.StaffAuthenticationGetUserByEmailResult, error) {
	req := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: aws.String(r.userPoolID),
		Filter:     aws.String(fmt.Sprintf("email = \"%s\"", email)),
		Limit:      aws.Int32(1),
	}
	res, err := r.cli.ListUsers(ctx, req)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	var user *types.UserType
	for _, cognitoUser := range res.Users {
		for _, attr := range cognitoUser.Attributes {
			if attr.Name == aws.String("email") && attr.Value == aws.String(email) {
				user = &cognitoUser
			}
		}
	}
	if user == nil {
		return &repository.StaffAuthenticationGetUserByEmailResult{
			Exist: false,
		}, nil
	}
	return &repository.StaffAuthenticationGetUserByEmailResult{
		AuthUID:     *user.Username,
		StaffClaims: marshaller.StaffUserAttributesToModel(dto.NewUserAttributesFromCognitoUser(user)),
		Exist:       true,
	}, nil
}

func (r *staffAuthentication) CreateUser(
	ctx context.Context,
	param repository.StaffAuthenticationCreateUserParam,
) (string, error) {
	authUID := uuid.UUIDBase64()
	emailAttr := &types.AttributeType{
		Name:  aws.String(string(types.UsernameAttributeTypeEmail)),
		Value: aws.String(param.Email),
	}
	attrs := []types.AttributeType{*emailAttr}
	req := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:             aws.String(r.userPoolID),
		Username:               aws.String(authUID),
		UserAttributes:         attrs,
		DesiredDeliveryMediums: []types.DeliveryMediumType{types.DeliveryMediumTypeEmail},
	}
	_, err := r.cli.AdminCreateUser(ctx, req)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}

	if param.Password.Valid {
		req := &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(r.userPoolID),
			Username:   aws.String(authUID),
			Password:   aws.String(param.Password.String),
			Permanent:  true,
		}
		_, err := r.cli.AdminSetUserPassword(ctx, req)
		if err != nil {
			return "", errors.InternalErr.Wrap(err)
		}
	}
	return authUID, nil
}

func (r *staffAuthentication) StoreClaims(
	ctx context.Context,
	authUID string,
	claims *model.StaffClaims,
) error {
	req := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId:     aws.String(r.userPoolID),
		Username:       aws.String(authUID),
		UserAttributes: marshaller.StaffClaimsToStaffCustomUserAttributes(claims).ToSlice(),
	}
	_, err := r.cli.AdminUpdateUserAttributes(ctx, req)
	if err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (r *staffAuthentication) CreateCustomToken(
	ctx context.Context,
	authUID string,
) (string, error) {
	return "", errors.InternalErr.New().WithDetail("can not create custom token in cognito")
}

func (r *staffAuthentication) CreateIDToken(
	ctx context.Context,
	authUID string,
	password string,
) (string, error) {
	req := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: types.AuthFlowTypeAdminUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": authUID,
			"PASSWORD": password,
		},
		ClientId:   aws.String(r.clientID),
		UserPoolId: aws.String(r.userPoolID),
	}
	res, err := r.cli.AdminInitiateAuth(ctx, req)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	if res == nil || res.AuthenticationResult == nil || res.AuthenticationResult.IdToken == nil {
		return "", errors.InternalErr.New().WithDetail("failed to auth")
	}

	return *res.AuthenticationResult.IdToken, nil
}
