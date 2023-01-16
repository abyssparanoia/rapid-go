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
	var userRole nullable.Type[model.UserRole]
	if userAttribute.UserRole.Valid {
		userRole = nullable.TypeFrom(model.NewUserRole(userAttribute.UserRole.String))
	}
	claims := model.NewClaims(
		userAttribute.AuthUID,
		userAttribute.TenantID,
		userAttribute.UserID,
		userRole,
	)
	return claims
}

func ClaimsToCustomUserAttributes(
	claims *model.Claims,
) *dto.CustomUserAttributes {
	ua := &dto.CustomUserAttributes{
		TenantID: claims.TenantID,
		UserID:   claims.UserID,
	}
	if claims.UserRole.Valid && claims.UserRole.Value.Valid() {
		ua.UserRole = null.StringFrom(claims.UserRole.Value.String())
	}
	return ua
}
