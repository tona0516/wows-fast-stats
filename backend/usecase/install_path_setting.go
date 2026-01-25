package usecase

import (
	"context"
	"wfs/backend/adapter"
	"wfs/backend/core"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type InstallPathSetting struct {
	prefStore      adapter.PrefStore
	wails          adapter.Wails
	validator      *ValidateInstallPathService
	gameClientFile string
}

func NewInstallPathSetting(i do.Injector) (*InstallPathSetting, error) {
	return &InstallPathSetting{
		prefStore:      do.MustInvoke[adapter.PrefStore](i),
		wails:          do.MustInvoke[adapter.Wails](i),
		validator:      do.MustInvoke[*ValidateInstallPathService](i),
		gameClientFile: "WorldOfWarships.exe",
	}, nil
}

func (s *InstallPathSetting) Invoke(ctx context.Context) error {
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

	pref, err := s.prefStore.Pref()
	if err != nil {
		return err
	}

	pref.InstallPath = selectedPath

	if err := s.prefStore.SetPref(pref); err != nil {
		return err
	}

	return nil
}
