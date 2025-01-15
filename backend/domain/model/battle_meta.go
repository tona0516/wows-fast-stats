package model

import "strings"

type BattleMeta struct {
	arenas map[int]string
	types  map[string]string
}

func NewBattleMeta(
	arenas map[int]string,
	types map[string]string,
) *BattleMeta {
	return &BattleMeta{
		arenas: arenas,
		types:  types,
	}
}

func (m BattleMeta) Arena(mapID int) string {
	return m.arenas[mapID]
}

func (m BattleMeta) Type(matchGroup string) string {
	return m.types[strings.ToUpper(matchGroup)]
}
