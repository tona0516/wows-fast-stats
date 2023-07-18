package domain

import "reflect"

type WGClansInfo struct {
	Status string                  `json:"status"`
	Data   map[int]WGClansInfoData `json:"data"`
	Error  WGError                 `json:"error"`
}

func (w WGClansInfo) GetStatus() string {
	return w.Status
}

func (w WGClansInfo) GetError() WGError {
	return w.Error
}

type WGClansInfoData struct {
	Tag string `json:"tag"`
}

func (w WGClansInfoData) Field() string {
	return fieldQuery(reflect.TypeOf(&w).Elem())
}
