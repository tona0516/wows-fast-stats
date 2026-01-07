package infra

import (
	"path/filepath"
	"wfs/internal/data"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type CacheStore interface {
	OwnIGN() (string, error)
	SetOwnIGN(ign string) error

	Warships() (data.Warships, error)
	SetWarships(data data.Warships) error

	BattleArenas() (map[int]string, error)
	SetBattleArenas(data map[int]string) error

	BattleTypes() (map[string]string, error)
	SetBattleTypes(data map[string]string) error
}

type cacheStore struct {
	dir              string
	ownIGNFile       string
	warshipsFile     string
	battleArenasFile string
	battleTypesFile  string
}

func NewCacheStore(dir string) CacheStore {
	return &cacheStore{
		dir:              dir,
		ownIGNFile:       "own_ign.txt",
		warshipsFile:     "warships.json",
		battleArenasFile: "battle_arenas.json",
		battleTypesFile:  "battle_types.json",
	}
}

func (s *cacheStore) OwnIGN() (string, error) {
	return readString(filepath.Join(s.dir, s.ownIGNFile))
}

func (s *cacheStore) SetOwnIGN(ign string) error {
	return writeString(filepath.Join(s.dir, s.ownIGNFile), ign)
}

func (s *cacheStore) Warships() (data.Warships, error) {
	return readJSON[data.Warships](filepath.Join(s.dir, s.warshipsFile))
}

func (s *cacheStore) SetWarships(data data.Warships) error {
	return writeJSON(filepath.Join(s.dir, s.warshipsFile), data)
}

func (s *cacheStore) BattleArenas() (map[int]string, error) {
	return readJSON[map[int]string](filepath.Join(s.dir, s.battleArenasFile))
}

func (s *cacheStore) SetBattleArenas(data map[int]string) error {
	return writeJSON(filepath.Join(s.dir, s.battleArenasFile), data)
}

func (s *cacheStore) BattleTypes() (map[string]string, error) {
	return readJSON[map[string]string](filepath.Join(s.dir, s.battleTypesFile))
}

func (s *cacheStore) SetBattleTypes(data map[string]string) error {
	return writeJSON(filepath.Join(s.dir, s.battleTypesFile), data)
}
