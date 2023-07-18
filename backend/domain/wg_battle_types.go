package domain

import "reflect"

type WGBattleTypes struct {
	Status string                       `json:"status"`
	Data   map[string]WGBattleTypesData `json:"data"`
	Error  WGError                      `json:"error"`
}

func (w WGBattleTypes) GetStatus() string {
	return w.Status
}

func (w WGBattleTypes) GetError() WGError {
	return w.Error
}

type WGBattleTypesData struct {
	Name string `json:"name"`
}

func (w WGBattleTypesData) Field() string {
	return fieldQuery(reflect.TypeOf(&w).Elem())
}
