package output

import "github.com/abyssparanoia/rapid-go/internal/domain/model"

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
