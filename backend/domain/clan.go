package domain

type Clans map[int]Clan

type Clan struct {
	Tag string `json:"tag"`
	ID  int    `json:"id"`
}
