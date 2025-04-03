package service

import (
	"context"
	"testing"
	"time"

	mock_cache "github.com/abyssparanoia/rapid-go/internal/domain/cache/mock"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/model/factory"
	mock_repository "github.com/abyssparanoia/rapid-go/internal/domain/repository/mock"
	"github.com/abyssparanoia/rapid-go/internal/pkg/id"
	"github.com/abyssparanoia/rapid-go/internal/pkg/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAssetService_CreatePresignedURL(t *testing.T) {
	t.Parallel()

	type args struct {
		assetType   model.AssetType
		contentType model.ContentType
		requestTime time.Time
	}

	type want struct {
		got            *AssetCreatePresignedURLResult
		expectedResult error
	}

	type testcase struct {
		args    args
		service Asset
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			mockID := id.Mock()
			uuid.MockUUIDBase64()
			testdata := factory.NewFactory()
			asset := testdata.Asset
			asset.ID = mockID
			presignedURL := "presignedURL"
			requestTime := testdata.RequestTime

			mockAssetRepo := mock_repository.NewMockAsset(ctrl)
			mockAssetPathCache := mock_cache.NewMockAssetPath(ctrl)
			mockAssetRepo.EXPECT().
				GenerateWritePresignedURL(
					gomock.Any(),
					asset.ContentType,
					asset.Path,
					asset.Expiration(),
				).
				Return(presignedURL, nil)
			mockAssetPathCache.EXPECT().
				Set(
					gomock.Any(),
					asset,
				).
				Return(nil)

			return testcase{
				args: args{
					assetType:   asset.Type,
					contentType: asset.ContentType,
					requestTime: requestTime,
				},
				service: &assetService{
					assetRepository: mockAssetRepo,
					assetPathCache:  mockAssetPathCache,
				},
				want: want{
					got: &AssetCreatePresignedURLResult{
						AssetID:      asset.ID,
						PresignedURL: presignedURL,
					},
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

			got, err := tc.service.CreatePresignedURL(ctx, tc.args.assetType, tc.args.contentType, tc.args.requestTime)
			if tc.want.expectedResult == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want.got, got)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}
