package service

import (
	"changeme/backend/domain"
	"changeme/backend/repo"
	"changeme/backend/vo"
	"crypto/md5"
	"fmt"
	"sort"
	"sync"
)

type StatsService struct{
    Parallels uint
}

func (s *StatsService) TempArenaInfoHash(installPath string) (string, error) {
    var result string
    local := repo.Local{}

    tempArenaInfo, err := local.TempArenaInfo(installPath)
    if err != nil {
        return result, err
    }

    md5 := md5.Sum([]byte(fmt.Sprintf("%x", tempArenaInfo)))
    result = fmt.Sprintf("%x", md5)
    return result, nil
}

func (s *StatsService) Battle(installPath string, appid string) (vo.Battle, error) {
    var result vo.Battle

	wargaming := repo.Wargaming{AppID: appid}
	numbers := repo.Numbers{}
	local := repo.Local{}
    unregistered := repo.Unregistered{}

	tempArenaInfo, err := local.TempArenaInfo(installPath)
	if err != nil {
		return result, err
	}

	accountListResult := make(chan vo.Result[vo.WGAccountList])
	encyclopediaInfoResult := make(chan vo.Result[vo.WGEncyclopediaInfo])
	accountInfoResult := make(chan vo.Result[vo.WGAccountInfo])
	shipStatsResult := make(chan vo.Result[map[int]vo.WGShipsStats])
	clanTagResult := make(chan vo.Result[map[int]string])
	shipInfoResult := make(chan vo.Result[map[int]vo.Warship])
	expectedStatsResult := make(chan vo.Result[vo.NSExpectedStats])
    battleArenasResult := make(chan vo.Result[vo.WGBattleArenas])
    battleTypesResult := make(chan vo.Result[vo.WGBattleTypes])

	go s.accountList(&wargaming, tempArenaInfo, accountListResult)
	go s.encyclopediaInfo(&wargaming, encyclopediaInfoResult)

	accountList := <-accountListResult
	if accountList.Error != nil {
		return result, accountList.Error
	}
	accountIDs := accountList.Value.AccountIDs()

	go s.accountInfo(&wargaming, accountIDs, accountInfoResult)
	go s.shipStats(&wargaming, accountIDs, shipStatsResult)
	go s.clanTag(&wargaming, accountIDs, clanTagResult)

	encyclopediaInfo := <-encyclopediaInfoResult
	if encyclopediaInfo.Error != nil {
		return result, encyclopediaInfo.Error
	}
	gameVersion := encyclopediaInfo.Value.Data.GameVersion

	go s.warship(&wargaming, gameVersion, shipInfoResult)
	go s.expectedStats(&numbers, gameVersion, expectedStatsResult)
    go s.battleArenas(&wargaming, gameVersion, battleArenasResult)
    go s.battleTypes(&wargaming, gameVersion, battleTypesResult)

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
	shipInfo := <-shipInfoResult
	if shipInfo.Error != nil {
		return result, shipInfo.Error
	}
	expectedStats := <-expectedStatsResult
	if expectedStats.Error != nil {
		return result, expectedStats.Error
	}
    unregisteredShipInfo, err := unregistered.Warship()
    if err != nil {
        return result, err
    }
    for k, v := range unregisteredShipInfo {
        if _, ok := shipInfo.Value[k]; !ok {
            shipInfo.Value[k] = v
        }
    }
    battleArenas := <-battleArenasResult
    if battleArenas.Error != nil {
        return result, battleArenas.Error
    }
    battleTypes := <-battleTypesResult
    if battleTypes.Error != nil {
        return result, battleTypes.Error
    }

	result = s.compose(
		tempArenaInfo,
		accountInfo.Value,
		accountList.Value,
		clanTag.Value,
		shipStats.Value,
		shipInfo.Value,
		expectedStats.Value,
        battleArenas.Value,
        battleTypes.Value,
	)

	return result, nil
}

func (s *StatsService) accountList(wargaming *repo.Wargaming, tempArenaInfo vo.TempArenaInfo, result chan vo.Result[(vo.WGAccountList)]) {
	accountNames := tempArenaInfo.AccountNames()

	accountList, err := wargaming.AccountList(accountNames)
	if err != nil {
		result <- vo.Result[vo.WGAccountList]{Value: accountList, Error: err}
		return
	}

	result <- vo.Result[vo.WGAccountList]{Value: accountList, Error: nil}
}

func (s *StatsService) encyclopediaInfo(wargaming *repo.Wargaming, result chan vo.Result[vo.WGEncyclopediaInfo]) {
	encyclopediaInfo, err := wargaming.EncyclopediaInfo()
	if err != nil {
		result <- vo.Result[vo.WGEncyclopediaInfo]{Value: encyclopediaInfo, Error: err}
		return
	}

	result <- vo.Result[vo.WGEncyclopediaInfo]{Value: encyclopediaInfo, Error: nil}
}

func (s *StatsService) accountInfo(wargaming *repo.Wargaming, accountIDs []int, result chan vo.Result[vo.WGAccountInfo]) {
	accountInfo, err := wargaming.AccountInfo(accountIDs)
	if err != nil {
		result <- vo.Result[vo.WGAccountInfo]{Value: accountInfo, Error: err}
		return
	}

	result <- vo.Result[vo.WGAccountInfo]{Value: accountInfo, Error: nil}
}

func (s *StatsService) shipStats(wargaming *repo.Wargaming, accountIDs []int, result chan vo.Result[map[int]vo.WGShipsStats]) {
	shipStatsMap := make(map[int]vo.WGShipsStats)
	limit := make(chan struct{}, s.Parallels)
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

func (s *StatsService) clanTag(wargaming *repo.Wargaming, accountIDs []int, result chan vo.Result[map[int]string]) {
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

func (s *StatsService) warship(wargaming *repo.Wargaming, gameVersion string, result chan vo.Result[map[int]vo.Warship]) {
	shipInfoMap := make(map[int]vo.Warship, 0)

	cache := repo.Cache[map[int]vo.Warship]{
		FileName: "shipinfo_" + gameVersion + ".bin",
	}

	object, err := cache.Deserialize()
	if err == nil {
		result <- vo.Result[map[int]vo.Warship]{Value: object, Error: nil}
		return
	}

	res, err := wargaming.EncyclopediaShips(1)
	if err != nil {
		result <- vo.Result[map[int]vo.Warship]{Value: shipInfoMap, Error: err}
		return
	}
	pageTotal := res.Meta.PageTotal

	var mu sync.Mutex
	limit := make(chan struct{}, s.Parallels)
	wg := sync.WaitGroup{}
	for i := 1; i <= pageTotal; i++ {
		limit <- struct{}{}
		wg.Add(1)
		go func(pageNo int) {
			defer func() {
				wg.Done()
				<-limit
			}()

			encyclopediaShips, err := wargaming.EncyclopediaShips(pageNo)
			if err != nil {
				result <- vo.Result[map[int]vo.Warship]{Value: shipInfoMap, Error: err}
				return
			}

			for shipID, shipInfo := range encyclopediaShips.Data {
				mu.Lock()
				shipInfoMap[shipID] = vo.Warship{
					Name:   shipInfo.Name,
					Tier:   shipInfo.Tier,
					Type:   shipInfo.Type,
					Nation: shipInfo.Nation,
				}
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()

	if err := cache.Serialize(shipInfoMap); err != nil {
		result <- vo.Result[map[int]vo.Warship]{Value: shipInfoMap, Error: err}
		return
	}

	result <- vo.Result[map[int]vo.Warship]{Value: shipInfoMap, Error: nil}
}

func (s *StatsService) expectedStats(numbers *repo.Numbers, gameVersion string, result chan vo.Result[vo.NSExpectedStats]) {
	cache := repo.Cache[vo.NSExpectedStats]{
		FileName: "expectedstats_" + gameVersion + ".bin",
	}

	object, err := cache.Deserialize()
	if err == nil {
		result <- vo.Result[vo.NSExpectedStats]{Value: object, Error: nil}
		return
	}

	expectedStats, err := numbers.ExpectedStats()
	if err != nil {
		result <- vo.Result[vo.NSExpectedStats]{Value: *expectedStats, Error: err}
		return
	}

	if err := cache.Serialize(*expectedStats); err != nil {
		result <- vo.Result[vo.NSExpectedStats]{Value: *expectedStats, Error: err}
		return
	}

	result <- vo.Result[vo.NSExpectedStats]{Value: *expectedStats, Error: err}
}

func (s *StatsService) battleArenas(wargaming *repo.Wargaming, gameVersion string, result chan vo.Result[vo.WGBattleArenas]) {
	cache := repo.Cache[vo.WGBattleArenas]{
		FileName: "battlearenas_" + gameVersion + ".bin",
	}

	object, err := cache.Deserialize()
	if err == nil {
		result <- vo.Result[vo.WGBattleArenas]{Value: object, Error: nil}
		return
	}

	battleArenas, err := wargaming.BattleArenas()
	if err != nil {
		result <- vo.Result[vo.WGBattleArenas]{Value: battleArenas, Error: err}
		return
	}

	if err := cache.Serialize(battleArenas); err != nil {
		result <- vo.Result[vo.WGBattleArenas]{Value: battleArenas, Error: err}
		return
	}

    result <- vo.Result[vo.WGBattleArenas]{Value: battleArenas, Error: err}
}

func (s *StatsService) battleTypes(wargaming *repo.Wargaming, gameVersion string, result chan vo.Result[vo.WGBattleTypes]) {
	cache := repo.Cache[vo.WGBattleTypes]{
		FileName: "battletypes_" + gameVersion + ".bin",
	}

	object, err := cache.Deserialize()
	if err == nil {
		result <- vo.Result[vo.WGBattleTypes]{Value: object, Error: nil}
		return
	}

	battleTypes, err := wargaming.BattleTypes()
	if err != nil {
		result <- vo.Result[vo.WGBattleTypes]{Value: battleTypes, Error: err}
		return
	}

	if err := cache.Serialize(battleTypes); err != nil {
		result <- vo.Result[vo.WGBattleTypes]{Value: battleTypes, Error: err}
		return
	}

    result <- vo.Result[vo.WGBattleTypes]{Value: battleTypes, Error: err}
}

func (s *StatsService) compose(
	tempArenaInfo vo.TempArenaInfo,
	accountInfo vo.WGAccountInfo,
	accountList vo.WGAccountList,
	clanTag map[int]string,
	shipStats map[int]vo.WGShipsStats,
	shipInfo map[int]vo.Warship,
	expectedStats vo.NSExpectedStats,
    battleArenas vo.WGBattleArenas,
    battleTypes vo.WGBattleTypes,
) vo.Battle {
    numbersURLGenerator := domain.NumbersURLGenerator{}

	friends := make(vo.Players, 0)
	enemies := make(vo.Players, 0)
	rating := domain.Rating{}

    var Ownship string

	for i := range tempArenaInfo.Vehicles {
		vehicle := tempArenaInfo.Vehicles[i]
		playerShipInfo := shipInfo[vehicle.ShipID]

		nickname := vehicle.Name
        if nickname == tempArenaInfo.PlayerName {
            Ownship = playerShipInfo.Name
        }
		accountID := accountList.AccountID(nickname)
		clan := clanTag[accountID]

		playerAccountInfo := accountInfo.Data[accountID]
        statsCalculator := domain.StatsCalculator{
            Player: domain.Stats{
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
                statsCalculator.SetShipStats(domain.Stats{
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

		expectedShipStats := expectedStats.Data[vehicle.ShipID]

		player := vo.Player{
			ShipInfo: vo.ShipInfo{
				Name:   playerShipInfo.Name,
				Nation: playerShipInfo.Nation,
				Tier:   playerShipInfo.Tier,
				Type:   playerShipInfo.Type,
                StatsURL: numbersURLGenerator.ShipPage(vehicle.ShipID, playerShipInfo.Name),
			},
			ShipStats: vo.ShipStats{
				Battles:   statsCalculator.Ship.Battles,
				AvgDamage: statsCalculator.ShipAvgDamage(),
				WinRate:   statsCalculator.ShipWinRate(),
                WinSurvivedRate: statsCalculator.ShipWinSurvivedRate(),
                LoseSurvivedRate: statsCalculator.ShipLoseSurvivedRate(),
				KdRate:    statsCalculator.ShipKdRate(),
                Exp: statsCalculator.ShipAvgExp(),
				PersonalRating: rating.PersonalRating(
					statsCalculator.ShipAvgDamage(),
					statsCalculator.ShipAvgFrags(),
					statsCalculator.ShipWinRate(),
					expectedShipStats.AverageDamageDealt,
					expectedShipStats.AverageFrags,
					expectedShipStats.WinRate,
				),
			},
			PlayerInfo: vo.PlayerInfo{
                ID: accountID,
				Name: nickname,
				Clan: clan,
                IsHidden: playerAccountInfo.HiddenProfile,
                StatsURL: numbersURLGenerator.PlayerPage(accountID, nickname),
			},
			PlayerStats: vo.PlayerStats{
				Battles:   statsCalculator.Player.Battles,
				AvgDamage: statsCalculator.PlayerAvgDamage(),
				WinRate:   statsCalculator.PlayerWinRate(),
                WinSurvivedRate: statsCalculator.PlayerWinSurvivedRate(),
                LoseSurvivedRate: statsCalculator.PlayerLoseSurvivedRate(),
				KdRate:    statsCalculator.PlayerKdRate(),
                Exp: statsCalculator.PlayerAvgExp(),
				AvgTier:   statsCalculator.PlayerAvgTier(accountID, shipInfo, shipStats),
                UsingShipTypeRate: statsCalculator.UsingShipTypeRate(accountID, shipInfo, shipStats),
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
        TeamAverage: friends.TeamAverage(),
    })
    teams = append(teams, vo.Team{
        Players: enemies,
        Name: "敵チーム",
        TeamAverage: enemies.TeamAverage(),
    })

    battle := vo.Battle {
        Meta: vo.Meta{
            Date: tempArenaInfo.FormattedDateTime(),
            Arena: tempArenaInfo.BattleArena(battleArenas),
            Type: tempArenaInfo.BattleType(battleTypes),
            OwnShip: Ownship,
        },
        Teams: teams,
    }

    return battle
}
