package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
)

func AdminContentTypeToModel(contentType admin_apiv1.ContentType) model.ContentType {
	switch contentType {
	case admin_apiv1.ContentType_CONTENT_TYPE_IMAGE_PNG:
		return model.ContentTypeImagePNG
	case admin_apiv1.ContentType_CONTENT_TYPE_IMAGE_JPEG:
		return model.ContentTypeImageJPEG
	case admin_apiv1.ContentType_CONTENT_TYPE_APPLICATION_ZIP:
		return model.ContentTypeApplicationZIP
	case admin_apiv1.ContentType_CONTENT_TYPE_APPLICATION_PDF:
		return model.ContentTypeApplicationPDF
	case admin_apiv1.ContentType_CONTENT_TYPE_TEXT_CSV:
		return model.ContentTypeTextCSV
	case admin_apiv1.ContentType_CONTENT_TYPE_UNSPECIFIED:
		fallthrough
	default:
		return model.ContentTypeUnknown
	}
}
