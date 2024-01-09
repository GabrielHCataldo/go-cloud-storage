package cstorage

type MimeType string
type storageSelected string

const (
	googleStorage storageSelected = "GOOGLE"
	awsStorage    storageSelected = "AWS"
)
const (
	MimeTypePdf  MimeType = "application/pdf"
	MimeTypeText MimeType = "text/plain"
	MimeTypeAvif MimeType = "image/avif"
	MimeTypeCss  MimeType = "text/css; charset=utf-8"
	MimeTypeGif  MimeType = "image/gif"
	MimeTypeHtml MimeType = "text/html; charset=utf-8"
	MimeTypeJpeg MimeType = "image/jpeg"
	MimeTypeJs   MimeType = "text/javascript; charset=utf-8"
	MimeTypeJson MimeType = "application/json"
	MimeTypePng  MimeType = "image/png"
	MimeTypeSvg  MimeType = "image/svg+xml"
	MimeTypeWasm MimeType = "application/wasm"
	MimeTypeWebp MimeType = "image/webp"
	MimeTypeXml  MimeType = "text/xml; charset=utf-8"
)

func (f MimeType) String() string {
	return string(f)
}

func (f MimeType) Extension() string {
	switch f {
	case MimeTypePdf:
		return ".pdf"
	case MimeTypeText:
		return ".txt"
	case MimeTypeAvif:
		return ".avif"
	case MimeTypeCss:
		return ".css"
	case MimeTypeGif:
		return ".gif"
	case MimeTypeHtml:
		return ".html"
	case MimeTypeJpeg:
		return ".jpeg"
	case MimeTypeJs:
		return ".js"
	case MimeTypeJson:
		return ".json"
	case MimeTypePng:
		return ".png"
	case MimeTypeWasm:
		return ".wasm"
	case MimeTypeWebp:
		return ".webp"
	case MimeTypeXml:
		return ".xml"
	case MimeTypeSvg:
		return ".svg"
	default:
		return ""
	}
}
