package repository

import "wfs/backend/domain"

type WargamingInterface interface {
	SetAppID(appid string)
	AccountInfo(accountIDs []int) (domain.WGAccountInfo, error)
	AccountList(accountNames []string) (domain.WGAccountList, error)
	AccountListForSearch(prefix string) (domain.WGAccountList, error)
	ClansAccountInfo(accountIDs []int) (domain.WGClansAccountInfo, error)
	ClansInfo(clanIDs []int) (domain.WGClansInfo, error)
	EncycShips(pageNo int) (domain.WGEncycShips, error)
	ShipsStats(accountID int) (domain.WGShipsStats, error)
	BattleArenas() (domain.WGBattleArenas, error)
	BattleTypes() (domain.WGBattleTypes, error)
	Test(appid string) (bool, error)
}
