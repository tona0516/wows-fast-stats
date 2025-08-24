package data

type AllPlayerShipsBadges map[int][]WGShipsBadgesData

type WGShipsBadges map[int][]WGShipsBadgesData

type WGShipsBadgesData struct {
	ShipID        int `json:"ship_id"`
	TopGradeClass int `json:"top_grade_class"`
}
