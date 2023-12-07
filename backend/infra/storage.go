package infra

import (
	"bytes"
	"encoding/gob"
	"errors"
	"wfs/backend/domain"

	"github.com/dgraph-io/badger/v4"
	"github.com/morikuni/failure"
)

const (
	dataVersionKey     = "data_version"
	userConfigKey      = "user_config"
	alertPlayersKey    = "alert_players"
	nsExpectedStatsKey = "ns_expected_stats"
)

type Storage struct {
	db *badger.DB
}

func NewStorage(db *badger.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) ReadDataVersion() (uint, error) {
	version, err := read[uint](s.db, dataVersionKey)
	if isErrKeyNotFound(err) {
		return 0, nil
	}

	return version, failure.Wrap(err)
}

func (s *Storage) WriteDataVersion(version uint) error {
	return write(s.db, dataVersionKey, version)
}

func (s *Storage) IsExistUserConfig() bool {
	_, err := read[domain.UserConfig](s.db, userConfigKey)
	return !isErrKeyNotFound(err)
}

func (s *Storage) ReadUserConfig() (domain.UserConfig, error) {
	config, err := read[domain.UserConfig](s.db, userConfigKey)
	if isErrKeyNotFound(err) {
		return domain.DefaultUserConfig, nil
	}

	return config, failure.Wrap(err)
}

func (s *Storage) WriteUserConfig(config domain.UserConfig) error {
	return write(s.db, userConfigKey, config)
}

func (s *Storage) IsExistAlertPlayers() bool {
	_, err := read[[]domain.AlertPlayer](s.db, alertPlayersKey)
	return !isErrKeyNotFound(err)
}

func (s *Storage) ReadAlertPlayers() ([]domain.AlertPlayer, error) {
	players, err := read[[]domain.AlertPlayer](s.db, alertPlayersKey)
	if isErrKeyNotFound(err) {
		return make([]domain.AlertPlayer, 0), nil
	}

	return players, failure.Wrap(err)
}

func (s *Storage) WriteAlertPlayers(players []domain.AlertPlayer) error {
	return write(s.db, alertPlayersKey, players)
}

func (s *Storage) ReadNSExpectedStats() (domain.NSExpectedStats, error) {
	return read[domain.NSExpectedStats](s.db, nsExpectedStatsKey)
}

func (s *Storage) WriteNSExpectedStats(nsExpectedStats domain.NSExpectedStats) error {
	return write(s.db, nsExpectedStatsKey, nsExpectedStats)
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

	return result, failure.Wrap(err)
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

	return failure.Wrap(err)
}

func isErrKeyNotFound(err error) bool {
	return err != nil && errors.Is(err, badger.ErrKeyNotFound)
}
