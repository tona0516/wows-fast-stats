package repository

import "wfs/backend/data"

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type WargamingInterface interface {
	AccountInfo(appID string, accountIDs []int) (data.WGAccountInfo, error)
	AccountList(appID string, accountNames []string) (data.WGAccountList, error)
	AccountListForSearch(appID string, prefix string) (data.WGAccountList, error)
	ClansAccountInfo(appID string, accountIDs []int) (data.WGClansAccountInfo, error)
	ClansInfo(appID string, clanIDs []int) (data.WGClansInfo, error)
	EncycShips(appID string, pageNo int) (data.WGEncycShips, int, error)
	ShipsStats(appID string, accountID int) (data.WGShipsStats, error)
	BattleArenas(appID string) (data.WGBattleArenas, error)
	BattleTypes(appID string) (data.WGBattleTypes, error)
	Test(appID string) (bool, error)
}
