package repository

import "wfs/backend/domain/model"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type AlertPlayerStore interface {
	IsExistV0() bool
	IsExistV1() bool
	GetV0() ([]model.AlertPlayer, error)
	GetV1() ([]model.AlertPlayer, error)
	SaveV1(players []model.AlertPlayer) error
	DeleteV0() error
}
