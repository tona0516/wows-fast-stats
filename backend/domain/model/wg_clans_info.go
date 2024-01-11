package model

type WGClansInfo map[int]WGClansInfoData

type WGClansInfoData struct {
	Tag string `json:"tag"`
}
