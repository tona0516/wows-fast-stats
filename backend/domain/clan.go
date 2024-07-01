package domain

type Clan struct {
	id          ClanID
	tag         ClanTag
	colorCode   ColorCode
	countryCode CountryCode
}

func NewClan[T ClanID](
	id ClanID,
	tag ClanTag,
	colorCode ColorCode,
	countryCode CountryCode,
) *Clan {
	return &Clan{
		id:          id,
		tag:         tag,
		colorCode:   colorCode,
		countryCode: countryCode,
	}
}

// Entity interface

func (c *Clan) ID() ClanID {
	return c.id
}

func (c *Clan) Equals(another *Clan) bool {
	return c.id.value == another.id.value
}
