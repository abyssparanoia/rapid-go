package model

import (
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/volatiletech/null/v8"
)

type StaffClaims struct {
	AuthUID   string
	TenantID  null.String
	StaffID   null.String
	StaffRole nullable.Type[StaffRole]
}

func NewStaffClaims(
	authUID string,
	tenantID null.String,
	staffID null.String,
	staffRole nullable.Type[StaffRole],
) *StaffClaims {
	return &StaffClaims{
		AuthUID:   authUID,
		TenantID:  tenantID,
		StaffID:   staffID,
		StaffRole: staffRole,
	}
}
