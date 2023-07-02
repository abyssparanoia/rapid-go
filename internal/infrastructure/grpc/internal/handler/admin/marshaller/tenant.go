package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TenantToPB(m *model.Tenant) *admin_apiv1.Tenant {
	if m == nil {
		return nil
	}
	return &admin_apiv1.Tenant{
		Id:        m.ID,
		Name:      m.Name,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}

func TenantsToPB(slice model.Tenants) []*admin_apiv1.Tenant {
	dsts := make([]*admin_apiv1.Tenant, len(slice))
	for idx, m := range slice {
		dsts[idx] = TenantToPB(m)
	}
	return dsts
}
