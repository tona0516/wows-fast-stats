package data

import (
	"reflect"
	"wfs/internal/util"
)

type WGEncycInfo struct {
	WGResponseCommon[WGEncycInfoData]
}

func (w WGEncycInfo) Field() string {
	return util.FieldQuery(reflect.TypeFor[WGEncycInfoData]())
}

type WGEncycInfoData struct {
	GameVersion string `json:"game_version"`
}
