package core

import (
	"slices"
)

type WGClansAccountInfo struct {
	WGResponseCommon[map[AccountID]WGClansAccountInfoData]
}

func (w WGClansAccountInfo) ClanIDs() []ClanID {
	clanIDs := make([]ClanID, 0)
	for _, v := range w.Data {
		if v.ClanID != 0 && !slices.Contains(clanIDs, v.ClanID) {
			clanIDs = append(clanIDs, v.ClanID)
		}
	}
	slices.Sort(clanIDs)
	return clanIDs
}

type WGClansAccountInfoData struct {
	ClanID ClanID `json:"clan_id"`
}
