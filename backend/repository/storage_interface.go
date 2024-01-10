package repository

import "wfs/backend/domain"

type StorageInterface interface {
	DataVersion() (uint, error)
	WriteDataVersion(version uint) error
	IsExistUserConfig() bool
	UserConfig() (domain.UserConfig, error)
	WriteUserConfig(config domain.UserConfig) error
	IsExistAlertPlayers() bool
	AlertPlayers() ([]domain.AlertPlayer, error)
	WriteAlertPlayers(players []domain.AlertPlayer) error
	ExpectedStats() (domain.ExpectedStats, error)
	WriteExpectedStats(expectedStats domain.ExpectedStats) error
	OwnIGN() (string, error)
	WriteOwnIGN(ign string) error
}
