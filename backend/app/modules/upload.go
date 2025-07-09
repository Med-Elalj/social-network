package modules

import (
	"image"
	"io"
	"strings"
)

func IsValidImage(file io.ReadSeeker, contentType string) bool {
	var err error
	if strings.Contains(contentType, "svg") {
		err = validateSecureSVG(file)
	} else {
		_, _, err = image.Decode(file)
	}
	_, _ = file.Seek(0, 0)
	return err == nil
}

func validateSecureSVG(file io.ReadSeeker) error {
	_, err := io.Copy(io.Discard, file)
	return err
}
