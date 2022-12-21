package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	modelv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/model/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserToPB(m *model.User) *modelv1.User {
	if m == nil {
		return nil
	}
	return &modelv1.User{
		Id:          m.ID,
		Role:        UserRoleToPB(m.Role),
		AuthUid:     m.AuthUID,
		DisplayName: m.DisplayName,
		ImageUrl:    m.ImagePath,
		Email:       m.Email,
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),

		Tenant: TenantToPB(m.Tenant),
	}
}
