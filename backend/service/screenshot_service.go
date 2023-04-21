package service

import (
	"changeme/backend/repo"
	"strconv"
	"time"
)

type ScreenshotService struct {}

func (s *ScreenshotService) Save(base64Data string) error {
    unixtime := time.Now().Unix()
    // TODO clearify filename
    screenshot := repo.Screenshot{FileName: strconv.FormatInt(unixtime, 10) + ".png"}
    return screenshot.Save(base64Data)
}
