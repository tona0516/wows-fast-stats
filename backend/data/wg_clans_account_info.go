package data

type WGClansAccountInfo map[int]WGClansAccountInfoData

type WGClansAccountInfoData struct {
	ClanID int `json:"clan_id"`
}
