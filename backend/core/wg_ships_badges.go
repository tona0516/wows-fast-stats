package core

type WGShipsBadges struct {
	WGResponseCommon[map[AccountID][]WGShipsBadgesData]
}

type WGShipsBadgesData struct {
	ShipID        ShipID `json:"ship_id"`
	TopGradeClass int    `json:"top_grade_class"`
}

type AllPlayerShipBadges map[AccountID]PlayerShipBadges

type PlayerShipBadges map[ShipID]WGShipsBadgesData

func (pb PlayerShipBadges) efficiencyBadge(shipID ShipID) EfficiencyBadge {
	if badge, ok := pb[shipID]; ok {
		switch badge.TopGradeClass {
		case 1:
			return EfficiencyBadgeExpert
		case 2:
			return EfficiencyBadgeFirst
		case 3:
			return EfficiencyBadgeSecond
		case 4:
			return EfficiencyBadgeThird
		}
	}

	return EfficiencyBadgeNone
}

func (pb PlayerShipBadges) efficiencyBadges() EfficiencyBadgeGroup {
	var badges EfficiencyBadgeGroup

	for _, b := range pb {
		switch b.TopGradeClass {
		case 1:
			badges.Expert++
		case 2:
			badges.First++
		case 3:
			badges.Second++
		case 4:
			badges.Third++
		}
	}

	return badges
}
