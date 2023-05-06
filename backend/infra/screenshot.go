package infra

import (
	"changeme/backend/apperr"
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/morikuni/failure"
)

type Screenshot struct{}

func (s *Screenshot) Save(path string, base64Data string) error {
	errCode := apperr.ScreenshotSave

	dir := filepath.Dir(path)
	_ = os.Mkdir(dir, 0o755)

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return failure.Translate(err, errCode)
	}

	f, err := os.Create(path)
	if err != nil {
		return failure.Translate(err, errCode)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return failure.Translate(err, errCode)
	}

	return nil
}
