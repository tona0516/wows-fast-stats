package service

import (
	"changeme/backend/apperr"
	"changeme/backend/infra"
	"changeme/backend/vo"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest
func TestPrepare_FetchCachable_Success(t *testing.T) {
	mockWargamingRepo := &mockWargamingRepo{}
	mockWargamingRepo.On("SetAppID", "").Return()
	mockWargamingRepo.On("EncyclopediaShips", 1).Return(vo.WGEncyclopediaShips{
		Status: "ok",
		Meta: struct {
			PageTotal int "json:\"page_total\""
			Page      int "json:\"page\""
		}{
			Page:      1,
			PageTotal: 2,
		},
		Data: map[int]vo.WGEncyclopediaShipsData{
			10: {Tier: 1, Type: "AirCarrier", Name: "Ship 1", Nation: "Japan"},
			20: {Tier: 2, Type: "Battleship", Name: "Ship 2", Nation: "USA"},
		},
	}, nil)
	mockWargamingRepo.On("EncyclopediaShips", 2).Return(vo.WGEncyclopediaShips{
		Status: "ok",
		Meta: struct {
			PageTotal int "json:\"page_total\""
			Page      int "json:\"page\""
		}{
			Page:      2,
			PageTotal: 2,
		},
		Data: map[int]vo.WGEncyclopediaShipsData{
			30: {Tier: 3, Type: "Cruiser", Name: "Ship 3", Nation: "Russia"},
			40: {Tier: 4, Type: "Destroyer", Name: "Ship 4", Nation: "United_Kingdom"},
		},
	}, nil)
	mockWargamingRepo.On("BattleArenas").Return(vo.WGBattleArenas{
		Status: "ok",
		Data: map[int]vo.WGBattleArenasData{
			1: {Name: "Map 1"},
		},
	}, nil)
	mockWargamingRepo.On("BattleTypes").Return(vo.WGBattleTypes{
		Status: "ok",
		Data: map[string]vo.WGBattleTypesData{
			"pvp": {Name: "ランダム戦"},
		},
	}, nil)

	mockUnregisteredRepo := &mockUnregisteredRepo{}
	mockUnregisteredRepo.On("Warship").Return(map[int]vo.Warship{
		50: {Tier: 5, Type: "Submarine", Name: "Ship 5", Nation: "Pan_Asia"},
	}, nil)

	mockNumbersRepo := &mockNumbersRepo{}
	mockNumbersRepo.On("ExpectedStats").Return(vo.NSExpectedStats{
		Time: int(time.Now().UnixMilli()),
		Data: map[int]vo.NSExpectedStatsData{
			10: {AverageDamageDealt: 10000, AverageFrags: 1, WinRate: 0.5},
		},
	}, nil)

	caches := infra.NewCaches("cache_test")
	defer os.RemoveAll(caches.Dir)

	prepare := NewPrepare(
		5,
		mockWargamingRepo,
		mockNumbersRepo,
		mockUnregisteredRepo,
		*caches,
	)

	errChan := make(chan error)
	go prepare.FetchCachable(vo.UserConfig{}, errChan)
	err := <-errChan

	assert.NoError(t, err)
	mockWargamingRepo.AssertCalled(t, "EncyclopediaShips", 1)
	mockWargamingRepo.AssertCalled(t, "EncyclopediaShips", 2)
	mockWargamingRepo.AssertCalled(t, "BattleArenas")
	mockWargamingRepo.AssertCalled(t, "BattleTypes")
	mockUnregisteredRepo.AssertCalled(t, "Warship")
	mockNumbersRepo.AssertCalled(t, "ExpectedStats")

	assert.FileExists(t, filepath.Join(caches.Dir, caches.Warship.Name+".bin"))
	assert.FileExists(t, filepath.Join(caches.Dir, caches.ExpectedStats.Name+".bin"))
	assert.FileExists(t, filepath.Join(caches.Dir, caches.BattleArenas.Name+".bin"))
	assert.FileExists(t, filepath.Join(caches.Dir, caches.BattleTypes.Name+".bin"))
}

//nolint:paralleltest
func TestPrepare_FetchCachable_Failure(t *testing.T) {
	mockWargamingRepo := &mockWargamingRepo{}
	mockWargamingRepo.On("SetAppID", "").Return()
	mockWargamingRepo.On("EncyclopediaShips", 1).Return(vo.WGEncyclopediaShips{
		Status: "ok",
		Meta: struct {
			PageTotal int "json:\"page_total\""
			Page      int "json:\"page\""
		}{
			Page:      1,
			PageTotal: 2,
		},
		Data: map[int]vo.WGEncyclopediaShipsData{
			10: {Tier: 1, Type: "AirCarrier", Name: "Ship 1", Nation: "Japan"},
			20: {Tier: 2, Type: "Battleship", Name: "Ship 2", Nation: "USA"},
		},
	}, nil)
	mockWargamingRepo.On("EncyclopediaShips", 2).Return(vo.WGEncyclopediaShips{
		Status: "ok",
		Meta: struct {
			PageTotal int "json:\"page_total\""
			Page      int "json:\"page\""
		}{
			Page:      2,
			PageTotal: 2,
		},
		Data: map[int]vo.WGEncyclopediaShipsData{
			30: {Tier: 3, Type: "Cruiser", Name: "Ship 3", Nation: "Russia"},
			40: {Tier: 4, Type: "Destroyer", Name: "Ship 4", Nation: "United_Kingdom"},
		},
	}, nil)
	mockWargamingRepo.On("BattleArenas").Return(vo.WGBattleArenas{
		Status: "ok",
		Data: map[int]vo.WGBattleArenasData{
			1: {Name: "Map 1"},
		},
	}, nil)
	mockWargamingRepo.On("BattleTypes").Return(vo.WGBattleTypes{
		Status: "ok",
		Data: map[string]vo.WGBattleTypesData{
			"pvp": {Name: "ランダム戦"},
		},
	}, nil)

	mockUnregisteredRepo := &mockUnregisteredRepo{}
	mockUnregisteredRepo.On("Warship").Return(map[int]vo.Warship{
		50: {Tier: 5, Type: "Submarine", Name: "Ship 5", Nation: "Pan_Asia"},
	}, nil)

	// Note: occur error
	expectedErr := errors.WithStack(apperr.Ns.Req)
	mockNumbersRepo := &mockNumbersRepo{}
	mockNumbersRepo.On("ExpectedStats").Return(vo.NSExpectedStats{}, expectedErr)

	caches := infra.NewCaches("cache_test")
	defer os.RemoveAll(caches.Dir)

	prepare := NewPrepare(
		5,
		mockWargamingRepo,
		mockNumbersRepo,
		mockUnregisteredRepo,
		*caches,
	)

	errChan := make(chan error)
	go prepare.FetchCachable(vo.UserConfig{}, errChan)
	err := <-errChan

	assert.EqualError(t, err, expectedErr.Error())
	mockWargamingRepo.AssertCalled(t, "EncyclopediaShips", 1)
	mockWargamingRepo.AssertCalled(t, "EncyclopediaShips", 2)
	mockWargamingRepo.AssertCalled(t, "BattleArenas")
	mockWargamingRepo.AssertCalled(t, "BattleTypes")
	mockUnregisteredRepo.AssertCalled(t, "Warship")
	mockNumbersRepo.AssertCalled(t, "ExpectedStats")

	assert.FileExists(t, filepath.Join(caches.Dir, caches.Warship.Name+".bin"))
	assert.NoFileExists(t, filepath.Join(caches.Dir, caches.ExpectedStats.Name+".bin"))
	assert.FileExists(t, filepath.Join(caches.Dir, caches.BattleArenas.Name+".bin"))
	assert.FileExists(t, filepath.Join(caches.Dir, caches.BattleTypes.Name+".bin"))
}
