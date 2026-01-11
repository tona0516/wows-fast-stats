package data

type PrefetchResult struct {
	Warships     Warships
	BattleArenas map[int]string
	BattleTypes  map[string]string
}
