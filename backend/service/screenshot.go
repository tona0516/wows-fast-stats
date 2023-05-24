package service

import (
	"changeme/backend/apperr"
	"changeme/backend/infra"
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Screenshot struct {
	screenshotRepo     infra.ScreenshotInterface
	saveFileDialogFunc func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error)
}

func NewScreenshot(
	screenshotRepo infra.ScreenshotInterface,
	saveFileDialogFunc func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error),
) *Screenshot {
	return &Screenshot{
		screenshotRepo:     screenshotRepo,
		saveFileDialogFunc: saveFileDialogFunc,
	}
}

func (s *Screenshot) SaveForAuto(filename string, base64Data string) error {
	return s.screenshotRepo.Save(filepath.Join("screenshot", filename), base64Data)
}

func (s *Screenshot) SaveWithDialog(ctx context.Context, filename string, base64Data string) error {
	path, err := s.saveFileDialogFunc(ctx, runtime.SaveDialogOptions{
		DefaultFilename: filename,
	})
	if err != nil {
		return errors.WithStack(apperr.SrvSs.SaveDialog.WithRaw(err))
	}
	if path == "" {
		return errors.WithStack(apperr.SrvSs.Canceled)
	}

	return s.screenshotRepo.Save(path, base64Data)
}
