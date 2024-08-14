package service

import (
	"context"
	"testing"
	"time"

	mock_cache "github.com/abyssparanoia/rapid-go/internal/domain/cache/mock"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	mock_repository "github.com/abyssparanoia/rapid-go/internal/domain/repository/mock"
	"github.com/abyssparanoia/rapid-go/internal/pkg/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAssetService_CreatePresignedURL(t *testing.T) {
	t.Parallel()

	mockUUID := uuid.MockUUIDBase64()
	asset := &model.Asset{
		Key:         mockUUID,
		ContentType: "image/png",
		AssetType:   model.AssetTypeUserImage,
		Path:        "private/user_images/" + mockUUID + ".png",
	}
	presignedURL := "presignedURL"

	type args struct {
		assetType   model.AssetType
		contentType string
	}

	type want struct {
		got            *AssetCreatePresignedURLResult
		expectedResult error
	}

	tests := map[string]struct {
		args    args
		service func(ctx context.Context, ctrl *gomock.Controller) Asset
		want    want
	}{
		"success": {
			args: args{
				assetType:   asset.AssetType,
				contentType: asset.ContentType,
			},
			service: func(_ context.Context, ctrl *gomock.Controller) Asset {
				mockAssetRepo := mock_repository.NewMockAsset(ctrl)
				mockAssetPathCache := mock_cache.NewMockAssetPath(ctrl)
				mockAssetRepo.EXPECT().
					GenerateWritePresignedURL(
						gomock.Any(),
						asset.ContentType,
						asset.Path,
						15*time.Minute,
					).
					Return(presignedURL, nil)
				mockAssetPathCache.EXPECT().
					Set(
						gomock.Any(),
						asset,
						24*time.Hour,
					).
					Return(nil)
				return &assetService{
					assetRepository: mockAssetRepo,
					assetPathCache:  mockAssetPathCache,
				}
			},
			want: want{
				got: &AssetCreatePresignedURLResult{
					AssetKey:     asset.Key,
					PresignedURL: presignedURL,
				},
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

			u := tc.service(ctx, ctrl)

			got, err := u.CreatePresignedURL(ctx, tc.args.assetType, tc.args.contentType)
			if tc.want.expectedResult == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want.got, got)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}
