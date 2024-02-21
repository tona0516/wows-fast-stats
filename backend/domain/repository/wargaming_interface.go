package repository

import "wfs/backend/domain/model"

//go:generate mockgen -source=$GOFILE -destination ../mock_$GOPACKAGE/$GOFILE -package mock_$GOPACKAGE
type WargamingInterface interface {
	AccountInfo(appID string, accountIDs []int) (model.WGAccountInfo, error)
	AccountList(appID string, accountNames []string) (model.WGAccountList, error)
	AccountListForSearch(appID string, prefix string) (model.WGAccountList, error)
	ClansAccountInfo(appID string, accountIDs []int) (model.WGClansAccountInfo, error)
	ClansInfo(appID string, clanIDs []int) (model.WGClansInfo, error)
	EncycShips(appID string, pageNo int) (model.WGEncycShips, int, error)
	ShipsStats(appID string, accountID int) (model.WGShipsStats, error)
	BattleArenas(appID string) (model.WGBattleArenas, error)
	BattleTypes(appID string) (model.WGBattleTypes, error)
	Test(appID string) (bool, error)
}
