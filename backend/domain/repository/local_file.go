package repository

import "wfs/backend/domain/model"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type LocalFile interface {
	SaveScreenshot(path string, base64Data string) error
	SaveTempArenaInfo(path string, tempArenaInfo model.TempArenaInfo) error
	ReadTempArenaInfo(installPath string) (model.TempArenaInfo, error)
}
