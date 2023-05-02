package service

import (
	"changeme/backend/domain"
	"changeme/backend/infra"
	"changeme/backend/vo"
	"crypto/md5"
	"fmt"
	"sort"
	"sync"
)

type Battle struct{
    Parallels uint
    UserConfig vo.UserConfig
}

func (b *Battle) TempArenaInfoHash() (string, error) {
    var result string
    tempArenaInfoRepo := infra.TempArenaInfo{}

    tempArenaInfo, err := tempArenaInfoRepo.Get(b.UserConfig.InstallPath)
    if err != nil {
        return result, err
    }

    md5 := md5.Sum([]byte(fmt.Sprintf("%x", tempArenaInfo)))
    result = fmt.Sprintf("%x", md5)
    return result, nil
}

func (b *Battle) Battle() (vo.Battle, error) {
    var result vo.Battle

	wargaming := infra.Wargaming{AppID: b.UserConfig.Appid}
    tempArenaInfoRepo := infra.TempArenaInfo{}

	tempArenaInfo, err := tempArenaInfoRepo.Get(b.UserConfig.InstallPath)
	if err != nil {
		return result, err
	}

    if b.UserConfig.SaveTempArenaInfo {
        if err := tempArenaInfoRepo.Save(tempArenaInfo); err != nil {
            return result, err
        }
    }

	accountListResult := make(chan vo.Result[vo.WGAccountList])
	encyclopediaInfoResult := make(chan vo.Result[vo.WGEncyclopediaInfo])
	accountInfoResult := make(chan vo.Result[vo.WGAccountInfo])
	shipStatsResult := make(chan vo.Result[map[int]vo.WGShipsStats])
	clanTagResult := make(chan vo.Result[map[int]string])

	go b.accountList(&wargaming, tempArenaInfo, accountListResult)
	go b.encyclopediaInfo(&wargaming, encyclopediaInfoResult)

	accountList := <-accountListResult
	if accountList.Error != nil {
		return result, accountList.Error
	}
	accountIDs := accountList.Value.AccountIDs()

	go b.accountInfo(&wargaming, accountIDs, accountInfoResult)
	go b.shipStats(&wargaming, accountIDs, shipStatsResult)
	go b.clanTag(&wargaming, accountIDs, clanTagResult)

    warshipCache := infra.Cache[map[int]vo.Warship]{Name: "warship"}
    warship, err := warshipCache.Deserialize()
    if err != nil {
        return result, err
    }

    expectedStatsCache := infra.Cache[vo.NSExpectedStats]{Name: "expectedstats"}
    expectedStats, err := expectedStatsCache.Deserialize()
    if err != nil {
        return result, err
    }

    battleArenasCache := infra.Cache[vo.WGBattleArenas]{Name: "battlearenas"}
    battleArenas, err := battleArenasCache.Deserialize()
    if err != nil {
        return result, err
    }

    battleTypesCache := infra.Cache[vo.WGBattleTypes]{Name: "battletypes"}
    battleTypes, err := battleTypesCache.Deserialize()
    if err != nil {
        return result, err
    }

	accountInfo := <-accountInfoResult
	if accountInfo.Error != nil {
		return result, accountInfo.Error
	}
	shipStats := <-shipStatsResult
	if shipStats.Error != nil {
		return result, shipStats.Error
	}
	clanTag := <-clanTagResult
	if clanTag.Error != nil {
		return result, clanTag.Error
	}

	result = b.compose(
		tempArenaInfo,
		accountInfo.Value,
		accountList.Value,
		clanTag.Value,
		shipStats.Value,
		warship,
		expectedStats,
        battleArenas,
        battleTypes,
	)

	return result, nil
}

func (b *Battle) accountList(wargaming *infra.Wargaming, tempArenaInfo vo.TempArenaInfo, result chan vo.Result[(vo.WGAccountList)]) {
	accountNames := tempArenaInfo.AccountNames()
	accountList, err := wargaming.AccountList(accountNames)
	result <- vo.Result[vo.WGAccountList]{Value: accountList, Error: err}
}

func (b *Battle) encyclopediaInfo(wargaming *infra.Wargaming, result chan vo.Result[vo.WGEncyclopediaInfo]) {
	encyclopediaInfo, err := wargaming.EncyclopediaInfo()
    result <- vo.Result[vo.WGEncyclopediaInfo]{Value: encyclopediaInfo, Error: err}
}

func (b *Battle) accountInfo(wargaming *infra.Wargaming, accountIDs []int, result chan vo.Result[vo.WGAccountInfo]) {
	accountInfo, err := wargaming.AccountInfo(accountIDs)
    result <- vo.Result[vo.WGAccountInfo]{Value: accountInfo, Error: err}
}

func (b *Battle) shipStats(wargaming *infra.Wargaming, accountIDs []int, result chan vo.Result[map[int]vo.WGShipsStats]) {
	shipStatsMap := make(map[int]vo.WGShipsStats)
	limit := make(chan struct{}, b.Parallels)
	wg := sync.WaitGroup{}
    var mu sync.Mutex
	for i := range accountIDs {
		limit <- struct{}{}
		wg.Add(1)
		go func(accountID int) {
			defer func() {
				wg.Done()
				<-limit
			}()

			shipStats, err := wargaming.ShipsStats(accountID)
			if err != nil {
				result <- vo.Result[map[int]vo.WGShipsStats]{Value: shipStatsMap, Error: err}
				return
			}

            mu.Lock()
			shipStatsMap[accountID] = shipStats
            mu.Unlock()
		}(accountIDs[i])
	}
	wg.Wait()

	result <- vo.Result[map[int]vo.WGShipsStats]{Value: shipStatsMap, Error: nil}
}

func (b *Battle) clanTag(wargaming *infra.Wargaming, accountIDs []int, result chan vo.Result[map[int]string]) {
	clanTagMap := make(map[int]string)

	clansAccountInfo, err := wargaming.ClansAccountInfo(accountIDs)
	if err != nil {
		result <- vo.Result[map[int]string]{Value: clanTagMap, Error: err}
		return
	}

	clanIDs := clansAccountInfo.ClanIDs()

	clansInfo, err := wargaming.ClansInfo(clanIDs)
	if err != nil {
		result <- vo.Result[map[int]string]{Value: clanTagMap, Error: err}
		return
	}

	for i := range accountIDs {
		accountID := accountIDs[i]
		clanID := clansAccountInfo.Data[accountID].ClanID
		clanTag := clansInfo.Data[clanID].Tag
		clanTagMap[accountID] = clanTag
	}

	result <- vo.Result[map[int]string]{Value: clanTagMap, Error: nil}
}

func (b *Battle) compose(
	tempArenaInfo vo.TempArenaInfo,
	accountInfo vo.WGAccountInfo,
	accountList vo.WGAccountList,
	clanTag map[int]string,
	shipStats map[int]vo.WGShipsStats,
	warships map[int]vo.Warship,
	expectedStats vo.NSExpectedStats,
    battleArenas vo.WGBattleArenas,
    battleTypes vo.WGBattleTypes,
) vo.Battle {
    urlGen := domain.NumbersURLGenerator{}

	friends := make(vo.Players, 0)
	enemies := make(vo.Players, 0)
	rating := domain.Rating{}

    var ownShip string

	for i := range tempArenaInfo.Vehicles {
		vehicle := tempArenaInfo.Vehicles[i]
		warship := warships[vehicle.ShipID]

		nickname := vehicle.Name
        if nickname == tempArenaInfo.PlayerName {
            ownShip = warship.Name
        }
		accountID := accountList.AccountID(nickname)
		clan := clanTag[accountID]

		playerAccountInfo := accountInfo.Data[accountID]
        stats := domain.Stats{
            Overall: domain.StatsFactor{
                Battles:         playerAccountInfo.Statistics.Pvp.Battles,
                SurvivedBattles: playerAccountInfo.Statistics.Pvp.SurvivedBattles,
                DamageDealt:     playerAccountInfo.Statistics.Pvp.DamageDealt,
                Frags:           playerAccountInfo.Statistics.Pvp.Frags,
                Wins:            playerAccountInfo.Statistics.Pvp.Wins,
                SurvivedWins:    playerAccountInfo.Statistics.Pvp.SurviveWins,
                Xp: playerAccountInfo.Statistics.Pvp.Xp,
            },
        }
		for _, v:= range shipStats[accountID].Data[accountID] {
			if v.ShipID == vehicle.ShipID {
                stats.SetShipStats(domain.StatsFactor{
                    Battles:         v.Pvp.Battles,
                    SurvivedBattles: v.Pvp.SurvivedBattles,
                    DamageDealt:     v.Pvp.DamageDealt,
                    Frags:           v.Pvp.Frags,
                    Wins:            v.Pvp.Wins,
                    SurvivedWins: v.Pvp.SurviveWins,
                    Xp: v.Pvp.Xp,
                })
				break
			}
		}

		expected := expectedStats.Data[vehicle.ShipID]

		player := vo.Player{
			ShipInfo: vo.ShipInfo{
				Name:   warship.Name,
				Nation: warship.Nation,
				Tier:   warship.Tier,
				Type:   warship.Type,
                StatsURL: urlGen.ShipPage(vehicle.ShipID, warship.Name),
			},
			ShipStats: vo.ShipStats{
				Battles:   stats.Ship.Battles,
				Damage: stats.ShipAvgDamage(),
				WinRate:   stats.ShipWinRate(),
                WinSurvivedRate: stats.ShipWinSurvivedRate(),
                LoseSurvivedRate: stats.ShipLoseSurvivedRate(),
				KdRate:    stats.ShipKdRate(),
                Exp: stats.ShipAvgExp(),
				PR: rating.PersonalRating(
                    domain.RatingFactor{
                        Damage: stats.ShipAvgDamage(),
                        Frags: stats.ShipAvgFrags(),
                        Wins: stats.ShipWinRate(),
                    },
                    domain.RatingFactor{
                        Damage: expected.AverageDamageDealt,
                        Frags: expected.AverageFrags,
                        Wins: expected.WinRate,
                    },
				),
			},
			PlayerInfo: vo.PlayerInfo{
                ID: accountID,
				Name: nickname,
				Clan: clan,
                IsHidden: playerAccountInfo.HiddenProfile,
                StatsURL: urlGen.PlayerPage(accountID, nickname),
			},
			OverallStats: vo.OverallStats{
				Battles:   stats.Overall.Battles,
				Damage: stats.OverallAvgDamage(),
				WinRate:   stats.OverallWinRate(),
                WinSurvivedRate: stats.OverallWinSurvivedRate(),
                LoseSurvivedRate: stats.OverallLoseSurvivedRate(),
				KdRate:    stats.OverallKdRate(),
                Exp: stats.OverallAvgExp(),
				AvgTier:   stats.OverallAvgTier(accountID, warships, shipStats),
                UsingShipTypeRate: stats.OverallUsingShipTypeRate(accountID, warships, shipStats),
                UsingTierRate: stats.OverallUsingTierRate(accountID, warships, shipStats),
			},
		}

		if vehicle.Relation == 0 || vehicle.Relation == 1 {
			friends = append(friends, player)
		} else {
			enemies = append(enemies, player)
		}
	}

	sort.Sort(friends)
	sort.Sort(enemies)

    teams := make([]vo.Team, 0)
    teams = append(teams, vo.Team{
        Players: friends,
        Name: "味方チーム",
    })
    teams = append(teams, vo.Team{
        Players: enemies,
        Name: "敵チーム",
    })

    battle := vo.Battle {
        Meta: vo.Meta{
            Date: tempArenaInfo.FormattedDateTime(),
            Arena: tempArenaInfo.BattleArena(battleArenas),
            Type: tempArenaInfo.BattleType(battleTypes),
            OwnShip: ownShip,
        },
        Teams: teams,
    }

    return battle
}
