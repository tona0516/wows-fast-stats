package infra

import "wfs/backend/vo"

type ConfigInterface interface {
	User() (vo.UserConfig, error)
	UpdateUser(config vo.UserConfig) error
	App() (vo.AppConfig, error)
	UpdateApp(config vo.AppConfig) error
	AlertPlayers() ([]vo.AlertPlayer, error)
	UpdateAlertPlayer(player vo.AlertPlayer) error
	RemoveAlertPlayer(accountID int) error
}
