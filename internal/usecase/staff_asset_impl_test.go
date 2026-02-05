package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/model/factory"
	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	mock_service "github.com/abyssparanoia/rapid-go/internal/domain/service/mock"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStaffAssetInteractor_CreatePresignedURL(t *testing.T) {
	t.Parallel()

	type args struct {
		contentType model.ContentType
		assetType   model.AssetType
		requestTime time.Time
	}

	type want struct {
		result         *output.StaffCreateAssetPresignedURL
		expectedResult error
	}

	type testcase struct {
		args    args
		usecase StaffAssetInteractor
		want    want
	}

	type testcaseFunc func(ctx context.Context, ctrl *gomock.Controller) testcase

	tests := map[string]testcaseFunc{
		"invalid argument": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			mockAssetService := mock_service.NewMockAsset(ctrl)

			return testcase{
				args: args{
					contentType: model.ContentTypeUnknown, // Invalid content type
					assetType:   model.AssetTypeUserImage,
					requestTime: testdata.RequestTime,
				},
				usecase: &staffAssetInteractor{
					assetService: mockAssetService,
				},
				want: want{
					expectedResult: errors.RequestInvalidArgumentErr,
				},
			}
		},
		"service error": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			contentType := model.ContentTypeImagePNG
			assetType := model.AssetTypeUserImage

			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				CreatePresignedURL(gomock.Any(), assetType, contentType, testdata.RequestTime).
				Return(nil, errors.InternalErr.New())

			return testcase{
				args: args{
					contentType: contentType,
					assetType:   assetType,
					requestTime: testdata.RequestTime,
				},
				usecase: &staffAssetInteractor{
					assetService: mockAssetService,
				},
				want: want{
					expectedResult: errors.InternalErr,
				},
			}
		},
		"success": func(ctx context.Context, ctrl *gomock.Controller) testcase {
			testdata := factory.NewFactory()
			contentType := model.ContentTypeImagePNG
			assetType := model.AssetTypeUserImage
			assetID := "test-asset-id"
			presignedURL := "https://example.com/presigned-url"

			serviceResult := &service.AssetCreatePresignedURLResult{
				AssetID:      assetID,
				PresignedURL: presignedURL,
			}

			mockAssetService := mock_service.NewMockAsset(ctrl)
			mockAssetService.EXPECT().
				CreatePresignedURL(gomock.Any(), assetType, contentType, testdata.RequestTime).
				Return(serviceResult, nil)

			return testcase{
				args: args{
					contentType: contentType,
					assetType:   assetType,
					requestTime: testdata.RequestTime,
				},
				usecase: &staffAssetInteractor{
					assetService: mockAssetService,
				},
				want: want{
					result: output.NewStaffCreateAssetPresignedURL(assetID, presignedURL),
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

			got, err := tc.usecase.CreatePresignedURL(ctx, input.NewStaffCreateAssetPresignedURL(
				tc.args.contentType,
				tc.args.assetType,
				tc.args.requestTime,
			))

			if tc.want.expectedResult == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want.result, got)
			} else {
				require.ErrorContains(t, err, tc.want.expectedResult.Error())
			}
		})
	}
}
