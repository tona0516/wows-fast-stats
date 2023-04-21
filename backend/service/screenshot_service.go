package service

import (
	"changeme/backend/repo"
)

type ScreenshotService struct {}

func (s *ScreenshotService) Save(filename string, base64Data string) error {
    screenshot := repo.Screenshot{FileName: filename}
    return screenshot.Save(base64Data)
}
