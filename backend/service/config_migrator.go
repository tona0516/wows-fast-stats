package service

import (
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"

	"github.com/morikuni/failure"
)

type ConfigMigrator struct {
	storage     repository.Storage
	userConfig  repository.UserConfigStore
	alertPlayer repository.AlertPlayerStore
}

func NewConfigMigrator(
	storage repository.Storage,
	userConfig repository.UserConfigStore,
	alertPlayer repository.AlertPlayerStore,
) *ConfigMigrator {
	return &ConfigMigrator{
		storage:     storage,
		userConfig:  userConfig,
		alertPlayer: alertPlayer,
	}
}

func (m *ConfigMigrator) ExecuteIfNeeded() error {
	if err := m.toV1(); err != nil {
		return failure.Translate(err, apperr.MigrationError)
	}

	if err := m.toV2(); err != nil {
		return failure.Translate(err, apperr.MigrationError)
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

	if m.userConfig.IsExistV0() && !m.userConfig.IsExistV1() {
		userConfig, err := m.userConfig.GetV0()
		if err != nil {
			return err
		}

		if err := m.userConfig.SaveV1(userConfig); err != nil {
			return err
		}

		_ = m.userConfig.DeleteV0()
	}

	if m.alertPlayer.IsExistV0() && !m.alertPlayer.IsExistV1() {
		alertPlayers, err := m.alertPlayer.GetV0()
		if err != nil {
			return err
		}

		if err := m.alertPlayer.SaveV1(alertPlayers); err != nil {
			return err
		}

		_ = m.alertPlayer.DeleteV0()
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

	update := func(config model.UserConfigV2) error {
		if err := m.userConfig.SaveV2(config); err != nil {
			return err
		}

		if err := m.storage.WriteDataVersion(2); err != nil {
			return err
		}

		return nil
	}

	v2, err := m.userConfig.GetV2()
	if err != nil {
		return err
	}

	// バージョンが存在しないかつバグが発生していない場合
	if v2.Version == 0 && v2.Display != (model.UCDisplay{}) {
		v2.Version = 2
		return update(v2)
	}

	// マイグレーションが必要な場合
	v1, err := m.userConfig.GetV1()
	if err != nil {
		return err
	}
	v2 = model.FromUserConfigV1(v1)

	return update(v2)
}
