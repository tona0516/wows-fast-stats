package data

type ShipBadge string

const (
	ShipBadgeExpert = "E"
	ShipBadgeFirst  = "1"
	ShipBadgeSecond = "2"
	ShipBadgeThird  = "3"
	ShipBadgeNone   = ""
)

func ShipBadges() []string {
	return []string{
		ShipBadgeExpert,
		ShipBadgeFirst,
		ShipBadgeSecond,
		ShipBadgeThird,
		ShipBadgeNone,
	}
}
