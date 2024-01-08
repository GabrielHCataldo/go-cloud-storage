package cstorage

type FileMimeType string
type storageSelected string

const (
	googleStorage storageSelected = "GOOGLE"
	awsStorage    storageSelected = "AWS"
)
const (
	FileMimeTypePdf  FileMimeType = "application/pdf"
	FileMimeTypeText FileMimeType = "text/plain"
	FileMimeTypeAvif FileMimeType = "image/avif"
	FileMimeTypeCss  FileMimeType = "text/css; charset=utf-8"
	FileMimeTypeGif  FileMimeType = "image/gif"
	FileMimeTypeHtml FileMimeType = "text/html; charset=utf-8"
	FileMimeTypeJpeg FileMimeType = "image/jpeg"
	FileMimeTypeJs   FileMimeType = "text/javascript; charset=utf-8"
	FileMimeTypeJson FileMimeType = "application/json"
	FileMimeTypePng  FileMimeType = "image/png"
	FileMimeTypeSvg  FileMimeType = "image/svg+xml"
	FileMimeTypeWasm FileMimeType = "application/wasm"
	FileMimeTypeWebp FileMimeType = "image/webp"
	FileMimeTypeXml  FileMimeType = "text/xml; charset=utf-8"
)

func (f FileMimeType) String() string {
	return string(f)
}

func (f FileMimeType) Extension() string {
	switch f {
	case FileMimeTypePdf:
		return ".pdf"
	case FileMimeTypeText:
		return ".txt"
	case FileMimeTypeAvif:
		return ".avif"
	case FileMimeTypeCss:
		return ".css"
	case FileMimeTypeGif:
		return ".gif"
	case FileMimeTypeHtml:
		return ".html"
	case FileMimeTypeJpeg:
		return ".jpeg"
	case FileMimeTypeJs:
		return ".js"
	case FileMimeTypeJson:
		return ".json"
	case FileMimeTypePng:
		return ".png"
	case FileMimeTypeWasm:
		return ".wasm"
	case FileMimeTypeWebp:
		return ".webp"
	case FileMimeTypeXml:
		return ".xml"
	case FileMimeTypeSvg:
		return ".svg"
	default:
		return ""
	}
}
