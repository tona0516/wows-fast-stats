package usecase

import (
	"wfs/backend/adapter"
	"wfs/backend/data"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type SavePref struct {
	configStore adapter.PrefStore
}

func NewSavePref(i do.Injector) (*SavePref, error) {
	return &SavePref{
		configStore: do.MustInvoke[adapter.PrefStore](i),
	}, nil
}

func (sp *SavePref) Invoke(config data.Pref) error {
	if err := sp.configStore.SetPref(config); err != nil {
		return failure.Wrap(err)
	}

	return nil
}
