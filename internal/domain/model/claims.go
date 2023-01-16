package model

import (
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/volatiletech/null/v8"
)

type Claims struct {
	AuthUID  string
	TenantID null.String
	UserID   null.String
	UserRole nullable.Type[UserRole]
}

func NewClaims(
	authUID string,
	tenantID null.String,
	userID null.String,
	userRole nullable.Type[UserRole],
) *Claims {
	return &Claims{
		AuthUID:  authUID,
		TenantID: tenantID,
		UserID:   userID,
		UserRole: userRole,
	}
}
