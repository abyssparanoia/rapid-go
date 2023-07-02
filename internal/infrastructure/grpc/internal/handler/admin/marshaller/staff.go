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
	dst := &admin_apiv1.Staff{
		Id:          m.ID,
		TenantId:    m.TenantID,
		Role:        StaffRoleToPB(m.Role),
		AuthUid:     m.AuthUID,
		DisplayName: m.DisplayName,
		ImageUrl:    m.ImagePath,
		Email:       m.Email,
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),
	}
	if m.Tenant != nil {
		dst.Tenant = TenantToPB(m.Tenant)
	}
	return dst
}
