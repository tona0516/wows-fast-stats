package infra

import (
	_ "embed"
	"encoding/json"
	"strings"
	"wfs/backend/domain"

	"github.com/morikuni/failure"
)

//go:embed resource/ships.json
var shipsByte []byte

type Unregistered struct{}

func NewUnregistered() *Unregistered {
	return &Unregistered{}
}

func (u *Unregistered) Warship() (map[int]domain.Warship, error) {
	var ships []unregisteredShip
	result := make(map[int]domain.Warship, 0)

	if err := json.Unmarshal(shipsByte, &ships); err != nil {
		return result, failure.Wrap(err)
	}

	for _, us := range ships {
		nation := strings.ToLower(us.Nation)
		if nation == "united_kingdom" {
			nation = "uk"
		}

		result[us.ID] = domain.Warship{
			Name:   us.En,
			Tier:   us.Level,
			Type:   domain.NewShipType(us.Species),
			Nation: domain.Nation(nation),
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
