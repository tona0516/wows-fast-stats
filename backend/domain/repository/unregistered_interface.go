package repository

import "wfs/backend/domain/model"

type UnregisteredInterface interface {
	Warship() (model.Warships, error)
}
