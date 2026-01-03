package data

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestPlayers_Sorter(t *testing.T) {
	t.Parallel()

	expected := Players{
		{ShipInfo: ShipInfo{Name: "CV2", Type: ShipTypeCV, Tier: 8, Nation: "usa"}},
		{ShipInfo: ShipInfo{Name: "CV1", Type: ShipTypeCV, Tier: 6, Nation: "japan"}},
		{ShipInfo: ShipInfo{Name: "BB3", Type: ShipTypeBB, Tier: 8, Nation: "japan"}},
		{ShipInfo: ShipInfo{Name: "BB2", Type: ShipTypeBB, Tier: 7, Nation: "japan"}},
		{ShipInfo: ShipInfo{Name: "BB1", Type: ShipTypeBB, Tier: 6, Nation: "japan"}},
		{ShipInfo: ShipInfo{Name: "CL3", Type: ShipTypeCL, Tier: 6, Nation: "ussr"}},
		{ShipInfo: ShipInfo{Name: "CL2", Type: ShipTypeCL, Tier: 6, Nation: "germany"}},
		{ShipInfo: ShipInfo{Name: "CL1", Type: ShipTypeCL, Tier: 6, Nation: "uk"}},
		{ShipInfo: ShipInfo{Name: "DD2", Type: ShipTypeDD, Tier: 8, Nation: "japan"}},
		{ShipInfo: ShipInfo{Name: "DD1", Type: ShipTypeDD, Tier: 8, Nation: "unspecified"}},
		{ShipInfo: ShipInfo{Name: "SS-A", Type: ShipTypeSS, Tier: 6, Nation: "japan"}},
		{ShipInfo: ShipInfo{Name: "SS-Z", Type: ShipTypeSS, Tier: 6, Nation: "japan"}},
	}

	for i := range 100 {
		actual := make(Players, len(expected))
		copy(actual, expected)
		//nolint:gosec
		rand.New(rand.NewSource(int64(i)))
		rand.Shuffle(len(actual), func(i, j int) { actual[i], actual[j] = actual[j], actual[i] })

		sort.Sort(actual)

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Sorter interface implementation failed. Expected %v, but got %v", expected, actual)
		}
	}
}
