package entity

import "wfs/backend/domain/model/vo"

type Clan struct {
	id          vo.ClanID
	tag         vo.ClanTag
	colorCode   vo.ColorCode
	countryCode vo.CountryCode
}

func NewClan[T vo.ClanID](
	id vo.ClanID,
	tag vo.ClanTag,
	colorCode vo.ColorCode,
	countryCode vo.CountryCode,
) *Clan {
	return &Clan{
		id:          id,
		tag:         tag,
		colorCode:   colorCode,
		countryCode: countryCode,
	}
}

// Entity interface

func (c *Clan) ID() vo.ClanID {
	return c.id
}

func (c *Clan) Equals(another *Clan) bool {
	return c.id.Value() == another.id.Value()
}
