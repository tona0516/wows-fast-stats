package response

import (
	"wfs/backend/data"
)

type WGBattleArenas struct {
	WGResponse[data.WGBattleArenas]
}

func (w WGBattleArenas) Field() string {
	return wgAPIField(&data.WGBattleArenasData{})
}
