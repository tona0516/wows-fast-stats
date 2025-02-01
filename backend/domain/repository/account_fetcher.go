package repository

import (
	"wfs/backend/domain/model"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type AccountFetcher interface {
	Search(prefix string) (model.Accounts, error)
	Fetch(playerNames []string) (model.Accounts, error)
}
