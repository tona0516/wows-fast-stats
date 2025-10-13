package repository

import "wfs/backend/domain"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type BlackListInterface interface {
	Load() (*domain.BlackList, error)
	Save(data domain.BlackList) error
}
