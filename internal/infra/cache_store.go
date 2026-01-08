package infra

import (
	"path/filepath"
	"wfs/internal/data"
)

type CacheStore struct {
	dir              string
	ownIGNFile       string
	warshipsFile     string
	battleArenasFile string
	battleTypesFile  string
}

func NewCacheStore(dir string) *CacheStore {
	return &CacheStore{
		dir:              dir,
		ownIGNFile:       "own_ign.txt",
		warshipsFile:     "warships.json",
		battleArenasFile: "battle_arenas.json",
		battleTypesFile:  "battle_types.json",
	}
}

func (s *CacheStore) OwnIGN() (string, error) {
	return readString(filepath.Join(s.dir, s.ownIGNFile))
}

func (s *CacheStore) SetOwnIGN(ign string) error {
	return writeString(filepath.Join(s.dir, s.ownIGNFile), ign)
}

func (s *CacheStore) Warships() (data.Warships, error) {
	return readJSON[data.Warships](filepath.Join(s.dir, s.warshipsFile))
}

func (s *CacheStore) SetWarships(data data.Warships) error {
	return writeJSON(filepath.Join(s.dir, s.warshipsFile), data)
}

func (s *CacheStore) BattleArenas() (map[int]string, error) {
	return readJSON[map[int]string](filepath.Join(s.dir, s.battleArenasFile))
}

func (s *CacheStore) SetBattleArenas(data map[int]string) error {
	return writeJSON(filepath.Join(s.dir, s.battleArenasFile), data)
}

func (s *CacheStore) BattleTypes() (map[string]string, error) {
	return readJSON[map[string]string](filepath.Join(s.dir, s.battleTypesFile))
}

func (s *CacheStore) SetBattleTypes(data map[string]string) error {
	return writeJSON(filepath.Join(s.dir, s.battleTypesFile), data)
}
