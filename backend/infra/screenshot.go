package infra

import (
	"changeme/backend/apperr"
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
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
		return errors.WithStack(apperr.Ss.Save.WithRaw(err))
	}

	f, err := os.Create(path)
	if err != nil {
		return errors.WithStack(apperr.Ss.Save.WithRaw(err))
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return errors.WithStack(apperr.Ss.Save.WithRaw(err))
	}

	return nil
}
