package vo

import "reflect"

type WGEncycInfo struct {
	Status string                 `json:"status"`
	Data   WGEncyclopediaInfoData `json:"data"`
	Error  WGError                `json:"error"`
}

func (w WGEncycInfo) GetStatus() string {
	return w.Status
}

func (w WGEncycInfo) GetError() WGError {
	return w.Error
}

type WGEncyclopediaInfoData struct {
	GameVersion string `json:"game_version"`
}

func (w WGEncyclopediaInfoData) Field() string {
	return fieldQuery(reflect.TypeOf(&w).Elem())
}
