package repository

import (
	"wfs/backend/domain"
)

type LocalFileInterface interface {
	SaveScreenshot(path string, base64Data string) error
	TempArenaInfo(installPath string) (domain.TempArenaInfo, error)
	SaveTempArenaInfo(tempArenaInfo domain.TempArenaInfo) error
}
