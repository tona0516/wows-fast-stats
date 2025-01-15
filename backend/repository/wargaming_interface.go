package repository

import "wfs/backend/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type WargamingInterface interface {
	AccountList(accountNames []string) (data.WGAccountList, error)
	AccountListForSearch(prefix string) (data.WGAccountList, error)
}
