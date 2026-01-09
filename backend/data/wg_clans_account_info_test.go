package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWGClansAccountInfo_ClanIDs(t *testing.T) {
	t.Parallel()

	w := WGClansAccountInfo{
		WGResponseCommon: WGResponseCommon[map[int]WGClansAccountInfoData]{
			Status: "",
			Error:  WGError{},
			Data: map[int]WGClansAccountInfoData{
				1: {ClanID: 123},
				2: {ClanID: 0},
				3: {ClanID: 456},
				5: {ClanID: 789},
				6: {ClanID: 123},
			},
		},
	}

	expectedIDs := []int{123, 456, 789}
	actualIDs := w.ClanIDs()

	assert.Equal(t, expectedIDs, actualIDs, "should be equal")
}
