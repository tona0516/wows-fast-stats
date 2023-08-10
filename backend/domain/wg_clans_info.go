package domain

import "reflect"

type WGClansInfo struct {
	WGResponseCommon[map[int]WGClansInfoData]
}

func (w WGClansInfo) Field() string {
	return fieldQuery(reflect.TypeOf(&WGClansInfoData{}).Elem())
}

type WGClansInfoData struct {
	Tag string `json:"tag"`
}
