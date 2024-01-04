package adapter

import "wfs/backend/domain"

type NumbersInterface interface {
	ExpectedStats() (domain.ExpectedStats, error)
}
