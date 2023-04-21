package service

import (
	"changeme/backend/repo"
	"context"
	"encoding/base64"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ScreenshotService struct {}

func (s *ScreenshotService) Save(filename string, base64Data string) error {
    screenshot := repo.Screenshot{FileName: filename}
    return screenshot.Save(base64Data)
}

func (s *ScreenshotService) SaveWithDialog(ctx context.Context, filename string, base64Data string) error {
    path, err := runtime.SaveFileDialog(ctx, runtime.SaveDialogOptions{
        DefaultFilename: filename,
    })
    if err != nil {
        return err
    }

    // TOOD integrate into repository
    data, err := base64.StdEncoding.DecodeString(base64Data)
    if err != nil {
        return err
    }

    f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)

    return err
}
