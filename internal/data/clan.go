package data

type Clans map[int]Clan

type Clan struct {
	ID       int    `json:"id"`
	Tag      string `json:"tag"`
	HexColor string `json:"hex_color"`
	Language string `json:"language"`
}
