package repository

import "wfs/backend/domain/model"

type ConfigV0Interface interface {
	User() (model.UserConfig, error)
	IsExistUser() bool
	DeleteUser() error
	AlertPlayers() ([]model.AlertPlayer, error)
	IsExistAlertPlayers() bool
	DeleteAlertPlayers() error
}
