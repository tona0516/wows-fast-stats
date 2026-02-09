package core

type TeamStats struct {
	TeamAverageStats TeamAverageStats `json:"teamAverageStats"`
}

func NewTeamStats(players Players, pattern StatsPattern) TeamStats {
	return TeamStats{
		TeamAverageStats: NewTeamAverageStats(players, pattern),
	}
}
