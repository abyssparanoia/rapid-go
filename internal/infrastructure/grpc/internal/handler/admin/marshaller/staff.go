package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func StaffToPB(m *model.Staff) *admin_apiv1.Staff {
	if m == nil {
		return nil
	}

	var tenant *admin_apiv1.TenantPartial
	if m.ReadonlyReference != nil && m.ReadonlyReference.Tenant != nil {
		tenant = TenantPartialToPB(m.ReadonlyReference.Tenant)
	}

	return &admin_apiv1.Staff{
		Id:          m.ID,
		Tenant:      tenant,
		Role:        StaffRoleToPB(m.Role),
		AuthUid:     m.AuthUID,
		DisplayName: m.DisplayName,
		ImageUrl:    m.ImageURL.String,
		Email:       m.Email,
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),
	}
}

func StaffPartialToPB(m *model.Staff) *admin_apiv1.StaffPartial {
	if m == nil {
		return nil
	}

	var tenant *admin_apiv1.TenantPartial
	if m.ReadonlyReference != nil && m.ReadonlyReference.Tenant != nil {
		tenant = TenantPartialToPB(m.ReadonlyReference.Tenant)
	}

	return &admin_apiv1.StaffPartial{
		Id:          m.ID,
		Tenant:      tenant,
		Role:        StaffRoleToPB(m.Role),
		AuthUid:     m.AuthUID,
		DisplayName: m.DisplayName,
		ImageUrl:    m.ImageURL.String,
		Email:       m.Email,
	}
}

func StaffsToPB(slice model.Staffs) []*admin_apiv1.Staff {
	dsts := make([]*admin_apiv1.Staff, len(slice))
	for idx, m := range slice {
		dsts[idx] = StaffToPB(m)
	}
	return dsts
}

func StaffsPartialToPB(slice model.Staffs) []*admin_apiv1.StaffPartial {
	dsts := make([]*admin_apiv1.StaffPartial, len(slice))
	for idx, m := range slice {
		dsts[idx] = StaffPartialToPB(m)
	}
	return dsts
}
