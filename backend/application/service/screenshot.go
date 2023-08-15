package service

import (
	"context"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/application/repository"

	"github.com/morikuni/failure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Screenshot struct {
	localFile          repository.LocalFileInterface
	SaveFileDialogFunc SaveFileDialogFunc
}

func NewScreenshot(
	localFile repository.LocalFileInterface) *Screenshot {
	return &Screenshot{
		localFile:          localFile,
		SaveFileDialogFunc: runtime.SaveFileDialog,
	}
}

func (s *Screenshot) SaveForAuto(filename string, base64Data string) error {
	err := s.localFile.SaveScreenshot(filepath.Join("screenshot", filename), base64Data)
	return failure.Wrap(err)
}

func (s *Screenshot) SaveWithDialog(ctx context.Context, filename string, base64Data string) error {
	path, err := s.SaveFileDialogFunc(ctx, runtime.SaveDialogOptions{
		DefaultFilename: filename,
	})

	if err != nil {
		return failure.New(apperr.WailsError, failure.Messagef("%s", err.Error()))
	}
	if path == "" {
		return failure.New(apperr.UserCanceled)
	}

	err = s.localFile.SaveScreenshot(path, base64Data)
	return failure.Wrap(err)
}
