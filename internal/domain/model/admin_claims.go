package model

import (
	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
)

type AdminClaims struct {
	AuthUID   string
	Email     string
	AdminID   null.String
	AdminRole nullable.Type[AdminRole]
}

func NewAdminClaims(
	authUID string,
	email string,
	adminID null.String,
	adminRole nullable.Type[AdminRole],
) *AdminClaims {
	return &AdminClaims{
		AuthUID:   authUID,
		Email:     email,
		AdminID:   adminID,
		AdminRole: adminRole,
	}
}
