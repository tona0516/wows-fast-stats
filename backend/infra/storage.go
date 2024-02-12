package infra

import (
	"bytes"
	"encoding/gob"
	"errors"
	"wfs/backend/domain/model"

	"github.com/dgraph-io/badger/v4"
	"github.com/morikuni/failure"
)

const (
	dataVersionKey   = "data_version"
	userConfigKey    = "user_config"
	alertPlayersKey  = "alert_players"
	expectedStatsKey = "expected_stats"
	ownIGNKey        = "own_ign"
)

type Storage struct {
	db *badger.DB
}

func NewStorage(db *badger.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) DataVersion() (uint, error) {
	version, err := read[uint](s.db, dataVersionKey)
	if isErrKeyNotFound(err) {
		return 0, nil
	}

	return version, err
}

func (s *Storage) WriteDataVersion(version uint) error {
	return write(s.db, dataVersionKey, version)
}

func (s *Storage) UserConfigV1() (model.UserConfigV1, error) {
	return read[model.UserConfigV1](s.db, userConfigKey)
}

func (s *Storage) WriteUserConfigV1(config model.UserConfigV1) error {
	return write(s.db, userConfigKey, config)
}

func (s *Storage) IsExistUserConfig() bool {
	_, err := read[model.UserConfig](s.db, userConfigKey)
	return !isErrKeyNotFound(err)
}

func (s *Storage) UserConfig() (model.UserConfig, error) {
	config, err := read[model.UserConfig](s.db, userConfigKey)
	if isErrKeyNotFound(err) {
		return model.DefaultUserConfig, nil
	}

	return config, err
}

func (s *Storage) WriteUserConfig(config model.UserConfig) error {
	return write(s.db, userConfigKey, config)
}

func (s *Storage) IsExistAlertPlayers() bool {
	_, err := read[[]model.AlertPlayer](s.db, alertPlayersKey)
	return !isErrKeyNotFound(err)
}

func (s *Storage) AlertPlayers() ([]model.AlertPlayer, error) {
	players, err := read[[]model.AlertPlayer](s.db, alertPlayersKey)
	if isErrKeyNotFound(err) {
		return make([]model.AlertPlayer, 0), nil
	}

	return players, err
}

func (s *Storage) WriteAlertPlayers(players []model.AlertPlayer) error {
	return write(s.db, alertPlayersKey, players)
}

func (s *Storage) ExpectedStats() (model.ExpectedStats, error) {
	return read[model.ExpectedStats](s.db, expectedStatsKey)
}

func (s *Storage) WriteExpectedStats(expectedStats model.ExpectedStats) error {
	return write(s.db, expectedStatsKey, expectedStats)
}

func (s *Storage) OwnIGN() (string, error) {
	return read[string](s.db, ownIGNKey)
}

func (s *Storage) WriteOwnIGN(ign string) error {
	return write(s.db, ownIGNKey, ign)
}

func read[T any](db *badger.DB, key string) (T, error) {
	var result T
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return failure.Wrap(err)
		}

		err = item.Value(func(val []byte) error {
			buf := bytes.NewBuffer(val)
			return gob.NewDecoder(buf).Decode(&result)
		})

		return failure.Wrap(err)
	})

	return result, err
}

func write[T any](db *badger.DB, key string, target T) error {
	err := db.Update(func(txn *badger.Txn) error {
		buf := bytes.NewBuffer(nil)

		if err := gob.NewEncoder(buf).Encode(&target); err != nil {
			return failure.Wrap(err)
		}

		entry := badger.NewEntry([]byte(key), buf.Bytes())
		return failure.Wrap(txn.SetEntry(entry))
	})

	return err
}

func isErrKeyNotFound(err error) bool {
	return err != nil && errors.Is(err, badger.ErrKeyNotFound)
}
