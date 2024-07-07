package vo

type ClanID struct {
	ValueObject[uint]
}

func NewClanID(value uint) ClanID {
	return ClanID{ValueObject[uint]{value}}
}

