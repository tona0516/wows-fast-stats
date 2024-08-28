package service

import (
	"testing"
	"wfs/backend/data"
	"wfs/backend/mock/repository"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestBattleHistory_AllKeys(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := repository.NewMockStorageInterface(ctrl)
	mockLogger := repository.NewMockLoggerInterface(ctrl)

	battleHistory := NewBattleHistory(mockStorage, mockLogger)

	expectedKeys := []string{"battle1", "battle2"}
	mockStorage.EXPECT().BattleHistoryKeys().Return(expectedKeys, nil)

	keys, err := battleHistory.AllKeys()

	require.NoError(t, err)
	assert.Equal(t, expectedKeys, keys)
}

func TestBattleHistory_Get(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := repository.NewMockStorageInterface(ctrl)
	mockLogger := repository.NewMockLoggerInterface(ctrl)

	battleHistory := NewBattleHistory(mockStorage, mockLogger)

	expectedBattle := data.Battle{
		Meta:  data.Meta{},
		Teams: []data.Team{},
	}
	mockStorage.EXPECT().BattleHistory("battle1").Return(expectedBattle, nil)

	battle, err := battleHistory.Get("battle1")

	require.NoError(t, err)
	assert.Equal(t, expectedBattle, battle)
}

func TestBattleHistory_Add(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := repository.NewMockStorageInterface(ctrl)
	mockLogger := repository.NewMockLoggerInterface(ctrl)

	battleHistory := NewBattleHistory(mockStorage, mockLogger)

	config := data.UserConfigV2{MaxBattleHistories: 2}
	mockStorage.EXPECT().UserConfigV2().Return(config, nil)

	keys := []string{"battle1", "battle2", "battle3"}
	mockStorage.EXPECT().BattleHistoryKeys().Return(keys, nil)
	mockStorage.EXPECT().DeleteBattleHistory("battle1").Return(nil)
	mockStorage.EXPECT().WriteBattleHistory(gomock.Any()).Return(nil)

	err := battleHistory.Add(data.Battle{
		Meta:  data.Meta{},
		Teams: []data.Team{},
	})

	require.NoError(t, err)
}

func TestBattleHistory_Add_ErrorInUserConfig(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := repository.NewMockStorageInterface(ctrl)
	mockLogger := repository.NewMockLoggerInterface(ctrl)

	battleHistory := NewBattleHistory(mockStorage, mockLogger)

	mockStorage.EXPECT().UserConfigV2().Return(data.UserConfigV2{}, badger.ErrKeyNotFound)

	err := battleHistory.Add(data.Battle{
		Meta:  data.Meta{},
		Teams: []data.Team{},
	})

	require.Error(t, err)
	assert.Equal(t, badger.ErrKeyNotFound, err)
}
