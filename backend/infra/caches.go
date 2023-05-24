package infra

import "changeme/backend/vo"

type Caches struct {
	Dir           string
	Warship       Cache[map[int]vo.Warship]
	ExpectedStats Cache[vo.NSExpectedStats]
	BattleArenas  Cache[vo.WGBattleArenas]
	BattleTypes   Cache[vo.WGBattleTypes]
}

func NewCaches(dir string) *Caches {
	return &Caches{
		Dir:           dir,
		Warship:       *NewCache[map[int]vo.Warship]("warship", dir),
		ExpectedStats: *NewCache[vo.NSExpectedStats]("expectedstats", dir),
		BattleArenas:  *NewCache[vo.WGBattleArenas]("battlearenas", dir),
		BattleTypes:   *NewCache[vo.WGBattleTypes]("battletypes", dir),
	}
}
