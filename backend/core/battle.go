package core

import (
	"sort"
)

type Battle struct {
	Teams []Team `json:"teams"`
}

func NewBattle(
	prefetchResult *PrefetchResult,
	tempArenaInfo TempArenaInfo,
	accountInfo WGAccountInfo,
	accountList WGAccountList,
	clans Clans,
	allPlayerShipsStats AllPlayerShipStats,
	allPlayerShipsBadges AllPlayerShipBadges,
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

		player := Player{
			PlayerInfo: PlayerInfo{
				ID:       accountID,
				Name:     nickname,
				Clan:     clan,
				IsHidden: accountInfo.Data[accountID].HiddenProfile,
				IsAlly:   vehicle.IsFriend(),
			},
			Warship: warship,
			PvPSolo: NewPlayerStats(
				StatsPatternPvPSolo,
				accountInfo.Data[accountID],
				allPlayerShipsStats[accountID],
				allPlayerShipsBadges[accountID],
				vehicle.ShipID,
				accountID,
				tempArenaInfo.Vehicles,
				prefetchResult.Warships,
			),
			PvPAll: NewPlayerStats(
				StatsPatternPvPAll,
				accountInfo.Data[accountID],
				allPlayerShipsStats[accountID],
				allPlayerShipsBadges[accountID],
				vehicle.ShipID,
				accountID,
				tempArenaInfo.Vehicles,
				prefetchResult.Warships,
			),
			RankSolo: NewPlayerStats(
				StatsPatternRankSolo,
				accountInfo.Data[accountID],
				allPlayerShipsStats[accountID],
				allPlayerShipsBadges[accountID],
				vehicle.ShipID,
				accountID,
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
		NewTeam(friends, true),
		NewTeam(enemies, false),
	}

	battle := Battle{
		Teams: teams,
	}

	return battle
}
