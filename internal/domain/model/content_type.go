package model

type ContentType string

const (
	ContentTypeUnknown        ContentType = "unknown"
	ContentTypeImagePNG       ContentType = "image/png"
	ContentTypeImageJPEG      ContentType = "image/jpeg"
	ContentTypeApplicationZIP ContentType = "application/zip"
	ContentTypeApplicationPDF ContentType = "application/pdf"
	ContentTypeTextCSV        ContentType = "text/csv"
)

func NewContentType(contentType string) ContentType {
	switch contentType {
	case ContentTypeImagePNG.String(),
		ContentTypeImageJPEG.String(),
		ContentTypeApplicationZIP.String(),
		ContentTypeApplicationPDF.String(),
		ContentTypeTextCSV.String():
		return ContentType(contentType)
	default:
		return ContentTypeUnknown
	}
}

func (m ContentType) String() string {
	return string(m)
}

func (m ContentType) Valid() bool {
	return m != ContentTypeUnknown && m != ""
}

func (m ContentType) Extension() string {
	switch m {
	case ContentTypeImagePNG:
		return "png"
	case ContentTypeImageJPEG:
		return "jpeg"
	case ContentTypeApplicationZIP:
		return "zip"
	case ContentTypeApplicationPDF:
		return "pdf"
	case ContentTypeTextCSV:
		return "csv"
	case ContentTypeUnknown:
		fallthrough
	default:
		panic("invalid content type")
	}
}
