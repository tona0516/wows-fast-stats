package repository

import "wfs/backend/domain"

type StorageInterface interface {
	ReadDataVersion() (uint, error)
	WriteDataVersion(version uint) error
	IsExistUserConfig() bool
	ReadUserConfig() (domain.UserConfig, error)
	WriteUserConfig(config domain.UserConfig) error
	IsExistAlertPlayers() bool
	ReadAlertPlayers() ([]domain.AlertPlayer, error)
	WriteAlertPlayers(players []domain.AlertPlayer) error
	ReadNSExpectedStats() (domain.NSExpectedStats, error)
	WriteNSExpectedStats(nsExpectedStats domain.NSExpectedStats) error
}
