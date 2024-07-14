package repository

import "wfs/backend/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type LocalFileInterface interface {
	SaveScreenshot(path string, base64Data string) error
	TempArenaInfo(installPath string) (data.TempArenaInfo, error)
	SaveTempArenaInfo(tempArenaInfo data.TempArenaInfo) error
	IGN() (string, error)
	WriteIGN(ign string) error
	ExpectedStats() (data.ExpectedStats, error)
	WriteExpectedStats(target data.ExpectedStats) error
	UserConfigV2() (data.UserConfigV2, error)
	WriteUserConfigV2(target data.UserConfigV2) error
	AlertPlayerV1() ([]data.AlertPlayer, error)
	WriteAlertPlayerV1(target []data.AlertPlayer) error
}
