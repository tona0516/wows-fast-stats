package service

import (
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
	return h.storage.BattleHistoryKeys()
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
	slices.Sort(keys)
	removeKeys := keys[:uint(len(keys))-maxHistories]
	for _, v := range removeKeys {
		_ = h.storage.DeleteBattleHistory(v)
	}

	return h.storage.WriteBattleHistory(battle)
}
