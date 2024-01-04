package adapter

import (
	"wfs/backend/domain"
)

type LocalFileInterface interface {
	User() (domain.UserConfig, error)
	AlertPlayers() ([]domain.AlertPlayer, error)
	SaveScreenshot(path string, base64Data string) error
	TempArenaInfo(installPath string) (domain.TempArenaInfo, error)
	SaveTempArenaInfo(tempArenaInfo domain.TempArenaInfo) error

	// for migration
	IsExistUser() bool
	DeleteUser() error
	IsExistAlertPlayers() bool
	DeleteAlertPlayers() error
}
