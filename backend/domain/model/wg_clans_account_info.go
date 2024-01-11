package model

import (
	"sort"

	"golang.org/x/exp/slices"
)

type WGClansAccountInfo map[int]WGClansAccountInfoData

func (w WGClansAccountInfo) ClanIDs() []int {
	clanIDs := make([]int, 0)
	for _, v := range w {
		if v.ClanID != 0 && !slices.Contains(clanIDs, v.ClanID) {
			clanIDs = append(clanIDs, v.ClanID)
		}
	}
	sort.Ints(clanIDs)
	return clanIDs
}

type WGClansAccountInfoData struct {
	ClanID int `json:"clan_id"`
}
