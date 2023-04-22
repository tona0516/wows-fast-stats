package vo

type WGBattleArenas struct {
	Status string `json:"status"`
	Data   map[int]struct {
		Name string `json:"Name"`
	} `json:"data"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field"`
		Value   string `json:"value"`
	} `json:"error"`
}