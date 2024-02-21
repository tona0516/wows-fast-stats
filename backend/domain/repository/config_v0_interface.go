package repository

import "wfs/backend/domain/model"

//go:generate mockgen -source=$GOFILE -destination ../mock_$GOPACKAGE/$GOFILE -package mock_$GOPACKAGE
type ConfigV0Interface interface {
	User() (model.UserConfig, error)
	IsExistUser() bool
	DeleteUser() error
	AlertPlayers() ([]model.AlertPlayer, error)
	IsExistAlertPlayers() bool
	DeleteAlertPlayers() error
}
