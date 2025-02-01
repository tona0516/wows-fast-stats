package repository

import "wfs/backend/domain/model"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type UserConfigStore interface {
	IsExistV0() bool
	IsExistV1() bool
	GetV0() (model.UserConfig, error)
	GetV1() (model.UserConfig, error)
	GetV2() (model.UserConfigV2, error)
	SaveV1(config model.UserConfig) error
	SaveV2(config model.UserConfigV2) error
	DeleteV0() error
}
