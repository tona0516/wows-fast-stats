package data

import (
	"sort"
)

type Battle struct {
	Meta  BattleMetaData `json:"metadata"`
	Teams []Team         `json:"teams"`
}

func NewBattle(
	tempArenaInfo TempArenaInfo,
	userData *UserData,
	startupResult *PrefetchResult,
) Battle {
	friends := make(Players, 0)
	enemies := make(Players, 0)

	for _, vehicle := range tempArenaInfo.Vehicles {
		nickname := vehicle.Name
		accountID := userData.AccountList.AccountID(nickname)
		clan := userData.Clans[accountID]

		warship, ok := startupResult.Warships[vehicle.ShipID]
		if !ok {
			warship = *NewUnknownWarship()
		}

		stats := NewPersonalStats(
			vehicle.ShipID,
			userData.AccountInfo.Data[accountID],
			userData.AllPlayerShipsStats.Player(accountID),
			userData.AllPlayerShipsBadges[accountID],
			startupResult.Warships,
			tempArenaInfo,
		)

		player := Player{
			PlayerInfo: PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: userData.AccountInfo.Data[accountID].HiddenProfile,
			},
			Warship: warship,
			PvPSolo: BuildPlayerStats(
				StatsPatternPvPSolo,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo,
				startupResult.Warships,
			),
			PvPAll: BuildPlayerStats(
				StatsPatternPvPAll,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo,
				startupResult.Warships,
			),
			RankSolo: BuildPlayerStats(
				StatsPatternRankSolo,
				stats,
				accountID,
				vehicle.ShipID,
				tempArenaInfo,
				startupResult.Warships,
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
		{
			Players: friends,
			PvPAll: TeamStats{
				TeamAverageStats: CalculateTeamAverageStats(friends, StatsPatternPvPAll),
				TeamThreatLevel:  CalculateTeamThreatLevel(friends, StatsPatternPvPAll),
			},
			PvPSolo: TeamStats{
				TeamAverageStats: CalculateTeamAverageStats(friends, StatsPatternPvPSolo),
				TeamThreatLevel:  CalculateTeamThreatLevel(friends, StatsPatternPvPSolo),
			},
			RankSolo: TeamStats{
				TeamAverageStats: CalculateTeamAverageStats(friends, StatsPatternRankSolo),
				TeamThreatLevel:  CalculateTeamThreatLevel(friends, StatsPatternRankSolo),
			},
		},
		{
			Players: enemies,
			PvPAll: TeamStats{
				TeamAverageStats: CalculateTeamAverageStats(enemies, StatsPatternPvPAll),
				TeamThreatLevel:  CalculateTeamThreatLevel(enemies, StatsPatternPvPAll),
			},
			PvPSolo: TeamStats{
				TeamAverageStats: CalculateTeamAverageStats(enemies, StatsPatternPvPSolo),
				TeamThreatLevel:  CalculateTeamThreatLevel(enemies, StatsPatternPvPSolo),
			},
			RankSolo: TeamStats{
				TeamAverageStats: CalculateTeamAverageStats(enemies, StatsPatternRankSolo),
				TeamThreatLevel:  CalculateTeamThreatLevel(enemies, StatsPatternRankSolo),
			},
		},
	}

	battle := Battle{
		Meta: BattleMetaData{
			Unixtime: tempArenaInfo.Unixtime(),
			Arena:    tempArenaInfo.BattleArena(startupResult.BattleArenas),
			Type:     tempArenaInfo.BattleType(startupResult.BattleTypes),
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

type TeamStats struct {
	TeamAverageStats TeamAverageStats `json:"team_average_stats"`
	TeamThreatLevel  TeamThreatLevel  `json:"team_threat_level"`
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

type TeamThreatLevel struct {
	Average            float64 `json:"average"`
	DissociationDegree float64 `json:"dissociation_degree"`
	Accuracy           float64 `json:"accuracy"`
}
