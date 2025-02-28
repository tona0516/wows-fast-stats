package service

import (
	"context"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/domain/repository"

	"github.com/morikuni/failure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Screenshot struct {
	localFile      repository.LocalFile
	SaveFileDialog saveFileDialogFunc
}

func NewScreenshot(
	localFile repository.LocalFile,
) *Screenshot {
	return &Screenshot{
		localFile:      localFile,
		SaveFileDialog: runtime.SaveFileDialog,
	}
}

func (s *Screenshot) SaveForAuto(filename string, base64Data string) error {
	return s.localFile.SaveScreenshot(filepath.Join("screenshot", filename), base64Data)
}

func (s *Screenshot) SaveWithDialog(ctx context.Context, filename string, base64Data string) (bool, error) {
	path, err := s.SaveFileDialog(ctx, runtime.SaveDialogOptions{
		DefaultFilename: filename,
	})
	if err != nil {
		return false, failure.Translate(err, apperr.WailsError)
	}
	if path == "" {
		return false, nil
	}

	if err := s.localFile.SaveScreenshot(path, base64Data); err != nil {
		return false, err
	}

	return true, nil
}
