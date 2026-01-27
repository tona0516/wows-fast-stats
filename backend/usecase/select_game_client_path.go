package usecase

import (
	"context"
	"wfs/backend/adapter"
	"wfs/backend/core"
	"wfs/backend/usecase/service"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type SelectGameClientPath struct {
	prefStore adapter.PrefStore
	wails     adapter.Wails
	validator *service.GameClientPathValidator
}

func NewSelectGameClientPath(i do.Injector) (*SelectGameClientPath, error) {
	return &SelectGameClientPath{
		prefStore: do.MustInvoke[adapter.PrefStore](i),
		wails:     do.MustInvoke[adapter.Wails](i),
		validator: do.MustInvoke[*service.GameClientPathValidator](i),
	}, nil
}

func (s *SelectGameClientPath) Invoke(ctx context.Context) error {
	selectedPath, err := s.wails.OpenDirectoryDialog(ctx)
	if err != nil {
		return failure.Translate(err, core.ErrWailsOpenDirectoryDialog)
	}

	if selectedPath == "" {
		return failure.New(core.ErrSelectFolderCancelled)
	}

	if err := s.validator.Validate(selectedPath); err != nil {
		return err
	}

	if err := s.prefStore.SetGameClientPath(selectedPath); err != nil {
		return err
	}

	return nil
}
