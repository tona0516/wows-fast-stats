package service

import (
	"fmt"
	"slices"
	"wfs/backend/data"
	"wfs/backend/repository"
)

type BattleHistory struct {
	storage repository.StorageInterface
	logger  repository.LoggerInterface
}

func NewBattleHistory(
	storage repository.StorageInterface,
	logger repository.LoggerInterface,
) *BattleHistory {
	return &BattleHistory{
		storage: storage,
		logger:  logger,
	}
}

func (h *BattleHistory) AllKeys() ([]string, error) {
	keys, error := h.storage.BattleHistoryKeys()
	slices.Sort(keys)
	slices.Reverse(keys)
	return keys, error
}

func (h *BattleHistory) Get(key string) (data.Battle, error) {
	return h.storage.BattleHistory(key)
}

func (h *BattleHistory) Add(battle data.Battle) error {
	config, err := h.storage.UserConfigV2()
	if err != nil {
		return err
	}

	keys, err := h.AllKeys()
	if err != nil {
		return err
	}

	maxHistories := config.MaxBattleHistories
	h.logger.Debug(fmt.Sprintf("MaxBattleHistories: %d", maxHistories), nil)

	slices.Sort(keys)

	historyLen := uint(len(keys))
	if historyLen > maxHistories {
		removeKeys := keys[:historyLen-maxHistories]
		for _, v := range removeKeys {
			_ = h.storage.DeleteBattleHistory(v)
		}
	}

	return h.storage.WriteBattleHistory(battle)
}
