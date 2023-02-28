package service

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/volatiletech/null/v8"
)

type staffService struct {
	staffRepository          repository.Staff
	authenticationRepository repository.Authentication
}

func NewStaff(
	staffRepository repository.Staff,
	authenticationRepository repository.Authentication,
) Staff {
	return &staffService{
		staffRepository,
		authenticationRepository,
	}
}

func (s *staffService) Create(
	ctx context.Context,
	param StaffCreateParam,
) (*model.Staff, error) {
	res, err := s.authenticationRepository.GetUserByEmail(ctx, param.Email)
	if err != nil {
		return nil, err
	}
	var authUID string
	// 存在してない場合、新規作成する
	if !res.Exist {
		authUID, err = s.authenticationRepository.CreateUser(
			ctx,
			repository.AuthenticationCreateUserParam{
				Email:    param.Email,
				Password: null.StringFrom(param.Password),
			},
		)
		if err != nil {
			return nil, err
		}
	} else {
		authUID = res.AuthUID
	}

	staff := model.NewStaff(
		param.TenantID,
		param.StaffRole,
		authUID,
		param.DisplayName,
		param.ImagePath,
		param.Email,
		param.RequestTime,
	)

	if _, err := s.staffRepository.Create(ctx, staff); err != nil {
		return nil, err
	}

	claims := model.NewClaims(
		authUID,
		null.StringFrom(param.TenantID),
		null.StringFrom(staff.ID),
		nullable.TypeFrom(staff.Role),
	)
	if err := s.authenticationRepository.StoreClaims(ctx, authUID, claims); err != nil {
		return nil, err
	}

	return staff, nil
}
