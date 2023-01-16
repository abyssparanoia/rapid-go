package service

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/volatiletech/null/v8"
)

type userService struct {
	userRepository           repository.User
	authenticationRepository repository.Authentication
}

func NewUser(
	userRepository repository.User,
	authenticationRepository repository.Authentication,
) User {
	return &userService{
		userRepository,
		authenticationRepository,
	}
}

func (s *userService) Create(
	ctx context.Context,
	param UserCreateParam,
) (*model.User, error) {
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

	user := model.NewUser(
		param.TenantID,
		param.UserRole,
		authUID,
		param.DisplayName,
		param.ImagePath,
		param.Email,
		param.RequestTime,
	)

	if _, err := s.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	claims := model.NewClaims(
		authUID,
		null.StringFrom(param.TenantID),
		null.StringFrom(user.ID),
		nullable.TypeFrom(user.Role),
	)
	if err := s.authenticationRepository.StoreClaims(ctx, authUID, claims); err != nil {
		return nil, err
	}

	return user, nil
}
