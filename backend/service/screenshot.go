package service

import (
	"changeme/backend/infra"
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Screenshot struct{
    screenshotRepo infra.Screenshot
}

func NewScreenshot(screenshotRepo infra.Screenshot) *Screenshot {
    return &Screenshot{
        screenshotRepo: screenshotRepo,
    }
}

func (s *Screenshot) SaveForAuto(filename string, base64Data string) error {
    return s.screenshotRepo.Save(filepath.Join("screenshot", filename), base64Data)
}

func (s *Screenshot) SaveWithDialog(ctx context.Context, filename string, base64Data string) error {
    path, err := runtime.SaveFileDialog(ctx, runtime.SaveDialogOptions{
        DefaultFilename: filename,
    })
    if err != nil {
        return errors.WithStack(err)
    }

    return s.screenshotRepo.Save(path, base64Data)
}
