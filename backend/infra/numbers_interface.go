package infra

import "wfs/backend/vo"

type NumbersInterface interface {
	ExpectedStats() (vo.NSExpectedStats, error)
}
