package dto

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/volatiletech/null/v8"
)

type AWSCognitoStaffClaims struct {
	Username      string `json:"cognito:username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	TenantID      string `json:"custom:tenant_id,omitempty"`
	StaffID       string `json:"custom:staff_id,omitempty"`
	StaffRole     string `json:"custom:staff_role,omitempty"`
	jwt.RegisteredClaims
}

const (
	staffAttributesTenantIDKey  = "custom:tenant_id"
	staffAttributesStaffIDKey   = "custom:staff_id"
	staffAttributesStaffRoleKey = "custom:staff_role"
)

type StaffUserAttributes struct {
	AuthUID string
	StaffCustomUserAttributes
}

func NewUserAttributesFromCognitoUser(cognitoUser *types.UserType) *StaffUserAttributes {
	userAttributes := &StaffUserAttributes{
		AuthUID: *cognitoUser.Username,
	}
	for _, attribute := range cognitoUser.Attributes {
		if *attribute.Name == staffAttributesTenantIDKey {
			userAttributes.TenantID = null.StringFrom(*attribute.Value)
		}
		if *attribute.Name == staffAttributesStaffIDKey {
			userAttributes.StaffID = null.StringFrom(*attribute.Value)
		}
		if *attribute.Name == staffAttributesStaffRoleKey {
			userAttributes.StaffRole = null.StringFrom(*attribute.Value)
		}
	}

	return userAttributes
}

func NewUserAttributesFromClaims(awsClaims *AWSCognitoStaffClaims) *StaffUserAttributes {
	userAttributes := &StaffUserAttributes{
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

type StaffCustomUserAttributes struct {
	TenantID  null.String
	StaffID   null.String
	StaffRole null.String
}

func (ua *StaffCustomUserAttributes) ToSlice() []types.AttributeType {
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
