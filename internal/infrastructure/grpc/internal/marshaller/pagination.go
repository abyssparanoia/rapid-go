package marshaller

import (
	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	modelv1 "github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/grpc/pb/mmg/model/v1"
)

func NewPagination(m *model.Pagination) *modelv1.Pagination {
	return &modelv1.Pagination{
		CurrentPage: m.CurrentPage,
		PrevPage:    m.PrevPage,
		NextPage:    m.NextPage,
		TotalPage:   m.TotalPage,
		TotalCount:  m.TotalCount,
		HasNext:     m.HasNext,
	}
}
