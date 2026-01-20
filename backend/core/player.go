package core

type AccountID int

type Player struct {
	PlayerInfo PlayerInfo  `json:"playerInfo"`
	Warship    Warship     `json:"warship"`
	PvPSolo    PlayerStats `json:"pvpSolo"`
	PvPAll     PlayerStats `json:"pvpAll"`
	RankSolo   PlayerStats `json:"rankSolo"`
}

type PlayerInfo struct {
	ID       AccountID `json:"id"`
	Name     string    `json:"name"`
	Clan     Clan      `json:"clan"`
	IsHidden bool      `json:"isHidden"`
}
