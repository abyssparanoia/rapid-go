package marshaller

import (
	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
)

func StaffClaimsToModel(authUID string, customClaim map[string]interface{}) *model.StaffClaims {
	var tenantID null.String
	if _tenantID, ok := customClaim["tenant_id"]; ok {
		tenantID = null.StringFrom(_tenantID.(string)) //nolint:errcheck
	}
	var staffID null.String
	if _staffID, ok := customClaim["staff_id"]; ok {
		staffID = null.StringFrom(_staffID.(string)) //nolint:errcheck
	}
	var staffRole nullable.Type[model.StaffRole]
	if _staffRole, ok := customClaim["staff_role"]; ok {
		staffRole = nullable.TypeFrom(model.NewStaffRole(_staffRole.(string))) //nolint:errcheck
	}

	claims := model.NewStaffClaims(
		authUID,
		tenantID,
		staffID,
		staffRole,
	)
	return claims
}

func StaffClaimsToMap(m *model.StaffClaims) map[string]interface{} {
	cmap := map[string]interface{}{}
	if m.TenantID.Valid {
		cmap["tenant_id"] = m.TenantID.String
	}
	if m.StaffID.Valid {
		cmap["staff_id"] = m.StaffID.String
	}
	if m.StaffRole.Valid && m.StaffRole.Value().Valid() {
		cmap["staff_role"] = m.StaffRole.Value().String()
	}
	return cmap
}

func AdminClaimsToModel(authUID string, email string, customClaim map[string]interface{}) *model.AdminClaims {
	var adminID null.String
	if _adminID, ok := customClaim["admin_id"]; ok {
		adminID = null.StringFrom(_adminID.(string)) //nolint:errcheck
	}
	var adminRole nullable.Type[model.AdminRole]
	if _adminRole, ok := customClaim["admin_role"]; ok {
		adminRole = nullable.TypeFrom(model.NewAdminRole(_adminRole.(string))) //nolint:errcheck
	}

	claims := model.NewAdminClaims(
		authUID,
		email,
		adminID,
		adminRole,
	)
	return claims
}

func AdminClaimsToMap(m *model.AdminClaims) map[string]interface{} {
	cmap := map[string]interface{}{}
	if m.AdminID.Valid {
		cmap["admin_id"] = m.AdminID.String
	}
	if m.AdminRole.Valid && m.AdminRole.Value().Valid() {
		cmap["admin_role"] = m.AdminRole.Value().String()
	}
	return cmap
}
