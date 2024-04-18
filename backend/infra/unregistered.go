package infra

import (
	_ "embed"
	"encoding/json"
	"strings"
	"wfs/backend/data"

	"github.com/morikuni/failure"
)

//go:embed resource/ships.json
var shipsByte []byte

type unregisteredShip struct {
	En      string `json:"en"`
	ID      int    `json:"id"`
	Ja      string `json:"ja"`
	Level   uint   `json:"level"`
	Nation  string `json:"nation"`
	Species string `json:"Species"`
}

type Unregistered struct{}

func NewUnregistered() *Unregistered {
	return &Unregistered{}
}

func (u *Unregistered) Warship() (data.Warships, error) {
	var ships []unregisteredShip
	result := make(data.Warships)

	if err := json.Unmarshal(shipsByte, &ships); err != nil {
		return result, failure.Wrap(err)
	}

	for _, us := range ships {
		result[us.ID] = data.Warship{
			Name:   us.Ja,
			Tier:   us.Level,
			Type:   data.NewShipType(us.Species),
			Nation: data.Nation(toOfficialNation(us.Nation)),
		}
	}

	return result, nil
}

func toOfficialNation(input string) string {
	nation := strings.ToLower(input)
	if nation == "united_kingdom" {
		nation = "uk"
	}

	return nation
}
