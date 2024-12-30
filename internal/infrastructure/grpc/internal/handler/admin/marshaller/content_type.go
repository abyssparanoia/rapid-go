package marshaller

import (
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
)

func AdminContentTypeToModel(contentType admin_apiv1.ContentType) string {
	switch contentType {
	case admin_apiv1.ContentType_CONTENT_TYPE_IMAGE_PNG:
		return "image/png"
	case admin_apiv1.ContentType_CONTENT_TYPE_IMAGE_JPEG:
		return "image/jpeg"
	case admin_apiv1.ContentType_CONTENT_TYPE_UNSPECIFIED:
		fallthrough
	default:
		return ""
	}
}
