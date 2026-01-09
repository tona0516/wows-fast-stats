package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShipType_New(t *testing.T) {
	t.Parallel()
	assert.Equal(t, ShipTypeCV, NewShipType("AirCarrier"))
	assert.Equal(t, ShipTypeBB, NewShipType("Battleship"))
	assert.Equal(t, ShipTypeCL, NewShipType("Cruiser"))
	assert.Equal(t, ShipTypeDD, NewShipType("Destroyer"))
	assert.Equal(t, ShipTypeSS, NewShipType("Submarine"))
	assert.Equal(t, ShipTypeAUX, NewShipType("Auxiliary"))
	assert.Equal(t, ShipTypeNONE, NewShipType("UnknownType"))
}

func TestShipType_Priority(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 0, ShipTypeCV.Priority())
	assert.Equal(t, 1, ShipTypeBB.Priority())
	assert.Equal(t, 2, ShipTypeCL.Priority())
	assert.Equal(t, 3, ShipTypeDD.Priority())
	assert.Equal(t, 4, ShipTypeSS.Priority())
	assert.Equal(t, 5, ShipTypeAUX.Priority())
	assert.Equal(t, 999, ShipTypeNONE.Priority())
}
