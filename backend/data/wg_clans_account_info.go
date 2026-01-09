package data

import (
	"reflect"
	"slices"
	"sort"
	"wfs/backend/util"
)

type WGClansAccountInfo struct {
	WGResponseCommon[map[int]WGClansAccountInfoData]
}

func (w WGClansAccountInfo) Field() string {
	return util.FieldQuery(reflect.TypeFor[WGClansAccountInfoData]())
}

func (w WGClansAccountInfo) ClanIDs() []int {
	clanIDs := make([]int, 0)
	for _, v := range w.Data {
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
