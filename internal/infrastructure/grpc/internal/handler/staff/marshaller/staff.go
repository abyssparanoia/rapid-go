package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	staff_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/staff_api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func StaffToPB(m *model.Staff) *staff_apiv1.Staff {
	if m == nil {
		return nil
	}

	var tenant *staff_apiv1.TenantPartial
	if m.ReadonlyReference != nil && m.ReadonlyReference.Tenant != nil {
		tenant = TenantPartialToPB(m.ReadonlyReference.Tenant)
	}

	return &staff_apiv1.Staff{
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

func StaffPartialToPB(m *model.Staff) *staff_apiv1.StaffPartial {
	if m == nil {
		return nil
	}

	var tenant *staff_apiv1.TenantPartial
	if m.ReadonlyReference != nil && m.ReadonlyReference.Tenant != nil {
		tenant = TenantPartialToPB(m.ReadonlyReference.Tenant)
	}

	return &staff_apiv1.StaffPartial{
		Id:          m.ID,
		Tenant:      tenant,
		Role:        StaffRoleToPB(m.Role),
		AuthUid:     m.AuthUID,
		DisplayName: m.DisplayName,
		ImageUrl:    m.ImageURL.String,
		Email:       m.Email,
	}
}

func StaffsToPB(slice model.Staffs) []*staff_apiv1.Staff {
	dsts := make([]*staff_apiv1.Staff, len(slice))
	for idx, m := range slice {
		dsts[idx] = StaffToPB(m)
	}
	return dsts
}

func StaffsPartialToPB(slice model.Staffs) []*staff_apiv1.StaffPartial {
	dsts := make([]*staff_apiv1.StaffPartial, len(slice))
	for idx, m := range slice {
		dsts[idx] = StaffPartialToPB(m)
	}
	return dsts
}

func StaffSortKeyToPB(s model.StaffSortKey) staff_apiv1.ListStaffsRequest_ListStaffsSortKey {
	switch s {
	case model.StaffSortKeyCreatedAtDesc:
		return staff_apiv1.ListStaffsRequest_LIST_STAFFS_SORT_KEY_CREATED_AT_DESC
	case model.StaffSortKeyCreatedAtAsc:
		return staff_apiv1.ListStaffsRequest_LIST_STAFFS_SORT_KEY_CREATED_AT_ASC
	case model.StaffSortKeyDisplayNameAsc:
		return staff_apiv1.ListStaffsRequest_LIST_STAFFS_SORT_KEY_DISPLAY_NAME_ASC
	case model.StaffSortKeyDisplayNameDesc:
		return staff_apiv1.ListStaffsRequest_LIST_STAFFS_SORT_KEY_DISPLAY_NAME_DESC
	case model.StaffSortKeyUnknown:
		return staff_apiv1.ListStaffsRequest_LIST_STAFFS_SORT_KEY_UNSPECIFIED
	default:
		return staff_apiv1.ListStaffsRequest_LIST_STAFFS_SORT_KEY_UNSPECIFIED
	}
}

func StaffSortKeyToModel(s staff_apiv1.ListStaffsRequest_ListStaffsSortKey) model.StaffSortKey {
	switch s {
	case staff_apiv1.ListStaffsRequest_LIST_STAFFS_SORT_KEY_CREATED_AT_DESC:
		return model.StaffSortKeyCreatedAtDesc
	case staff_apiv1.ListStaffsRequest_LIST_STAFFS_SORT_KEY_CREATED_AT_ASC:
		return model.StaffSortKeyCreatedAtAsc
	case staff_apiv1.ListStaffsRequest_LIST_STAFFS_SORT_KEY_DISPLAY_NAME_ASC:
		return model.StaffSortKeyDisplayNameAsc
	case staff_apiv1.ListStaffsRequest_LIST_STAFFS_SORT_KEY_DISPLAY_NAME_DESC:
		return model.StaffSortKeyDisplayNameDesc
	case staff_apiv1.ListStaffsRequest_LIST_STAFFS_SORT_KEY_UNSPECIFIED:
		return model.StaffSortKeyUnknown
	default:
		return model.StaffSortKeyUnknown
	}
}
