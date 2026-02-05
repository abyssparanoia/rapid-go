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

func TestStaffMeTenantInteractor_Get(t *testing.T) {
	t.Parallel()

	type args struct {
		tenantID    string
		staffID     string
		requestTime time.Time
	}

	type want struct {
		tenant         *model.Tenant
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase StaffMeTenantInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			mockTransactable := mock_repository.NewMockTransactable(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)

			return testcase{
				args: args{
					// Empty args to trigger validation error
				},
				usecase: &staffMeTenantInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					assetService:     mockAssetService,
				},
				want: want{
					expectedResult: errors.RequestInvalidArgumentErr,
				},
			}
		},
		"not found": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			tenant := testdata.Tenant
			mockID := id.Mock()
			tenant.ID = mockID
			requestTime := time.Now()

			mockTransactable := mock_repository.NewMockTransactable(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockTenantRepo.EXPECT().
				Get(gomock.Any(),
					repository.GetTenantQuery{
						ID: null.StringFrom(tenant.ID),
						BaseGetOptions: repository.BaseGetOptions{
							OrFail:  true,
							Preload: true,
						},
					}).
				Return(nil, errors.TenantNotFoundErr)
			mockAssetService := mock_service.NewMockAsset(ctrl)

			return testcase{
				args: args{
					tenantID:    tenant.ID,
					staffID:     "staff-id",
					requestTime: requestTime,
				},
				usecase: &staffMeTenantInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					assetService:     mockAssetService,
				},
				want: want{
					expectedResult: errors.TenantNotFoundErr,
				},
			}
		},
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			tenant := testdata.Tenant
			mockID := id.Mock()
			tenant.ID = mockID
			requestTime := time.Now()

			mockTransactable := mock_repository.NewMockTransactable(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockTenantRepo.EXPECT().
				Get(gomock.Any(),
					repository.GetTenantQuery{
						ID: null.StringFrom(tenant.ID),
						BaseGetOptions: repository.BaseGetOptions{
							OrFail:  true,
							Preload: true,
						},
					}).
				Return(tenant, nil)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				BatchSetTenantURLs(gomock.Any(), gomock.Any()).
				Return(nil)

			return testcase{
				args: args{
					tenantID:    tenant.ID,
					staffID:     "staff-id",
					requestTime: requestTime,
				},
				usecase: &staffMeTenantInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					assetService:     mockAssetService,
				},
				want: want{
					tenant: tenant,
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

			got, err := tc.usecase.Get(ctx, input.NewStaffGetMeTenant(
				tc.args.tenantID,
				tc.args.staffID,
				tc.args.requestTime,
			))

			if tc.want.expectedResult == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want.tenant, got)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}

func TestStaffMeTenantInteractor_Update(t *testing.T) {
	t.Parallel()

	type args struct {
		tenantID    string
		staffID     string
		name        null.String
		requestTime time.Time
	}

	type want struct {
		tenant         *model.Tenant
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase StaffMeTenantInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			mockTransactable := mock_repository.NewMockTransactable(ctrl)
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockAssetService := mock_service.NewMockAsset(ctrl)

			return testcase{
				args: args{
					// Empty args to trigger validation error
				},
				usecase: &staffMeTenantInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					assetService:     mockAssetService,
				},
				want: want{
					expectedResult: errors.RequestInvalidArgumentErr,
				},
			}
		},
		"not found": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			tenant := testdata.Tenant
			mockID := id.Mock()
			tenant.ID = mockID
			requestTime := time.Now()

			mockTransactable := mock_repository.TestMockTransactable()
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockTenantRepo.EXPECT().
				Get(gomock.Any(),
					repository.GetTenantQuery{
						ID: null.StringFrom(tenant.ID),
						BaseGetOptions: repository.BaseGetOptions{
							OrFail:    true,
							ForUpdate: true,
						},
					}).
				Return(nil, errors.TenantNotFoundErr)
			mockAssetService := mock_service.NewMockAsset(ctrl)

			return testcase{
				args: args{
					tenantID:    tenant.ID,
					staffID:     "staff-id",
					name:        null.StringFrom("Updated Name"),
					requestTime: requestTime,
				},
				usecase: &staffMeTenantInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					assetService:     mockAssetService,
				},
				want: want{
					expectedResult: errors.TenantNotFoundErr,
				},
			}
		},
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			tenant := testdata.Tenant
			mockID := id.Mock()
			tenant.ID = mockID
			requestTime := time.Now()

			mockTransactable := mock_repository.TestMockTransactable()
			mockTenantRepo := mock_repository.NewMockTenant(ctrl)
			mockTenantRepo.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Times(2).
				Return(tenant, nil)
			mockTenantRepo.EXPECT().
				Update(gomock.Any(), gomock.Any()).
				Return(nil)
			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				BatchSetTenantURLs(gomock.Any(), gomock.Any()).
				Return(nil)

			return testcase{
				args: args{
					tenantID:    tenant.ID,
					staffID:     "staff-id",
					name:        null.StringFrom("Updated Name"),
					requestTime: requestTime,
				},
				usecase: &staffMeTenantInteractor{
					transactable:     mockTransactable,
					tenantRepository: mockTenantRepo,
					assetService:     mockAssetService,
				},
				want: want{
					tenant: tenant,
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

			got, err := tc.usecase.Update(ctx, input.NewStaffUpdateMeTenant(
				tc.args.tenantID,
				tc.args.staffID,
				tc.args.name,
				tc.args.requestTime,
			))

			if tc.want.expectedResult == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want.tenant, got)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}
