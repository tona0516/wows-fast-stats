package infra

import (
	"github.com/dgraph-io/badger/v4"
)

const (
	dataVersionKey = "data_version"
	ownIGNKey      = "own_ign"
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

func (s *Storage) OwnIGN() (string, error) {
	return read[string](s.db, ownIGNKey)
}

func (s *Storage) WriteOwnIGN(ign string) error {
	return write(s.db, ownIGNKey, ign)
}
