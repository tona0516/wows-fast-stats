package data

import (
	"reflect"
	"wfs/internal/util"
)

type AllPlayerShipsBadges map[int][]WGShipsBadgesData

type WGShipsBadges struct {
	WGResponseCommon[map[int][]WGShipsBadgesData]
}

func (w WGShipsBadges) Field() string {
	return util.FieldQuery(reflect.TypeFor[WGShipsBadgesData]())
}

type WGShipsBadgesData struct {
	ShipID        int `json:"ship_id"`
	TopGradeClass int `json:"top_grade_class"`
}
