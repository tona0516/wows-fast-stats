package repository

import "wfs/backend/domain/model"

type LocalFileInterface interface {
	SaveScreenshot(path string, base64Data string) error
	TempArenaInfo(installPath string) (model.TempArenaInfo, error)
	SaveTempArenaInfo(tempArenaInfo model.TempArenaInfo) error
}
