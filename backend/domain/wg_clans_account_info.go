package domain

import (
	"reflect"
	"sort"

	"golang.org/x/exp/slices"
)

type WGClansAccountInfo struct {
	Status string                         `json:"status"`
	Data   map[int]WGClansAccountInfoData `json:"data"`
	Error  WGError                        `json:"error"`
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

func (w WGClansAccountInfo) GetStatus() string {
	return w.Status
}

func (w WGClansAccountInfo) GetError() WGError {
	return w.Error
}

type WGClansAccountInfoData struct {
	ClanID int `json:"clan_id"`
}

func (w WGClansAccountInfoData) Field() string {
	return fieldQuery(reflect.TypeOf(&w).Elem())
}
