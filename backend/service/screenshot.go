package service

import (
	"changeme/backend/infra"
	"context"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Screenshot struct{}

func (s *Screenshot) SaveForAuto(filename string, base64Data string) error {
    screenshot := infra.Screenshot{}
    return screenshot.Save(filepath.Join("screenshot", filename), base64Data)
}

func (s *Screenshot) SaveWithDialog(ctx context.Context, filename string, base64Data string) error {
    path, err := runtime.SaveFileDialog(ctx, runtime.SaveDialogOptions{
        DefaultFilename: filename,
    })
    if err != nil {
        return err
    }

    screenshot := infra.Screenshot{}
    return screenshot.Save(path, base64Data)
}
