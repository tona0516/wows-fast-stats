package repository

import "wfs/backend/domain"

type UnregisteredInterface interface {
	Warship() (map[int]domain.Warship, error)
}
