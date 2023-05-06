package infra

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/morikuni/failure"
)

//go:embed resource/ships.json
var shipsByte []byte

type Unregistered struct{}

func (u *Unregistered) Warship() (map[int]vo.Warship, error) {
	errCode := apperr.UnregisteredWarship

	var ships []unregisteredShip
	result := make(map[int]vo.Warship, 0)

	if err := json.Unmarshal(shipsByte, &ships); err != nil {
		return result, failure.Translate(err, errCode)
	}

	for _, us := range ships {
		nation := strings.ToLower(us.Nation)
		if nation == "united_kingdom" {
			nation = "uk"
		}

		result[us.ID] = vo.Warship{
			Name:   us.En,
			Tier:   us.Level,
			Type:   us.Species,
			Nation: nation,
		}
	}

	return result, nil
}

type unregisteredShip struct {
	ID      int    `json:"id"`
	En      string `json:"en"`
	Level   uint   `json:"level"`
	Nation  string `json:"nation"`
	Species string `json:"Species"`
}
