package repository

import "wfs/backend/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type WargamingInterface interface {
	AccountInfo(accountIDs []int) (data.WGAccountInfo, error)
	AccountList(accountNames []string) (data.WGAccountList, error)
	AccountListForSearch(prefix string) (data.WGAccountList, error)
	ShipsStats(accountID int) (data.WGShipsStats, error)
	BattleArenas() (data.WGBattleArenas, error)
	BattleTypes() (data.WGBattleTypes, error)
}
