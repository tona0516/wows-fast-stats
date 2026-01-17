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
