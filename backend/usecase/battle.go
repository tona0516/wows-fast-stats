package usecase

import (
	"sort"
	"sync"
	"wfs/backend/apperr"
	"wfs/backend/domain"
	"wfs/backend/repository"

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

	warship          domain.Warships
	allExpectedStats domain.ExpectedStats
	battleArenas     domain.WGBattleArenas
	battleTypes      domain.WGBattleTypes
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

func (b *Battle) Get(userConfig domain.UserConfig) (domain.Battle, error) {
	b.wargaming.SetAppID(userConfig.Appid)
	var result domain.Battle

	// Fetch on-memory stored data
	warshipResult := make(chan domain.Result[domain.Warships])
	allExpectedStatsResult := make(chan domain.Result[domain.ExpectedStats])
	battleArenasResult := make(chan domain.Result[domain.WGBattleArenas])
	battleTypesResult := make(chan domain.Result[domain.WGBattleTypes])
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
	accountInfoResult := make(chan domain.Result[domain.WGAccountInfo])
	shipStatsResult := make(chan domain.Result[domain.AllPlayerShipsStats])
	clanResult := make(chan domain.Result[domain.Clans])
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

func (b *Battle) getTempArenaInfo(userConfig domain.UserConfig) (domain.TempArenaInfo, error) {
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

func (b *Battle) fetchWarships(channel chan domain.Result[domain.Warships]) {
	warships := make(domain.Warships)
	var result domain.Result[domain.Warships]

	var mu sync.Mutex
	addToResult := func(data domain.WGEncycShips) {
		for shipID, warship := range data {
			mu.Lock()
			warships[shipID] = domain.Warship{
				Name:      warship.Name,
				Tier:      warship.Tier,
				Type:      domain.NewShipType(warship.Type),
				Nation:    domain.Nation(warship.Nation),
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

func (b *Battle) fetchExpectedStats(channel chan domain.Result[domain.ExpectedStats]) {
	var result domain.Result[domain.ExpectedStats]

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

func (b *Battle) fetchBattleArenas(channel chan domain.Result[domain.WGBattleArenas]) {
	battleArenas, err := b.wargaming.BattleArenas()
	channel <- domain.Result[domain.WGBattleArenas]{Value: battleArenas, Error: err}
}

func (b *Battle) fetchBattleTypes(channel chan domain.Result[domain.WGBattleTypes]) {
	battleTypes, err := b.wargaming.BattleTypes()
	channel <- domain.Result[domain.WGBattleTypes]{Value: battleTypes, Error: err}
}

func (b *Battle) fetchAccountInfo(accountIDs []int, channel chan domain.Result[domain.WGAccountInfo]) {
	accountInfo, err := b.wargaming.AccountInfo(accountIDs)
	channel <- domain.Result[domain.WGAccountInfo]{Value: accountInfo, Error: err}
}

func (b *Battle) fetchAllPlayerShipsStats(accountIDs []int, channel chan domain.Result[domain.AllPlayerShipsStats]) {
	shipStatsMap := make(domain.AllPlayerShipsStats)
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

	channel <- domain.Result[domain.AllPlayerShipsStats]{Value: shipStatsMap, Error: err}
}

func (b *Battle) fetchClanTag(accountIDs []int, channel chan domain.Result[domain.Clans]) {
	clans := make(domain.Clans)
	var result domain.Result[domain.Clans]

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
		clans[accountID] = domain.Clan{Tag: clanTag, ID: clanID}
	}

	result.Value = clans
	channel <- result
}

func (b *Battle) compose(
	tempArenaInfo domain.TempArenaInfo,
	accountInfo domain.WGAccountInfo,
	accountList domain.WGAccountList,
	clans domain.Clans,
	allPlayerShipsStats domain.AllPlayerShipsStats,
	warships domain.Warships,
	allExpectedStats domain.ExpectedStats,
	battleArenas domain.WGBattleArenas,
	battleTypes domain.WGBattleTypes,
) domain.Battle {
	friends := make(domain.Players, 0)
	enemies := make(domain.Players, 0)

	var ownShip string

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := accountList.AccountID(nickname)
		clan := clans[accountID]

		warship, ok := warships[vehicle.ShipID]
		if !ok {
			warship = domain.Warship{
				Name:   "Unknown",
				Tier:   0,
				Type:   domain.ShipTypeNONE,
				Nation: "",
			}
		}
		if nickname == tempArenaInfo.PlayerName {
			ownShip = warship.Name
		}

		stats := domain.NewStats(
			vehicle.ShipID,
			accountInfo[accountID],
			allPlayerShipsStats.Player(accountID),
			allExpectedStats,
			warships,
		)

		player := domain.Player{
			PlayerInfo: domain.PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: accountInfo[accountID].HiddenProfile,
			},
			ShipInfo: domain.ShipInfo{
				ID:        vehicle.ShipID,
				Name:      warship.Name,
				Nation:    warship.Nation,
				Tier:      warship.Tier,
				Type:      warship.Type,
				IsPremium: warship.IsPremium,
				AvgDamage: allExpectedStats[vehicle.ShipID].AverageDamageDealt,
			},
			PvPSolo: playerStats(domain.StatsPatternPvPSolo, stats),
			PvPAll:  playerStats(domain.StatsPatternPvPAll, stats),
		}

		if vehicle.IsFriend() {
			friends = append(friends, player)
		} else {
			enemies = append(enemies, player)
		}
	}

	sort.Sort(friends)
	sort.Sort(enemies)

	teams := []domain.Team{
		{Players: friends},
		{Players: enemies},
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
) domain.PlayerStats {
	return domain.PlayerStats{
		ShipStats: domain.ShipStats{
			Battles:            stats.Battles(domain.StatsCategoryShip, statsPattern),
			Damage:             stats.AvgDamage(domain.StatsCategoryShip, statsPattern),
			MaxDamage:          stats.MaxDamage(domain.StatsCategoryShip, statsPattern),
			WinRate:            stats.WinRate(domain.StatsCategoryShip, statsPattern),
			WinSurvivedRate:    stats.WinSurvivedRate(domain.StatsCategoryShip, statsPattern),
			LoseSurvivedRate:   stats.LoseSurvivedRate(domain.StatsCategoryShip, statsPattern),
			KdRate:             stats.KdRate(domain.StatsCategoryShip, statsPattern),
			Kill:               stats.AvgKill(domain.StatsCategoryShip, statsPattern),
			Exp:                stats.AvgExp(domain.StatsCategoryShip, statsPattern),
			PR:                 stats.PR(domain.StatsCategoryShip, statsPattern),
			MainBatteryHitRate: stats.MainBatteryHitRate(statsPattern),
			TorpedoesHitRate:   stats.TorpedoesHitRate(statsPattern),
			PlanesKilled:       stats.PlanesKilled(statsPattern),
			PlatoonRate:        stats.PlatoonRate(domain.StatsCategoryShip),
		},
		OverallStats: domain.OverallStats{
			Battles:           stats.Battles(domain.StatsCategoryOverall, statsPattern),
			Damage:            stats.AvgDamage(domain.StatsCategoryOverall, statsPattern),
			MaxDamage:         stats.MaxDamage(domain.StatsCategoryOverall, statsPattern),
			WinRate:           stats.WinRate(domain.StatsCategoryOverall, statsPattern),
			WinSurvivedRate:   stats.WinSurvivedRate(domain.StatsCategoryOverall, statsPattern),
			LoseSurvivedRate:  stats.LoseSurvivedRate(domain.StatsCategoryOverall, statsPattern),
			KdRate:            stats.KdRate(domain.StatsCategoryOverall, statsPattern),
			Kill:              stats.AvgKill(domain.StatsCategoryOverall, statsPattern),
			Exp:               stats.AvgExp(domain.StatsCategoryOverall, statsPattern),
			PR:                stats.PR(domain.StatsCategoryOverall, statsPattern),
			AvgTier:           stats.AvgTier(statsPattern),
			UsingShipTypeRate: stats.UsingShipTypeRate(statsPattern),
			UsingTierRate:     stats.UsingTierRate(statsPattern),
			PlatoonRate:       stats.PlatoonRate(domain.StatsCategoryOverall),
		},
	}
}
