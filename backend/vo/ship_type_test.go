package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShipType(t *testing.T) {
	t.Parallel()
	assert.Equal(t, CV, NewShipType("AirCarrier"))
	assert.Equal(t, BB, NewShipType("Battleship"))
	assert.Equal(t, CL, NewShipType("Cruiser"))
	assert.Equal(t, DD, NewShipType("Destroyer"))
	assert.Equal(t, SS, NewShipType("Submarine"))
	assert.Equal(t, AUX, NewShipType("Auxiliary"))
	assert.Equal(t, NONE, NewShipType("UnknownType"))
}

func TestShipType_Priority(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 0, CV.Priority())
	assert.Equal(t, 1, BB.Priority())
	assert.Equal(t, 2, CL.Priority())
	assert.Equal(t, 3, DD.Priority())
	assert.Equal(t, 4, SS.Priority())
	assert.Equal(t, 5, AUX.Priority())
	assert.Equal(t, 999, NONE.Priority())
}
