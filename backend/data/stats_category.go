package data

type StatsCategory string

const (
	StatsCategoryShip    = "ship"
	StatsCategoryOverall = "overall"
)

func StatsCategories() []string {
	return []string{
		StatsCategoryShip,
		StatsCategoryOverall,
	}
}
