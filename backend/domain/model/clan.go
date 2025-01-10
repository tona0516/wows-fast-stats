package model

type Clans map[int]Clan

type Clan struct {
	ID          int
	Tag         string
	Description string
	HexColor    string
	Lang        string
}
