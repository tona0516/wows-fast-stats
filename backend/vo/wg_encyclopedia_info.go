package vo

type WGEncyclopediaInfo struct {
	Status string `json:"status"`
	Data   struct {
		GameVersion string `json:"game_version"`
	} `json:"data"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field"`
		Value   string `json:"value"`
	} `json:"error"`
}

func (w WGEncyclopediaInfo) GetStatus() string {
	return w.Status
}

func (w WGEncyclopediaInfo) GetError() WGError {
	return WGError(w.Error)
}
