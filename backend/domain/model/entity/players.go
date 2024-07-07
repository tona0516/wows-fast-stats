package entity

import (
	"errors"
	"sort"
	"wfs/backend/domain/model/vo"
)

type Players struct {
	values []Player
}

func NewPlayers(values []Player) *Players {
	return &Players{values}
}

func (ps *Players) Player(id vo.PlayerID) (Player, error) {
	for _, p := range ps.values {
		if p.ID() == id {
			return p, nil
		}
	}

	return Player{}, errors.New("player_not_found")
}

func (ps *Players) Sorted() Players {
	result := make([]Player, 0)
	result = append(result, ps.values...)

	sort.Slice(ps.values, func(i, j int) bool {
		one := result[i].warship
		another := result[j].warship

		return one.less(&another)
	})

	return *NewPlayers(result)
}
