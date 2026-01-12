package data

type WGBattleArenas struct {
	WGResponseCommon[map[int]WGBattleArenasData]
}

type WGBattleArenasData struct {
	Name string `json:"name"`
}
