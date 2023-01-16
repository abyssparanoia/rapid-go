package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/volatiletech/null/v8"
)

func ClaimsToModel(authUID string, customClaim map[string]interface{}) *model.Claims {

	var tenantID null.String
	if _tenantID, ok := customClaim["tenant_id"]; ok {
		tenantID = null.StringFrom(_tenantID.(string))
	}
	var userID null.String
	if _userID, ok := customClaim["service_user_id"]; ok {
		userID = null.StringFrom(_userID.(string))
	}
	var userRole nullable.Type[model.UserRole]
	if _userRole, ok := customClaim["user_role"]; ok {
		userRole = nullable.TypeFrom(model.NewUserRole(_userRole.(string)))
	}

	claims := model.NewClaims(
		authUID,
		tenantID,
		userID,
		userRole,
	)
	return claims
}

func ClaimsToMap(m *model.Claims) map[string]interface{} {
	cmap := map[string]interface{}{}
	if m.TenantID.Valid {
		cmap["tenant_id"] = m.TenantID.String
	}
	if m.UserID.Valid {
		cmap["service_user_id"] = m.UserID.String
	}
	if m.UserRole.Valid && m.UserRole.Value.Valid() {
		cmap["user_role"] = m.UserRole.Value.String()
	}
	return cmap
}
