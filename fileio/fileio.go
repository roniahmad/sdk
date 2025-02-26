package fileio

import (
	"encoding/base64"
	"image"
	"os"
	"path/filepath"
	"slices"

	"github.com/gabriel-vasile/mimetype"
	"github.com/nfnt/resize"
	"github.com/roniahmad/sdk"
)

func CreateDirectory(path string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exePath := filepath.Dir(exe)
	dir := filepath.Join(exePath, path)

	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0751); err != nil {
			return "", sdk.ErrCreateDirectory
		}
	}

	return dir, nil
}

func IsDocumentAllowed(path string, allowedDocs []string) (string, error) {
	contentType, err := mimetype.DetectFile(path)
	if err != nil {
		return "", err
	}

	if !slices.Contains(allowedDocs, contentType.String()) {
		return "", sdk.ErrUnsupportedDocument
	}

	return contentType.String(), nil
}

func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Thumbnail(img image.Image) image.Image {
	resizedImage := resize.Resize(72, 72, img, resize.Lanczos2)
	return resizedImage
}
