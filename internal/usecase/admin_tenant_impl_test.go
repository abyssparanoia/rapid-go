package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/model/factory"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	mock_repository "github.com/abyssparanoia/rapid-go/internal/domain/repository/mock"
	"github.com/abyssparanoia/rapid-go/internal/pkg/id"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"go.uber.org/mock/gomock"
)

func TestAdminAdminTenantInteractor_Get(t *testing.T) {
	t.Parallel()
	testdata := factory.NewFactory()
	tenant := testdata.Tenant

	type args struct {
		tenantID string
	}

	type want struct {
		tenant         *model.Tenant
		expectedResult error
	}

	tests := map[string]struct {
		args    args
		usecase func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor
		want    want
	}{
		"invalid argument": {
			args: args{},
			usecase: func(_ context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				return &adminTenantInteractor{
					tenantRepository: mockTenantRepo,
				}
			},
			want: want{
				expectedResult: errors.RequestInvalidArgumentErr,
			},
		},
		"not found": {
			args: args{
				tenantID: tenant.ID,
			},
			usecase: func(_ context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				mockTenantRepo.EXPECT().
					Get(
						gomock.Any(),
						repository.GetTenantQuery{
							ID: null.StringFrom(tenant.ID),
							BaseGetOptions: repository.BaseGetOptions{
								OrFail: true,
							},
						}).
					Return(nil, errors.TenantNotFoundErr)
				return &adminTenantInteractor{
					tenantRepository: mockTenantRepo,
				}
			},
			want: want{
				expectedResult: errors.TenantNotFoundErr,
			},
		},
		"success": {
			args: args{
				tenantID: tenant.ID,
			},
			usecase: func(_ context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
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

				return &adminTenantInteractor{
					tenantRepository: mockTenantRepo,
				}
			},
			want: want{
				tenant: tenant,
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

			u := tc.usecase(ctx, ctrl)

			got, err := u.Get(ctx, input.NewAdminGetTenant(
				tc.args.tenantID,
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

func TestAdminAdminTenantInteractor_List(t *testing.T) {
	t.Parallel()
	testdata := factory.NewFactory()
	tenant := testdata.Tenant

	type args struct {
		page  uint64
		limit uint64
	}

	type want struct {
		output         *output.ListTenants
		expectedResult error
	}

	tests := map[string]struct {
		args    args
		usecase func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor
		want    want
	}{
		"success": {
			args: args{
				page:  2,
				limit: 30,
			},
			usecase: func(_ context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				mockTenantRepo.EXPECT().
					List(
						gomock.Any(),
						repository.ListTenantsQuery{
							BaseListOptions: repository.BaseListOptions{
								Page:  null.Uint64From(2),
								Limit: null.Uint64From(30),
							},
						}).
					Return(model.Tenants{tenant}, nil)

				mockTenantRepo.EXPECT().
					Count(
						gomock.Any(),
						repository.ListTenantsQuery{
							BaseListOptions: repository.BaseListOptions{
								Page:  null.Uint64From(2),
								Limit: null.Uint64From(30),
							},
						},
					).
					Return(uint64(60), nil)

				return &adminTenantInteractor{
					tenantRepository: mockTenantRepo,
				}
			},
			want: want{
				output: output.NewAdminListTenants(
					model.Tenants{tenant},
					model.NewPagination(2, 30, 60),
				),
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

			u := tc.usecase(ctx, ctrl)

			got, err := u.List(ctx, input.NewAdminListTenants(
				tc.args.page,
				tc.args.limit,
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

func TestAdminAdminTenantInteractor_Create(t *testing.T) {
	t.Parallel()
	testdata := factory.NewFactory()
	requestTime := testdata.RequestTime
	tenant := testdata.Tenant
	mockID := id.Mock()
	tenant.ID = mockID

	type args struct {
		name        string
		requestTime time.Time
	}

	type want struct {
		tenant         *model.Tenant
		expectedResult error
	}

	tests := map[string]struct {
		args    args
		usecase func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor
		want    want
	}{
		"invalid argument": {
			args: args{},
			usecase: func(_ context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				return &adminTenantInteractor{
					tenantRepository: mockTenantRepo,
				}
			},
			want: want{
				expectedResult: errors.RequestInvalidArgumentErr,
			},
		},
		"success": {
			args: args{
				name:        tenant.Name,
				requestTime: requestTime,
			},
			usecase: func(_ context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				mockTenantRepo.EXPECT().
					Create(gomock.Any(), tenant).
					Return(nil)

				return &adminTenantInteractor{
					tenantRepository: mockTenantRepo,
				}
			},
			want: want{
				tenant: tenant,
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

			u := tc.usecase(ctx, ctrl)

			got, err := u.Create(ctx, input.NewAdminCreateTenant(
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

func TestAdminAdminTenantInteractor_Update(t *testing.T) {
	t.Parallel()
	testdata := factory.NewFactory()
	requestTime := testdata.RequestTime
	tenant := testdata.Tenant

	type args struct {
		tenantID    string
		name        null.String
		requestTime time.Time
	}

	type want struct {
		tenant         *model.Tenant
		expectedResult error
	}

	tests := map[string]struct {
		args    args
		usecase func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor
		want    want
	}{
		"invalid argument": {
			args: args{},
			usecase: func(_ context.Context, _ *gomock.Controller) AdminTenantInteractor {
				return &adminTenantInteractor{}
			},
			want: want{
				expectedResult: errors.RequestInvalidArgumentErr,
			},
		},
		"success": {
			args: args{
				tenantID:    tenant.ID,
				name:        null.StringFrom(tenant.Name),
				requestTime: requestTime,
			},
			usecase: func(_ context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
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
					Return(tenant, nil)
				mockTenantRepo.EXPECT().
					Update(gomock.Any(), tenant).
					Return(nil)

				return &adminTenantInteractor{
					transactable:     mock_repository.TestMockTransactable(),
					tenantRepository: mockTenantRepo,
				}
			},
			want: want{
				tenant: tenant,
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

			u := tc.usecase(ctx, ctrl)

			got, err := u.Update(ctx, input.NewAdminUpdateTenant(
				tc.args.tenantID,
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

func TestAdminAdminTenantInteractor_Delete(t *testing.T) {
	t.Parallel()
	testdata := factory.NewFactory()
	tenant := testdata.Tenant

	type args struct {
		tenantID string
	}

	type want struct {
		expectedResult error
	}

	tests := map[string]struct {
		args    args
		usecase func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor
		want    want
	}{
		"invalid argument": {
			args: args{},
			usecase: func(_ context.Context, _ *gomock.Controller) AdminTenantInteractor {
				return &adminTenantInteractor{}
			},
			want: want{
				expectedResult: errors.RequestInvalidArgumentErr,
			},
		},
		"success": {
			args: args{
				tenantID: tenant.ID,
			},
			usecase: func(_ context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				mockTenantRepo.EXPECT().
					Delete(gomock.Any(), tenant.ID).
					Return(nil)

				return &adminTenantInteractor{
					tenantRepository: mockTenantRepo,
				}
			},
			want: want{},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := tc.usecase(ctx, ctrl)

			err := u.Delete(ctx, input.NewAdminDeleteTenant(
				tc.args.tenantID,
			))
			if tc.want.expectedResult == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}
