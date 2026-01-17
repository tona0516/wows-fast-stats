package usecase

import (
	"wfs/backend/adapter"
	"wfs/backend/core"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type LoadPref struct {
	configStore adapter.PrefStore
}

func NewLoadPref(i do.Injector) (*LoadPref, error) {
	return &LoadPref{
		configStore: do.MustInvoke[adapter.PrefStore](i),
	}, nil
}

func (lp *LoadPref) Invoke() (core.Pref, error) {
	config, err := lp.configStore.Pref()
	if err != nil {
		if failure.Is(err, core.ErrJSONNotFound) {
			return core.DefaultPref(), nil
		}

		return core.Pref{}, err
	}

	return config, nil
}
