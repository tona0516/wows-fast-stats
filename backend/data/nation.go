package data

type Nation string

//nolint:gochecknoglobals
var (
	nations = []string{
		"japan",
		"usa",
		"ussr",
		"germany",
		"uk",
		"france",
		"italy",
		"pan_asia",
		"europe",
		"netherlands",
		"commonwealth",
		"pan_america",
		"spain",
	}
)

func (n Nation) Priority() int {
	for i, nation := range nations {
		if nation == string(n) {
			return i
		}
	}

	return 999
}
