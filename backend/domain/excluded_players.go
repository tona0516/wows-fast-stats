package domain

type ExcludedPlayers map[int]bool

func (ep ExcludedPlayers) IDs() []int {
	ids := make([]int, 0, len(ep))
	for id := range ep {
		ids = append(ids, id)
	}

	return ids
}

func (ep ExcludedPlayers) Add(playerID int) {
	ep[playerID] = true
}

func (ep ExcludedPlayers) Remove(playerID int) {
	delete(ep, playerID)
}
