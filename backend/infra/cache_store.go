package infra

import (
	"path/filepath"
	"wfs/backend/config"
	"wfs/backend/data"

	"github.com/samber/do/v2"
)

type CacheStore struct {
	dir              string
	ownIGNFile       string
	warshipsFile     string
	battleArenasFile string
	battleTypesFile  string
}

func NewCacheStore(i do.Injector) (*CacheStore, error) {
	config := do.MustInvoke[config.Config](i)
	return &CacheStore{
		dir:              config.LocalFile.CacheDir,
		ownIGNFile:       "own_ign.txt",
		warshipsFile:     "warships.json",
		battleArenasFile: "battle_arenas.json",
		battleTypesFile:  "battle_types.json",
	}, nil
}

func (s *CacheStore) OwnIGN() (string, error) {
	return readString(filepath.Join(s.dir, s.ownIGNFile))
}

func (s *CacheStore) SetOwnIGN(ign string) {
	_ = writeString(filepath.Join(s.dir, s.ownIGNFile), ign)
}

func (s *CacheStore) Warships() (data.Warships, error) {
	return readJSON[data.Warships](filepath.Join(s.dir, s.warshipsFile))
}

func (s *CacheStore) SetWarships(data data.Warships) {
	_ = writeJSON(filepath.Join(s.dir, s.warshipsFile), data)
}

func (s *CacheStore) BattleArenas() (map[int]string, error) {
	return readJSON[map[int]string](filepath.Join(s.dir, s.battleArenasFile))
}

func (s *CacheStore) SetBattleArenas(data map[int]string) {
	_ = writeJSON(filepath.Join(s.dir, s.battleArenasFile), data)
}

func (s *CacheStore) BattleTypes() (map[string]string, error) {
	return readJSON[map[string]string](filepath.Join(s.dir, s.battleTypesFile))
}

func (s *CacheStore) SetBattleTypes(data map[string]string) {
	_ = writeJSON(filepath.Join(s.dir, s.battleTypesFile), data)
}
