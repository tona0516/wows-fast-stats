package model

type WGBattleArenas map[int]WGBattleArenasData

type WGBattleArenasData struct {
	Name string `json:"name"`
}
