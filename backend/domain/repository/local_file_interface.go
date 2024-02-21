package repository

import "wfs/backend/domain/model"

//go:generate mockgen -source=$GOFILE -destination ../mock_$GOPACKAGE/$GOFILE -package mock_$GOPACKAGE
type LocalFileInterface interface {
	SaveScreenshot(path string, base64Data string) error
	TempArenaInfo(installPath string) (model.TempArenaInfo, error)
	SaveTempArenaInfo(tempArenaInfo model.TempArenaInfo) error
}
