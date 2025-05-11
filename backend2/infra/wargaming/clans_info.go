package wargaming

import (
	"reflect"
)

type ClansInfoResponse struct {
	ResponseCommon[ClansInfo]
}

func (w ClansInfoResponse) Field() string {
	return fieldQuery(reflect.TypeOf(&ClansInfoData{}).Elem())
}

type ClansInfo map[int]ClansInfoData

type ClansInfoData struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
}
