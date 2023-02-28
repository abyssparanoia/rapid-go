package model

import (
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/volatiletech/null/v8"
)

type Claims struct {
	AuthUID   string
	TenantID  null.String
	StaffID   null.String
	StaffRole nullable.Type[StaffRole]
}

func NewClaims(
	authUID string,
	tenantID null.String,
	userID null.String,
	userRole nullable.Type[StaffRole],
) *Claims {
	return &Claims{
		AuthUID:   authUID,
		TenantID:  tenantID,
		StaffID:   userID,
		StaffRole: userRole,
	}
}
