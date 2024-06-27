package infra

import (
	"bytes"
	"encoding/gob"
	"errors"
	"wfs/backend/data"

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

func (s *Storage) UserConfig() (data.UserConfig, error) {
	config, err := read[data.UserConfig](s.db, userConfigKey)
	if isErrKeyNotFound(err) {
		return data.DefaultUserConfig(), nil
	}

	return config, err
}

func (s *Storage) WriteUserConfig(config data.UserConfig) error {
	return write(s.db, userConfigKey, config)
}

func (s *Storage) IsExistUserConfig() bool {
	_, err := read[data.UserConfig](s.db, userConfigKey)
	return !isErrKeyNotFound(err)
}

func (s *Storage) UserConfigV2() (data.UserConfigV2, error) {
	config, err := read[data.UserConfigV2](s.db, userConfigKey)
	if isErrKeyNotFound(err) {
		return data.DefaultUserConfigV2(), nil
	}

	return config, err
}

func (s *Storage) WriteUserConfigV2(config data.UserConfigV2) error {
	return write(s.db, userConfigKey, config)
}

func (s *Storage) IsExistAlertPlayers() bool {
	_, err := read[[]data.AlertPlayer](s.db, alertPlayersKey)
	return !isErrKeyNotFound(err)
}

func (s *Storage) AlertPlayers() ([]data.AlertPlayer, error) {
	players, err := read[[]data.AlertPlayer](s.db, alertPlayersKey)
	if isErrKeyNotFound(err) {
		return []data.AlertPlayer{}, nil
	}

	return players, err
}

func (s *Storage) WriteAlertPlayers(players []data.AlertPlayer) error {
	if len(players) == 0 {
		return delete(s.db, alertPlayersKey)
	}
	return write(s.db, alertPlayersKey, players)
}

func (s *Storage) ExpectedStats() (data.ExpectedStats, error) {
	return read[data.ExpectedStats](s.db, expectedStatsKey)
}

func (s *Storage) WriteExpectedStats(expectedStats data.ExpectedStats) error {
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

func delete(db *badger.DB, key string) error {
	return db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return failure.Wrap(err)
	})
}

func isErrKeyNotFound(err error) bool {
	return err != nil && errors.Is(err, badger.ErrKeyNotFound)
}
