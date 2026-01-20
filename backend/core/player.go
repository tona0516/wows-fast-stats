package core

type AccountID int

type Player struct {
	PlayerInfo PlayerInfo  `json:"player_info"`
	Warship    Warship     `json:"warship"`
	PvPSolo    PlayerStats `json:"pvp_solo"`
	PvPAll     PlayerStats `json:"pvp_all"`
	RankSolo   PlayerStats `json:"rank_solo"`
}

type PlayerInfo struct {
	ID       AccountID `json:"id"`
	Name     string    `json:"name"`
	Clan     Clan      `json:"clan"`
	IsHidden bool      `json:"is_hidden"`
}
