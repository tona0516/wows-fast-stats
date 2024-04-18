package repository

import "wfs/backend/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type StorageInterface interface {
	DataVersion() (uint, error)
	WriteDataVersion(version uint) error
	UserConfig() (data.UserConfig, error)
	WriteUserConfig(config data.UserConfig) error
	IsExistUserConfig() bool
	UserConfigV2() (data.UserConfigV2, error)
	WriteUserConfigV2(config data.UserConfigV2) error
	IsExistAlertPlayers() bool
	AlertPlayers() ([]data.AlertPlayer, error)
	WriteAlertPlayers(players []data.AlertPlayer) error
	ExpectedStats() (data.ExpectedStats, error)
	WriteExpectedStats(expectedStats data.ExpectedStats) error
	OwnIGN() (string, error)
	WriteOwnIGN(ign string) error
}
