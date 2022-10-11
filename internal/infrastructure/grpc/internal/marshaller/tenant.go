package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	modelv1 "github.com/abyssparanoia/rapid-go/schema/proto/pb/rapid/model/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TenantToPB(m *model.Tenant) *modelv1.Tenant {
	if m == nil {
		return nil
	}
	return &modelv1.Tenant{
		Id:        m.ID,
		Name:      m.Name,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}

func TenantsToPB(slice []*model.Tenant) []*modelv1.Tenant {
	dsts := make([]*modelv1.Tenant, len(slice))
	for idx, m := range slice {
		dsts[idx] = TenantToPB(m)
	}
	return dsts
}
