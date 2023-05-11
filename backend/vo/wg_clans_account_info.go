package vo

import (
	"reflect"

	"golang.org/x/exp/slices"
)

type WGClansAccountInfo struct {
	Status string                         `json:"status"`
	Data   map[int]WGClansAccountInfoData `json:"data"`
	Error  WGError                        `json:"error"`
}

func (w WGClansAccountInfo) ClanIDs() []int {
	clanIDs := make([]int, 0)
	for i := range w.Data {
		clanID := w.Data[i].ClanID
		if clanID == 0 {
			continue
		}

		if slices.Contains(clanIDs, clanID) {
			continue
		}

		clanIDs = append(clanIDs, clanID)
	}

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
