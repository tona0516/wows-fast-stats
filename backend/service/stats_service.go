package service

import (
	"changeme/backend/domain"
	"changeme/backend/repo"
	"changeme/backend/vo"
	"math"
	"sort"
	"sync"
)

type StatsService struct {
	InstallPath string
	AppID       string
	Parallels   uint
}

func (s *StatsService) GetsStats() (*vo.Team, error) {
	wargaming := repo.Wargaming{AppID: s.AppID}
	numbers := repo.Numbers{}
	local := repo.Local{}

	tempArenaInfo, err := local.GetTempArenaInfo(s.InstallPath)
	if err != nil {
		return nil, err
	}

	accountListResult := make(chan vo.Result[vo.WGAccountList])
	encyclopediaInfoResult := make(chan vo.Result[vo.WGEncyclopediaInfo])
	accountInfoResult := make(chan vo.Result[vo.WGAccountInfo])
	shipStatsResult := make(chan vo.Result[map[int]vo.WGShipsStats])
	clanTagResult := make(chan vo.Result[map[int]string])
	shipInfoResult := make(chan vo.Result[map[int]vo.ShipInfo])
	expectedStatsResult := make(chan vo.Result[vo.NSExpectedStats])

	go s.fetchAccountList(&wargaming, tempArenaInfo, accountListResult)
	go s.fetchEncyclopediaInfo(&wargaming, encyclopediaInfoResult)

	accountList := <-accountListResult
	if accountList.Error != nil {
		return nil, accountList.Error
	}
	accountIDs := accountList.Value.AccountIDs()

	go s.fetchAccountInfo(&wargaming, accountIDs, accountInfoResult)
	go s.fetchShipStats(&wargaming, accountIDs, shipStatsResult)
	go s.fetchClanTag(&wargaming, accountIDs, clanTagResult)

	encyclopediaInfo := <-encyclopediaInfoResult
	if encyclopediaInfo.Error != nil {
		return nil, encyclopediaInfo.Error
	}
	gameVersion := encyclopediaInfo.Value.Data.GameVersion

	go s.fetchShipInfo(&wargaming, gameVersion, shipInfoResult)
	go s.fetchExpectedStats(&numbers, gameVersion, expectedStatsResult)

	accountInfo := <-accountInfoResult
	if accountInfo.Error != nil {
		return nil, accountInfo.Error
	}
	shipStats := <-shipStatsResult
	if shipStats.Error != nil {
		return nil, shipStats.Error
	}
	clanTag := <-clanTagResult
	if clanTag.Error != nil {
		return nil, clanTag.Error
	}
	shipInfo := <-shipInfoResult
	if shipInfo.Error != nil {
		return nil, shipInfo.Error
	}
	expectedStats := <-expectedStatsResult
	if expectedStats.Error != nil {
		return nil, expectedStats.Error
	}

	team := s.compose(
		tempArenaInfo,
		accountInfo.Value,
		accountList.Value,
		clanTag.Value,
		shipStats.Value,
		shipInfo.Value,
		expectedStats.Value,
	)

	return &team, nil
}

func (s *StatsService) fetchAccountList(wargaming *repo.Wargaming, tempArenaInfo vo.TempArenaInfo, result chan vo.Result[(vo.WGAccountList)]) {
	accountNames := tempArenaInfo.AccountNames()

	accountList, err := wargaming.GetAccountList(accountNames)
	if err != nil {
		result <- vo.Result[vo.WGAccountList]{Value: accountList, Error: err}
		return
	}

	result <- vo.Result[vo.WGAccountList]{Value: accountList, Error: nil}
}

func (s *StatsService) fetchEncyclopediaInfo(wargaming *repo.Wargaming, result chan vo.Result[vo.WGEncyclopediaInfo]) {
	encyclopediaInfo, err := wargaming.GetEncyclopediaInfo()
	if err != nil {
		result <- vo.Result[vo.WGEncyclopediaInfo]{Value: encyclopediaInfo, Error: err}
		return
	}

	result <- vo.Result[vo.WGEncyclopediaInfo]{Value: encyclopediaInfo, Error: nil}
}

func (s *StatsService) fetchAccountInfo(wargaming *repo.Wargaming, accountIDs []int, result chan vo.Result[vo.WGAccountInfo]) {
	accountInfo, err := wargaming.GetAccountInfo(accountIDs)
	if err != nil {
		result <- vo.Result[vo.WGAccountInfo]{Value: accountInfo, Error: err}
		return
	}

	result <- vo.Result[vo.WGAccountInfo]{Value: accountInfo, Error: nil}
}

func (s *StatsService) fetchShipStats(wargaming *repo.Wargaming, accountIDs []int, result chan vo.Result[map[int]vo.WGShipsStats]) {
	shipStatsMap := make(map[int]vo.WGShipsStats)
	limit := make(chan struct{}, s.Parallels)
	wg := sync.WaitGroup{}
	for i := range accountIDs {
		limit <- struct{}{}
		wg.Add(1)
		go func(accountID int) {
			defer func() {
				wg.Done()
				<-limit
			}()

			shipStats, err := wargaming.GetShipsStats(accountID)
			if err != nil {
				result <- vo.Result[map[int]vo.WGShipsStats]{Value: shipStatsMap, Error: err}
				return
			}

			shipStatsMap[accountID] = shipStats
		}(accountIDs[i])
	}
	wg.Wait()

	result <- vo.Result[map[int]vo.WGShipsStats]{Value: shipStatsMap, Error: nil}
}

func (s *StatsService) fetchClanTag(wargaming *repo.Wargaming, accountIDs []int, result chan vo.Result[map[int]string]) {
	clanTagMap := make(map[int]string)

	clansAccountInfo, err := wargaming.GetClansAccountInfo(accountIDs)
	if err != nil {
		result <- vo.Result[map[int]string]{Value: clanTagMap, Error: err}
		return
	}

	clanIDs := clansAccountInfo.ClanIDs()

	clansInfo, err := wargaming.GetClansInfo(clanIDs)
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

func (s *StatsService) fetchShipInfo(wargaming *repo.Wargaming, gameVersion string, result chan vo.Result[map[int]vo.ShipInfo]) {
	shipInfoMap := make(map[int]vo.ShipInfo, 0)

	cache := repo.Cache[map[int]vo.ShipInfo]{
		FileName: "shipinfo_" + gameVersion + ".bin",
	}

	object, err := cache.Deserialize()
	if err == nil {
		result <- vo.Result[map[int]vo.ShipInfo]{Value: object, Error: nil}
		return
	}

	res, err := wargaming.GetEncyclopediaShips(1)
	if err != nil {
		result <- vo.Result[map[int]vo.ShipInfo]{Value: shipInfoMap, Error: err}
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

			encyclopediaShips, err := wargaming.GetEncyclopediaShips(pageNo)
			if err != nil {
				result <- vo.Result[map[int]vo.ShipInfo]{Value: shipInfoMap, Error: err}
				return
			}

			for shipID, shipInfo := range encyclopediaShips.Data {
				mu.Lock()
				shipInfoMap[shipID] = vo.ShipInfo{
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
		result <- vo.Result[map[int]vo.ShipInfo]{Value: shipInfoMap, Error: err}
		return
	}

	result <- vo.Result[map[int]vo.ShipInfo]{Value: shipInfoMap, Error: nil}
}

func (s *StatsService) fetchExpectedStats(numbers *repo.Numbers, gameVersion string, result chan vo.Result[vo.NSExpectedStats]) {
	cache := repo.Cache[vo.NSExpectedStats]{
		FileName: "expectedstats_" + gameVersion + ".bin",
	}

	object, err := cache.Deserialize()
	if err == nil {
		result <- vo.Result[vo.NSExpectedStats]{Value: object, Error: nil}
		return
	}

	expectedStats, err := numbers.Get()
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

func (s *StatsService) compose(
	tempArenaInfo vo.TempArenaInfo,
	accountInfo vo.WGAccountInfo,
	accountList vo.WGAccountList,
	clanTag map[int]string,
	shipStats map[int]vo.WGShipsStats,
	shipInfo map[int]vo.ShipInfo,
	expectedStats vo.NSExpectedStats,
) (vo.Team) {
	friends := make([]vo.Player, 0)
	enemies := make([]vo.Player, 0)
	rating := domain.Rating{}

	for i := range tempArenaInfo.Vehicles {
		vehicle := tempArenaInfo.Vehicles[i]
		playerShipInfo := shipInfo[vehicle.ShipID]

		nickname := vehicle.Name
		accountID := accountList.AccountID(nickname)
		clan := clanTag[accountID]

		var summaryStats domain.SummaryStats
		playerAccountInfo := accountInfo.Data[accountID]
		for k := range shipStats[accountID].Data[accountID] {
			playerShipStats := shipStats[accountID].Data[accountID][k]
			if playerShipStats.ShipID == vehicle.ShipID {
				summaryStats = domain.SummaryStats{
					Player: domain.Stats{
						Battles:         playerAccountInfo.Statistics.Pvp.Battles,
						SurvivedBattles: playerAccountInfo.Statistics.Pvp.SurvivedBattles,
						DamageDealt:     playerAccountInfo.Statistics.Pvp.DamageDealt,
						Xp:              playerAccountInfo.Statistics.Pvp.Xp,
						Frags:           playerAccountInfo.Statistics.Pvp.Frags,
						Wins:            playerAccountInfo.Statistics.Pvp.Wins,
					},
					Ship: domain.Stats{
						Battles:         playerShipStats.Pvp.Battles,
						SurvivedBattles: playerShipStats.Pvp.SurvivedBattles,
						DamageDealt:     playerShipStats.Pvp.DamageDealt,
						Xp:              playerShipStats.Pvp.Xp,
						Frags:           playerShipStats.Pvp.Frags,
						Wins:            playerShipStats.Pvp.Wins,
					},
				}
				break
			}
		}

		expectedShipStats := expectedStats.Data[vehicle.ShipID]

		player := vo.Player{
			ShipInfo: vo.PlayerShipInfo{
				Name:   playerShipInfo.Name,
				Nation: playerShipInfo.Nation,
				Tier:   playerShipInfo.Tier,
				Type:   playerShipInfo.Type,
			},
			ShipStats: vo.PlayerShipStats{
				Battles:   summaryStats.Player.Battles,
				AvgDamage: int(summaryStats.ShipAvgDamage()),
				AvgExp:    int(summaryStats.ShipAvgExp()),
				WinRate:   math.Round(summaryStats.ShipWinRate()*10) / 10,
				KdRate:    summaryStats.ShipKdRate(),
				CombatPower: rating.CombatPower(
					summaryStats.ShipAvgDamage(),
					summaryStats.ShipKdRate(),
					summaryStats.ShipAvgExp(),
					playerShipInfo.Tier,
					playerShipInfo.Type,
				),
				PersonalRating: rating.PersonalRating(
					summaryStats.ShipAvgDamage(),
					summaryStats.ShipAvgFrags(),
					summaryStats.ShipWinRate(),
					expectedShipStats.AverageDamageDealt,
					expectedShipStats.AverageFrags,
					expectedShipStats.WinRate,
				),
			},
			PlayerInfo: vo.PlayerPlayerInfo{
				Name: nickname,
				Clan: clan,
			},
			PlayerStats: vo.PlayerPlayerStats{
				Battles:   summaryStats.Player.Battles,
				AvgDamage: int(summaryStats.PlayerAvgDamage()),
				AvgExp:    int(summaryStats.PlayerAvgExp()),
				WinRate:   math.Round(summaryStats.PlayerWinRate()*10) / 10,
				KdRate:    summaryStats.PlayerKdRate(),
				AvgTier:   math.Round(summaryStats.PlayerAvgTier(accountID, shipInfo, shipStats)*10) / 10,
			},
		}

		if vehicle.Relation == 0 || vehicle.Relation == 1 {
			friends = append(friends, player)
		} else {
			enemies = append(enemies, player)
		}
	}

	sort.Slice(friends, func(i, j int) bool {
		one := friends[i].ShipInfo
		second := friends[j].ShipInfo
		if one.Type != second.Type {
			return one.Type < second.Type
		}
		if one.Tier != second.Tier {
			return one.Tier > second.Tier
		}
		if one.Nation != second.Nation {
			return one.Nation < second.Nation
		}
		return one.Name < second.Name
	})
	sort.Slice(enemies, func(i, j int) bool {
		one := enemies[i].ShipInfo
		second := enemies[j].ShipInfo
		if one.Type != second.Type {
			return one.Type < second.Type
		}
		if one.Tier != second.Tier {
			return one.Tier > second.Tier
		}
		if one.Nation != second.Nation {
			return one.Nation < second.Nation
		}
		return one.Name < second.Name
	})
	return vo.Team{
        Friends: friends,
        Enemies: enemies,
    }
}
