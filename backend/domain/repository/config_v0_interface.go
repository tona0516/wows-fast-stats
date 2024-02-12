package repository

import "wfs/backend/domain/model"

type ConfigV0Interface interface {
	UserV1() (model.UserConfigV1, error)
	IsExistUser() bool
	DeleteUser() error
	AlertPlayers() ([]model.AlertPlayer, error)
	IsExistAlertPlayers() bool
	DeleteAlertPlayers() error
}
