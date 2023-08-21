package domain

type StatsPattern string

//nolint:gochecknoglobals
var StatsPatterns = []string{
	StatsPatternPvPAll,
	StatsPatternPvPSolo,
}

const (
	StatsPatternPvPAll  = "pvp_all"
	StatsPatternPvPSolo = "pvp_solo"
)
