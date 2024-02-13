package repository

import "wfs/backend/domain/model"

type StorageInterface interface {
	DataVersion() (uint, error)
	WriteDataVersion(version uint) error
	UserConfig() (model.UserConfig, error)
	WriteUserConfig(config model.UserConfig) error
	IsExistUserConfig() bool
	UserConfigV2() (model.UserConfigV2, error)
	WriteUserConfigV2(config model.UserConfigV2) error
	IsExistAlertPlayers() bool
	AlertPlayers() ([]model.AlertPlayer, error)
	WriteAlertPlayers(players []model.AlertPlayer) error
	ExpectedStats() (model.ExpectedStats, error)
	WriteExpectedStats(expectedStats model.ExpectedStats) error
	OwnIGN() (string, error)
	WriteOwnIGN(ign string) error
}
