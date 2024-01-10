package repository

import (
	"wfs/backend/domain"
)

type ConfigV0Interface interface {
	User() (domain.UserConfig, error)
	IsExistUser() bool
	DeleteUser() error
	AlertPlayers() ([]domain.AlertPlayer, error)
	IsExistAlertPlayers() bool
	DeleteAlertPlayers() error
}
