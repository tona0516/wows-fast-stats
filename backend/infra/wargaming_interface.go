package infra

import "wfs/backend/vo"

type WargamingInterface interface {
	SetAppID(appid string)
	AccountInfo(accountIDs []int) (vo.WGAccountInfo, error)
	AccountList(accountNames []string) (vo.WGAccountList, error)
	AccountListForSearch(prefix string) (vo.WGAccountList, error)
	ClansAccountInfo(accountIDs []int) (vo.WGClansAccountInfo, error)
	ClansInfo(clanIDs []int) (vo.WGClansInfo, error)
	EncycShips(pageNo int) (vo.WGEncycShips, error)
	ShipsStats(accountID int) (vo.WGShipsStats, error)
	EncycInfo() (vo.WGEncycInfo, error)
	BattleArenas() (vo.WGBattleArenas, error)
	BattleTypes() (vo.WGBattleTypes, error)
}
