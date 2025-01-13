package response

import (
	"reflect"
)

type WGClansInfoResponse struct {
	WGResponseCommon[WGClansInfo]
}

func (w WGClansInfoResponse) Field() string {
	return fieldQuery(reflect.TypeOf(&WGClansInfoData{}).Elem())
}

type WGClansInfo map[int]WGClansInfoData

type WGClansInfoData struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
}
