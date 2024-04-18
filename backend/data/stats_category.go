package data

type StatsCategory string

//nolint:gochecknoglobals
var StatsCategories = []string{
	StatsCategoryShip,
	StatsCategoryOverall,
}

const (
	StatsCategoryShip    = "ship"
	StatsCategoryOverall = "overall"
)
