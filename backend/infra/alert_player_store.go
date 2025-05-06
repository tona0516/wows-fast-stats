package infra

import (
	"os"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"

	"github.com/dgraph-io/badger/v4"
	"github.com/morikuni/failure"
)

type AlertPlayerStore struct {
	db      *badger.DB
	keyName string
	v0Path  string
}

func NewAlertPlayerStore(db *badger.DB) *AlertPlayerStore {
	return &AlertPlayerStore{
		db:      db,
		keyName: "alert_player",
		v0Path:  filepath.Join("config", "alert_player.json"),
	}
}

func (s *AlertPlayerStore) IsExistV0() bool {
	_, err := os.Stat(s.v0Path)

	return err == nil
}

func (s *AlertPlayerStore) IsExistV1() bool {
	_, err := read[[]model.AlertPlayer](s.db, s.keyName)

	return !isErrKeyNotFound(err)
}

func (s *AlertPlayerStore) GetV0() ([]model.AlertPlayer, error) {
	players, err := readJSON(s.v0Path, []model.AlertPlayer{})
	if err != nil && failure.Is(err, apperr.FileNotExist) {
		return []model.AlertPlayer{}, nil
	}

	return players, err
}

func (s *AlertPlayerStore) GetV1() ([]model.AlertPlayer, error) {
	players, err := read[[]model.AlertPlayer](s.db, s.keyName)
	if isErrKeyNotFound(err) {
		return []model.AlertPlayer{}, nil
	}

	return players, err
}

func (s *AlertPlayerStore) SaveV1(players []model.AlertPlayer) error {
	if len(players) == 0 {
		return remove(s.db, s.keyName)
	}

	return write(s.db, s.keyName, players)
}

func (s *AlertPlayerStore) DeleteV0() error {
	return os.RemoveAll(s.v0Path)
}
