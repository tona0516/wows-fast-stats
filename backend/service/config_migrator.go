package service

import (
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/repository"

	"github.com/morikuni/failure"
)

type ConfigMigrator struct {
	configV0 repository.ConfigV0Interface
	storage  repository.StorageInterface
	logger   repository.LoggerInterface
}

func NewConfigMigrator(
	configV0 repository.ConfigV0Interface,
	storage repository.StorageInterface,
	logger repository.LoggerInterface,
) *ConfigMigrator {
	return &ConfigMigrator{
		configV0: configV0,
		storage:  storage,
		logger:   logger,
	}
}

func (m *ConfigMigrator) ExecuteIfNeeded() error {
	if err := m.toV1(); err != nil {
		return failure.New(apperr.MigrationError, failure.Messagef("%s", err.Error()))
	}

	if err := m.toV2(); err != nil {
		return failure.New(apperr.MigrationError, failure.Messagef("%s", err.Error()))
	}

	return nil
}

//nolint:cyclop
func (m *ConfigMigrator) toV1() error {
	version, err := m.storage.DataVersion()
	if err != nil {
		return err
	}

	if version > 0 {
		return nil
	}

	if m.configV0.IsExistUser() && !m.storage.IsExistUserConfig() {
		userConfig, err := m.configV0.User()
		if err != nil {
			return err
		}

		if err := m.storage.WriteUserConfig(userConfig); err != nil {
			return err
		}

		_ = m.configV0.DeleteUser()
	}

	if m.configV0.IsExistAlertPlayers() && !m.storage.IsExistAlertPlayers() {
		alertPlayers, err := m.configV0.AlertPlayers()
		if err != nil {
			return err
		}

		if err := m.storage.WriteAlertPlayers(alertPlayers); err != nil {
			return err
		}

		_ = m.configV0.DeleteAlertPlayers()
	}

	if err := m.storage.WriteDataVersion(1); err != nil {
		return err
	}

	return nil
}

func (m *ConfigMigrator) toV2() error {
	version, err := m.storage.DataVersion()
	if err != nil {
		return err
	}

	if version > 1 {
		return nil
	}

	update := func(config data.UserConfigV2) error {
		if err := m.storage.WriteUserConfigV2(config); err != nil {
			return err
		}

		if err := m.storage.WriteDataVersion(2); err != nil {
			return err
		}

		return nil
	}

	v2, err := m.storage.UserConfigV2()
	if err != nil {
		return err
	}

	// バージョンが存在しないかつバグが発生していない場合
	if v2.Version == 0 && v2.Display != (data.UCDisplay{}) {
		v2.Version = 2
		return update(v2)
	}

	// マイグレーションが必要な場合
	v1, err := m.storage.UserConfig()
	if err != nil {
		return err
	}
	v2 = data.FromUserConfigV1(v1)

	return update(v2)
}
