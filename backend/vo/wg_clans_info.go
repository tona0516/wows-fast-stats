package vo

type WGClansInfo struct {
	Status string `json:"status"`
	Data   map[int]struct {
		Tag string `json:"tag"`
	} `json:"data"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field"`
		Value   string `json:"value"`
	} `json:"error"`
}

func (w WGClansInfo) GetStatus() string {
	return w.Status
}

func (w WGClansInfo) GetError() WGError {
	return WGError(w.Error)
}
