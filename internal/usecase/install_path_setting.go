package usecase

import (
	"context"
	"wfs/internal/infra"

	"github.com/morikuni/failure"
)

type OpenDirectoryDialogFunc func(ctx context.Context) (string, error)

type InstallPathSetting struct {
	localStorage        infra.LocalStorage
	openDirectoryDialog OpenDirectoryDialogFunc
}

func NewInstallPathSetting(
	localStorage infra.LocalStorage,
	openDirectoryDialogFunc OpenDirectoryDialogFunc,
) *InstallPathSetting {
	return &InstallPathSetting{
		localStorage:        localStorage,
		openDirectoryDialog: openDirectoryDialogFunc,
	}
}

func (s *InstallPathSetting) Invoke(ctx context.Context) (bool, error) {
	selected, err := s.openDirectoryDialog(ctx)
	if err != nil {
		return false, failure.Wrap(err)
	}

	if selected == "" {
		return false, nil
	}

	config, err := s.localStorage.UserConfig()
	if err != nil {
		return false, failure.Wrap(err)
	}

	config.InstallPath = selected

	if err = s.localStorage.SetUserConfig(config); err != nil {
		return false, failure.Wrap(err)
	}

	return true, nil
}
