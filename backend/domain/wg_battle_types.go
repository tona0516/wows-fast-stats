package domain

type WGBattleTypes map[string]WGBattleTypesData

type WGBattleTypesData struct {
	Name string `json:"name"`
}
