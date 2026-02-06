package output

import "github.com/abyssparanoia/rapid-go/internal/domain/model"

type AdminCreateStaff struct {
	Staff    *model.Staff
	Password string
}

func NewAdminCreateStaff(staff *model.Staff, password string) *AdminCreateStaff {
	return &AdminCreateStaff{
		Staff:    staff,
		Password: password,
	}
}

type ListStaffs struct {
	Staffs     model.Staffs
	Pagination *model.Pagination
}

func NewAdminListStaffs(
	staffs model.Staffs,
	pagination *model.Pagination,
) *ListStaffs {
	return &ListStaffs{
		Staffs:     staffs,
		Pagination: pagination,
	}
}

func NewStaffListStaffs(
	staffs model.Staffs,
	pagination *model.Pagination,
) *ListStaffs {
	return &ListStaffs{
		Staffs:     staffs,
		Pagination: pagination,
	}
}
