package infra

import (
	"os"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"

	"github.com/dgraph-io/badger/v4"
	"github.com/morikuni/failure"
)

type UserConfigStore struct {
	db      *badger.DB
	keyName string
	v0Path  string
}

func NewUserConfigStore(db *badger.DB) *UserConfigStore {
	return &UserConfigStore{
		db:      db,
		keyName: "user_config",
		v0Path:  filepath.Join("config", "user.json"),
	}
}

func (s *UserConfigStore) IsExistV0() bool {
	_, err := os.Stat(s.v0Path)
	return err == nil
}

func (s *UserConfigStore) IsExistV1() bool {
	_, err := read[model.UserConfig](s.db, s.keyName)
	return !isErrKeyNotFound(err)
}

func (s *UserConfigStore) GetV0() (model.UserConfig, error) {
	config, err := readJSON(s.v0Path, model.DefaultUserConfig())
	if err != nil && failure.Is(err, apperr.FileNotExist) {
		return model.DefaultUserConfig(), nil
	}

	return config, err
}

func (s *UserConfigStore) GetV1() (model.UserConfig, error) {
	config, err := read[model.UserConfig](s.db, s.keyName)
	if isErrKeyNotFound(err) {
		return model.DefaultUserConfig(), nil
	}

	return config, err
}

func (s *UserConfigStore) GetV2() (model.UserConfigV2, error) {
	config, err := read[model.UserConfigV2](s.db, s.keyName)
	if isErrKeyNotFound(err) {
		return model.DefaultUserConfigV2(), nil
	}

	return config, err
}

func (s *UserConfigStore) SaveV1(config model.UserConfig) error {
	return write(s.db, s.keyName, config)
}

func (s *UserConfigStore) SaveV2(config model.UserConfigV2) error {
	return write(s.db, s.keyName, config)
}

func (s *UserConfigStore) DeleteV0() error {
	return os.RemoveAll(s.v0Path)
}
