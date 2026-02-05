package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	staff_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/staff_api/v1"
)

func StaffContentTypeToModel(contentType staff_apiv1.ContentType) model.ContentType {
	switch contentType {
	case staff_apiv1.ContentType_CONTENT_TYPE_IMAGE_PNG:
		return model.ContentTypeImagePNG
	case staff_apiv1.ContentType_CONTENT_TYPE_IMAGE_JPEG:
		return model.ContentTypeImageJPEG
	case staff_apiv1.ContentType_CONTENT_TYPE_APPLICATION_ZIP:
		return model.ContentTypeApplicationZIP
	case staff_apiv1.ContentType_CONTENT_TYPE_APPLICATION_PDF:
		return model.ContentTypeApplicationPDF
	case staff_apiv1.ContentType_CONTENT_TYPE_TEXT_CSV:
		return model.ContentTypeTextCSV
	case staff_apiv1.ContentType_CONTENT_TYPE_UNSPECIFIED:
		fallthrough
	default:
		return model.ContentTypeUnknown
	}
}
