package vo

type WGBattleTypes struct {
	Status string `json:"status"`
	Data   map[string]struct {
		Name string `json:"Name"`
	} `json:"data"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field"`
		Value   string `json:"value"`
	} `json:"error"`
}
