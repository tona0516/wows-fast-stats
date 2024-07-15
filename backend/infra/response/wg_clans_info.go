package response

import (
	"wfs/backend/data"
)

type WGClansInfo struct {
	WGResponse[data.WGClansInfo]
}

func (w WGClansInfo) Field() string {
	return wgAPIField(&data.WGClansInfoData{})
}
