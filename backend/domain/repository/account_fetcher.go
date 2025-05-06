package repository

import (
	"wfs/backend/domain/model"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type AccountFetcher interface {
	FetchByPrefix(prefix string) (model.Accounts, error)
	FetchByNames(names []string) (model.Accounts, error)
}
