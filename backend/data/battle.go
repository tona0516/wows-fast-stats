package data

import (
	"math"
	"sort"
)

type Battle struct {
	Meta  BattleMetaData `json:"metadata"`
	Teams []Team         `json:"teams"`
}

func NewBattle(
	prefetchResult *PrefetchResult,
	tempArenaInfo TempArenaInfo,
	accountInfo WGAccountInfo,
	accountList WGAccountList,
	clans Clans,
	allPlayerShipsStats AllPlayerShipsStats,
	allPlayerShipsBadges AllPlayerShipsBadges,
) Battle {
	friends := make(Players, 0)
	enemies := make(Players, 0)

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := accountList.AccountID(nickname)
		clan := clans[accountID]

		warship, ok := prefetchResult.Warships[vehicle.ShipID]
		if !ok {
			warship = *NewUnknownWarship()
		}

		stats := NewPersonalStats(
			vehicle.ShipID,
			accountInfo.Data[accountID],
			allPlayerShipsStats.Player(accountID),
			allPlayerShipsBadges[accountID],
			prefetchResult.Warships,
			tempArenaInfo,
		)

		player := Player{
			PlayerInfo: PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: accountInfo.Data[accountID].HiddenProfile,
			},
			Warship: warship,
			PvPSolo: NewPlayerStats(
				StatsPatternPvPSolo,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo.Vehicles,
				prefetchResult.Warships,
			),
			PvPAll: NewPlayerStats(
				StatsPatternPvPAll,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo.Vehicles,
				prefetchResult.Warships,
			),
			RankSolo: NewPlayerStats(
				StatsPatternRankSolo,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo.Vehicles,
				prefetchResult.Warships,
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

	teams := []Team{
		NewTeam(friends),
		NewTeam(enemies),
	}

	battle := Battle{
		Meta: BattleMetaData{
			Unixtime: tempArenaInfo.Unixtime(),
			Arena:    tempArenaInfo.BattleArena(prefetchResult.BattleArenas),
			Type:     tempArenaInfo.BattleType(prefetchResult.BattleTypes),
		},
		Teams: teams,
	}

	return battle
}

type BattleMetaData struct {
	Unixtime int64  `json:"unixtime"`
	Arena    string `json:"arena"`
	Type     string `json:"type"`
}

type Team struct {
	Players  Players   `json:"players"`
	PvPSolo  TeamStats `json:"pvp_solo"`
	PvPAll   TeamStats `json:"pvp_all"`
	RankSolo TeamStats `json:"rank_solo"`
}

func NewTeam(players Players) Team {
	return Team{
		Players:  players,
		PvPAll:   NewTeamStats(players, StatsPatternPvPAll),
		PvPSolo:  NewTeamStats(players, StatsPatternPvPSolo),
		RankSolo: NewTeamStats(players, StatsPatternRankSolo),
	}
}

type TeamStats struct {
	TeamAverageStats TeamAverageStats `json:"team_average_stats"`
	TeamThreatLevel  TeamThreatLevel  `json:"team_threat_level"`
}

func NewTeamStats(players Players, pattern StatsPattern) TeamStats {
	return TeamStats{
		TeamAverageStats: NewTeamAverageStats(players, pattern),
		TeamThreatLevel:  NewTeamThreatLevel(players, pattern),
	}
}

type TeamAverageStats struct {
	ShipPR         float64 `json:"ship_pr"`
	ShipDamage     float64 `json:"ship_damage"`
	ShipWinRate    float64 `json:"ship_win_rate"`
	ShipBattles    uint    `json:"ship_battles"`
	OverallPR      float64 `json:"overall_pr"`
	OverallDamage  float64 `json:"overall_damage"`
	OverallWinRate float64 `json:"overall_win_rate"`
	OverallBattles uint    `json:"overall_battles"`
}

func NewTeamAverageStats(players Players, statsPattern StatsPattern) TeamAverageStats {
	var shipPRSum, shipDamageSum, shipWinRateSum float64
	var shipBattlesSum, shipStatsCount uint

	var overallPRSum, overallDamageSum, overallWinRateSum float64
	var overallBattlesSum, overallStatsCount uint

	for _, player := range players {
		var shipStats ShipStats
		var overallStats OverallStats
		switch statsPattern {
		case StatsPatternPvPSolo:
			shipStats = player.PvPSolo.ShipStats
			overallStats = player.PvPSolo.OverallStats
		case StatsPatternPvPAll:
			shipStats = player.PvPAll.ShipStats
			overallStats = player.PvPAll.OverallStats
		case StatsPatternRankSolo:
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

	return TeamAverageStats{
		ShipPR:         safeDivide(shipPRSum, shipStatsCount),
		ShipDamage:     safeDivide(shipDamageSum, shipStatsCount),
		ShipWinRate:    safeDivide(shipWinRateSum, shipStatsCount),
		ShipBattles:    uint(safeDivide(float64(shipBattlesSum), shipStatsCount)),
		OverallPR:      safeDivide(overallPRSum, overallStatsCount),
		OverallDamage:  safeDivide(overallDamageSum, overallStatsCount),
		OverallWinRate: safeDivide(overallWinRateSum, overallStatsCount),
		OverallBattles: uint(safeDivide(float64(overallBattlesSum), overallStatsCount)),
	}
}

type TeamThreatLevel struct {
	Average            float64 `json:"average"`
	DissociationDegree float64 `json:"dissociation_degree"`
	Accuracy           float64 `json:"accuracy"`
}

func NewTeamThreatLevel(players Players, statsPattern StatsPattern) TeamThreatLevel {
	if len(players) == 0 {
		return TeamThreatLevel{}
	}

	scores := make([]float64, 0, len(players))
	for _, player := range players {
		if player.PlayerInfo.ID == 0 {
			continue
		}

		if player.PlayerInfo.IsHidden {
			continue
		}

		var score float64
		switch statsPattern {
		case StatsPatternPvPSolo:
			score = player.PvPSolo.OverallStats.ThreatLevel.Modified
		case StatsPatternPvPAll:
			score = player.PvPAll.OverallStats.ThreatLevel.Modified
		case StatsPatternRankSolo:
			score = player.RankSolo.OverallStats.ThreatLevel.Modified
		default:
			continue
		}

		scores = append(scores, score)
	}

	if len(scores) == 0 {
		return TeamThreatLevel{}
	}

	maxScore := maxThreatLevelScore(scores)
	mean := geometricMeanThreatLevel(scores)

	return TeamThreatLevel{
		Average:            mean,
		DissociationDegree: (maxScore/mean - 1) * 100,
		Accuracy:           float64(len(scores)) / float64(len(players)) * 100,
	}
}

func maxThreatLevelScore(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	maxValue := values[0]
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

func geometricMeanThreatLevel(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	product := 1.0
	for _, value := range values {
		product *= value
	}
	return math.Pow(product, 1/float64(len(values)))
}
