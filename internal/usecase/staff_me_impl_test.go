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
	"github.com/abyssparanoia/rapid-go/internal/domain/service"
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
			staff := testdata.Staff
			mockID := id.Mock()
			requestTime := testdata.RequestTime

			// Create tenant with mock ID (matches what implementation will create)
			tenant := model.NewTenant("Test Tenant", requestTime)
			tenant.ID = mockID

			staff.ID = mockID
			staff.TenantID = mockID
			staff.ImagePath = "/path/to/image"

			mockTransactable := mock_repository.TestMockTransactable()
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockTenantRepo.EXPECT().
				Create(gomock.Any(), tenant).
				Return(nil)
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
			mockStaffService.EXPECT().
				Create(gomock.Any(),
					service.StaffCreateParam{
						TenantID:    tenant.ID,
						Email:       "test@example.com",
						Password:    "",
						StaffRole:   model.StaffRoleAdmin,
						DisplayName: "Test User",
						ImagePath:   "/path/to/image",
						RequestTime: requestTime,
					}).
				Return(staff, nil)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				GetWithValidate(gomock.Any(), model.AssetTypeUserImage, "image-asset-id").
				Return("/path/to/image", nil)
			mockAssetService.EXPECT().
				BatchSetStaffURLs(gomock.Any(), model.Staffs{staff}).
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
				BatchSetStaffURLs(gomock.Any(), model.Staffs{staff}).
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
			updatedStaff := &model.Staff{} //nolint:exhaustruct
			factory.CloneValue(staff, updatedStaff)
			mockID := id.Mock()
			staff.ID = mockID
			updatedStaff.ID = mockID
			updatedStaff.DisplayName = "Updated Name"
			updatedStaff.ImagePath = "/path/to/image"
			requestTime := testdata.RequestTime
			updatedStaff.UpdatedAt = requestTime

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
				Return(staff, nil)
			mockStaffRepo.EXPECT().
				Update(gomock.Any(), updatedStaff).
				Return(nil)
			mockStaffRepo.EXPECT().
				Get(gomock.Any(),
					repository.GetStaffQuery{
						ID: null.StringFrom(staff.ID),
						BaseGetOptions: repository.BaseGetOptions{
							OrFail:  true,
							Preload: true,
						},
					}).
				Return(updatedStaff, nil)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				GetWithValidate(gomock.Any(), model.AssetTypeUserImage, "image-asset-id").
				Return("/path/to/image", nil)
			mockAssetService.EXPECT().
				BatchSetStaffURLs(gomock.Any(), model.Staffs{updatedStaff}).
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
					staff: updatedStaff,
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
