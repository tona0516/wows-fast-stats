package data

type Clans map[int]Clan

type Clan struct {
	Tag      string `json:"tag"`
	ID       int    `json:"id"`
	HexColor string `json:"hex_color"`
	Language string `json:"language"`
}
