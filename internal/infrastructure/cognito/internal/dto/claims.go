package dto

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt/v4"
	"github.com/volatiletech/null/v8"
)

type AWSCognitoClaims struct {
	Username      string `json:"cognito:username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	TenantID      string `json:"custom:tenant_id,omitepmty"`
	UserID        string `json:"custom:user_id,omitepmty"`
	UserRole      string `json:"custom:user_role,omitepmty"`
	jwt.StandardClaims
}

const (
	attributesTenantIDKey = "custom:tenant_id"
	attributesUserIDKey   = "custom:user_id"
	attributesUserRoleKey = "custom:user_role"
)

type UserAttributes struct {
	AuthUID string
	CustomUserAttributes
}

func NewUserAttributesFromCognitoUser(cognitoUser *cognitoidentityprovider.UserType) *UserAttributes {
	userAttributes := &UserAttributes{
		AuthUID: *cognitoUser.Username,
	}
	for _, attribute := range cognitoUser.Attributes {
		if *attribute.Name == attributesTenantIDKey {
			userAttributes.TenantID = null.StringFrom(*attribute.Value)
		}
		if *attribute.Name == attributesUserIDKey {
			userAttributes.UserID = null.StringFrom(*attribute.Value)
		}
		if *attribute.Name == attributesUserRoleKey {
			userAttributes.UserRole = null.StringFrom(*attribute.Value)
		}
	}

	return userAttributes
}

func NewUserAttributesFromClaims(awsClaims *AWSCognitoClaims) *UserAttributes {
	userAttributes := &UserAttributes{
		AuthUID: awsClaims.Username,
	}
	if awsClaims.TenantID != "" {
		userAttributes.TenantID = null.StringFrom(awsClaims.TenantID)
	}
	if awsClaims.UserID != "" {
		userAttributes.UserID = null.StringFrom(awsClaims.UserID)
	}
	if awsClaims.UserRole != "" {
		userAttributes.UserRole = null.StringFrom(awsClaims.UserRole)
	}
	return userAttributes
}

type CustomUserAttributes struct {
	TenantID null.String
	UserID   null.String
	UserRole null.String
}

func (ua *CustomUserAttributes) ToSlice() []*cognitoidentityprovider.AttributeType {
	attrs := []*cognitoidentityprovider.AttributeType{}
	if ua.TenantID.Valid {
		attrs = append(attrs, NewUserAttribute("custom:tenant_id", ua.TenantID.String))
	}
	if ua.UserID.Valid {
		attrs = append(attrs, NewUserAttribute("custom:user_id", ua.UserID.String))
	}
	if ua.UserRole.Valid {
		attrs = append(attrs, NewUserAttribute("custom:user_role", ua.UserRole.String))
	}
	return attrs
}

func NewUserAttribute(name, value string) *cognitoidentityprovider.AttributeType {
	attr := &cognitoidentityprovider.AttributeType{
		Name:  aws.String(name),
		Value: aws.String(value),
	}
	return attr
}
