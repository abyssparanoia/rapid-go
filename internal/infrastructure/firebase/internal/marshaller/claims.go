package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

func ClaimsToModel(authUID string, customClaim map[string]interface{}) *model.Claims {
	claims := model.NewClaims(authUID)
	if tenantID, ok := customClaim["tenant_id"]; ok {
		claims.SetTenantID(tenantID.(string))
	}
	if userID, ok := customClaim["service_user_id"]; ok {
		claims.SetUserID(userID.(string))
	}
	if userRole, ok := customClaim["user_role"]; ok {
		claims.SetUserRole(model.NewUserRole(userRole.(string)))
	}

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
	if m.UserRole.Valid() {
		cmap["user_role"] = m.UserID.String
	}
	return cmap
}
