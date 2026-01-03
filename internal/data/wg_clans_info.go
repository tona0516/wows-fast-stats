package data

import (
	"reflect"
	"wfs/internal/util"
)

type WGClansInfo struct {
	WGResponseCommon[map[int]WGClansInfoData]
}

func (w WGClansInfo) Field() string {
	return util.FieldQuery(reflect.TypeFor[WGClansInfoData]())
}

func (w WGClansInfo) ToArray() []WGClansInfoData {
	array := make([]WGClansInfoData, 0)
	for _, v := range w.Data {
		array = append(array, v)
	}
	return array
}

type WGClansInfoData struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
}
