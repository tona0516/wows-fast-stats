package repository

import "wfs/backend/domain/model"

type WargamingInterface interface {
	SetAppID(appid string)
	AccountInfo(accountIDs []int) (model.WGAccountInfo, error)
	AccountList(accountNames []string) (model.WGAccountList, error)
	AccountListForSearch(prefix string) (model.WGAccountList, error)
	ClansAccountInfo(accountIDs []int) (model.WGClansAccountInfo, error)
	ClansInfo(clanIDs []int) (model.WGClansInfo, error)
	EncycShips(pageNo int) (model.WGEncycShips, int, error)
	ShipsStats(accountID int) (model.WGShipsStats, error)
	BattleArenas() (model.WGBattleArenas, error)
	BattleTypes() (model.WGBattleTypes, error)
	Test(appid string) (bool, error)
}
