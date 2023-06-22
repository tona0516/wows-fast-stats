package infra

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"wfs/backend/apperr"
)

type Screenshot struct{}

func NewScreenshot() *Screenshot {
	return &Screenshot{}
}

func (s *Screenshot) Save(path string, base64Data string) error {
	dir := filepath.Dir(path)
	_ = os.Mkdir(dir, 0o755)

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return apperr.New(apperr.DecodeBase64, err)
	}

	f, err := os.Create(path)
	if err != nil {
		return apperr.New(apperr.WriteFile, err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return apperr.New(apperr.WriteFile, err)
	}

	return nil
}
