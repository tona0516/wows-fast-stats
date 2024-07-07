package vo

import (
	"errors"
	"strings"
)

var _nations = []string{
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

type Nation struct {
	ValueObject[string]
}

func NewNation(value string) (Nation, error) {
	if len(strings.TrimSpace(value)) == 0 {
		return Nation{}, errors.New("empty_nation")
	}

	return Nation{ValueObject[string]{value}}, nil
}

func (n Nation) Priority() int {
	for i, nation := range _nations {
		if nation == n.value {
			return i + 1
		}
	}

	return 999
}
