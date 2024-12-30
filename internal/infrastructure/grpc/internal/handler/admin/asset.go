package admin

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/admin/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *AdminHandler) CreateAssetPresignedURL(ctx context.Context, req *admin_apiv1.CreateAssetPresignedURLRequest) (*admin_apiv1.CreateAssetPresignedURLResponse, error) {
	got, err := h.assetInteractor.CreatePresignedURL(
		ctx,
		input.NewAdminCreateAssetPresignedURL(
			marshaller.AdminContentTypeToModel(req.GetContentType()),
			marshaller.AdminAssetTypeToModel(req.GetAssetType()),
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}
	return &admin_apiv1.CreateAssetPresignedURLResponse{
		AssetId:      got.AssetID,
		PresignedUrl: got.PresignedURL,
	}, nil
}
