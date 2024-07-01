package domain

type Player struct {
	id       PlayerID
	name     PlayerName
	clan     Clan
	isHidden bool
}

func NewPlayer(
	id PlayerID,
	name PlayerName,
	clan Clan,
	isHidden bool,
) *Player {
	return &Player{
		id:       id,
		name:     name,
		clan:     clan,
		isHidden: isHidden,
	}
}

// Entity interface

func (p *Player) ID() PlayerID {
	return p.id
}

func (p *Player) Equals(another *Player) bool {
	return p.id.value == another.id.value
}
