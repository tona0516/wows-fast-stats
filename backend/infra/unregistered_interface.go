package infra

import "wfs/backend/vo"

type UnregisteredInterface interface {
	Warship() (map[int]vo.Warship, error)
}
