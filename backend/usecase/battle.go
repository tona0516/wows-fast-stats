package usecase

import (
	"sort"
	"sync"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"

	"github.com/morikuni/failure"
)

type Battle struct {
	parallels     uint
	wargaming     repository.WargamingInterface
	numbers       repository.NumbersInterface
	unregistered  repository.UnregisteredInterface
	localFile     repository.LocalFileInterface
	storage       repository.StorageInterface
	logger        repository.LoggerInterface
	isFirstBattle bool

	warship          model.Warships
	allExpectedStats model.ExpectedStats
	battleArenas     model.WGBattleArenas
	battleTypes      model.WGBattleTypes
}

func NewBattle(
	parallels uint,
	wargaming repository.WargamingInterface,
	localFile repository.LocalFileInterface,
	numbers repository.NumbersInterface,
	unregistered repository.UnregisteredInterface,
	storage repository.StorageInterface,
	logger repository.LoggerInterface,
) *Battle {
	return &Battle{
		parallels:     parallels,
		wargaming:     wargaming,
		localFile:     localFile,
		numbers:       numbers,
		unregistered:  unregistered,
		storage:       storage,
		logger:        logger,
		isFirstBattle: true,
	}
}

func (b *Battle) Get(userConfig model.UserConfig) (model.Battle, error) {
	b.wargaming.SetAppID(userConfig.Appid)
	var result model.Battle

	// Fetch on-memory stored data
	warshipResult := make(chan model.Result[model.Warships])
	allExpectedStatsResult := make(chan model.Result[model.ExpectedStats])
	battleArenasResult := make(chan model.Result[model.WGBattleArenas])
	battleTypesResult := make(chan model.Result[model.WGBattleTypes])
	if b.isFirstBattle {
		go b.fetchWarships(warshipResult)
		go b.fetchExpectedStats(allExpectedStatsResult)
		go b.fetchBattleArenas(battleArenasResult)
		go b.fetchBattleTypes(battleTypesResult)
	}

	// Get tempArenaInfo.json
	tempArenaInfo, err := b.getTempArenaInfo(userConfig)
	if err != nil {
		return result, err
	}

	// persist own ign for reporting
	_ = b.storage.WriteOwnIGN(tempArenaInfo.PlayerName)

	// Get Account ID list
	accountList, err := b.wargaming.AccountList(tempArenaInfo.AccountNames())
	if err != nil {
		return result, err
	}
	accountIDs := accountList.AccountIDs()

	// Fetch each stats
	accountInfoResult := make(chan model.Result[model.WGAccountInfo])
	shipStatsResult := make(chan model.Result[model.AllPlayerShipsStats])
	clanResult := make(chan model.Result[model.Clans])
	go b.fetchAccountInfo(accountIDs, accountInfoResult)
	go b.fetchAllPlayerShipsStats(accountIDs, shipStatsResult)
	go b.fetchClanTag(accountIDs, clanResult)

	errs := make([]error, 0)

	if b.isFirstBattle {
		warship := <-warshipResult
		b.warship = warship.Value
		errs = append(errs, warship.Error)

		expectedStats := <-allExpectedStatsResult
		b.allExpectedStats = expectedStats.Value
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
		b.allExpectedStats,
		b.battleArenas,
		b.battleTypes,
	)

	b.isFirstBattle = false

	return result, nil
}

func (b *Battle) getTempArenaInfo(userConfig model.UserConfig) (model.TempArenaInfo, error) {
	tempArenaInfo, err := b.localFile.TempArenaInfo(userConfig.InstallPath)
	if err != nil {
		return tempArenaInfo, err
	}

	if userConfig.SaveTempArenaInfo {
		if err := b.localFile.SaveTempArenaInfo(tempArenaInfo); err != nil {
			return tempArenaInfo, err
		}
	}

	return tempArenaInfo, nil
}

func (b *Battle) fetchWarships(channel chan model.Result[model.Warships]) {
	warships := make(model.Warships)
	var result model.Result[model.Warships]

	var mu sync.Mutex
	addToResult := func(data model.WGEncycShips) {
		for shipID, warship := range data {
			mu.Lock()
			warships[shipID] = model.Warship{
				Name:      warship.Name,
				Tier:      warship.Tier,
				Type:      model.NewShipType(warship.Type),
				Nation:    model.Nation(warship.Nation),
				IsPremium: warship.IsPremium,
			}
			mu.Unlock()
		}
	}

	first := 1
	encycShips, pageTotal, err := b.wargaming.EncycShips(first)
	if err != nil {
		result.Error = err
		channel <- result
		return
	}
	addToResult(encycShips)

	pages := makeRange(first+1, pageTotal+1)
	err = doParallel(b.parallels, pages, func(page int) error {
		res, _, err := b.wargaming.EncycShips(page)
		if err != nil {
			return err
		}

		addToResult(res)
		return nil
	})
	if err != nil {
		result.Error = err
		channel <- result
		return
	}

	unregisteredShipInfo, err := b.unregistered.Warship()
	if err != nil {
		result.Error = err
		channel <- result
		return
	}
	for k, v := range unregisteredShipInfo {
		if _, ok := warships[k]; !ok {
			warships[k] = v
		}
	}

	result.Value = warships
	channel <- result
}

func (b *Battle) fetchExpectedStats(channel chan model.Result[model.ExpectedStats]) {
	var result model.Result[model.ExpectedStats]

	expectedStats, errFetch := b.numbers.ExpectedStats()
	if errFetch == nil {
		_ = b.storage.WriteExpectedStats(expectedStats)
		result.Value = expectedStats
		channel <- result
		return
	}

	b.logger.Warn(failure.New(apperr.FailSafeProccess), nil)

	expectedStats, errCache := b.storage.ExpectedStats()
	if errCache == nil {
		result.Value = expectedStats
		channel <- result
		return
	}

	result.Error = errFetch
	channel <- result
}

func (b *Battle) fetchBattleArenas(channel chan model.Result[model.WGBattleArenas]) {
	battleArenas, err := b.wargaming.BattleArenas()
	channel <- model.Result[model.WGBattleArenas]{Value: battleArenas, Error: err}
}

func (b *Battle) fetchBattleTypes(channel chan model.Result[model.WGBattleTypes]) {
	battleTypes, err := b.wargaming.BattleTypes()
	channel <- model.Result[model.WGBattleTypes]{Value: battleTypes, Error: err}
}

func (b *Battle) fetchAccountInfo(accountIDs []int, channel chan model.Result[model.WGAccountInfo]) {
	accountInfo, err := b.wargaming.AccountInfo(accountIDs)
	channel <- model.Result[model.WGAccountInfo]{Value: accountInfo, Error: err}
}

func (b *Battle) fetchAllPlayerShipsStats(accountIDs []int, channel chan model.Result[model.AllPlayerShipsStats]) {
	shipStatsMap := make(model.AllPlayerShipsStats)
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

	channel <- model.Result[model.AllPlayerShipsStats]{Value: shipStatsMap, Error: err}
}

func (b *Battle) fetchClanTag(accountIDs []int, channel chan model.Result[model.Clans]) {
	clans := make(model.Clans)
	var result model.Result[model.Clans]

	clansAccountInfo, err := b.wargaming.ClansAccountInfo(accountIDs)
	if err != nil {
		result.Error = err
		channel <- result
		return
	}

	clanIDs := clansAccountInfo.ClanIDs()
	clansInfo, err := b.wargaming.ClansInfo(clanIDs)
	if err != nil {
		result.Error = err
		channel <- result
		return
	}

	for _, accountID := range accountIDs {
		clanID := clansAccountInfo[accountID].ClanID
		clanTag := clansInfo[clanID].Tag
		clans[accountID] = model.Clan{Tag: clanTag, ID: clanID}
	}

	result.Value = clans
	channel <- result
}

func (b *Battle) compose(
	tempArenaInfo model.TempArenaInfo,
	accountInfo model.WGAccountInfo,
	accountList model.WGAccountList,
	clans model.Clans,
	allPlayerShipsStats model.AllPlayerShipsStats,
	warships model.Warships,
	allExpectedStats model.ExpectedStats,
	battleArenas model.WGBattleArenas,
	battleTypes model.WGBattleTypes,
) model.Battle {
	friends := make(model.Players, 0)
	enemies := make(model.Players, 0)

	var ownShip string

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := accountList.AccountID(nickname)
		clan := clans[accountID]

		warship, ok := warships[vehicle.ShipID]
		if !ok {
			warship = model.Warship{
				Name:   "Unknown",
				Tier:   0,
				Type:   model.ShipTypeNONE,
				Nation: "",
			}
		}
		if nickname == tempArenaInfo.PlayerName {
			ownShip = warship.Name
		}

		stats := model.NewStats(
			vehicle.ShipID,
			accountInfo[accountID],
			allPlayerShipsStats.Player(accountID),
			allExpectedStats,
			warships,
		)

		player := model.Player{
			PlayerInfo: model.PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: accountInfo[accountID].HiddenProfile,
			},
			ShipInfo: model.ShipInfo{
				ID:        vehicle.ShipID,
				Name:      warship.Name,
				Nation:    warship.Nation,
				Tier:      warship.Tier,
				Type:      warship.Type,
				IsPremium: warship.IsPremium,
				AvgDamage: allExpectedStats[vehicle.ShipID].AverageDamageDealt,
			},
			PvPSolo: playerStats(model.StatsPatternPvPSolo, stats),
			PvPAll:  playerStats(model.StatsPatternPvPAll, stats),
		}

		if vehicle.IsFriend() {
			friends = append(friends, player)
		} else {
			enemies = append(enemies, player)
		}
	}

	sort.Sort(friends)
	sort.Sort(enemies)

	teams := []model.Team{
		{Players: friends},
		{Players: enemies},
	}

	battle := model.Battle{
		Meta: model.Meta{
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
	statsPattern model.StatsPattern,
	stats *model.Stats,
) model.PlayerStats {
	return model.PlayerStats{
		ShipStats: model.ShipStats{
			Battles:            stats.Battles(model.StatsCategoryShip, statsPattern),
			Damage:             stats.AvgDamage(model.StatsCategoryShip, statsPattern),
			MaxDamage:          stats.MaxDamage(model.StatsCategoryShip, statsPattern),
			WinRate:            stats.WinRate(model.StatsCategoryShip, statsPattern),
			WinSurvivedRate:    stats.WinSurvivedRate(model.StatsCategoryShip, statsPattern),
			LoseSurvivedRate:   stats.LoseSurvivedRate(model.StatsCategoryShip, statsPattern),
			KdRate:             stats.KdRate(model.StatsCategoryShip, statsPattern),
			Kill:               stats.AvgKill(model.StatsCategoryShip, statsPattern),
			Exp:                stats.AvgExp(model.StatsCategoryShip, statsPattern),
			PR:                 stats.PR(model.StatsCategoryShip, statsPattern),
			MainBatteryHitRate: stats.MainBatteryHitRate(statsPattern),
			TorpedoesHitRate:   stats.TorpedoesHitRate(statsPattern),
			PlanesKilled:       stats.PlanesKilled(statsPattern),
			PlatoonRate:        stats.PlatoonRate(model.StatsCategoryShip),
		},
		OverallStats: model.OverallStats{
			Battles:           stats.Battles(model.StatsCategoryOverall, statsPattern),
			Damage:            stats.AvgDamage(model.StatsCategoryOverall, statsPattern),
			MaxDamage:         stats.MaxDamage(model.StatsCategoryOverall, statsPattern),
			WinRate:           stats.WinRate(model.StatsCategoryOverall, statsPattern),
			WinSurvivedRate:   stats.WinSurvivedRate(model.StatsCategoryOverall, statsPattern),
			LoseSurvivedRate:  stats.LoseSurvivedRate(model.StatsCategoryOverall, statsPattern),
			KdRate:            stats.KdRate(model.StatsCategoryOverall, statsPattern),
			Kill:              stats.AvgKill(model.StatsCategoryOverall, statsPattern),
			Exp:               stats.AvgExp(model.StatsCategoryOverall, statsPattern),
			PR:                stats.PR(model.StatsCategoryOverall, statsPattern),
			AvgTier:           stats.AvgTier(statsPattern),
			UsingShipTypeRate: stats.UsingShipTypeRate(statsPattern),
			UsingTierRate:     stats.UsingTierRate(statsPattern),
			PlatoonRate:       stats.PlatoonRate(model.StatsCategoryOverall),
		},
	}
}
