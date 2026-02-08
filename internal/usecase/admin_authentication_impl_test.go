package usecase

import (
	"context"
	"testing"

	"github.com/aarondl/null/v9"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	mock_repository "github.com/abyssparanoia/rapid-go/internal/domain/repository/mock"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAdminAuthenticationInteractor_VerifyAdminIDToken(t *testing.T) {
	t.Parallel()

	type args struct {
		idToken string
	}

	type want struct {
		claims         *model.AdminClaims
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase AdminAuthenticationInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			mockAdminAuthRepo := mock_repository.NewMockAdminAuthentication(ctrl)

			return testcase{
				args: args{
					idToken: "", // Empty token causes validation error
				},
				usecase: &adminAuthenticationInteractor{
					adminAuthenticationRepository: mockAdminAuthRepo,
				},
				want: want{
					expectedResult: errors.RequestInvalidArgumentErr,
				},
			}
		},
		"repository error": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			idToken := "test-id-token"

			mockAdminAuthRepo := mock_repository.NewMockAdminAuthentication(ctrl)
			mockAdminAuthRepo.EXPECT().
				VerifyIDToken(gomock.Any(), idToken).
				Return(nil, errors.InvalidIDTokenErr.New())

			return testcase{
				args: args{
					idToken: idToken,
				},
				usecase: &adminAuthenticationInteractor{
					adminAuthenticationRepository: mockAdminAuthRepo,
				},
				want: want{
					expectedResult: errors.InvalidIDTokenErr,
				},
			}
		},
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			idToken := "test-id-token"
			adminID := "test-admin-id"
			claims := model.NewAdminClaims(
				"test-auth-uid",
				"test@example.com",
				null.StringFrom(adminID),
				nullable.TypeFrom(model.AdminRoleRoot),
			)

			mockAdminAuthRepo := mock_repository.NewMockAdminAuthentication(ctrl)
			mockAdminAuthRepo.EXPECT().
				VerifyIDToken(gomock.Any(), idToken).
				Return(claims, nil)

			return testcase{
				args: args{
					idToken: idToken,
				},
				usecase: &adminAuthenticationInteractor{
					adminAuthenticationRepository: mockAdminAuthRepo,
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

			got, err := tc.usecase.VerifyAdminIDToken(ctx, input.NewVerifyIDToken(
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
