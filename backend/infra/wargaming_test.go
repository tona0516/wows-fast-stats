package infra

import (
	"changeme/backend/vo"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWargaming_AccountInfo(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{}, &mockLogger{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGAccountInfo]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(vo.WGAccountInfo{}, nil)
	wargaming.accountInfoClient = mockAPIClient

	result, err := wargaming.AccountInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, vo.WGAccountInfo{}, result)
}

func TestWargaming_AccountList(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{}, &mockLogger{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGAccountList]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(vo.WGAccountList{}, nil)
	wargaming.accountListClient = mockAPIClient

	result, err := wargaming.AccountList([]string{"player_1", "player_2"})

	assert.NoError(t, err)
	assert.Equal(t, vo.WGAccountList{}, result)
}

func TestWargaming_ClansAccountInfo(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{}, &mockLogger{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGClansAccountInfo]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(vo.WGClansAccountInfo{}, nil)
	wargaming.clansAccountInfoClient = mockAPIClient

	result, err := wargaming.ClansAccountInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, vo.WGClansAccountInfo{}, result)
}

func TestWargaming_ClansInfo(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{}, &mockLogger{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGClansInfo]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(vo.WGClansInfo{}, nil)
	wargaming.clansInfoClient = mockAPIClient

	result, err := wargaming.ClansInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, vo.WGClansInfo{}, result)
}

func TestWargaming_ShipsStats(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{}, &mockLogger{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGShipsStats]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(vo.WGShipsStats{}, nil)
	wargaming.shipsStatsClient = mockAPIClient

	result, err := wargaming.ShipsStats(123)

	assert.NoError(t, err)
	assert.Equal(t, vo.WGShipsStats{}, result)
}

func TestWargaming_EncycShips(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{}, &mockLogger{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGEncycShips]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(vo.WGEncycShips{}, nil)
	wargaming.encycShipsClient = mockAPIClient

	result, err := wargaming.EncycShips(1)

	assert.NoError(t, err)
	assert.Equal(t, vo.WGEncycShips{}, result)
}

func TestWargaming_EncycInfo(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{}, &mockLogger{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGEncycInfo]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(vo.WGEncycInfo{}, nil)
	wargaming.encycInfoClient = mockAPIClient

	result, err := wargaming.EncycInfo()

	assert.NoError(t, err)
	assert.Equal(t, vo.WGEncycInfo{}, result)
}

func TestWargaming_BattleArena(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{}, &mockLogger{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGBattleArenas]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(vo.WGBattleArenas{}, nil)
	wargaming.battleArenasClient = mockAPIClient

	result, err := wargaming.BattleArenas()

	assert.NoError(t, err)
	assert.Equal(t, vo.WGBattleArenas{}, result)
}

func TestWargaming_BattleTypes(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{}, &mockLogger{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGBattleTypes]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(vo.WGBattleTypes{}, nil)
	wargaming.battleTypesClient = mockAPIClient

	result, err := wargaming.BattleTypes()

	assert.NoError(t, err)
	assert.Equal(t, vo.WGBattleTypes{}, result)
}

func TestWargaming_AccountInfo_異常系_リトライなし(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{}, &mockLogger{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGAccountInfo]{}
	response := vo.WGAccountInfo{
		Status: "error",
		Error: vo.WGError{
			Message: "INVALID_APPLICATION_ID",
		},
	}
	mockAPIClient.On("GetRequest", mock.Anything).Return(response, nil)
	wargaming.accountInfoClient = mockAPIClient

	_, err := wargaming.AccountInfo([]int{123, 456})

	assert.EqualError(t, err, fmt.Sprintf("I100 AccountInfo %s", response.Error.Message))
	mockAPIClient.AssertNumberOfCalls(t, "GetRequest", 1)
}

func TestWargaming_AccountInfo_正常系_最大リトライ(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{}, &mockLogger{})
	wargaming.SetAppID("your-app-id")

	messages := []string{
		"REQUEST_LIMIT_EXCEEDED",
		"SOURCE_NOT_AVAILABLE",
	}

	for _, v := range messages {
		mockAPIClient := &mockAPIClient[vo.WGAccountInfo]{}
		errorResponse := vo.WGAccountInfo{
			Error: vo.WGError{
				Message: v,
			},
		}
		successResponse := vo.WGAccountInfo{}
		mockAPIClient.On("GetRequest", mock.Anything).Return(errorResponse, nil).Times(3)
		mockAPIClient.On("GetRequest", mock.Anything).Return(successResponse, nil)
		wargaming.accountInfoClient = mockAPIClient

		_, err := wargaming.AccountInfo([]int{123, 456})

		assert.NoError(t, err)
		mockAPIClient.AssertNumberOfCalls(t, "GetRequest", 4)
	}
}

func TestWargaming_AccountInfo_異常系_最大リトライ(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{}, &mockLogger{})
	wargaming.SetAppID("your-app-id")

	messages := []string{
		"REQUEST_LIMIT_EXCEEDED",
		"SOURCE_NOT_AVAILABLE",
	}

	for _, v := range messages {
		mockAPIClient := &mockAPIClient[vo.WGAccountInfo]{}
		response := vo.WGAccountInfo{
			Error: vo.WGError{
				Message: v,
			},
		}
		mockAPIClient.On("GetRequest", mock.Anything).Return(response, nil)
		wargaming.accountInfoClient = mockAPIClient

		_, err := wargaming.AccountInfo([]int{123, 456})

		assert.EqualError(t, err, fmt.Sprintf("I100 AccountInfo %s", v))
		mockAPIClient.AssertNumberOfCalls(t, "GetRequest", 4)
	}
}
