package core

type AccountID int

type Player struct {
	PlayerInfo PlayerInfo  `json:"player_info"`
	Warship    Warship     `json:"warship"`
	PvPSolo    PlayerStats `json:"pvp_solo"`
	PvPAll     PlayerStats `json:"pvp_all"`
	RankSolo   PlayerStats `json:"rank_solo"`
}

type PlayerInfo struct {
	ID       AccountID `json:"id"`
	Name     string    `json:"name"`
	Clan     Clan      `json:"clan"`
	IsHidden bool      `json:"is_hidden"`
}

type PlayerStats struct {
	ShipStats    ShipStats    `json:"ship"`
	OverallStats OverallStats `json:"overall"`
}

func NewPlayerStats(
	statsPattern StatsPattern,
	accountInfoData WGAccountInfoData,
	playerShipStats PlayerShipStats,
	playerShipBadges PlayerShipBadges,
	shipID ShipID,
	accountID AccountID,
	vehicles []Vehicle,
	warships Warships,
) PlayerStats {
	shipStats := NewShipStats(
		statsPattern,
		playerShipStats,
		playerShipBadges,
		shipID,
		warships,
	)

	overallStats := NewOverallStats(
		statsPattern,
		accountInfoData,
		playerShipStats,
		playerShipBadges,
		shipID,
		accountID,
		vehicles,
		warships,
	)

	return PlayerStats{
		ShipStats:    shipStats,
		OverallStats: overallStats,
	}
}
