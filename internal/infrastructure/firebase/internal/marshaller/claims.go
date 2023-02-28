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
	var staffID null.String
	if _staffID, ok := customClaim["staff_id"]; ok {
		staffID = null.StringFrom(_staffID.(string))
	}
	var staffRole nullable.Type[model.StaffRole]
	if _staffRole, ok := customClaim["staff_role"]; ok {
		staffRole = nullable.TypeFrom(model.NewStaffRole(_staffRole.(string)))
	}

	claims := model.NewClaims(
		authUID,
		tenantID,
		staffID,
		staffRole,
	)
	return claims
}

func ClaimsToMap(m *model.Claims) map[string]interface{} {
	cmap := map[string]interface{}{}
	if m.TenantID.Valid {
		cmap["tenant_id"] = m.TenantID.String
	}
	if m.StaffID.Valid {
		cmap["staff_id"] = m.StaffID.String
	}
	if m.StaffRole.Valid && m.StaffRole.Value.Valid() {
		cmap["staff_role"] = m.StaffRole.Value.String()
	}
	return cmap
}
