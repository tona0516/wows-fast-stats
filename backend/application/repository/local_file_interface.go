package repository

import (
	"wfs/backend/domain"
)

type LocalFileInterface interface {
	User() (domain.UserConfig, error)
	UpdateUser(config domain.UserConfig) error
	AlertPlayers() ([]domain.AlertPlayer, error)
	UpdateAlertPlayer(player domain.AlertPlayer) error
	RemoveAlertPlayer(accountID int) error
	SaveScreenshot(path string, base64Data string) error
	TempArenaInfo(installPath string) (domain.TempArenaInfo, error)
	SaveTempArenaInfo(tempArenaInfo domain.TempArenaInfo) error
	CachedNSExpectedStats() (domain.NSExpectedStats, error)
	SaveNSExpectedStats(expectedStats domain.NSExpectedStats) error
}
