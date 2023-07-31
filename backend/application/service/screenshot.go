package service

import (
	"context"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/application/repository"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Screenshot struct {
	localFile          repository.LocalFileInterface
	saveFileDialogFunc func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error)
}

func NewScreenshot(
	localFile repository.LocalFileInterface,
	saveFileDialogFunc func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error),
) *Screenshot {
	return &Screenshot{
		localFile:          localFile,
		saveFileDialogFunc: saveFileDialogFunc,
	}
}

func (s *Screenshot) SaveForAuto(filename string, base64Data string) error {
	return s.localFile.SaveScreenshot(filepath.Join("screenshot", filename), base64Data)
}

func (s *Screenshot) SaveWithDialog(ctx context.Context, filename string, base64Data string) error {
	path, err := s.saveFileDialogFunc(ctx, runtime.SaveDialogOptions{
		DefaultFilename: filename,
	})

	if err != nil {
		return apperr.New(apperr.ErrShowDialog, err)
	}
	if path == "" {
		return apperr.New(apperr.ErrUserCanceled, nil)
	}

	return s.localFile.SaveScreenshot(path, base64Data)
}
