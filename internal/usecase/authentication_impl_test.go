package usecase

import (
	"context"
	"testing"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/model/factory"
	mock_repository "github.com/abyssparanoia/rapid-go/internal/domain/repository/mock"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthenticationInteractor_VerifyStaffIDToken(t *testing.T) {
	t.Parallel()

	type args struct {
		idToken string
	}

	type want struct {
		claims         *model.StaffClaims
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase AuthenticationInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			mockStaffAuthRepo := mock_repository.NewMockStaffAuthentication(ctrl)

			return testcase{
				args: args{
					idToken: "", // Empty token causes validation error
				},
				usecase: &staffAuthenticationInteractor{
					staffAuthenticationRepository: mockStaffAuthRepo,
				},
				want: want{
					expectedResult: errors.RequestInvalidArgumentErr,
				},
			}
		},
		"repository error": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			idToken := "test-id-token"

			mockStaffAuthRepo := mock_repository.NewMockStaffAuthentication(ctrl)
			mockStaffAuthRepo.EXPECT().
				VerifyIDToken(gomock.Any(), idToken).
				Return(nil, errors.InvalidIDTokenErr.New())

			return testcase{
				args: args{
					idToken: idToken,
				},
				usecase: &staffAuthenticationInteractor{
					staffAuthenticationRepository: mockStaffAuthRepo,
				},
				want: want{
					expectedResult: errors.InvalidIDTokenErr,
				},
			}
		},
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			idToken := "test-id-token"
			claims := model.NewStaffClaims(
				"test-auth-uid",
				"test@example.com",
				null.StringFrom(testdata.Tenant.ID),
				null.StringFrom(testdata.Staff.ID),
				nullable.TypeFrom(model.StaffRoleAdmin),
			)

			mockStaffAuthRepo := mock_repository.NewMockStaffAuthentication(ctrl)
			mockStaffAuthRepo.EXPECT().
				VerifyIDToken(gomock.Any(), idToken).
				Return(claims, nil)

			return testcase{
				args: args{
					idToken: idToken,
				},
				usecase: &staffAuthenticationInteractor{
					staffAuthenticationRepository: mockStaffAuthRepo,
				},
				want: want{
					claims: claims,
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

			got, err := tc.usecase.VerifyStaffIDToken(ctx, input.NewVerifyIDToken(
				tc.args.idToken,
			))

			if tc.want.expectedResult == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want.claims, got)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}
