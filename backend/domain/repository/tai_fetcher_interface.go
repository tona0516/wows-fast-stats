package repository

import (
	"wfs/backend/domain/model"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type TAIFetcherInterface interface {
	Get(installPath string) (model.TempArenaInfo, error)
}
