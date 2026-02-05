package dto

import (
	"github.com/aarondl/null/v8"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
)

type AWSCognitoAdminClaims struct {
	Username      string `json:"cognito:username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AdminID       string `json:"custom:admin_id,omitempty"`
	AdminRole     string `json:"custom:admin_role,omitempty"`
	jwt.RegisteredClaims
}

const (
	adminAttributesAdminIDKey   = "custom:admin_id"
	adminAttributesAdminRoleKey = "custom:admin_role"
)

type AdminUserAttributes struct {
	AuthUID string
	Email   string
	AdminCustomUserAttributes
}

func NewAdminUserAttributesFromCognitoUser(cognitoUser *types.UserType) *AdminUserAttributes {
	userAttributes := &AdminUserAttributes{
		AuthUID: *cognitoUser.Username,
	}

	for _, attribute := range cognitoUser.Attributes {
		if *attribute.Name == "email" {
			userAttributes.Email = *attribute.Value
		}
		if *attribute.Name == adminAttributesAdminIDKey {
			userAttributes.AdminID = null.StringFrom(*attribute.Value)
		}
		if *attribute.Name == adminAttributesAdminRoleKey {
			userAttributes.AdminRole = null.StringFrom(*attribute.Value)
		}
	}

	return userAttributes
}

func NewAdminUserAttributesFromClaims(awsClaims *AWSCognitoAdminClaims) *AdminUserAttributes {
	userAttributes := &AdminUserAttributes{
		AuthUID: awsClaims.Username,
		Email:   awsClaims.Email,
	}
	if awsClaims.AdminID != "" {
		userAttributes.AdminID = null.StringFrom(awsClaims.AdminID)
	}
	if awsClaims.AdminRole != "" {
		userAttributes.AdminRole = null.StringFrom(awsClaims.AdminRole)
	}
	return userAttributes
}

type AdminCustomUserAttributes struct {
	AdminID   null.String
	AdminRole null.String
}

func (ua *AdminCustomUserAttributes) ToSlice() []types.AttributeType {
	attrs := []types.AttributeType{}
	if ua.AdminID.Valid {
		attrs = append(attrs, *NewUserAttribute("custom:admin_id", ua.AdminID.String))
	}
	if ua.AdminRole.Valid {
		attrs = append(attrs, *NewUserAttribute("custom:admin_role", ua.AdminRole.String))
	}
	return attrs
}
