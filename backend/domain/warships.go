package domain

import "errors"

type Warships struct {
	values map[ShipID]Warship
}

func NewWarships(values map[ShipID]Warship) *Warships {
	return &Warships{values}
}

func (w *Warships) Get(id ShipID) (Warship, error) {
	warship, ok := w.values[id]
	if !ok {
		return Warship{}, errors.New("undefined_warship")
	}
	return warship, nil
}
