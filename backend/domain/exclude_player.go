package domain

import "golang.org/x/exp/slices"

type ExcludePlayer struct {
    playerIDs []int
}

func NewExcludePlayer() *ExcludePlayer {
    return &ExcludePlayer{
        playerIDs: make([]int, 0),
    }
}

func (e *ExcludePlayer) Get() []int{
    return e.playerIDs
}

func (e *ExcludePlayer) Add(playerID int) {
    if !slices.Contains(e.playerIDs, playerID) {
		e.playerIDs = append(e.playerIDs, playerID)
	}
}

func (e *ExcludePlayer) Remove(playerID int) {
    index := slices.Index(e.playerIDs, playerID)
	if index != -1 {
		e.playerIDs = append(e.playerIDs[:index], e.playerIDs[index+1:]...)
	}
}
