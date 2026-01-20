package core

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
