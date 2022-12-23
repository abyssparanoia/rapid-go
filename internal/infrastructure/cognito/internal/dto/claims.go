package dto

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt/v4"
	"github.com/volatiletech/null/v8"
)

type AWSCognitoClaims struct {
	ClientID string `json:"client_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type UserAttributes struct {
	TenantID null.String
	UserID   null.String
	UserRole null.String
}

func (ua *UserAttributes) ToSlice() []*cognitoidentityprovider.AttributeType {
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
	attr := &cognitoidentityprovider.AttributeType{}
	attr.SetName(name).
		SetValue(value)
	return attr
}
