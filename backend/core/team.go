package core

type Team struct {
	Players  Players   `json:"players"`
	PvPSolo  TeamStats `json:"pvpSolo"`
	PvPAll   TeamStats `json:"pvpAll"`
	RankSolo TeamStats `json:"rankSolo"`
}

func NewTeam(players Players) Team {
	return Team{
		Players:  players,
		PvPAll:   NewTeamStats(players, StatsPatternPvPAll),
		PvPSolo:  NewTeamStats(players, StatsPatternPvPSolo),
		RankSolo: NewTeamStats(players, StatsPatternRankSolo),
	}
}
