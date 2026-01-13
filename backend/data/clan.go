package data

type ClanID int

type Clans map[AccountID]Clan

type Clan struct {
	ID       ClanID `json:"id"`
	Tag      string `json:"tag"`
	HexColor string `json:"hex_color"`
	Language string `json:"language"`
}
