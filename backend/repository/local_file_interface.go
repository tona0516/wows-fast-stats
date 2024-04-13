package repository

import "wfs/backend/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type LocalFileInterface interface {
	SaveScreenshot(path string, base64Data string) error
	TempArenaInfo(installPath string) (data.TempArenaInfo, error)
	SaveTempArenaInfo(tempArenaInfo data.TempArenaInfo) error
}
