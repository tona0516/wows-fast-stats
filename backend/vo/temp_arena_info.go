package vo

import "strings"

type TempArenaInfo struct {
	Vehicles []struct {
		ShipID   int    `json:"shipId"`
		Relation int    `json:"relation"`
		ID       int    `json:"id"`
		Name     string `json:"name"`
	} `json:"vehicles"`
}

func (t *TempArenaInfo) AccountNames() []string {
	accountNames := make([]string, 0)
	for i := range t.Vehicles {
		vehicle := t.Vehicles[i]
		if strings.HasPrefix(vehicle.Name, ":") && strings.HasSuffix(vehicle.Name, ":") {
			continue
		}

		accountNames = append(accountNames, vehicle.Name)
	}

	return accountNames
}
