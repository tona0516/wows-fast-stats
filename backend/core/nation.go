package core

type Nation string

var Nations = []Nation{
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

func (n Nation) Priority() int {
	for i, nation := range Nations {
		if nation == n {
			return i
		}
	}

	return 999
}
