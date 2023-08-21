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
		result[us.ID] = domain.Warship{
			Name:   us.Ja,
			Tier:   us.Level,
			Type:   domain.NewShipType(us.Species),
			Nation: domain.Nation(toOfficialNation(us.Nation)),
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

type unregisteredShip struct {
	En      string `json:"en"`
	ID      int    `json:"id"`
	Ja      string `json:"ja"`
	Level   uint   `json:"level"`
	Nation  string `json:"nation"`
	Species string `json:"Species"`
}
