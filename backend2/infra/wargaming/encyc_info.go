package wargaming

import (
	"reflect"
)

type EncycInfo struct {
	ResponseCommon[EncycInfoData]
}

func (w EncycInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&EncycInfoData{}).Elem())
}

type EncycInfoData struct {
	GameVersion string `json:"game_version"`
}
