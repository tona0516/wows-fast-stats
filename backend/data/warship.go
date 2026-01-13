package data

type ShipID int

type Warships map[ShipID]Warship

type Warship struct {
	ID            ShipID        `json:"id"`
	Name          string        `json:"name"`
	Tier          uint          `json:"tier"`
	Type          ShipType      `json:"type"`
	Nation        Nation        `json:"nation"`
	IsPremium     bool          `json:"isPremium"`
	ServerAverage ServerAverage `json:"serverAverage"`
	DamageRatings []RatingValue `json:"damageRatings"`
}

func NewWarship(
	id ShipID,
	name string,
	tier uint,
	shipType ShipType,
	nation Nation,
	isPremium bool,
	serverAverage ServerAverage,
) *Warship {
	return &Warship{
		ID:            id,
		Name:          name,
		Tier:          tier,
		Type:          shipType,
		Nation:        nation,
		IsPremium:     isPremium,
		ServerAverage: serverAverage,
		DamageRatings: NewDamageRatings(serverAverage.Damage),
	}
}

func NewUnknownWarship() *Warship {
	return &Warship{
		Name:          "Unknown",
		Type:          ShipTypeNONE,
		DamageRatings: make([]RatingValue, 0),
	}
}

type ServerAverage struct {
	Damage  float64 `json:"damage"`
	Frags   float64 `json:"frags"`
	WinRate float64 `json:"winRate"`
}
