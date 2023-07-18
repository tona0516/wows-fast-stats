package domain

import "reflect"

type WGBattleArenas struct {
	Status string                     `json:"status"`
	Data   map[int]WGBattleArenasData `json:"data"`
	Error  WGError                    `json:"error"`
}

type WGBattleArenasData struct {
	Name string `json:"name"`
}

func (w WGBattleArenas) GetStatus() string {
	return w.Status
}

func (w WGBattleArenas) GetError() WGError {
	return w.Error
}

func (w WGBattleArenasData) Field() string {
	return fieldQuery(reflect.TypeOf(&w).Elem())
}
