package vo

import (
	"reflect"
	"sort"

	"github.com/samber/lo"
)

type WGClansAccountInfo struct {
	Status string                         `json:"status"`
	Data   map[int]WGClansAccountInfoData `json:"data"`
	Error  WGError                        `json:"error"`
}

func (w WGClansAccountInfo) ClanIDs() []int {
	clansAccounts := lo.Values(w.Data)
	clanIDs := lo.FilterMap(clansAccounts, func(clansAccount WGClansAccountInfoData, _ int) (int, bool) {
		return clansAccount.ClanID, clansAccount.ClanID != 0
	})
	clanIDs = lo.Uniq(clanIDs)
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
