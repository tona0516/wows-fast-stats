package repository

import "wfs/backend/domain/model"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type LocalFileInterface interface {
	SaveScreenshot(path string, base64Data string) error
	SaveTempArenaInfo(tempArenaInfo model.TempArenaInfo) error
}
