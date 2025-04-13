package wargaming

import (
	"reflect"
)

type ClansAccountInfoResponse struct {
	ResponseCommon[ClansAccountInfo]
}

func (w ClansAccountInfoResponse) Field() string {
	return fieldQuery(reflect.TypeOf(&ClansAccountInfoData{}).Elem())
}

type ClansAccountInfo map[int]ClansAccountInfoData

type ClansAccountInfoData struct {
	ClanID int `json:"clan_id"`
}
