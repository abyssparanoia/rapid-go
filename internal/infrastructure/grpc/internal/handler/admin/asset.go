package admin

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/marshaller"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"

	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *AdminHandler) AdminCreateAssetPresignedURL(ctx context.Context, req *admin_apiv1.AdminCreateAssetPresignedURLRequest) (*admin_apiv1.AdminCreateAssetPresignedURLResponse, error) {
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
	return &admin_apiv1.AdminCreateAssetPresignedURLResponse{
		Path:         got.Path,
		PresignedUrl: got.PresignedURL,
	}, nil
}
