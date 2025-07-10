package modules

import (
	"image"
	"io"
	"net/http"
)

func IsValidImage(file io.ReadSeeker) bool {
	//read first 512 bytes without consuming the stream
	header := make([]byte, 512)
	n, _ := file.Read(header)

	mediaType := http.DetectContentType(header[:n])

	file.Seek(0, io.SeekStart)

	switch mediaType {
	case "image/jpeg", "image/png", "image/gif":
		if _, _, err := image.Decode(file); err != nil {
			file.Seek(0, io.SeekStart)
			return false
		}
		file.Seek(0, io.SeekStart)
		return true
	default:
		return false
	}
}
