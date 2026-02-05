package core

type Team struct {
	Players  Players   `json:"players"`
	PvPSolo  TeamStats `json:"pvpSolo"`
	PvPAll   TeamStats `json:"pvpAll"`
	RankSolo TeamStats `json:"rankSolo"`
	IsAlly   bool      `json:"isAlly"`
}

func NewTeam(players Players, isAlly bool) Team {
	return Team{
		Players:  players,
		PvPAll:   NewTeamStats(players, StatsPatternPvPAll),
		PvPSolo:  NewTeamStats(players, StatsPatternPvPSolo),
		RankSolo: NewTeamStats(players, StatsPatternRankSolo),
		IsAlly:   isAlly,
	}
}
