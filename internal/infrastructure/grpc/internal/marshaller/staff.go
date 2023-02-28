package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	modelv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/model/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func StaffToPB(m *model.Staff) *modelv1.Staff {
	if m == nil {
		return nil
	}
	dst := &modelv1.Staff{
		Id:          m.ID,
		Role:        StaffRoleToPB(m.Role),
		AuthUid:     m.AuthUID,
		DisplayName: m.DisplayName,
		ImageUrl:    m.ImagePath,
		Email:       m.Email,
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),
	}
	if m.Tenant == nil {
		dst.OneofTenant = &modelv1.Staff_TenantId{
			TenantId: m.TenantID,
		}
	} else {
		dst.OneofTenant = &modelv1.Staff_Tenant{
			Tenant: TenantToPB(m.Tenant),
		}
	}
	return dst
}
