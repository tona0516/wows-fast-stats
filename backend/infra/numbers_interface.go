package infra

import "changeme/backend/vo"

type NumbersInterface interface {
	ExpectedStats() (vo.NSExpectedStats, error)
}
