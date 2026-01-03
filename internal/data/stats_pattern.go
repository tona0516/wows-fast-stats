package data

type StatsPattern string

const (
	StatsPatternPvPAll   = "pvp_all"
	StatsPatternPvPSolo  = "pvp_solo"
	StatsPatternRankSolo = "rank_solo"
)

func StatsPatterns() []string {
	return []string{
		StatsPatternPvPAll,
		StatsPatternPvPSolo,
		StatsPatternRankSolo,
	}
}
