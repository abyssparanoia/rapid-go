package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito/internal/dto"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/volatiletech/null/v8"
)

func UserAttributesToModel(
	userAttribute *dto.UserAttributes,
) *model.Claims {
	var staffRole nullable.Type[model.StaffRole]
	if userAttribute.StaffRole.Valid {
		staffRole = nullable.TypeFrom(model.NewStaffRole(userAttribute.StaffRole.String))
	}
	claims := model.NewClaims(
		userAttribute.AuthUID,
		userAttribute.TenantID,
		userAttribute.StaffID,
		staffRole,
	)
	return claims
}

func ClaimsToCustomUserAttributes(
	claims *model.Claims,
) *dto.CustomUserAttributes {
	ua := &dto.CustomUserAttributes{
		TenantID: claims.TenantID,
		StaffID:  claims.StaffID,
	}
	if claims.StaffRole.Valid && claims.StaffRole.Value.Valid() {
		ua.StaffRole = null.StringFrom(claims.StaffRole.Value.String())
	}
	return ua
}
