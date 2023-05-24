package infra

import "changeme/backend/vo"

type WargamingInterface interface {
	SetAppID(appid string)
	AccountInfo(accountIDs []int) (vo.WGAccountInfo, error)
	AccountList(accountNames []string) (vo.WGAccountList, error)
	ClansAccountInfo(accountIDs []int) (vo.WGClansAccountInfo, error)
	ClansInfo(clanIDs []int) (vo.WGClansInfo, error)
	EncyclopediaShips(pageNo int) (vo.WGEncyclopediaShips, error)
	ShipsStats(accountID int) (vo.WGShipsStats, error)
	EncyclopediaInfo() (vo.WGEncyclopediaInfo, error)
	BattleArenas() (vo.WGBattleArenas, error)
	BattleTypes() (vo.WGBattleTypes, error)
}
