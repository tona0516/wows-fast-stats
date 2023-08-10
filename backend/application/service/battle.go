package service

import (
	"sort"
	"sync"

	"wfs/backend/application/repository"
	"wfs/backend/application/vo"
	"wfs/backend/domain"
)

type Battle struct {
	parallels     uint
	wargaming     repository.WargamingInterface
	numbers       repository.NumbersInterface
	unregistered  repository.UnregisteredInterface
	localFile     repository.LocalFileInterface
	isFirstBattle bool

	warship       map[int]domain.Warship
	expectedStats domain.NSExpectedStats
	battleArenas  domain.WGBattleArenas
	battleTypes   domain.WGBattleTypes
}

func NewBattle(
	parallels uint,
	wargaming repository.WargamingInterface,
	localFile repository.LocalFileInterface,
	numbers repository.NumbersInterface,
	unregistered repository.UnregisteredInterface,
) *Battle {
	return &Battle{
		parallels:     parallels,
		wargaming:     wargaming,
		localFile:     localFile,
		numbers:       numbers,
		unregistered:  unregistered,
		isFirstBattle: true,
	}
}

func (b *Battle) Battle(userConfig domain.UserConfig) (domain.Battle, error) {
	b.wargaming.SetAppID(userConfig.Appid)
	var result domain.Battle

	// Fetch on-memory stored data
	warshipResult := make(chan vo.Result[map[int]domain.Warship])
	expectedStatsResult := make(chan vo.Result[domain.NSExpectedStats])
	battleArenasResult := make(chan vo.Result[domain.WGBattleArenas])
	battleTypesResult := make(chan vo.Result[domain.WGBattleTypes])
	if b.isFirstBattle {
		go b.fetchWarship(warshipResult)
		go b.fetchExpectedStats(expectedStatsResult)
		go b.fetchBattleArenas(battleArenasResult)
		go b.fetchBattleTypes(battleTypesResult)
	}

	// Get tempArenaInfo.json
	tempArenaInfo, err := b.localFile.TempArenaInfo(userConfig.InstallPath)
	if err != nil {
		return result, err
	}
	if userConfig.SaveTempArenaInfo {
		if err := b.localFile.SaveTempArenaInfo(tempArenaInfo); err != nil {
			return result, err
		}
	}

	// Get Account ID list
	accountNames := tempArenaInfo.AccountNames()
	accountList, err := b.wargaming.AccountList(accountNames)
	if err != nil {
		return result, err
	}
	accountIDs := accountList.AccountIDs()

	// Fetch each stats
	accountInfoResult := make(chan vo.Result[domain.WGAccountInfo])
	shipStatsResult := make(chan vo.Result[map[int]domain.WGShipsStats])
	clanResult := make(chan vo.Result[map[int]domain.Clan])
	go b.accountInfo(accountIDs, accountInfoResult)
	go b.shipStats(accountIDs, shipStatsResult)
	go b.clanTag(accountIDs, clanResult)

	errs := make([]error, 0)

	if b.isFirstBattle {
		warship := <-warshipResult
		b.warship = warship.Value
		errs = append(errs, warship.Error)

		expectedStats := <-expectedStatsResult
		b.expectedStats = expectedStats.Value
		errs = append(errs, expectedStats.Error)

		battleArenas := <-battleArenasResult
		b.battleArenas = battleArenas.Value
		errs = append(errs, battleArenas.Error)

		battleTypes := <-battleTypesResult
		b.battleTypes = battleTypes.Value
		errs = append(errs, battleTypes.Error)
	}

	accountInfo := <-accountInfoResult
	errs = append(errs, accountInfo.Error)

	shipStats := <-shipStatsResult
	errs = append(errs, shipStats.Error)

	clan := <-clanResult
	errs = append(errs, clan.Error)

	for _, err := range errs {
		if err != nil {
			return result, err
		}
	}

	result = b.compose(
		tempArenaInfo,
		accountInfo.Value,
		accountList,
		clan.Value,
		shipStats.Value,
		b.warship,
		b.expectedStats,
		b.battleArenas,
		b.battleTypes,
	)

	b.isFirstBattle = false

	return result, nil
}

func (b *Battle) fetchWarship(result chan vo.Result[map[int]domain.Warship]) {
	warships := make(map[int]domain.Warship, 0)

	var mu sync.Mutex
	addToResult := func(data map[int]domain.WGEncycShipsData) {
		for shipID, warship := range data {
			mu.Lock()
			warships[shipID] = domain.Warship{
				Name:   warship.Name,
				Tier:   warship.Tier,
				Type:   domain.NewShipType(warship.Type),
				Nation: domain.Nation(warship.Nation),
			}
			mu.Unlock()
		}
	}

	first := 1
	res, err := b.wargaming.EncycShips(first)
	if err != nil {
		result <- vo.Result[map[int]domain.Warship]{Error: err}
		return
	}
	addToResult(res.Data)

	pages := makeRange(first+1, res.Meta.PageTotal+1)
	err = doParallel(b.parallels, pages, func(page int) error {
		res, err := b.wargaming.EncycShips(page)
		if err != nil {
			return err
		}

		addToResult(res.Data)
		return nil
	})
	if err != nil {
		result <- vo.Result[map[int]domain.Warship]{Error: err}
		return
	}

	unregisteredShipInfo, err := b.unregistered.Warship()
	if err != nil {
		result <- vo.Result[map[int]domain.Warship]{Error: err}
		return
	}
	for k, v := range unregisteredShipInfo {
		if _, ok := warships[k]; !ok {
			warships[k] = v
		}
	}

	result <- vo.Result[map[int]domain.Warship]{Value: warships}
}

func (b *Battle) fetchExpectedStats(result chan vo.Result[domain.NSExpectedStats]) {
	expectedStats, err := b.numbers.ExpectedStats()
	result <- vo.Result[domain.NSExpectedStats]{Value: expectedStats, Error: err}
}

func (b *Battle) fetchBattleArenas(result chan vo.Result[domain.WGBattleArenas]) {
	battleArenas, err := b.wargaming.BattleArenas()
	result <- vo.Result[domain.WGBattleArenas]{Value: battleArenas, Error: err}
}

func (b *Battle) fetchBattleTypes(result chan vo.Result[domain.WGBattleTypes]) {
	battleTypes, err := b.wargaming.BattleTypes()
	result <- vo.Result[domain.WGBattleTypes]{Value: battleTypes, Error: err}
}

func (b *Battle) accountInfo(accountIDs []int, result chan vo.Result[domain.WGAccountInfo]) {
	accountInfo, err := b.wargaming.AccountInfo(accountIDs)
	result <- vo.Result[domain.WGAccountInfo]{Value: accountInfo, Error: err}
}

func (b *Battle) shipStats(accountIDs []int, result chan vo.Result[map[int]domain.WGShipsStats]) {
	shipStatsMap := make(map[int]domain.WGShipsStats)
	var mu sync.Mutex
	err := doParallel(b.parallels, accountIDs, func(accountID int) error {
		shipStats, err := b.wargaming.ShipsStats(accountID)
		if err != nil {
			return err
		}

		mu.Lock()
		shipStatsMap[accountID] = shipStats
		mu.Unlock()

		return nil
	})

	result <- vo.Result[map[int]domain.WGShipsStats]{Value: shipStatsMap, Error: err}
}

func (b *Battle) clanTag(accountIDs []int, result chan vo.Result[map[int]domain.Clan]) {
	clanMap := make(map[int]domain.Clan)

	clansAccountInfo, err := b.wargaming.ClansAccountInfo(accountIDs)
	if err != nil {
		result <- vo.Result[map[int]domain.Clan]{Error: err}

		return
	}

	clanIDs := clansAccountInfo.ClanIDs()
	clansInfo, err := b.wargaming.ClansInfo(clanIDs)
	if err != nil {
		result <- vo.Result[map[int]domain.Clan]{Error: err}

		return
	}

	for _, accountID := range accountIDs {
		clanID := clansAccountInfo.Data[accountID].ClanID
		clanTag := clansInfo.Data[clanID].Tag
		clanMap[accountID] = domain.Clan{Tag: clanTag, ID: clanID}
	}

	result <- vo.Result[map[int]domain.Clan]{Value: clanMap}
}

func (b *Battle) compose(
	tempArenaInfo domain.TempArenaInfo,
	accountInfo domain.WGAccountInfo,
	accountList domain.WGAccountList,
	clan map[int]domain.Clan,
	shipStats map[int]domain.WGShipsStats,
	warships map[int]domain.Warship,
	expectedStats domain.NSExpectedStats,
	battleArenas domain.WGBattleArenas,
	battleTypes domain.WGBattleTypes,
) domain.Battle {
	friends := make(domain.Players, 0)
	enemies := make(domain.Players, 0)

	var ownShip string

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := accountList.AccountID(nickname)
		clan := clan[accountID]

		warship, ok := warships[vehicle.ShipID]
		if !ok {
			warship = domain.Warship{
				Name:   "Unknown",
				Tier:   0,
				Type:   domain.NONE,
				Nation: "",
			}
		}
		if nickname == tempArenaInfo.PlayerName {
			ownShip = warship.Name
		}

		stats := domain.NewStats(
			vehicle.ShipID,
			accountInfo.Data[accountID],
			shipStats[accountID].Data[accountID],
			expectedStats.Data,
			nickname,
		)

		player := domain.Player{
			PlayerInfo: domain.PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: accountInfo.Data[accountID].HiddenProfile,
			},
			ShipInfo: domain.ShipInfo{
				ID:        vehicle.ShipID,
				Name:      warship.Name,
				Nation:    warship.Nation,
				Tier:      warship.Tier,
				Type:      warship.Type,
				AvgDamage: expectedStats.Data[vehicle.ShipID].AverageDamageDealt,
			},
			PvPSolo: playerStats(domain.StatsPatternPvPSolo, stats, warships),
			PvPAll:  playerStats(domain.StatsPatternPvPAll, stats, warships),
		}

		if vehicle.Relation == 0 || vehicle.Relation == 1 {
			friends = append(friends, player)
		} else {
			enemies = append(enemies, player)
		}
	}

	sort.Sort(friends)
	sort.Sort(enemies)

	teams := []domain.Team{
		{
			Players: friends,
			Name:    "味方チーム",
		},
		{
			Players: enemies,
			Name:    "敵チーム",
		},
	}

	battle := domain.Battle{
		Meta: domain.Meta{
			Unixtime: tempArenaInfo.Unixtime(),
			Arena:    tempArenaInfo.BattleArena(battleArenas),
			Type:     tempArenaInfo.BattleType(battleTypes),
			OwnShip:  ownShip,
		},
		Teams: teams,
	}

	return battle
}

func playerStats(
	statsPattern domain.StatsPattern,
	stats *domain.Stats,
	warships map[int]domain.Warship,
) domain.PlayerStats {
	return domain.PlayerStats{
		ShipStats: domain.ShipStats{
			Battles:            stats.Battles(domain.StatsCategoryShip, statsPattern),
			Damage:             stats.AvgDamage(domain.StatsCategoryShip, statsPattern),
			WinRate:            stats.WinRate(domain.StatsCategoryShip, statsPattern),
			WinSurvivedRate:    stats.WinSurvivedRate(domain.StatsCategoryShip, statsPattern),
			LoseSurvivedRate:   stats.LoseSurvivedRate(domain.StatsCategoryShip, statsPattern),
			KdRate:             stats.KdRate(domain.StatsCategoryShip, statsPattern),
			Kill:               stats.AvgKill(domain.StatsCategoryShip, statsPattern),
			Exp:                stats.AvgExp(domain.StatsCategoryShip, statsPattern),
			PR:                 stats.PR(domain.StatsCategoryShip, statsPattern),
			MainBatteryHitRate: stats.MainBatteryHitRate(domain.StatsCategoryShip, statsPattern),
			TorpedoesHitRate:   stats.TorpedoesHitRate(domain.StatsCategoryShip, statsPattern),
			PlanesKilled:       stats.PlanesKilled(domain.StatsCategoryShip, statsPattern),
		},
		OverallStats: domain.OverallStats{
			Battles:           stats.Battles(domain.StatsCategoryOverall, statsPattern),
			Damage:            stats.AvgDamage(domain.StatsCategoryOverall, statsPattern),
			WinRate:           stats.WinRate(domain.StatsCategoryOverall, statsPattern),
			WinSurvivedRate:   stats.WinSurvivedRate(domain.StatsCategoryOverall, statsPattern),
			LoseSurvivedRate:  stats.LoseSurvivedRate(domain.StatsCategoryOverall, statsPattern),
			KdRate:            stats.KdRate(domain.StatsCategoryOverall, statsPattern),
			Kill:              stats.AvgKill(domain.StatsCategoryOverall, statsPattern),
			Exp:               stats.AvgExp(domain.StatsCategoryOverall, statsPattern),
			PR:                stats.PR(domain.StatsCategoryOverall, statsPattern),
			AvgTier:           stats.AvgTier(statsPattern, warships),
			UsingShipTypeRate: stats.UsingShipTypeRate(statsPattern, warships),
			UsingTierRate:     stats.UsingTierRate(statsPattern, warships),
		},
	}
}
