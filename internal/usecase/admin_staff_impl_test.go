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
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAdminStaffInteractor_Get(t *testing.T) {
	t.Parallel()

	type args struct {
		staffID     string
		requestTime time.Time
	}

	type want struct {
		staff          *model.Staff
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase AdminStaffInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			mockStaffRepo := mock_repository.NewMockStaff(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			return testcase{
				args: args{},
				usecase: &adminStaffInteractor{
					staffRepository:  mockStaffRepo,
					tenantRepository: mockTenantRepo,
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
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)

			return testcase{
				args: args{
					staffID:     staff.ID,
					requestTime: testdata.RequestTime,
				},
				usecase: &adminStaffInteractor{
					staffRepository:  mockStaffRepo,
					tenantRepository: mockTenantRepo,
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
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				BatchSetStaffURLs(gomock.Any(), model.Staffs{staff}).
				Return(nil)

			return testcase{
				args: args{
					staffID:     staff.ID,
					requestTime: testdata.RequestTime,
				},
				usecase: &adminStaffInteractor{
					staffRepository:  mockStaffRepo,
					tenantRepository: mockTenantRepo,
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
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tc := tc(ctx, ctrl)

			got, err := tc.usecase.Get(ctx, input.NewAdminGetStaff(
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

func TestAdminStaffInteractor_List(t *testing.T) {
	t.Parallel()

	type args struct {
		tenantID    string
		page        uint64
		limit       uint64
		requestTime time.Time
	}

	type want struct {
		output         *output.ListStaffs
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase AdminStaffInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			staff := testdata.Staff
			tenant := testdata.Tenant

			mockStaffRepo := mock_repository.NewMockStaff(ctrl)
			mockStaffRepo.EXPECT().
				List(
					gomock.Any(),
					repository.ListStaffQuery{
						TenantID: null.StringFrom(tenant.ID),
						BaseListOptions: repository.BaseListOptions{
							Page:    null.Uint64From(2),
							Limit:   null.Uint64From(30),
							Preload: true,
						},
						SortKey: nullable.TypeFrom(model.StaffSortKeyCreatedAtDesc),
					}).
				Return(model.Staffs{staff}, nil)

			mockStaffRepo.EXPECT().
				Count(
					gomock.Any(),
					repository.ListStaffQuery{
						TenantID: null.StringFrom(tenant.ID),
						BaseListOptions: repository.BaseListOptions{
							Page:    null.Uint64From(2),
							Limit:   null.Uint64From(30),
							Preload: true,
						},
						SortKey: nullable.TypeFrom(model.StaffSortKeyCreatedAtDesc),
					},
				).
				Return(uint64(60), nil)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				BatchSetStaffURLs(gomock.Any(), model.Staffs{staff}).
				Return(nil)

			return testcase{
				args: args{
					tenantID:    tenant.ID,
					page:        2,
					limit:       30,
					requestTime: testdata.RequestTime,
				},
				usecase: &adminStaffInteractor{
					staffRepository:  mockStaffRepo,
					tenantRepository: mockTenantRepo,
					staffService:     mockStaffService,
					assetService:     mockAssetService,
				},
				want: want{
					output: output.NewAdminListStaffs(
						model.Staffs{staff},
						model.NewPagination(2, 30, 60),
					),
				},
			}
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tc := tc(ctx, ctrl)

			got, err := tc.usecase.List(ctx, input.NewAdminListStaffs(
				tc.args.tenantID,
				tc.args.page,
				tc.args.limit,
				nullable.Type[model.StaffSortKey]{}, // Use empty nullable for default
				tc.args.requestTime,
			))
			if tc.want.expectedResult == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want.output, got)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}

func TestAdminStaffInteractor_Create(t *testing.T) {
	t.Parallel()

	type args struct {
		tenantID     string
		email        string
		displayName  string
		role         model.StaffRole
		imageAssetID string
		requestTime  time.Time
	}

	type want struct {
		staff          *model.Staff
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase AdminStaffInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			mockStaffRepo := mock_repository.NewMockStaff(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			return testcase{
				args: args{},
				usecase: &adminStaffInteractor{
					staffRepository:  mockStaffRepo,
					tenantRepository: mockTenantRepo,
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
			requestTime := testdata.RequestTime
			staff := testdata.Staff
			tenant := testdata.Tenant
			mockID := id.Mock()
			staff.ID = mockID

			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockTenantRepo.EXPECT().
				Get(gomock.Any(),
					repository.GetTenantQuery{
						ID: null.StringFrom(tenant.ID),
						BaseGetOptions: repository.BaseGetOptions{
							OrFail: true,
						},
					}).
				Return(tenant, nil)

			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				GetWithValidate(gomock.Any(), model.AssetTypeUserImage, staff.ImagePath).
				Return(staff.ImagePath, nil)

			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockStaffService.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(staff, nil)

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

			mockAssetService.EXPECT().
				BatchSetStaffURLs(gomock.Any(), model.Staffs{staff}).
				Return(nil)

			return testcase{
				args: args{
					tenantID:     tenant.ID,
					email:        staff.Email,
					displayName:  staff.DisplayName,
					role:         staff.Role,
					imageAssetID: staff.ImagePath,
					requestTime:  requestTime,
				},
				usecase: &adminStaffInteractor{
					transactable:     mock_repository.TestMockTransactable(),
					staffRepository:  mockStaffRepo,
					tenantRepository: mockTenantRepo,
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
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tc := tc(ctx, ctrl)

			got, err := tc.usecase.Create(ctx, input.NewAdminCreateStaff(
				tc.args.tenantID,
				tc.args.email,
				tc.args.displayName,
				tc.args.role,
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

func TestAdminStaffInteractor_Update(t *testing.T) {
	t.Parallel()

	type args struct {
		staffID      string
		displayName  null.String
		role         nullable.Type[model.StaffRole]
		imageAssetID null.String
		requestTime  time.Time
	}

	type want struct {
		staff          *model.Staff
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase AdminStaffInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			return testcase{
				args:    args{},
				usecase: &adminStaffInteractor{},
				want: want{
					expectedResult: errors.RequestInvalidArgumentErr,
				},
			}
		},
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			requestTime := testdata.RequestTime
			staff := testdata.Staff

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
				Update(gomock.Any(), staff).
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
				Return(staff, nil)

			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				GetWithValidate(gomock.Any(), model.AssetTypeUserImage, staff.ImagePath).
				Return(staff.ImagePath, nil)
			mockAssetService.EXPECT().
				BatchSetStaffURLs(gomock.Any(), model.Staffs{staff}).
				Return(nil)

			return testcase{
				args: args{
					staffID:      staff.ID,
					displayName:  null.StringFrom(staff.DisplayName),
					role:         nullable.TypeFrom(staff.Role),
					imageAssetID: null.StringFrom(staff.ImagePath),
					requestTime:  requestTime,
				},
				usecase: &adminStaffInteractor{
					transactable:     mock_repository.TestMockTransactable(),
					staffRepository:  mockStaffRepo,
					tenantRepository: mockTenantRepo,
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
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tc := tc(ctx, ctrl)

			got, err := tc.usecase.Update(ctx, input.NewAdminUpdateStaff(
				tc.args.staffID,
				tc.args.displayName,
				tc.args.role,
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
