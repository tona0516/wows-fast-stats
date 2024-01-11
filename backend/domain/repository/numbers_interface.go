package repository

import "wfs/backend/domain/model"

type NumbersInterface interface {
	ExpectedStats() (model.ExpectedStats, error)
}
