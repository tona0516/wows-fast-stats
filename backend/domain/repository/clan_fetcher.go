package repository

import (
	"wfs/backend/domain/model"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type ClanFetcher interface {
	Fetch(accountIDs []int) (model.Clans, error)
}
