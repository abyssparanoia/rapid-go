package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	staff_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/staff_api/v1"
)

func NewPagination(m *model.Pagination) *staff_apiv1.Pagination {
	return &staff_apiv1.Pagination{
		CurrentPage: m.CurrentPage,
		PrevPage:    m.PrevPage,
		NextPage:    m.NextPage,
		TotalPage:   m.TotalPage,
		TotalCount:  m.TotalCount,
		HasNext:     m.HasNext,
	}
}
