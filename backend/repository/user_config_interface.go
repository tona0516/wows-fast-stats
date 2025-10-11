package repository

import "wfs/backend/domain"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type UserConfigInterface interface {
	Load() (*domain.UserConfig, error)
	Save(data domain.UserConfig) error
}
