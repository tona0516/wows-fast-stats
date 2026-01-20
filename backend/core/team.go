package core

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
