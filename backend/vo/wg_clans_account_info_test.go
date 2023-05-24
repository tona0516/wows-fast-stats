package vo_test

import (
	"changeme/backend/vo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWGClansAccountInfo_ClanIDs(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	w := vo.WGClansAccountInfo{
		Data: map[int]vo.WGClansAccountInfoData{
			1: {ClanID: 123},
			2: {ClanID: 0},
			3: {ClanID: 456},
			5: {ClanID: 789},
			6: {ClanID: 123},
		},
	}

	expectedIDs := []int{123, 456, 789}
	actualIDs := w.ClanIDs()

	assert.Equal(actualIDs, expectedIDs, "should be equal")
}
