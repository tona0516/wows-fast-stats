package core

type BattleMetadata struct {
	Unixtime int64  `json:"unixtime"`
	Arena    string `json:"arena"`
	Type     string `json:"type"`
}
