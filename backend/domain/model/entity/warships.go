package entity

import (
	"errors"
	"wfs/backend/domain/model/vo"
)

type Warships struct {
	values map[vo.ShipID]Warship
}

func NewWarships(values map[vo.ShipID]Warship) *Warships {
	return &Warships{values}
}

func (w *Warships) Warship(id vo.ShipID) (Warship, error) {
	warship, ok := w.values[id]
	if !ok {
		return Warship{}, errors.New("undefined_warship")
	}
	return warship, nil
}
