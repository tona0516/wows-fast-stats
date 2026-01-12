package data

type WGBattleTypes struct {
	WGResponseCommon[map[string]WGBattleTypesData]
}

type WGBattleTypesData struct {
	Name string `json:"name"`
}
