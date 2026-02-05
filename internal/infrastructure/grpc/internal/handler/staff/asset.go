package staff

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/staff/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
	staff_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/staff_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *StaffHandler) CreateAssetPresignedURL(ctx context.Context, req *staff_apiv1.CreateAssetPresignedURLRequest) (*staff_apiv1.CreateAssetPresignedURLResponse, error) {
	got, err := h.assetInteractor.CreatePresignedURL(
		ctx,
		input.NewStaffCreateAssetPresignedURL(
			marshaller.StaffContentTypeToModel(req.GetContentType()),
			marshaller.StaffAssetTypeToModel(req.GetAssetType()),
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}
	return &staff_apiv1.CreateAssetPresignedURLResponse{
		AssetId:      got.AssetID,
		PresignedUrl: got.PresignedURL,
	}, nil
}
