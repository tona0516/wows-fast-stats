package usecase

import (
	"errors"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"

	"github.com/dgraph-io/badger/v4"
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
		userConfig, err := m.configV0.UserV1()
		if err != nil {
			return err
		}

		if err := m.storage.WriteUserConfigV1(userConfig); err != nil {
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

	v1, err := m.storage.UserConfigV1()
	if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
		return err
	}

	if err == nil {
		v2 := model.FromUserConfigV1(v1)

		if err := m.storage.WriteUserConfig(v2); err != nil {
			return err
		}
	}

	if err := m.storage.WriteDataVersion(2); err != nil {
		return err
	}

	return nil
}
