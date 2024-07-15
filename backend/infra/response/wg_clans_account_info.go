package response

import (
	"wfs/backend/data"
)

type WGClansAccountInfo struct {
	WGResponse[data.WGClansAccountInfo]
}

func (w WGClansAccountInfo) Field() string {
	return wgAPIField(&data.WGClansAccountInfoData{})
}
