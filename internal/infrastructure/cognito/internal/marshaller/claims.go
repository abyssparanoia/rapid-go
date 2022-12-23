package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito/internal/dto"
	"github.com/volatiletech/null/v8"
)

func ClaimsToModel(
	username string,
) *model.Claims {
	claims := model.NewClaims(username)
	return claims
}

func ClaimsToUserAttributes(
	claims *model.Claims,
) *dto.UserAttributes {
	ua := &dto.UserAttributes{
		TenantID: claims.TenantID,
		UserID:   claims.UserID,
	}
	if claims.UserRole.Valid && claims.UserRole.Value.Valid() {
		ua.UserRole = null.StringFrom(claims.UserRole.Value.String())
	}
	return ua
}
