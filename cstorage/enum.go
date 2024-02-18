package cstorage

type MimeType string

//goland:noinspection GoUnusedConst
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
	MimeTypeMp4  MimeType = "video/mp4"
	MimeTypeMp3  MimeType = "audio/mp3"
	MimeTypeSvg  MimeType = "image/svg+xml"
	MimeTypeWasm MimeType = "application/wasm"
	MimeTypeWebp MimeType = "image/webp"
	MimeTypeXml  MimeType = "text/xml; charset=utf-8"
)

func (f MimeType) String() string {
	return string(f)
}
