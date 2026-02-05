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
	"github.com/abyssparanoia/rapid-go/internal/pkg/id"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTaskAdminInteractor_Create(t *testing.T) {
	t.Parallel()

	type args struct {
		email       string
		displayName string
		password    string
		requestTime time.Time
	}

	type want struct {
		result         *output.TaskCreateAdmin
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase TaskAdminInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			mockAdminRepo := mock_repository.NewMockAdmin(ctrl)
			mockAdminAuthRepo := mock_repository.NewMockAdminAuthentication(ctrl)

			return testcase{
				args: args{
					// Empty email will cause validation error
					email:       "",
					displayName: "Test Admin",
					password:    "password123",
					requestTime: time.Now(),
				},
				usecase: &taskAdminInteractor{
					transactable:                  mock_repository.TestMockTransactable(),
					adminRepository:               mockAdminRepo,
					adminAuthenticationRepository: mockAdminAuthRepo,
				},
				want: want{
					expectedResult: errors.RequestInvalidArgumentErr,
				},
			}
		},
		"cognito create user error": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			requestTime := testdata.RequestTime

			mockAdminAuthRepo := mock_repository.NewMockAdminAuthentication(ctrl)
			mockAdminAuthRepo.EXPECT().
				CreateUser(
					gomock.Any(),
					repository.AdminAuthenticationCreateUserParam{
						Email:    "test@example.com",
						Password: null.StringFrom("password123"),
					},
				).
				Return("", errors.InternalErr.New())

			mockAdminRepo := mock_repository.NewMockAdmin(ctrl)

			return testcase{
				args: args{
					email:       "test@example.com",
					displayName: "Test Admin",
					password:    "password123",
					requestTime: requestTime,
				},
				usecase: &taskAdminInteractor{
					transactable:                  mock_repository.TestMockTransactable(),
					adminRepository:               mockAdminRepo,
					adminAuthenticationRepository: mockAdminAuthRepo,
				},
				want: want{
					expectedResult: errors.InternalErr,
				},
			}
		},
		"admin repository create error": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			requestTime := testdata.RequestTime
			authUID := "test-auth-uid"
			email := "test@example.com"
			displayName := "Test Admin"
			mockID := id.Mock() // Mock ID generation

			mockAdminAuthRepo := mock_repository.NewMockAdminAuthentication(ctrl)
			mockAdminAuthRepo.EXPECT().
				CreateUser(
					gomock.Any(),
					repository.AdminAuthenticationCreateUserParam{
						Email:    email,
						Password: null.StringFrom("password123"),
					},
				).
				Return(authUID, nil)

			// Create expected admin object with mock ID
			admin := model.NewAdmin(
				model.AdminRoleRoot,
				authUID,
				email,
				displayName,
				requestTime,
			)
			admin.ID = mockID // Set to mock ID

			mockAdminRepo := mock_repository.NewMockAdmin(ctrl)
			mockAdminRepo.EXPECT().
				Create(gomock.Any(), admin).
				Return(errors.InternalErr.New())

			return testcase{
				args: args{
					email:       email,
					displayName: displayName,
					password:    "password123",
					requestTime: requestTime,
				},
				usecase: &taskAdminInteractor{
					transactable:                  mock_repository.TestMockTransactable(),
					adminRepository:               mockAdminRepo,
					adminAuthenticationRepository: mockAdminAuthRepo,
				},
				want: want{
					expectedResult: errors.InternalErr,
				},
			}
		},
		"cognito store claims error": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			requestTime := testdata.RequestTime
			authUID := "test-auth-uid"
			email := "test@example.com"
			displayName := "Test Admin"
			mockID := id.Mock() // Mock ID generation

			mockAdminAuthRepo := mock_repository.NewMockAdminAuthentication(ctrl)
			mockAdminAuthRepo.EXPECT().
				CreateUser(
					gomock.Any(),
					repository.AdminAuthenticationCreateUserParam{
						Email:    email,
						Password: null.StringFrom("password123"),
					},
				).
				Return(authUID, nil)

			// Create expected admin object with mock ID
			admin := model.NewAdmin(
				model.AdminRoleRoot,
				authUID,
				email,
				displayName,
				requestTime,
			)
			admin.ID = mockID // Set to mock ID

			// Create expected claims object
			claims := model.NewAdminClaims(
				authUID,
				email,
				null.StringFrom(admin.ID),
				nullable.TypeFrom(admin.Role),
			)

			mockAdminAuthRepo.EXPECT().
				StoreClaims(gomock.Any(), authUID, claims).
				Return(errors.InternalErr.New())

			mockAdminRepo := mock_repository.NewMockAdmin(ctrl)
			mockAdminRepo.EXPECT().
				Create(gomock.Any(), admin).
				Return(nil)

			return testcase{
				args: args{
					email:       email,
					displayName: displayName,
					password:    "password123",
					requestTime: requestTime,
				},
				usecase: &taskAdminInteractor{
					transactable:                  mock_repository.TestMockTransactable(),
					adminRepository:               mockAdminRepo,
					adminAuthenticationRepository: mockAdminAuthRepo,
				},
				want: want{
					expectedResult: errors.InternalErr,
				},
			}
		},
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			requestTime := testdata.RequestTime
			authUID := "test-auth-uid"
			email := "test@example.com"
			displayName := "Test Admin"
			password := "password123"
			mockID := id.Mock() // Mock ID generation

			mockAdminAuthRepo := mock_repository.NewMockAdminAuthentication(ctrl)
			mockAdminAuthRepo.EXPECT().
				CreateUser(
					gomock.Any(),
					repository.AdminAuthenticationCreateUserParam{
						Email:    email,
						Password: null.StringFrom(password),
					},
				).
				Return(authUID, nil)

			// Create expected admin object with mock ID
			admin := model.NewAdmin(
				model.AdminRoleRoot,
				authUID,
				email,
				displayName,
				requestTime,
			)
			admin.ID = mockID // Set to mock ID

			// Create expected claims object
			claims := model.NewAdminClaims(
				authUID,
				email,
				null.StringFrom(admin.ID),
				nullable.TypeFrom(admin.Role),
			)

			mockAdminAuthRepo.EXPECT().
				StoreClaims(gomock.Any(), authUID, claims).
				Return(nil)

			mockAdminRepo := mock_repository.NewMockAdmin(ctrl)
			mockAdminRepo.EXPECT().
				Create(gomock.Any(), admin).
				Return(nil)

			return testcase{
				args: args{
					email:       email,
					displayName: displayName,
					password:    password,
					requestTime: requestTime,
				},
				usecase: &taskAdminInteractor{
					transactable:                  mock_repository.TestMockTransactable(),
					adminRepository:               mockAdminRepo,
					adminAuthenticationRepository: mockAdminAuthRepo,
				},
				want: want{
					result: &output.TaskCreateAdmin{
						AdminID:  admin.ID,
						AuthUID:  authUID,
						Password: password,
					},
				},
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tc := tc(ctx, ctrl)

			got, err := tc.usecase.Create(ctx, input.NewTaskCreateAdmin(
				tc.args.email,
				tc.args.displayName,
				tc.args.password,
				tc.args.requestTime,
			))

			if tc.want.expectedResult == nil {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NotEmpty(t, got.AdminID)
				require.Equal(t, tc.want.result.AuthUID, got.AuthUID)
				require.Equal(t, tc.want.result.Password, got.Password)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}
