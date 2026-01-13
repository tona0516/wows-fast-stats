package data

type WGClansInfo struct {
	WGResponseCommon[map[ClanID]WGClansInfoData]
}

type WGClansInfoData struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
}
