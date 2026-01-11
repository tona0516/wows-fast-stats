package usecase

import (
	"wfs/backend/adapter"
	"wfs/backend/data"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type SavePref struct {
	configStore adapter.ConfigStore
}

func NewSavePref(i do.Injector) (*SavePref, error) {
	return &SavePref{
		configStore: do.MustInvoke[adapter.ConfigStore](i),
	}, nil
}

func (sp *SavePref) Invoke(config data.UserConfig) error {
	if err := sp.configStore.SetUserConfig(config); err != nil {
		return failure.Wrap(err)
	}

	return nil
}
