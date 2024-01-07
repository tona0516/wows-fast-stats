package usecase

import (
	"wfs/backend/adapter"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
)

type ConfigMigrator struct {
	configV0 adapter.ConfigV0Interface
	storage  adapter.StorageInterface
	logger   adapter.LoggerInterface
}

func NewConfigMigrator(
	configV0 adapter.ConfigV0Interface,
	storage adapter.StorageInterface,
	logger adapter.LoggerInterface,
) *ConfigMigrator {
	return &ConfigMigrator{
		configV0: configV0,
		storage:  storage,
		logger:   logger,
	}
}

func (m *ConfigMigrator) Execute() error {
	if err := m.toV1(); err != nil {
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
