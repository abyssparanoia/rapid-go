package service

import (
	"context"
	"testing"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/model/factory"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	mock_repository "github.com/abyssparanoia/rapid-go/internal/domain/repository/mock"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/pkg/ulid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func TestUserService_Create(t *testing.T) {
	mockPassword := "password"
	testdata := factory.NewFactory()
	requestTime := testdata.RequestTime
	tenant := testdata.Tenant
	mockULID := ulid.Mock()
	user := testdata.User
	user.ID = mockULID
	user.Tenant = nil

	claims := model.NewClaims(
		user.AuthUID,
		null.StringFrom(tenant.ID),
		null.StringFrom(user.ID),
		nullable.TypeFrom(user.Role),
	)

	type args struct {
		tenantID    string
		email       string
		password    string
		userRole    model.UserRole
		displayName string
		imagePath   string
		requestTime time.Time
	}

	type want struct {
		user           *model.User
		expectedResult error
	}

	tests := map[string]struct {
		args    args
		service func(ctx context.Context, ctrl *gomock.Controller) User
		want    want
	}{
		"success": {
			args: args{
				tenantID:    tenant.ID,
				email:       user.Email,
				password:    mockPassword,
				userRole:    user.Role,
				displayName: user.DisplayName,
				imagePath:   user.ImagePath,
				requestTime: requestTime,
			},
			service: func(ctx context.Context, ctrl *gomock.Controller) User {
				mockUserRepo := mock_repository.NewMockUser(ctrl)
				mockAuthenticationRepo := mock_repository.NewMockAuthentication(ctrl)

				mockAuthenticationRepo.EXPECT().
					GetUserByEmail(
						gomock.Any(),
						user.Email,
					).
					Return(&repository.AuthenticationGetUserByEmailResult{}, nil)
				mockAuthenticationRepo.EXPECT().
					CreateUser(
						gomock.Any(),
						repository.AuthenticationCreateUserParam{
							Email:    user.Email,
							Password: null.StringFrom(mockPassword),
						},
					).
					Return(user.AuthUID, nil)
				mockUserRepo.EXPECT().
					Create(gomock.Any(), user).
					Return(user, nil)
				mockAuthenticationRepo.EXPECT().
					StoreClaims(
						gomock.Any(),
						user.AuthUID,
						claims,
					).
					Return(nil)

				return &userService{
					userRepository:           mockUserRepo,
					authenticationRepository: mockAuthenticationRepo,
				}
			},
			want: want{
				user: user,
			},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := tc.service(ctx, ctrl)

			got, err := s.Create(ctx, UserCreateParam{
				TenantID:    tc.args.tenantID,
				Email:       tc.args.email,
				Password:    tc.args.password,
				UserRole:    tc.args.userRole,
				DisplayName: tc.args.displayName,
				ImagePath:   tc.args.imagePath,
				RequestTime: tc.args.requestTime,
			})
			if tc.want.expectedResult == nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.want.user, got)
			} else {
				assert.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}
