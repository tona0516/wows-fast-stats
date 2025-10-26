package data

import (
	"reflect"
	"wfs/backend/util"
)

type AllPlayerShipsBadges map[int][]WGShipsBadgesData

type WGShipsBadges struct {
	WGResponseCommon[map[int][]WGShipsBadgesData]
}

func (w WGShipsBadges) Field() string {
	return util.FieldQuery(reflect.TypeOf(&WGShipsBadgesData{}).Elem())
}

type WGShipsBadgesData struct {
	ShipID        int `json:"ship_id"`
	TopGradeClass int `json:"top_grade_class"`
}
