package infra

import (
	"changeme/backend/vo"
)

type Caches struct {
    Warships Cache[map[int]vo.Warship]
    ExpectedStats Cache[vo.NSExpectedStats]
    BattleArenas Cache[vo.WGBattleArenas]
    BattleTypes Cache[vo.WGBattleTypes]
}

func NewCaches(gameVersion string) *Caches {
    return &Caches{
		Warships: Cache[map[int]vo.Warship]{
			Prefix:   "warships",
			GameVersion: gameVersion,
		},
		ExpectedStats: Cache[vo.NSExpectedStats]{
			Prefix:   "expectedstats",
			GameVersion: gameVersion,
		},
		BattleArenas: Cache[vo.WGBattleArenas]{
			Prefix:   "battlearenas",
			GameVersion: gameVersion,
		},
		BattleTypes: Cache[vo.WGBattleTypes]{
			Prefix:   "battletypes",
			GameVersion: gameVersion,
		},
	}
}

func (c *Caches) RemoveOld() {
    _ = c.Warships.RemoveOld()
    _ = c.ExpectedStats.RemoveOld()
    _ = c.BattleArenas.RemoveOld()
    _ = c.BattleTypes.RemoveOld()
}
