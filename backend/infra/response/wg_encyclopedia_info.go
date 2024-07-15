package response

import (
	"wfs/backend/data"
)

type WGEncycInfo struct {
	WGResponse[data.WGEncycInfoData]
}

func (w WGEncycInfo) Field() string {
	return wgAPIField(&data.WGEncycInfoData{})
}
