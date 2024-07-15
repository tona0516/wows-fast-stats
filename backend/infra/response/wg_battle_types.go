package response

import (
	"wfs/backend/data"
)

type WGBattleTypes struct {
	WGResponse[data.WGBattleTypes]
}

func (w WGBattleTypes) Field() string {
	return wgAPIField(&data.WGBattleTypesData{})
}
