package infra

import (
	"reflect"
)

type WGEncycInfo struct {
	WGResponseCommon[WGEncycInfoData]
}

func (w WGEncycInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&WGEncycInfoData{}).Elem())
}

type WGEncycInfoData struct {
	GameVersion string `json:"game_version"`
}
