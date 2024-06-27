package data

type Nation string

func nations() []string {
	return []string{
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
}

func (n Nation) Priority() int {
	nations := nations()
	for i, nation := range nations {
		if nation == string(n) {
			return i
		}
	}

	return 999
}
