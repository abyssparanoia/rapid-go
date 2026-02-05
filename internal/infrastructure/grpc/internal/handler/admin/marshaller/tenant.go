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

func TenantPartialToPB(m *model.Tenant) *admin_apiv1.TenantPartial {
	if m == nil {
		return nil
	}
	return &admin_apiv1.TenantPartial{
		Id:   m.ID,
		Name: m.Name,
	}
}

func TenantsToPB(slice model.Tenants) []*admin_apiv1.Tenant {
	dsts := make([]*admin_apiv1.Tenant, len(slice))
	for idx, m := range slice {
		dsts[idx] = TenantToPB(m)
	}
	return dsts
}

func TenantSortKeyToPB(s model.TenantSortKey) admin_apiv1.ListTenantsRequest_ListTenantsSortKey {
	switch s {
	case model.TenantSortKeyCreatedAtDesc:
		return admin_apiv1.ListTenantsRequest_LIST_TENANTS_SORT_KEY_CREATED_AT_DESC
	case model.TenantSortKeyCreatedAtAsc:
		return admin_apiv1.ListTenantsRequest_LIST_TENANTS_SORT_KEY_CREATED_AT_ASC
	case model.TenantSortKeyNameAsc:
		return admin_apiv1.ListTenantsRequest_LIST_TENANTS_SORT_KEY_NAME_ASC
	case model.TenantSortKeyNameDesc:
		return admin_apiv1.ListTenantsRequest_LIST_TENANTS_SORT_KEY_NAME_DESC
	case model.TenantSortKeyUnknown:
		return admin_apiv1.ListTenantsRequest_LIST_TENANTS_SORT_KEY_UNSPECIFIED
	default:
		return admin_apiv1.ListTenantsRequest_LIST_TENANTS_SORT_KEY_UNSPECIFIED
	}
}

func TenantSortKeyToModel(s admin_apiv1.ListTenantsRequest_ListTenantsSortKey) model.TenantSortKey {
	switch s {
	case admin_apiv1.ListTenantsRequest_LIST_TENANTS_SORT_KEY_CREATED_AT_DESC:
		return model.TenantSortKeyCreatedAtDesc
	case admin_apiv1.ListTenantsRequest_LIST_TENANTS_SORT_KEY_CREATED_AT_ASC:
		return model.TenantSortKeyCreatedAtAsc
	case admin_apiv1.ListTenantsRequest_LIST_TENANTS_SORT_KEY_NAME_ASC:
		return model.TenantSortKeyNameAsc
	case admin_apiv1.ListTenantsRequest_LIST_TENANTS_SORT_KEY_NAME_DESC:
		return model.TenantSortKeyNameDesc
	case admin_apiv1.ListTenantsRequest_LIST_TENANTS_SORT_KEY_UNSPECIFIED:
		return model.TenantSortKeyUnknown
	default:
		return model.TenantSortKeyUnknown
	}
}
