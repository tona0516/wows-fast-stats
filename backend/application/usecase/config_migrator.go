package usecase

import (
	"wfs/backend/adapter"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
)

type ConfigMigrator struct {
	localFile adapter.LocalFileInterface
	storage   adapter.StorageInterface
	logger    adapter.LoggerInterface
}

func NewConfigMigrator(
	localFile adapter.LocalFileInterface,
	storage adapter.StorageInterface,
	logger adapter.LoggerInterface,
) *ConfigMigrator {
	return &ConfigMigrator{
		localFile: localFile,
		storage:   storage,
		logger:    logger,
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
	version, err := m.storage.ReadDataVersion()
	if err != nil {
		return err
	}

	if version > 0 {
		return nil
	}

	if m.localFile.IsExistUser() && !m.storage.IsExistUserConfig() {
		userConfig, err := m.localFile.User()
		if err != nil {
			return err
		}

		if err := m.storage.WriteUserConfig(userConfig); err != nil {
			return err
		}

		_ = m.localFile.DeleteUser()
	}

	if m.localFile.IsExistAlertPlayers() && !m.storage.IsExistAlertPlayers() {
		alertPlayers, err := m.localFile.AlertPlayers()
		if err != nil {
			return err
		}

		if err := m.storage.WriteAlertPlayers(alertPlayers); err != nil {
			return err
		}

		_ = m.localFile.DeleteAlertPlayers()
	}

	if err := m.storage.WriteDataVersion(1); err != nil {
		return err
	}

	return nil
}
