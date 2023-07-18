package repository

import (
	"wfs/backend/application/vo"
	"wfs/backend/domain"
)

type LocalFileInterface interface {
	User() (domain.UserConfig, error)
	UpdateUser(config domain.UserConfig) error
	App() (vo.AppConfig, error)
	UpdateApp(config vo.AppConfig) error
	AlertPlayers() ([]domain.AlertPlayer, error)
	UpdateAlertPlayer(player domain.AlertPlayer) error
	RemoveAlertPlayer(accountID int) error
	SaveScreenshot(path string, base64Data string) error
	TempArenaInfo(installPath string) (domain.TempArenaInfo, error)
	SaveTempArenaInfo(tempArenaInfo domain.TempArenaInfo) error
}
