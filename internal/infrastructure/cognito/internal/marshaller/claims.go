package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito/internal/dto"
	"github.com/volatiletech/null/v8"
)

func UserAttributesToModel(
	userAttribute *dto.UserAttributes,
) *model.Claims {
	claims := model.NewClaims(userAttribute.AuthUID)
	if userAttribute.TenantID.Valid {
		claims.SetTenantID(userAttribute.TenantID.String)
	}
	if userAttribute.UserID.Valid {
		claims.SetUserID(userAttribute.UserID.String)
	}
	if userAttribute.UserRole.Valid {
		claims.SetUserRole(model.NewUserRole(userAttribute.UserRole.String))
	}
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
