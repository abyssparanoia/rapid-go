package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
)

func AdminAssetTypeToModel(assetType admin_apiv1.AssetType) model.AssetType {
	switch assetType {
	case admin_apiv1.AssetType_ASSET_TYPE_USER_IMAGE:
		return model.AssetTypeUserImage
	case admin_apiv1.AssetType_ASSET_TYPE_UNSPECIFIED:
		fallthrough
	default:
		return model.AssetTypeUnknown
	}
}
