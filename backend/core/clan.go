package core

type ClanID int

type Clans map[AccountID]Clan

type Clan struct {
	ID        ClanID `json:"id"`
	Tag       string `json:"tag"`
	ColorCode string `json:"colorCode"`
	Language  string `json:"language"`
}
