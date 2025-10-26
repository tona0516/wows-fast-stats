package data

import (
	"reflect"
	"wfs/backend/util"
)

type WGEncycInfo struct {
	WGResponseCommon[WGEncycInfoData]
}

func (w WGEncycInfo) Field() string {
	return util.FieldQuery(reflect.TypeOf(&WGEncycInfoData{}).Elem())
}

type WGEncycInfoData struct {
	GameVersion string `json:"game_version"`
}
