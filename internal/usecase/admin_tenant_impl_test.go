package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model/factory"
	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/repository"
	mock_repository "github.com/playground-live/moala-meet-and-greet-back/internal/domain/repository/mock"
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/errors"
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/ulid"
	"github.com/playground-live/moala-meet-and-greet-back/internal/usecase/input"
	"github.com/playground-live/moala-meet-and-greet-back/internal/usecase/output"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func TestAdminAdminTenantInteractor_Get(t *testing.T) {
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
			usecase: func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				return &adminTenantInteractor{
					tenantRepository: mockTenantRepo,
				}
			},
			want: want{
				expectedResult: errors.RequestInvalidArgumentErr.New(),
			},
		},
		"not found": {
			args: args{
				tenantID: tenant.ID,
			},
			usecase: func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				mockTenantRepo.EXPECT().
					Get(gomock.Any(), tenant.ID, true).
					Return(nil, errors.NotFoundErr.New())
				return &adminTenantInteractor{
					tenantRepository: mockTenantRepo,
				}
			},
			want: want{
				expectedResult: errors.NotFoundErr.New(),
			},
		},
		"success": {
			args: args{
				tenantID: tenant.ID,
			},
			usecase: func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				mockTenantRepo.EXPECT().
					Get(gomock.Any(), tenant.ID, true).
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
				assert.NoError(t, err)
				assert.Equal(t, tc.want.tenant, got)
			} else {
				assert.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}

func TestAdminAdminTenantInteractor_List(t *testing.T) {
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
			usecase: func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				mockTenantRepo.EXPECT().
					List(gomock.Any(), repository.ListTenantsQuery{
						Page:  null.Uint64From(2),
						Limit: null.Uint64From(30),
					}).
					Return([]*model.Tenant{tenant}, nil)

				mockTenantRepo.EXPECT().
					Count(gomock.Any(), repository.CountTenantsQuery{}).
					Return(uint64(60), nil)

				return &adminTenantInteractor{
					tenantRepository: mockTenantRepo,
				}
			},
			want: want{
				output: output.NewAdminListTenants(
					[]*model.Tenant{tenant},
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
				assert.NoError(t, err)
				assert.Equal(t, tc.want.output, got)
			} else {
				assert.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}

func TestAdminAdminTenantInteractor_Create(t *testing.T) {
	testdata := factory.NewFactory()
	tenant := testdata.Tenant
	mockULID := ulid.Mock()
	tenant.ID = mockULID

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
			usecase: func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				return &adminTenantInteractor{
					tenantRepository: mockTenantRepo,
				}
			},
			want: want{
				expectedResult: errors.RequestInvalidArgumentErr.New(),
			},
		},
		"success": {
			args: args{
				name:        tenant.Name,
				requestTime: tenant.CreatedAt,
			},
			usecase: func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				mockTenantRepo.EXPECT().
					Create(gomock.Any(), tenant).
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

			got, err := u.Create(ctx, input.NewAdminCreateTenant(
				tc.args.name,
				tc.args.requestTime,
			))
			if tc.want.expectedResult == nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.want.tenant, got)
			} else {
				assert.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}

func TestAdminAdminTenantInteractor_Update(t *testing.T) {
	testdata := factory.NewFactory()
	tenant := testdata.Tenant

	type args struct {
		tenantID    string
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
			usecase: func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				return &adminTenantInteractor{}
			},
			want: want{
				expectedResult: errors.RequestInvalidArgumentErr.New(),
			},
		},
		"success": {
			args: args{
				tenantID:    tenant.ID,
				name:        tenant.Name,
				requestTime: tenant.UpdatedAt,
			},
			usecase: func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				mockTenantRepo := mock_repository.NewMockTenant(ctrl)
				mockTenantRepo.EXPECT().
					Get(gomock.Any(), tenant.ID, true).
					Return(tenant, nil)
				mockTenantRepo.EXPECT().
					Update(gomock.Any(), tenant).
					Return(tenant, nil)

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
				assert.NoError(t, err)
				assert.Equal(t, tc.want.tenant, got)
			} else {
				assert.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}

func TestAdminAdminTenantInteractor_Delete(t *testing.T) {
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
			usecase: func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
				return &adminTenantInteractor{}
			},
			want: want{
				expectedResult: errors.RequestInvalidArgumentErr.New(),
			},
		},
		"success": {
			args: args{
				tenantID: tenant.ID,
			},
			usecase: func(ctx context.Context, ctrl *gomock.Controller) AdminTenantInteractor {
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
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}
