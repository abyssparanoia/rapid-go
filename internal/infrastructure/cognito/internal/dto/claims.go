package dto

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/volatiletech/null/v8"
)

type AWSCognitoClaims struct {
	Username      string `json:"cognito:username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	TenantID      string `json:"custom:tenant_id,omitempty"`
	StaffID       string `json:"custom:staff_id,omitempty"`
	StaffRole     string `json:"custom:staff_role,omitempty"`
	jwt.RegisteredClaims
}

const (
	attributesTenantIDKey  = "custom:tenant_id"
	attributesStaffIDKey   = "custom:staff_id"
	attributesStaffRoleKey = "custom:staff_role"
)

type UserAttributes struct {
	AuthUID string
	CustomUserAttributes
}

func NewUserAttributesFromCognitoUser(cognitoUser *types.UserType) *UserAttributes {
	userAttributes := &UserAttributes{
		AuthUID: *cognitoUser.Username,
	}
	for _, attribute := range cognitoUser.Attributes {
		if *attribute.Name == attributesTenantIDKey {
			userAttributes.TenantID = null.StringFrom(*attribute.Value)
		}
		if *attribute.Name == attributesStaffIDKey {
			userAttributes.StaffID = null.StringFrom(*attribute.Value)
		}
		if *attribute.Name == attributesStaffRoleKey {
			userAttributes.StaffRole = null.StringFrom(*attribute.Value)
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
	if awsClaims.StaffID != "" {
		userAttributes.StaffID = null.StringFrom(awsClaims.StaffID)
	}
	if awsClaims.StaffRole != "" {
		userAttributes.StaffRole = null.StringFrom(awsClaims.StaffRole)
	}
	return userAttributes
}

type CustomUserAttributes struct {
	TenantID  null.String
	StaffID   null.String
	StaffRole null.String
}

func (ua *CustomUserAttributes) ToSlice() []types.AttributeType {
	attrs := []types.AttributeType{}
	if ua.TenantID.Valid {
		attrs = append(attrs, *NewUserAttribute("custom:tenant_id", ua.TenantID.String))
	}
	if ua.StaffID.Valid {
		attrs = append(attrs, *NewUserAttribute("custom:staff_id", ua.StaffID.String))
	}
	if ua.StaffRole.Valid {
		attrs = append(attrs, *NewUserAttribute("custom:staff_role", ua.StaffRole.String))
	}
	return attrs
}

func NewUserAttribute(name, value string) *types.AttributeType {
	attr := &types.AttributeType{
		Name:  aws.String(name),
		Value: aws.String(value),
	}
	return attr
}
