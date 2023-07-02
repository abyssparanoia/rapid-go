package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
)

func NewPagination(m *model.Pagination) *admin_apiv1.Pagination {
	return &admin_apiv1.Pagination{
		CurrentPage: m.CurrentPage,
		PrevPage:    m.PrevPage,
		NextPage:    m.NextPage,
		TotalPage:   m.TotalPage,
		TotalCount:  m.TotalCount,
		HasNext:     m.HasNext,
	}
}
