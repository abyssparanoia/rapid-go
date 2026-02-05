package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	staff_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/staff_api/v1"
)

func StaffAssetTypeToModel(assetType staff_apiv1.AssetType) model.AssetType {
	switch assetType {
	case staff_apiv1.AssetType_ASSET_TYPE_USER_IMAGE:
		return model.AssetTypeUserImage
	case staff_apiv1.AssetType_ASSET_TYPE_UNSPECIFIED:
		fallthrough
	default:
		return model.AssetTypeUnknown
	}
}
