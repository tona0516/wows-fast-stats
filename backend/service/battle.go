package service

import (
	"context"
	"sort"
	"strings"
	"wfs/backend/data"
	"wfs/backend/domain/model"
	domainRepository "wfs/backend/domain/repository"
	"wfs/backend/repository"
)

type Battle struct {
	wargaming      repository.WargamingInterface
	localFile      repository.LocalFileInterface
	warshipFetcher domainRepository.WarshipFetcherInterface
	clanFetcher    domainRepository.ClanFetcherInterface
	taiFetcher     domainRepository.TAIFetcherInterface
	rawStatFetcher domainRepository.RawStatFetcherInterface
	storage        repository.StorageInterface
	logger         repository.LoggerInterface
	eventsEmitFunc eventEmitFunc

	isFirstBattle bool
	battleArenas  data.WGBattleArenas
	battleTypes   data.WGBattleTypes
}

func NewBattle(
	wargaming repository.WargamingInterface,
	localFile repository.LocalFileInterface,
	warshipFetcher domainRepository.WarshipFetcherInterface,
	clanFetcher domainRepository.ClanFetcherInterface,
	taiFetcher domainRepository.TAIFetcherInterface,
	rawStatFetcher domainRepository.RawStatFetcherInterface,
	storage repository.StorageInterface,
	logger repository.LoggerInterface,
	eventsEmitFunc eventEmitFunc,
) *Battle {
	return &Battle{
		wargaming:      wargaming,
		localFile:      localFile,
		warshipFetcher: warshipFetcher,
		clanFetcher:    clanFetcher,
		taiFetcher:     taiFetcher,
		rawStatFetcher: rawStatFetcher,
		storage:        storage,
		logger:         logger,
		eventsEmitFunc: eventsEmitFunc,
		isFirstBattle:  true,
	}
}

func (b *Battle) Get(appCtx context.Context, userConfig data.UserConfigV2) (data.Battle, error) {
	var result data.Battle

	// Fetch on-memory stored data
	battleArenasResult := make(chan data.Result[data.WGBattleArenas])
	battleTypesResult := make(chan data.Result[data.WGBattleTypes])
	if b.isFirstBattle {
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
	b.logger.SetOwnIGN(tempArenaInfo.PlayerName)

	warshipResult := make(chan data.Result[model.Warships])
	go b.fetchWarships(warshipResult)

	// Get Account ID list
	accountList, err := b.wargaming.AccountList(tempArenaInfo.AccountNames())
	if err != nil {
		return result, err
	}
	accountIDs := accountList.AccountIDs()

	// Fetch each stats
	rawStatsResult := make(chan data.Result[model.RawStats])
	clanResult := make(chan data.Result[model.Clans])
	go b.fetchRawStats(accountIDs, rawStatsResult)
	go b.fetchClans(accountIDs, clanResult)

	errs := make([]error, 0)

	if b.isFirstBattle {
		battleArenas := <-battleArenasResult
		b.battleArenas = battleArenas.Value
		errs = append(errs, battleArenas.Error)

		battleTypes := <-battleTypesResult
		b.battleTypes = battleTypes.Value
		errs = append(errs, battleTypes.Error)
	}

	warship := <-warshipResult
	errs = append(errs, warship.Error)

	clan := <-clanResult
	errs = append(errs, clan.Error)

	rawStats := <-rawStatsResult
	errs = append(errs, rawStats.Error)

	for _, err := range errs {
		if err != nil {
			return result, err
		}
	}

	result = b.compose(
		tempArenaInfo,
		accountList,
		warship.Value,
		clan.Value,
		rawStats.Value,
		b.battleArenas,
		b.battleTypes,
	)

	b.isFirstBattle = false

	return result, nil
}

func (b *Battle) getTempArenaInfo(userConfig data.UserConfigV2) (model.TempArenaInfo, error) {
	tempArenaInfo, err := b.taiFetcher.Get(userConfig.InstallPath)
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

func (b *Battle) fetchWarships(channel chan data.Result[model.Warships]) {
	warships, err := b.warshipFetcher.Fetch()
	channel <- data.Result[model.Warships]{
		Value: warships,
		Error: err,
	}
}

func (b *Battle) fetchClans(accountIDs []int, channel chan data.Result[model.Clans]) {
	clans, err := b.clanFetcher.Fetch(accountIDs)
	channel <- data.Result[model.Clans]{
		Value: clans,
		Error: err,
	}
}

func (b *Battle) fetchBattleArenas(channel chan data.Result[data.WGBattleArenas]) {
	battleArenas, err := b.wargaming.BattleArenas()
	channel <- data.Result[data.WGBattleArenas]{Value: battleArenas, Error: err}
}

func (b *Battle) fetchBattleTypes(channel chan data.Result[data.WGBattleTypes]) {
	battleTypes, err := b.wargaming.BattleTypes()
	channel <- data.Result[data.WGBattleTypes]{Value: battleTypes, Error: err}
}

func (b *Battle) fetchRawStats(accountIDs []int, channel chan data.Result[model.RawStats]) {
	rawStats, err := b.rawStatFetcher.Fetch(accountIDs)
	channel <- data.Result[model.RawStats]{Value: rawStats, Error: err}
}

func (b *Battle) compose(
	tempArenaInfo model.TempArenaInfo,
	accountList data.WGAccountList,
	warships model.Warships,
	clans model.Clans,
	rawStats model.RawStats,
	battleArenas data.WGBattleArenas,
	battleTypes data.WGBattleTypes,
) data.Battle {
	friends := make(data.Players, 0)
	enemies := make(data.Players, 0)

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

		stats := data.NewStats(
			vehicle.ShipID,
			rawStats[accountID],
			warships,
		)

		player := data.Player{
			PlayerInfo: data.PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: rawStats[accountID].Overall.IsHidden,
			},
			Warship:  warship,
			PvPSolo:  playerStats(data.StatsPatternPvPSolo, stats, accountID, vehicle.ShipID, tempArenaInfo, warships),
			PvPAll:   playerStats(data.StatsPatternPvPAll, stats, accountID, vehicle.ShipID, tempArenaInfo, warships),
			RankSolo: playerStats(data.StatsPatternRankSolo, stats, accountID, vehicle.ShipID, tempArenaInfo, warships),
		}

		if vehicle.IsFriend() {
			friends = append(friends, player)
		} else {
			enemies = append(enemies, player)
		}
	}

	sort.Sort(friends)
	sort.Sort(enemies)

	teams := []data.Team{
		{Players: friends},
		{Players: enemies},
	}

	battle := data.Battle{
		Meta: data.Meta{
			Unixtime: tempArenaInfo.Unixtime(),
			Arena:    battleArenas[tempArenaInfo.MapID].Name,
			Type:     strings.ReplaceAll(battleTypes[strings.ToUpper(tempArenaInfo.MatchGroup)].Name, " ", ""),
			OwnShip:  ownShip,
		},
		Teams: teams,
	}

	return battle
}

func playerStats(
	statsPattern data.StatsPattern,
	stats *data.Stats,
	accountID int,
	shipID int,
	tempArenaInfo model.TempArenaInfo,
	warships model.Warships,
) data.PlayerStats {
	threatLevel := model.CalculateThreatLevel(
		accountID,
		tempArenaInfo,
		warships,
		shipID,
		stats.Battles(data.StatsCategoryShip, statsPattern),
		stats.AvgDamage(data.StatsCategoryShip, statsPattern),
		stats.WinRate(data.StatsCategoryShip, statsPattern),
		stats.SurvivedRate(data.StatsCategoryShip, statsPattern),
		stats.PlanesKilled(data.StatsCategoryShip),
		stats.Battles(data.StatsCategoryOverall, statsPattern),
		stats.AvgDamage(data.StatsCategoryOverall, statsPattern),
		stats.WinRate(data.StatsCategoryOverall, statsPattern),
		stats.AvgKill(data.StatsCategoryOverall, statsPattern),
		stats.KdRate(data.StatsCategoryOverall, statsPattern),
	)

	return data.PlayerStats{
		ShipStats: data.ShipStats{
			Battles:            stats.Battles(data.StatsCategoryShip, statsPattern),
			Damage:             stats.AvgDamage(data.StatsCategoryShip, statsPattern),
			MaxDamage:          stats.MaxDamage(data.StatsCategoryShip, statsPattern),
			WinRate:            stats.WinRate(data.StatsCategoryShip, statsPattern),
			WinSurvivedRate:    stats.WinSurvivedRate(data.StatsCategoryShip, statsPattern),
			LoseSurvivedRate:   stats.LoseSurvivedRate(data.StatsCategoryShip, statsPattern),
			KdRate:             stats.KdRate(data.StatsCategoryShip, statsPattern),
			Kill:               stats.AvgKill(data.StatsCategoryShip, statsPattern),
			Exp:                stats.AvgExp(data.StatsCategoryShip, statsPattern),
			PR:                 stats.PR(data.StatsCategoryShip, statsPattern),
			MainBatteryHitRate: stats.MainBatteryHitRate(statsPattern),
			TorpedoesHitRate:   stats.TorpedoesHitRate(statsPattern),
			PlanesKilled:       stats.PlanesKilled(statsPattern),
			PlatoonRate:        stats.PlatoonRate(data.StatsCategoryShip),
		},
		OverallStats: data.OverallStats{
			Battles:           stats.Battles(data.StatsCategoryOverall, statsPattern),
			Damage:            stats.AvgDamage(data.StatsCategoryOverall, statsPattern),
			MaxDamage:         stats.MaxDamage(data.StatsCategoryOverall, statsPattern),
			WinRate:           stats.WinRate(data.StatsCategoryOverall, statsPattern),
			WinSurvivedRate:   stats.WinSurvivedRate(data.StatsCategoryOverall, statsPattern),
			LoseSurvivedRate:  stats.LoseSurvivedRate(data.StatsCategoryOverall, statsPattern),
			KdRate:            stats.KdRate(data.StatsCategoryOverall, statsPattern),
			Kill:              stats.AvgKill(data.StatsCategoryOverall, statsPattern),
			Exp:               stats.AvgExp(data.StatsCategoryOverall, statsPattern),
			PR:                stats.PR(data.StatsCategoryOverall, statsPattern),
			AvgTier:           stats.AvgTier(statsPattern),
			UsingShipTypeRate: stats.UsingShipTypeRate(statsPattern),
			UsingTierRate:     stats.UsingTierRate(statsPattern),
			PlatoonRate:       stats.PlatoonRate(data.StatsCategoryOverall),
			ThreatLevel:       threatLevel,
		},
	}
}
