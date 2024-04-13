package repository

import "wfs/backend/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type ConfigV0Interface interface {
	User() (data.UserConfig, error)
	IsExistUser() bool
	DeleteUser() error
	AlertPlayers() ([]data.AlertPlayer, error)
	IsExistAlertPlayers() bool
	DeleteAlertPlayers() error
}
