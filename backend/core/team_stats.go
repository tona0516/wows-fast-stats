package core

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
