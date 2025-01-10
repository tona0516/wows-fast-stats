package data

type WGClansInfo map[int]WGClansInfoData

type WGClansInfoData struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
}
