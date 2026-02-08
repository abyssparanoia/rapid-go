package marshaller

import (
	"github.com/aarondl/null/v9"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito/internal/dto"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
)

func AdminUserAttributesToModel(
	userAttribute *dto.AdminUserAttributes,
) *model.AdminClaims {
	var adminRole nullable.Type[model.AdminRole]
	if userAttribute.AdminRole.Valid {
		adminRole = nullable.TypeFrom(model.NewAdminRole(userAttribute.AdminRole.String))
	}
	claims := model.NewAdminClaims(
		userAttribute.AuthUID,
		userAttribute.Email,
		userAttribute.AdminID,
		adminRole,
	)
	return claims
}

func AdminClaimsToAdminCustomUserAttributes(
	claims *model.AdminClaims,
) *dto.AdminCustomUserAttributes {
	ua := &dto.AdminCustomUserAttributes{
		AdminID: claims.AdminID,
	}
	if claims.AdminRole.Valid && claims.AdminRole.Value().Valid() {
		ua.AdminRole = null.StringFrom(claims.AdminRole.Value().String())
	}
	return ua
}
