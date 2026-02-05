package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/model/factory"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	mock_repository "github.com/abyssparanoia/rapid-go/internal/domain/repository/mock"
	mock_service "github.com/abyssparanoia/rapid-go/internal/domain/service/mock"
	"github.com/abyssparanoia/rapid-go/internal/pkg/id"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStaffMeInteractor_SignUp(t *testing.T) {
	t.Parallel()

	type args struct {
		authUID      string
		email        string
		tenantName   string
		displayName  string
		imageAssetID string
		requestTime  time.Time
	}

	type want struct {
		staff          *model.Staff
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase StaffMeInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			mockTransactable := mock_repository.NewMockTransactable(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffRepo := mock_repository.NewMockStaff(ctrl)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)

			return testcase{
				args: args{
					// Empty args to trigger validation error
				},
				usecase: &staffMeInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					staffRepository:  mockStaffRepo,
					staffService:     mockStaffService,
					assetService:     mockAssetService,
				},
				want: want{
					expectedResult: errors.RequestInvalidArgumentErr,
				},
			}
		},
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			tenant := testdata.Tenant
			staff := testdata.Staff
			mockID := id.Mock()
			tenant.ID = mockID
			staff.ID = mockID
			staff.TenantID = tenant.ID
			requestTime := time.Now()

			mockTransactable := mock_repository.TestMockTransactable()
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockTenantRepo.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(nil)
			mockStaffRepo := mock_repository.NewMockStaff(ctrl)
			mockStaffRepo.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Return(staff, nil)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockStaffService.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(staff, nil)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				GetWithValidate(gomock.Any(), gomock.Any(), gomock.Any()).
				Return("/path/to/image", nil)
			mockAssetService.EXPECT().
				BatchSetStaffURLs(gomock.Any(), gomock.Any()).
				Return(nil)

			return testcase{
				args: args{
					authUID:      "auth-uid",
					email:        "test@example.com",
					tenantName:   "Test Tenant",
					displayName:  "Test User",
					imageAssetID: "image-asset-id",
					requestTime:  requestTime,
				},
				usecase: &staffMeInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					staffRepository:  mockStaffRepo,
					staffService:     mockStaffService,
					assetService:     mockAssetService,
				},
				want: want{
					staff: staff,
				},
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tc := tc(ctx, ctrl)

			got, err := tc.usecase.SignUp(ctx, input.NewStaffSignUp(
				tc.args.authUID,
				tc.args.email,
				tc.args.tenantName,
				tc.args.displayName,
				tc.args.imageAssetID,
				tc.args.requestTime,
			))

			if tc.want.expectedResult == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want.staff, got)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}

func TestStaffMeInteractor_Get(t *testing.T) {
	t.Parallel()

	type args struct {
		tenantID    string
		staffID     string
		requestTime time.Time
	}

	type want struct {
		staff          *model.Staff
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase StaffMeInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			mockTransactable := mock_repository.NewMockTransactable(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffRepo := mock_repository.NewMockStaff(ctrl)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)

			return testcase{
				args: args{
					// Empty args to trigger validation error
				},
				usecase: &staffMeInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					staffRepository:  mockStaffRepo,
					staffService:     mockStaffService,
					assetService:     mockAssetService,
				},
				want: want{
					expectedResult: errors.RequestInvalidArgumentErr,
				},
			}
		},
		"not found": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			staff := testdata.Staff
			mockID := id.Mock()
			staff.ID = mockID
			requestTime := time.Now()

			mockTransactable := mock_repository.NewMockTransactable(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffRepo := mock_repository.NewMockStaff(ctrl)
			mockStaffRepo.EXPECT().
				Get(gomock.Any(),
					repository.GetStaffQuery{
						ID: null.StringFrom(staff.ID),
						BaseGetOptions: repository.BaseGetOptions{
							OrFail:  true,
							Preload: true,
						},
					}).
				Return(nil, errors.StaffNotFoundErr)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)

			return testcase{
				args: args{
					tenantID:    staff.TenantID,
					staffID:     staff.ID,
					requestTime: requestTime,
				},
				usecase: &staffMeInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					staffRepository:  mockStaffRepo,
					staffService:     mockStaffService,
					assetService:     mockAssetService,
				},
				want: want{
					expectedResult: errors.StaffNotFoundErr,
				},
			}
		},
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			staff := testdata.Staff
			mockID := id.Mock()
			staff.ID = mockID
			requestTime := time.Now()

			mockTransactable := mock_repository.NewMockTransactable(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffRepo := mock_repository.NewMockStaff(ctrl)
			mockStaffRepo.EXPECT().
				Get(gomock.Any(),
					repository.GetStaffQuery{
						ID: null.StringFrom(staff.ID),
						BaseGetOptions: repository.BaseGetOptions{
							OrFail:  true,
							Preload: true,
						},
					}).
				Return(staff, nil)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				BatchSetStaffURLs(gomock.Any(), gomock.Any()).
				Return(nil)

			return testcase{
				args: args{
					tenantID:    staff.TenantID,
					staffID:     staff.ID,
					requestTime: requestTime,
				},
				usecase: &staffMeInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					staffRepository:  mockStaffRepo,
					staffService:     mockStaffService,
					assetService:     mockAssetService,
				},
				want: want{
					staff: staff,
				},
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tc := tc(ctx, ctrl)

			got, err := tc.usecase.Get(ctx, input.NewStaffGetMe(
				tc.args.tenantID,
				tc.args.staffID,
				tc.args.requestTime,
			))

			if tc.want.expectedResult == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want.staff, got)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}

func TestStaffMeInteractor_Update(t *testing.T) {
	t.Parallel()

	type args struct {
		tenantID     string
		staffID      string
		displayName  null.String
		imageAssetID null.String
		requestTime  time.Time
	}

	type want struct {
		staff          *model.Staff
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase StaffMeInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			mockTransactable := mock_repository.NewMockTransactable(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffRepo := mock_repository.NewMockStaff(ctrl)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)

			return testcase{
				args: args{
					// Empty args to trigger validation error
				},
				usecase: &staffMeInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					staffRepository:  mockStaffRepo,
					staffService:     mockStaffService,
					assetService:     mockAssetService,
				},
				want: want{
					expectedResult: errors.RequestInvalidArgumentErr,
				},
			}
		},
		"not found": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			staff := testdata.Staff
			mockID := id.Mock()
			staff.ID = mockID
			requestTime := time.Now()

			mockTransactable := mock_repository.TestMockTransactable()
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffRepo := mock_repository.NewMockStaff(ctrl)
			mockStaffRepo.EXPECT().
				Get(gomock.Any(),
					repository.GetStaffQuery{
						ID: null.StringFrom(staff.ID),
						BaseGetOptions: repository.BaseGetOptions{
							OrFail:    true,
							ForUpdate: true,
						},
					}).
				Return(nil, errors.StaffNotFoundErr)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)

			return testcase{
				args: args{
					tenantID:     staff.TenantID,
					staffID:      staff.ID,
					displayName:  null.StringFrom("Updated Name"),
					imageAssetID: null.String{},
					requestTime:  requestTime,
				},
				usecase: &staffMeInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					staffRepository:  mockStaffRepo,
					staffService:     mockStaffService,
					assetService:     mockAssetService,
				},
				want: want{
					expectedResult: errors.StaffNotFoundErr,
				},
			}
		},
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			staff := testdata.Staff
			mockID := id.Mock()
			staff.ID = mockID
			requestTime := time.Now()

			mockTransactable := mock_repository.TestMockTransactable()
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffRepo := mock_repository.NewMockStaff(ctrl)
			mockStaffRepo.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Times(2).
				Return(staff, nil)
			mockStaffRepo.EXPECT().
				Update(gomock.Any(), gomock.Any()).
				Return(nil)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				GetWithValidate(gomock.Any(), gomock.Any(), gomock.Any()).
				Return("/path/to/image", nil)
			mockAssetService.EXPECT().
				BatchSetStaffURLs(gomock.Any(), gomock.Any()).
				Return(nil)

			return testcase{
				args: args{
					tenantID:     staff.TenantID,
					staffID:      staff.ID,
					displayName:  null.StringFrom("Updated Name"),
					imageAssetID: null.StringFrom("image-asset-id"),
					requestTime:  requestTime,
				},
				usecase: &staffMeInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					staffRepository:  mockStaffRepo,
					staffService:     mockStaffService,
					assetService:     mockAssetService,
				},
				want: want{
					staff: staff,
				},
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tc := tc(ctx, ctrl)

			got, err := tc.usecase.Update(ctx, input.NewStaffUpdateMe(
				tc.args.tenantID,
				tc.args.staffID,
				tc.args.displayName,
				tc.args.imageAssetID,
				tc.args.requestTime,
			))

			if tc.want.expectedResult == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want.staff, got)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}
