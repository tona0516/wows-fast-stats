package vo

type PlayerID struct {
	ValueObject[uint]
}

func NewPlayerID(value uint) PlayerID {
	return PlayerID{ValueObject[uint]{value}}
}

