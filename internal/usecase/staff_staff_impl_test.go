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

func TestStaffStaffInteractor_Get(t *testing.T) {
	t.Parallel()

	type args struct {
		tenantID      string
		staffID       string
		targetStaffID string
		requestTime   time.Time
	}

	type want struct {
		staff          *model.Staff
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase StaffStaffInteractor
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
				args: args{},
				usecase: &staffStaffInteractor{
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

			mockTransactable := mock_repository.NewMockTransactable(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)

			return testcase{
				args: args{
					tenantID:      staff.TenantID,
					staffID:       staff.ID,
					targetStaffID: staff.ID,
					requestTime:   testdata.RequestTime,
				},
				usecase: &staffStaffInteractor{
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

			mockTransactable := mock_repository.NewMockTransactable(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				BatchSetStaffURLs(gomock.Any(), model.Staffs{staff}, gomock.Any()).
				Return(nil)

			return testcase{
				args: args{
					tenantID:      staff.TenantID,
					staffID:       staff.ID,
					targetStaffID: staff.ID,
					requestTime:   testdata.RequestTime,
				},
				usecase: &staffStaffInteractor{
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

			got, err := tc.usecase.Get(ctx, input.NewStaffGetStaff(
				tc.args.tenantID,
				tc.args.staffID,
				tc.args.targetStaffID,
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

func TestStaffStaffInteractor_List(t *testing.T) {
	t.Parallel()

	type args struct {
		tenantID    string
		staffID     string
		page        uint64
		limit       uint64
		sortKey     model.StaffSortKey
		requestTime time.Time
	}

	type want struct {
		output         *output.ListStaffs
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase StaffStaffInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			staff := testdata.Staff
			tenant := testdata.Tenant
			mockID := id.Mock()
			staff.ID = mockID
			tenant.ID = mockID
			staff.TenantID = mockID

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
					}).
				Return(uint64(60), nil)

			mockTransactable := mock_repository.NewMockTransactable(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockStaffService := mock_service.NewMockStaff(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				BatchSetStaffURLs(gomock.Any(), model.Staffs{staff}, gomock.Any()).
				Return(nil)

			return testcase{
				args: args{
					tenantID:    tenant.ID,
					staffID:     staff.ID,
					page:        2,
					limit:       30,
					sortKey:     model.StaffSortKeyCreatedAtDesc,
					requestTime: testdata.RequestTime,
				},
				usecase: &staffStaffInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					staffRepository:  mockStaffRepo,
					staffService:     mockStaffService,
					assetService:     mockAssetService,
				},
				want: want{
					output: output.NewStaffListStaffs(
						model.Staffs{staff},
						model.NewPagination(2, 30, 60),
					),
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

			got, err := tc.usecase.List(ctx, input.NewStaffListStaffs(
				tc.args.tenantID,
				tc.args.staffID,
				tc.args.page,
				tc.args.limit,
				nullable.TypeFrom(tc.args.sortKey),
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
