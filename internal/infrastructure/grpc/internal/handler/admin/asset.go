package admin

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/admin/marshaller"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"

	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *AdminHandler) CreateAssetPresignedURL(ctx context.Context, req *admin_apiv1.CreateAssetPresignedURLRequest) (*admin_apiv1.CreateAssetPresignedURLResponse, error) {
	got, err := h.assetInteractor.CreatePresignedURL(
		ctx,
		input.NewAdminCreateAssetPresignedURL(
			req.GetContentType(),
			marshaller.AdminAssetTypeToModel(req.GetAssetType()),
		),
	)
	if err != nil {
		return nil, err
	}
	return &admin_apiv1.CreateAssetPresignedURLResponse{
		Path:         got.Path,
		PresignedUrl: got.PresignedURL,
	}, nil
}
