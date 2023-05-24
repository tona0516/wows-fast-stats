package infra

import "changeme/backend/vo"

type UnregisteredInterface interface {
	Warship() (map[int]vo.Warship, error)
}
