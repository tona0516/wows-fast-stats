package infra

import (
	"reflect"
)

type WGClansAccountInfoResponse struct {
	WGResponseCommon[WGClansAccountInfo]
}

func (w WGClansAccountInfoResponse) Field() string {
	return fieldQuery(reflect.TypeOf(&WGClansAccountInfoData{}).Elem())
}

type WGClansAccountInfo map[int]WGClansAccountInfoData

type WGClansAccountInfoData struct {
	ClanID int `json:"clan_id"`
}
