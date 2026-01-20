package core

type TeamStats struct {
	TeamAverageStats TeamAverageStats `json:"teamAverageStats"`
	TeamThreatLevel  TeamThreatLevel  `json:"teamThreatLevel"`
}

func NewTeamStats(players Players, pattern StatsPattern) TeamStats {
	return TeamStats{
		TeamAverageStats: NewTeamAverageStats(players, pattern),
		TeamThreatLevel:  NewTeamThreatLevel(players, pattern),
	}
}
