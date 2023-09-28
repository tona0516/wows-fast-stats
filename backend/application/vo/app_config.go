package vo

type WindowConfig struct {
	X      int `json:"x"`
	Y      int `json:"Y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type AppConfig struct {
	Window WindowConfig `json:"window"`
}
