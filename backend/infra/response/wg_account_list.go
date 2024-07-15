package response

import (
	"wfs/backend/data"
)

type WGAccountList struct {
	WGResponse[data.WGAccountList]
}

func (w WGAccountList) Field() string {
	return wgAPIField(&data.WGAccountListData{})
}
