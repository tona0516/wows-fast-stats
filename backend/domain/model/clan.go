package model

type Clans map[int]Clan

type Clan struct {
	ID          int    `json:"id"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	HexColor    string `json:"hex_color"`
	Lang        string `json:"lang"`
}
