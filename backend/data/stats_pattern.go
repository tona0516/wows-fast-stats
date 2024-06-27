package data

type StatsPattern string

const (
	StatsPatternPvPAll  = "pvp_all"
	StatsPatternPvPSolo = "pvp_solo"
)

func StatsPatterns() []string {
	return []string{
		StatsPatternPvPAll,
		StatsPatternPvPSolo,
	}
}
