package service

import (
	"context"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"
)

type Battle struct {
	localFile         repository.LocalFile
	warshipStore      repository.WarshipFetcher
	clanFetcher       repository.ClanFetcher
	rawStatFetcher    repository.RawStatFetcher
	battleMetaFetcher repository.BattleMetaFetcher
	accountFetcher    repository.AccountFetcher
	logger            repository.Logger
}

func NewBattle(
	localFile repository.LocalFile,
	warshipFetcher repository.WarshipFetcher,
	clanFetcher repository.ClanFetcher,
	rawStatFetcher repository.RawStatFetcher,
	battleMetaFetcher repository.BattleMetaFetcher,
	accountFetcher repository.AccountFetcher,
	logger repository.Logger,
) *Battle {
	return &Battle{
		localFile:         localFile,
		warshipStore:      warshipFetcher,
		clanFetcher:       clanFetcher,
		rawStatFetcher:    rawStatFetcher,
		battleMetaFetcher: battleMetaFetcher,
		accountFetcher:    accountFetcher,
		logger:            logger,
	}
}

func (b *Battle) Get(appCtx context.Context, userConfig model.UserConfigV2) (model.Battle, error) {
	var result model.Battle

	// Get tempArenaInfo.json
	tempArenaInfo, err := b.getTempArenaInfo(userConfig)
	if err != nil {
		return result, err
	}

	// persist own ign for reporting
	b.logger.SetOwnIGN(tempArenaInfo.PlayerName)

	warshipResult := make(chan model.Result[model.Warships])
	go b.fetchWarships(warshipResult)

	battleMetaResult := make(chan model.Result[model.BattleMeta])
	go b.fetchBattleMeta(battleMetaResult)

	// Get Account ID list
	accountList, err := b.accountFetcher.Fetch(tempArenaInfo.AccountNames())
	if err != nil {
		return result, err
	}
	accountIDs := make([]int, 0, len(accountList))
	for _, id := range accountList {
		accountIDs = append(accountIDs, id)
	}

	// Fetch each stats
	rawStatsResult := make(chan model.Result[model.RawStats])
	clanResult := make(chan model.Result[model.Clans])
	go b.fetchRawStats(accountIDs, rawStatsResult)
	go b.fetchClans(accountIDs, clanResult)

	errs := make([]error, 0)
	warship := <-warshipResult
	errs = append(errs, warship.Error)

	battleMeta := <-battleMetaResult
	errs = append(errs, battleMeta.Error)

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
		battleMeta.Value,
	)

	return result, nil
}

func (b *Battle) getTempArenaInfo(userConfig model.UserConfigV2) (model.TempArenaInfo, error) {
	tempArenaInfo, err := b.localFile.ReadTempArenaInfo(userConfig.InstallPath)
	if err != nil {
		return tempArenaInfo, err
	}

	if userConfig.SaveTempArenaInfo {
		path := filepath.Join("temp_arena_info", "tempArenaInfo_"+strconv.FormatInt(tempArenaInfo.Unixtime(), 10)+".json")
		if err := b.localFile.SaveTempArenaInfo(path, tempArenaInfo); err != nil {
			return tempArenaInfo, err
		}
	}

	return tempArenaInfo, nil
}

func (b *Battle) fetchWarships(channel chan model.Result[model.Warships]) {
	warships, err := b.warshipStore.Fetch()
	channel <- model.Result[model.Warships]{
		Value: warships,
		Error: err,
	}
}

func (b *Battle) fetchClans(accountIDs []int, channel chan model.Result[model.Clans]) {
	clans, err := b.clanFetcher.Fetch(accountIDs)
	channel <- model.Result[model.Clans]{
		Value: clans,
		Error: err,
	}
}

func (b *Battle) fetchRawStats(accountIDs []int, channel chan model.Result[model.RawStats]) {
	rawStats, err := b.rawStatFetcher.Fetch(accountIDs)
	channel <- model.Result[model.RawStats]{Value: rawStats, Error: err}
}

func (b *Battle) fetchBattleMeta(channel chan model.Result[model.BattleMeta]) {
	battleMeta, err := b.battleMetaFetcher.Fetch()
	channel <- model.Result[model.BattleMeta]{Value: battleMeta, Error: err}
}

func (b *Battle) compose(
	tempArenaInfo model.TempArenaInfo,
	accounts model.Accounts,
	warships model.Warships,
	clans model.Clans,
	rawStats model.RawStats,
	battleMeta model.BattleMeta,
) model.Battle {
	friends := make(model.Players, 0)
	enemies := make(model.Players, 0)
	var ownShip string

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := accounts[nickname]
		clan := clans[accountID]
		rawStat := rawStats[accountID]
		shipID := vehicle.ShipID

		warship, ok := warships[shipID]
		if !ok {
			warship = model.Warship{
				Name: "Unknown",
				Type: model.ShipTypeNONE,
			}
		}
		if nickname == tempArenaInfo.PlayerName {
			ownShip = warship.Name
		}

		stats := model.NewStats(shipID, rawStat, warships)

		player := model.Player{
			PlayerInfo: model.PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: rawStat.Overall.IsHidden,
			},
			Warship:  warship,
			PvPSolo:  playerStats(model.StatsPatternPvPSolo, stats, accountID, shipID, tempArenaInfo, warships),
			PvPAll:   playerStats(model.StatsPatternPvPAll, stats, accountID, shipID, tempArenaInfo, warships),
			RankSolo: playerStats(model.StatsPatternRankSolo, stats, accountID, shipID, tempArenaInfo, warships),
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
			Arena:    battleMeta.Arena(tempArenaInfo.MapID),
			Type:     strings.ReplaceAll(battleMeta.Type(tempArenaInfo.MatchGroup), " ", ""),
			OwnShip:  ownShip,
		},
		Teams: teams,
	}

	return battle
}

func playerStats(
	statsPattern model.StatsPattern,
	stats *model.Stats,
	accountID int,
	shipID int,
	tempArenaInfo model.TempArenaInfo,
	warships model.Warships,
) model.PlayerStats {
	threatLevel := model.CalculateThreatLevel(
		accountID,
		tempArenaInfo,
		warships,
		shipID,
		stats.Battles(model.StatsCategoryShip, statsPattern),
		stats.AvgDamage(model.StatsCategoryShip, statsPattern),
		stats.WinRate(model.StatsCategoryShip, statsPattern),
		stats.SurvivedRate(model.StatsCategoryShip, statsPattern),
		stats.PlanesKilled(model.StatsCategoryShip),
		stats.Battles(model.StatsCategoryOverall, statsPattern),
		stats.AvgDamage(model.StatsCategoryOverall, statsPattern),
		stats.WinRate(model.StatsCategoryOverall, statsPattern),
		stats.AvgKill(model.StatsCategoryOverall, statsPattern),
		stats.KdRate(model.StatsCategoryOverall, statsPattern),
	)

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
			ThreatLevel:       threatLevel,
		},
	}
}
