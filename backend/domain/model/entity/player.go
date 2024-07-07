package entity

import "wfs/backend/domain/model/vo"

type Player struct {
	id       vo.PlayerID
	name     vo.PlayerName
	clan     Clan
	isHidden bool
	warship  Warship
	stat     vo.ShipStat
}

func NewPlayer(
	id vo.PlayerID,
	name vo.PlayerName,
	clan Clan,
	isHidden bool,
	warship Warship,
) *Player {
	return &Player{
		id:       id,
		name:     name,
		clan:     clan,
		isHidden: isHidden,
		warship:  warship,
	}
}

// Entity interface

func (p *Player) ID() vo.PlayerID {
	return p.id
}

func (p *Player) Equals(another *Player) bool {
	return p.id.Value() == another.id.Value()
}
