package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	staff_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/staff_api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TenantToPB(m *model.Tenant) *staff_apiv1.Tenant {
	if m == nil {
		return nil
	}
	return &staff_apiv1.Tenant{
		Id:        m.ID,
		Name:      m.Name,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}

func TenantPartialToPB(m *model.Tenant) *staff_apiv1.TenantPartial {
	if m == nil {
		return nil
	}
	return &staff_apiv1.TenantPartial{
		Id:   m.ID,
		Name: m.Name,
	}
}

func TenantsToPB(slice model.Tenants) []*staff_apiv1.Tenant {
	dsts := make([]*staff_apiv1.Tenant, len(slice))
	for idx, m := range slice {
		dsts[idx] = TenantToPB(m)
	}
	return dsts
}
