package response

import (
	"wfs/backend/data"
)

type WGAccountInfo struct {
	WGResponse[data.WGAccountInfo]
}

func (w WGAccountInfo) Field() string {
	return wgAPIField(&data.WGAccountInfoData{})
}
