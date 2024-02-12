package repository

import "wfs/backend/domain/model"

type StorageInterface interface {
	DataVersion() (uint, error)
	WriteDataVersion(version uint) error
	UserConfigV1() (model.UserConfigV1, error)
	WriteUserConfigV1(config model.UserConfigV1) error
	IsExistUserConfig() bool
	UserConfig() (model.UserConfig, error)
	WriteUserConfig(config model.UserConfig) error
	IsExistAlertPlayers() bool
	AlertPlayers() ([]model.AlertPlayer, error)
	WriteAlertPlayers(players []model.AlertPlayer) error
	ExpectedStats() (model.ExpectedStats, error)
	WriteExpectedStats(expectedStats model.ExpectedStats) error
	OwnIGN() (string, error)
	WriteOwnIGN(ign string) error
}
