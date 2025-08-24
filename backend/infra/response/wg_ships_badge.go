package response

import (
	"reflect"
	"wfs/backend/data"
)

type WGShipsBadges struct {
	WGResponseCommon[data.WGShipsBadges]
}

func (w WGShipsBadges) Field() string {
	return fieldQuery(reflect.TypeOf(&data.WGShipsBadgesData{}).Elem())
}
