package repository

import "wfs/backend/domain"

type NumbersInterface interface {
	ExpectedStats() (domain.NSExpectedStats, error)
}
