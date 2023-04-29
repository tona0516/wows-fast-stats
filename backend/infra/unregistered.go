package infra

import (
	"changeme/backend/vo"
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed resource/ships.json
var shipsByte []byte

type Unregistered struct {
}

func (u *Unregistered) Warship() (map[int]vo.Warship, error) {
    var ships []unregisteredShip
    result := make(map[int]vo.Warship, 0)

    err := json.Unmarshal(shipsByte, &ships)
	if err != nil {
		return result, err
	}

    for _, us := range ships {
        nation := strings.ToLower(us.Nation)
        if nation == "united_kingdom" {
            nation = "uk"
        }

        result[us.Id] = vo.Warship{
            Name: us.En,
            Tier: us.Level,
            Type: us.Species,
            Nation: nation,
        }
    }

    return result, nil
}

type unregisteredShip struct {
    Id int `json:"id"`
    En string `json:"en"`
    Level uint `json:"level"`
    Nation string `json:"nation"`
    Species string `json:"Species"`
  }