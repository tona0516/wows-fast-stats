package infra

import (
	"fmt"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/vo"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWargaming_AccountInfo(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGAccountInfo]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(APIResponse[vo.WGAccountInfo]{}, nil)
	wargaming.accountInfoClient = mockAPIClient

	result, err := wargaming.AccountInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, vo.WGAccountInfo{}, result)
}

func TestWargaming_AccountList(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGAccountList]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(APIResponse[vo.WGAccountList]{}, nil)
	wargaming.accountListClient = mockAPIClient

	result, err := wargaming.AccountList([]string{"player_1", "player_2"})

	assert.NoError(t, err)
	assert.Equal(t, vo.WGAccountList{}, result)
}

func TestWargaming_ClansAccountInfo(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGClansAccountInfo]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(APIResponse[vo.WGClansAccountInfo]{}, nil)
	wargaming.clansAccountInfoClient = mockAPIClient

	result, err := wargaming.ClansAccountInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, vo.WGClansAccountInfo{}, result)
}

func TestWargaming_ClansInfo(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGClansInfo]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(APIResponse[vo.WGClansInfo]{}, nil)
	wargaming.clansInfoClient = mockAPIClient

	result, err := wargaming.ClansInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, vo.WGClansInfo{}, result)
}

func TestWargaming_ShipsStats(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGShipsStats]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(APIResponse[vo.WGShipsStats]{}, nil)
	wargaming.shipsStatsClient = mockAPIClient

	result, err := wargaming.ShipsStats(123)

	assert.NoError(t, err)
	assert.Equal(t, vo.WGShipsStats{}, result)
}

func TestWargaming_EncycShips(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGEncycShips]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(APIResponse[vo.WGEncycShips]{}, nil)
	wargaming.encycShipsClient = mockAPIClient

	result, err := wargaming.EncycShips(1)

	assert.NoError(t, err)
	assert.Equal(t, vo.WGEncycShips{}, result)
}

func TestWargaming_BattleArena(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGBattleArenas]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(APIResponse[vo.WGBattleArenas]{}, nil)
	wargaming.battleArenasClient = mockAPIClient

	result, err := wargaming.BattleArenas()

	assert.NoError(t, err)
	assert.Equal(t, vo.WGBattleArenas{}, result)
}

func TestWargaming_BattleTypes(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{})
	wargaming.SetAppID("your-app-id")

	mockAPIClient := &mockAPIClient[vo.WGBattleTypes]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(APIResponse[vo.WGBattleTypes]{}, nil)
	wargaming.battleTypesClient = mockAPIClient

	result, err := wargaming.BattleTypes()

	assert.NoError(t, err)
	assert.Equal(t, vo.WGBattleTypes{}, result)
}

func TestWargaming_Test(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{})

	mockAPIClient := &mockAPIClient[vo.WGEncycInfo]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(APIResponse[vo.WGEncycInfo]{}, nil)
	wargaming.encycInfoClient = mockAPIClient

	valid, err := wargaming.Test("hoge")
	assert.True(t, valid)
	assert.NoError(t, err)
}

func TestWargaming_AccountInfo_異常系_リトライなし(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{})
	wargaming.SetAppID("your-app-id")

	message := "INVALID_APPLICATION_ID"
	mockAPIClient := &mockAPIClient[vo.WGAccountInfo]{}

	response := APIResponse[vo.WGAccountInfo]{
		Body: vo.WGAccountInfo{
			Status: "error",
			Error:  vo.WGError{Message: message},
		},
		BodyString: `{
            \"status\":\"error\",
            \"error\":{
                \"field\":null,
                \"message\":\"INVALID_APPLICATION_ID\",
                \"code\":407,\"value\":null
            }
        }`,
	}
	mockAPIClient.On("GetRequest", mock.Anything).Return(response, nil)
	wargaming.accountInfoClient = mockAPIClient

	_, err := wargaming.AccountInfo([]int{123, 456})

	assert.EqualError(t, err, apperr.New(apperr.WargamingAPIError, errors.New(response.BodyString)).Error())
	mockAPIClient.AssertNumberOfCalls(t, "GetRequest", 1)
}

func TestWargaming_AccountInfo_正常系_最大リトライ(t *testing.T) {
	t.Parallel()

	wargaming := NewWargaming(vo.WGConfig{})
	wargaming.SetAppID("your-app-id")

	messages := []string{
		"REQUEST_LIMIT_EXCEEDED",
		"SOURCE_NOT_AVAILABLE",
	}

	for _, v := range messages {
		mockAPIClient := &mockAPIClient[vo.WGAccountInfo]{}
		errorResponse := APIResponse[vo.WGAccountInfo]{
			Body: vo.WGAccountInfo{Error: vo.WGError{Message: v}},
		}
		successResponse := APIResponse[vo.WGAccountInfo]{
			Body: vo.WGAccountInfo{},
		}
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

	wargaming := NewWargaming(vo.WGConfig{})
	wargaming.SetAppID("your-app-id")

	messages := []string{
		"REQUEST_LIMIT_EXCEEDED",
		"SOURCE_NOT_AVAILABLE",
	}

	for _, message := range messages {
		mockAPIClient := &mockAPIClient[vo.WGAccountInfo]{}
		response := APIResponse[vo.WGAccountInfo]{
			Body: vo.WGAccountInfo{
				Error: vo.WGError{
					Message: message,
				},
			},
			BodyString: fmt.Sprintf(`{
                \"status\":\"error\",
                \"error\":{
                    \"field\":null,
                    \"message\":\"%s\",
                    \"code\":407,\"value\":null
                }
            }`, message),
		}
		mockAPIClient.On("GetRequest", mock.Anything).Return(response, nil)
		wargaming.accountInfoClient = mockAPIClient

		_, err := wargaming.AccountInfo([]int{123, 456})

		assert.EqualError(t, err, apperr.New(
			apperr.WargamingAPITemporaryUnavaillalble,
			errors.New(response.BodyString),
		).Error())
		mockAPIClient.AssertNumberOfCalls(t, "GetRequest", 4)
	}
}
