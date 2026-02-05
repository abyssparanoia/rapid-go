package usecase

import (
	"context"
	"testing"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	mock_repository "github.com/abyssparanoia/rapid-go/internal/domain/repository/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDebugInteractor_CreateStaffIDToken(t *testing.T) {
	t.Parallel()

	type args struct {
		authUID  string
		password string
	}

	type want struct {
		token          string
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase DebugInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"repository error": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			authUID := "test-auth-uid"
			password := "password123"

			mockStaffAuthRepo := mock_repository.NewMockStaffAuthentication(ctrl)
			mockStaffAuthRepo.EXPECT().
				CreateIDToken(gomock.Any(), authUID, password).
				Return("", errors.InvalidIDTokenErr.New())

			return testcase{
				args: args{
					authUID:  authUID,
					password: password,
				},
				usecase: &debugInteractor{
					staffAuthenticationRepository: mockStaffAuthRepo,
				},
				want: want{
					expectedResult: errors.InvalidIDTokenErr,
				},
			}
		},
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			authUID := "test-auth-uid"
			password := "password123"
			token := "test-id-token"

			mockStaffAuthRepo := mock_repository.NewMockStaffAuthentication(ctrl)
			mockStaffAuthRepo.EXPECT().
				CreateIDToken(gomock.Any(), authUID, password).
				Return(token, nil)

			return testcase{
				args: args{
					authUID:  authUID,
					password: password,
				},
				usecase: &debugInteractor{
					staffAuthenticationRepository: mockStaffAuthRepo,
				},
				want: want{
					token: token,
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

			got, err := tc.usecase.CreateStaffIDToken(ctx, tc.args.authUID, tc.args.password)

			if tc.want.expectedResult == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want.token, got)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}
