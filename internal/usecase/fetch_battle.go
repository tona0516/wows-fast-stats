package usecase

import (
	"context"
	"sort"
	"wfs/internal/data"
	"wfs/internal/infra"
	"wfs/internal/service"
	"wfs/internal/util"
	"wfs/internal/yamibuka"

	"golang.org/x/sync/errgroup"
)

type FetchBattle struct {
	userDataFetcher    *service.UserDataFetcher
	nonUserDataFetcher *service.NonUserDataFetcher
	localStorage       infra.LocalStorage
	logger             infra.Logger
	eventsEmitFunc     eventEmitFunc
}

func NewFetchBattle(
	userDataFetcher *service.UserDataFetcher,
	nonUserDataFetcher *service.NonUserDataFetcher,
	localStorage infra.LocalStorage,
	logger infra.Logger,
	eventsEmitFunc eventEmitFunc,
) *FetchBattle {
	return &FetchBattle{
		userDataFetcher:    userDataFetcher,
		nonUserDataFetcher: nonUserDataFetcher,
		localStorage:       localStorage,
		logger:             logger,
		eventsEmitFunc:     eventsEmitFunc,
	}
}

func (b *FetchBattle) Invoke(
	ctx context.Context,
	tempArenaInfo data.TempArenaInfo,
) {
	_ = b.localStorage.SetOwnIGN(tempArenaInfo.PlayerName)
	b.logger.SetOwnIGN(tempArenaInfo.PlayerName)

	eg := errgroup.Group{}

	var userData *service.UserData
	eg.Go(func() error {
		var err error
		userData, err = b.userDataFetcher.Fetch(tempArenaInfo.AccountNames())
		return err
	})

	var nonUserData *service.NonUserData
	eg.Go(func() error {
		var err error
		nonUserData, err = b.nonUserDataFetcher.Fetch()
		return err
	})

	if err := eg.Wait(); err != nil {
		b.eventsEmitFunc(ctx, EventErr, err)
		return
	}

	result := b.compose(
		tempArenaInfo,
		userData,
		nonUserData,
	)

	b.eventsEmitFunc(ctx, EventFetchDone, result)
}

func (b *FetchBattle) compose(
	tempArenaInfo data.TempArenaInfo,
	userData *service.UserData,
	nonUserData *service.NonUserData,
) data.Battle {
	friends := make(data.Players, 0)
	enemies := make(data.Players, 0)

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := userData.AccountList.AccountID(nickname)
		clan := userData.Clans[accountID]

		warship, ok := nonUserData.Warships[vehicle.ShipID]
		if !ok {
			warship = *data.NewUnknownWarship()
		}

		stats := data.NewPersonalStats(
			vehicle.ShipID,
			userData.AccountInfo.Data[accountID],
			userData.AllPlayerShipsStats.Player(accountID),
			userData.AllPlayerShipsBadges[accountID],
			nonUserData.Warships,
			tempArenaInfo,
		)

		player := data.Player{
			PlayerInfo: data.PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: userData.AccountInfo.Data[accountID].HiddenProfile,
			},
			Warship: warship,
			PvPSolo: playerStats(
				data.StatsPatternPvPSolo,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo,
				nonUserData.Warships,
			),
			PvPAll: playerStats(
				data.StatsPatternPvPAll,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo,
				nonUserData.Warships,
			),
			RankSolo: playerStats(
				data.StatsPatternRankSolo,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo,
				nonUserData.Warships,
			),
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
		{
			Players: friends,
			PvPAll: data.TeamStats{
				TeamAverageStats: calculateTeamAverageStats(friends, data.StatsPatternPvPAll),
				TeamThreatLevel:  yamibuka.CalculateTeamThreatLevel(friends, data.StatsPatternPvPAll),
			},
			PvPSolo: data.TeamStats{
				TeamAverageStats: calculateTeamAverageStats(friends, data.StatsPatternPvPSolo),
				TeamThreatLevel:  yamibuka.CalculateTeamThreatLevel(friends, data.StatsPatternPvPSolo),
			},
			RankSolo: data.TeamStats{
				TeamAverageStats: calculateTeamAverageStats(friends, data.StatsPatternRankSolo),
				TeamThreatLevel:  yamibuka.CalculateTeamThreatLevel(friends, data.StatsPatternRankSolo),
			},
		},
		{
			Players: enemies,
			PvPAll: data.TeamStats{
				TeamAverageStats: calculateTeamAverageStats(enemies, data.StatsPatternPvPAll),
				TeamThreatLevel:  yamibuka.CalculateTeamThreatLevel(enemies, data.StatsPatternPvPAll),
			},
			PvPSolo: data.TeamStats{
				TeamAverageStats: calculateTeamAverageStats(enemies, data.StatsPatternPvPSolo),
				TeamThreatLevel:  yamibuka.CalculateTeamThreatLevel(enemies, data.StatsPatternPvPSolo),
			},
			RankSolo: data.TeamStats{
				TeamAverageStats: calculateTeamAverageStats(enemies, data.StatsPatternRankSolo),
				TeamThreatLevel:  yamibuka.CalculateTeamThreatLevel(enemies, data.StatsPatternRankSolo),
			},
		},
	}

	battle := data.Battle{
		Meta: data.BattleMetaData{
			Unixtime: tempArenaInfo.Unixtime(),
			Arena:    tempArenaInfo.BattleArena(nonUserData.BattleArenas),
			Type:     tempArenaInfo.BattleType(nonUserData.BattleTypes),
		},
		Teams: teams,
	}

	return battle
}

func playerStats(
	statsPattern data.StatsPattern,
	stats *data.PersonalStats,
	accountID int,
	shipID int,
	tempArenaInfo data.TempArenaInfo,
	warships data.Warships,
) data.PlayerStats {
	threatLevel := yamibuka.CalculateThreatLevel(yamibuka.NewThreatLevelFactor(
		accountID,
		tempArenaInfo,
		warships,
		shipID,
		stats.Battles(data.StatsCategoryShip, statsPattern),
		stats.AvgDamage(data.StatsCategoryShip, statsPattern).Value,
		stats.WinRate(data.StatsCategoryShip, statsPattern).Value,
		stats.SurvivedRate(data.StatsCategoryShip, statsPattern).All,
		stats.PlanesKilled(data.StatsCategoryShip),
		stats.Battles(data.StatsCategoryOverall, statsPattern),
		stats.AvgDamage(data.StatsCategoryOverall, statsPattern).Value,
		stats.WinRate(data.StatsCategoryOverall, statsPattern).Value,
		stats.AvgKill(data.StatsCategoryOverall, statsPattern),
		stats.KdRate(data.StatsCategoryOverall, statsPattern),
	))

	return data.PlayerStats{
		ShipStats: data.ShipStats{
			Battles:         stats.Battles(data.StatsCategoryShip, statsPattern),
			Damage:          stats.AvgDamage(data.StatsCategoryShip, statsPattern),
			MaxDamage:       stats.MaxDamage(data.StatsCategoryShip, statsPattern),
			WinRate:         stats.WinRate(data.StatsCategoryShip, statsPattern),
			SurvivedRate:    stats.SurvivedRate(data.StatsCategoryShip, statsPattern),
			KdRate:          stats.KdRate(data.StatsCategoryShip, statsPattern),
			Kill:            stats.AvgKill(data.StatsCategoryShip, statsPattern),
			Exp:             stats.AvgExp(data.StatsCategoryShip, statsPattern),
			PR:              stats.PR(data.StatsCategoryShip, statsPattern),
			HitRate:         stats.HitRate(statsPattern),
			PlanesKilled:    stats.PlanesKilled(statsPattern),
			PlatoonRate:     stats.PlatoonRate(data.StatsCategoryShip),
			EfficiencyBadge: stats.EfficiencyBadge(),
		},
		OverallStats: data.OverallStats{
			Battles:           stats.Battles(data.StatsCategoryOverall, statsPattern),
			Damage:            stats.AvgDamage(data.StatsCategoryOverall, statsPattern),
			MaxDamage:         stats.MaxDamage(data.StatsCategoryOverall, statsPattern),
			WinRate:           stats.WinRate(data.StatsCategoryOverall, statsPattern),
			SurvivedRate:      stats.SurvivedRate(data.StatsCategoryOverall, statsPattern),
			KdRate:            stats.KdRate(data.StatsCategoryOverall, statsPattern),
			Kill:              stats.AvgKill(data.StatsCategoryOverall, statsPattern),
			Exp:               stats.AvgExp(data.StatsCategoryOverall, statsPattern),
			PR:                stats.PR(data.StatsCategoryOverall, statsPattern),
			AvgTier:           stats.AvgTier(statsPattern),
			UsingShipTypeRate: stats.UsingShipTypeRate(statsPattern),
			UsingTierRate:     stats.UsingTierRate(statsPattern),
			PlatoonRate:       stats.PlatoonRate(data.StatsCategoryOverall),
			EfficiencyBadge:   stats.EfficiencyBadges(),
			ThreatLevel:       threatLevel,
		},
	}
}

func calculateTeamAverageStats(
	players data.Players,
	statsPattern data.StatsPattern,
) data.TeamAverageStats {
	var shipPRSum, shipDamageSum, shipWinRateSum float64
	var shipBattlesSum, shipStatsCount uint

	var overallPRSum, overallDamageSum, overallWinRateSum float64
	var overallBattlesSum, overallStatsCount uint

	for _, player := range players {
		var shipStats data.ShipStats
		var overallStats data.OverallStats
		switch statsPattern {
		case data.StatsPatternPvPSolo:
			shipStats = player.PvPSolo.ShipStats
			overallStats = player.PvPSolo.OverallStats
		case data.StatsPatternPvPAll:
			shipStats = player.PvPAll.ShipStats
			overallStats = player.PvPAll.OverallStats
		case data.StatsPatternRankSolo:
			shipStats = player.RankSolo.ShipStats
			overallStats = player.RankSolo.OverallStats
		}

		if shipStats.Battles > 0 {
			shipPRSum += shipStats.PR.Value
			shipDamageSum += shipStats.Damage.Value
			shipWinRateSum += shipStats.WinRate.Value
			shipBattlesSum += shipStats.Battles

			shipStatsCount++
		}

		if overallStats.Battles > 0 {
			overallPRSum += overallStats.PR.Value
			overallDamageSum += overallStats.Damage.Value
			overallWinRateSum += overallStats.WinRate.Value
			overallBattlesSum += overallStats.Battles

			overallStatsCount++
		}
	}

	return data.TeamAverageStats{
		ShipPR:         util.SafeDivide(shipPRSum, shipStatsCount),
		ShipDamage:     util.SafeDivide(shipDamageSum, shipStatsCount),
		ShipWinRate:    util.SafeDivide(shipWinRateSum, shipStatsCount),
		ShipBattles:    uint(util.SafeDivide(float64(shipBattlesSum), shipStatsCount)),
		OverallPR:      util.SafeDivide(overallPRSum, overallStatsCount),
		OverallDamage:  util.SafeDivide(overallDamageSum, overallStatsCount),
		OverallWinRate: util.SafeDivide(overallWinRateSum, overallStatsCount),
		OverallBattles: uint(util.SafeDivide(float64(overallBattlesSum), overallStatsCount)),
	}
}
