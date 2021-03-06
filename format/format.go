package format

import (
	"errors"
	"github.com/die-net/fotomat/vips"
	"net/http"
)

var (
	ErrInvalidOperation = errors.New("Invalid operation")
	ErrUnknownFormat    = errors.New("Unknown image format")
)

type Format int

const (
	Unknown Format = iota
	Jpeg
	Png
	Gif
	Webp
)

var formatInfo = []struct {
	mime      string
	loadFile  func(filename string) (*vips.Image, error)
	loadBytes func([]byte) (*vips.Image, error)
}{
	{mime: "application/octet-stream", loadFile: nil, loadBytes: nil},
	{mime: "image/jpeg", loadFile: vips.Jpegload, loadBytes: vips.JpegloadBuffer},
	{mime: "image/png", loadFile: vips.Pngload, loadBytes: vips.PngloadBuffer},
	{mime: "image/gif", loadFile: vips.Gifload, loadBytes: vips.GifloadBuffer},
	{mime: "image/webp", loadFile: vips.Webpload, loadBytes: vips.WebploadBuffer},
}

func DetectFormat(blob []byte) Format {
	mime := http.DetectContentType(blob)

	for format, info := range formatInfo {
		if info.mime == mime {
			return Format(format)
		}
	}

	return Unknown
}

func (format Format) String() string {
	return formatInfo[format].mime
}

func (format Format) CanLoadFile() bool {
	return formatInfo[format].loadFile != nil
}

func (format Format) CanLoadBytes() bool {
	return formatInfo[format].loadBytes != nil
}

func (format Format) LoadFile(filename string) (*vips.Image, error) {
	loadFile := formatInfo[format].loadFile
	if loadFile == nil {
		return nil, ErrInvalidOperation
	}

	return loadFile(filename)
}

func (format Format) LoadBytes(blob []byte) (*vips.Image, error) {
	loadBytes := formatInfo[format].loadBytes
	if loadBytes == nil {
		return nil, ErrInvalidOperation
	}

	return loadBytes(blob)
}
