package service

import (
	"context"
	"testing"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/model/factory"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	mock_repository "github.com/abyssparanoia/rapid-go/internal/domain/repository/mock"
	"github.com/abyssparanoia/rapid-go/internal/pkg/id"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func TestStaffService_Create(t *testing.T) {
	t.Parallel()
	mockPassword := "password"
	testdata := factory.NewFactory()
	requestTime := testdata.RequestTime
	tenant := testdata.Tenant
	mockID := id.Mock()
	staff := testdata.Staff
	staff.ID = mockID
	staff.Tenant = nil

	claims := model.NewStaffClaims(
		staff.AuthUID,
		null.StringFrom(tenant.ID),
		null.StringFrom(staff.ID),
		nullable.TypeFrom(staff.Role),
	)

	type args struct {
		tenantID    string
		email       string
		password    string
		staffRole   model.StaffRole
		displayName string
		imagePath   string
		requestTime time.Time
	}

	type want struct {
		staff          *model.Staff
		expectedResult error
	}

	tests := map[string]struct {
		args    args
		service func(ctx context.Context, ctrl *gomock.Controller) Staff
		want    want
	}{
		"success": {
			args: args{
				tenantID:    tenant.ID,
				email:       staff.Email,
				password:    mockPassword,
				staffRole:   staff.Role,
				displayName: staff.DisplayName,
				imagePath:   staff.ImagePath,
				requestTime: requestTime,
			},
			service: func(ctx context.Context, ctrl *gomock.Controller) Staff {
				mockStaffRepo := mock_repository.NewMockStaff(ctrl)
				mockStaffAuthenticationRepo := mock_repository.NewMockStaffAuthentication(ctrl)

				mockStaffAuthenticationRepo.EXPECT().
					GetUserByEmail(
						gomock.Any(),
						staff.Email,
					).
					Return(&repository.StaffAuthenticationGetUserByEmailResult{}, nil)
				mockStaffAuthenticationRepo.EXPECT().
					CreateUser(
						gomock.Any(),
						repository.StaffAuthenticationCreateUserParam{
							Email:    staff.Email,
							Password: null.StringFrom(mockPassword),
						},
					).
					Return(staff.AuthUID, nil)
				mockStaffRepo.EXPECT().
					Create(gomock.Any(), staff).
					Return(staff, nil)
				mockStaffAuthenticationRepo.EXPECT().
					StoreClaims(
						gomock.Any(),
						staff.AuthUID,
						claims,
					).
					Return(nil)

				return &staffService{
					staffRepository:               mockStaffRepo,
					staffAuthenticationRepository: mockStaffAuthenticationRepo,
				}
			},
			want: want{
				staff: staff,
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

			got, err := s.Create(ctx, StaffCreateParam{
				TenantID:    tc.args.tenantID,
				Email:       tc.args.email,
				Password:    tc.args.password,
				StaffRole:   tc.args.staffRole,
				DisplayName: tc.args.displayName,
				ImagePath:   tc.args.imagePath,
				RequestTime: tc.args.requestTime,
			})
			if tc.want.expectedResult == nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.want.staff, got)
			} else {
				assert.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}
