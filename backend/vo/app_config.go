package vo

type WindowConfig struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type AppConfig struct {
	Window WindowConfig `json:"window"`
}
