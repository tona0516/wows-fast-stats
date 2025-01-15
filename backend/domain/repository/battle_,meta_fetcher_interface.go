package repository

import (
	"wfs/backend/domain/model"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type BattleMetaFetcherInterface interface {
	Fetch() (model.BattleMeta, error)
}
