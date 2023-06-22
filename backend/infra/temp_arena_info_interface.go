package infra

import "wfs/backend/vo"

type TempArenaInfoInterface interface {
	Get(installPath string) (vo.TempArenaInfo, error)
	Save(tempArenaInfo vo.TempArenaInfo) error
}
