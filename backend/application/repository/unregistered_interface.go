package repository

import "wfs/backend/domain"

type UnregisteredInterface interface {
	Warship() (domain.Warships, error)
}
