package infra

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wfs/backend/domain"
)

func TestUnregistered_Warship(t *testing.T) {
	t.Parallel()

	// テスト用ships.json
	mockShipsJSON := `[
		{"id": 1, "en": "Ship 1", "level": 11, "nation": "United_Kingdom", "species": "AirCarrier"},
        {"id": 2, "en": "Ship 2", "level": 10, "nation": "USA", "species": "Battleship"},
        {"id": 3, "en": "Ship 3", "level": 9, "nation": "Japan", "species": "Cruiser"},
        {"id": 4, "en": "Ship 4", "level": 8, "nation": "Pan_Asia", "species": "Destroyer"},
        {"id": 5, "en": "Ship 5", "level": 7, "nation": "Commonwealth", "species": "Submarine"},
        {"id": 6, "en": "Ship 6", "level": 6, "nation": "Events", "species": "Auxiliary"}
	]`
	shipsByte = []byte(mockShipsJSON)

	// テスト
	unregistered := NewUnregistered()
	actual, err := unregistered.Warship()

	// アサーション
	expected := map[int]domain.Warship{
		1: {
			Name:   "Ship 1",
			Tier:   11,
			Type:   domain.CV,
			Nation: "uk",
		},
		2: {
			Name:   "Ship 2",
			Tier:   10,
			Type:   domain.BB,
			Nation: "usa",
		},
		3: {
			Name:   "Ship 3",
			Tier:   9,
			Type:   domain.CL,
			Nation: "japan",
		},
		4: {
			Name:   "Ship 4",
			Tier:   8,
			Type:   domain.DD,
			Nation: "pan_asia",
		},
		5: {
			Name:   "Ship 5",
			Tier:   7,
			Type:   domain.SS,
			Nation: "commonwealth",
		},
		6: {
			Name:   "Ship 6",
			Tier:   6,
			Type:   domain.AUX,
			Nation: "events",
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
