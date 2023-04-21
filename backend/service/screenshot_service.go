package service

import (
	"changeme/backend/repo"
	"context"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ScreenshotService struct{}

func (s *ScreenshotService) SaveForAuto(filename string, base64Data string) error {
    screenshot := repo.Screenshot{}
    return screenshot.Save(filepath.Join("screenshot", filename), base64Data)
}

func (s *ScreenshotService) SaveWithDialog(ctx context.Context, filename string, base64Data string) error {
    path, err := runtime.SaveFileDialog(ctx, runtime.SaveDialogOptions{
        DefaultFilename: filename,
    })
    if err != nil {
        return err
    }

    screenshot := repo.Screenshot{}
    return screenshot.Save(path, base64Data)
}
