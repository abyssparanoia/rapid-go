package model

import (
	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
)

type StaffClaims struct {
	AuthUID   string
	Email     string
	TenantID  null.String
	StaffID   null.String
	StaffRole nullable.Type[StaffRole]
}

func NewStaffClaims(
	authUID string,
	email string,
	tenantID null.String,
	staffID null.String,
	staffRole nullable.Type[StaffRole],
) *StaffClaims {
	return &StaffClaims{
		AuthUID:   authUID,
		Email:     email,
		TenantID:  tenantID,
		StaffID:   staffID,
		StaffRole: staffRole,
	}
}
