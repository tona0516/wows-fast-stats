package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBattle_正常系(t *testing.T) {
	t.Parallel()

	// Arrange
	prefetchResult := &PrefetchResult{
		Warships: Warships{
			123: {ID: 123, Name: "Yamato", Type: ShipTypeBB, Tier: 10, Nation: "japan"},
			456: {ID: 456, Name: "Iowa", Type: ShipTypeBB, Tier: 10, Nation: "usa"},
			789: {ID: 789, Name: "Conqueror", Type: ShipTypeBB, Tier: 10, Nation: "uk"},
			101: {ID: 101, Name: "Republique", Type: ShipTypeBB, Tier: 10, Nation: "france"},
		},
		BattleArenas: map[int]string{
			1: "Strait",
			2: "Sleeping Giant",
		},
		BattleTypes: map[string]string{
			"PVE": "Operation",
			"PVP": "Randoms",
		},
	}

	tempArenaInfo := TempArenaInfo{
		Vehicles: []Vehicle{
			{ShipID: 123, Name: "Player1", Relation: 0},
			{ShipID: 456, Name: "Player2", Relation: 0},
			{ShipID: 789, Name: "Enemy1", Relation: 2},
			{ShipID: 101, Name: "Enemy2", Relation: 2},
		},
		DateTime:   "01.01.2025 12:00:00",
		MapID:      1,
		MatchGroup: "PVP",
	}

	accountInfo := WGAccountInfo{
		WGResponseCommon: WGResponseCommon[map[AccountID]WGAccountInfoData]{
			Data: map[AccountID]WGAccountInfoData{
				1: {HiddenProfile: false},
				2: {HiddenProfile: false},
				3: {HiddenProfile: true},
				4: {HiddenProfile: false},
			},
		},
	}

	accountList := WGAccountList{
		WGResponseCommon: WGResponseCommon[[]WGAccountListData]{
			Data: []WGAccountListData{
				{NickName: "Player1", AccountID: 1},
				{NickName: "Player2", AccountID: 2},
				{NickName: "Enemy1", AccountID: 3},
				{NickName: "Enemy2", AccountID: 4},
			},
		},
	}

	clans := Clans{
		1: {Tag: "TAG1"},
		2: {Tag: "TAG2"},
		3: {},
		4: {Tag: "TAG3"},
	}

	allPlayerShipsStats := AllPlayerShipStats{
		1: {},
		2: {},
		3: {},
		4: {},
	}

	allPlayerShipsBadges := AllPlayerShipBadges{
		1: {},
		2: {},
		3: {},
		4: {},
	}

	// Act
	battle := NewBattle(
		prefetchResult,
		tempArenaInfo,
		accountInfo,
		accountList,
		clans,
		allPlayerShipsStats,
		allPlayerShipsBadges,
	)

	// Assert
	assert.NotNil(t, battle)
	assert.Equal(t, 2, len(battle.Teams))

	// メタデータの確認
	assert.Equal(t, "Strait", battle.Metadata.Arena)
	assert.Equal(t, "Randoms", battle.Metadata.Type)

	// 友軍チーム（Team 0）の確認
	assert.Equal(t, 2, len(battle.Teams[0].Players))
	assert.Equal(t, "Player1", battle.Teams[0].Players[0].PlayerInfo.Name)
	assert.Equal(t, "Player2", battle.Teams[0].Players[1].PlayerInfo.Name)

	// 敵軍チーム（Team 1）の確認
	assert.Equal(t, 2, len(battle.Teams[1].Players))
	assert.Equal(t, "Enemy1", battle.Teams[1].Players[0].PlayerInfo.Name)
	assert.Equal(t, "Enemy2", battle.Teams[1].Players[1].PlayerInfo.Name)

	// プレイヤー情報の確認
	assert.Equal(t, AccountID(1), battle.Teams[0].Players[0].PlayerInfo.ID)
	assert.Equal(t, "TAG1", battle.Teams[0].Players[0].PlayerInfo.Clan.Tag)
	assert.False(t, battle.Teams[0].Players[0].PlayerInfo.IsHidden)

	assert.Equal(t, AccountID(3), battle.Teams[1].Players[0].PlayerInfo.ID)
	assert.True(t, battle.Teams[1].Players[0].PlayerInfo.IsHidden)
}

func TestNewBattle_未知の戦艦(t *testing.T) {
	t.Parallel()

	// Arrange
	prefetchResult := &PrefetchResult{
		Warships: Warships{
			123: {ID: 123, Name: "Yamato", Type: ShipTypeBB, Tier: 10, Nation: "japan"},
		},
		BattleArenas: map[int]string{
			1: "Strait",
		},
		BattleTypes: map[string]string{
			"PVP": "Randoms",
		},
	}

	tempArenaInfo := TempArenaInfo{
		Vehicles: []Vehicle{
			{ShipID: 123, Name: "Player1", Relation: 0},
			{ShipID: 999, Name: "Player2", Relation: 0}, // 存在しないShipID
		},
		DateTime:   "01.01.2025 12:00:00",
		MapID:      1,
		MatchGroup: "PVP",
	}

	accountInfo := WGAccountInfo{
		WGResponseCommon: WGResponseCommon[map[AccountID]WGAccountInfoData]{
			Data: map[AccountID]WGAccountInfoData{
				1: {HiddenProfile: false},
				2: {HiddenProfile: false},
			},
		},
	}

	accountList := WGAccountList{
		WGResponseCommon: WGResponseCommon[[]WGAccountListData]{
			Data: []WGAccountListData{
				{NickName: "Player1", AccountID: 1},
				{NickName: "Player2", AccountID: 2},
			},
		},
	}

	clans := Clans{
		1: {Tag: "TAG1"},
		2: {},
	}

	allPlayerShipsStats := AllPlayerShipStats{
		1: {},
		2: {},
	}

	allPlayerShipsBadges := AllPlayerShipBadges{
		1: {},
		2: {},
	}

	// Act
	battle := NewBattle(
		prefetchResult,
		tempArenaInfo,
		accountInfo,
		accountList,
		clans,
		allPlayerShipsStats,
		allPlayerShipsBadges,
	)

	// Assert
	assert.Equal(t, 2, len(battle.Teams[0].Players))
	// 存在しない戦艦は UnknownWarship に置き換わる
	assert.Equal(t, "Yamato", battle.Teams[0].Players[0].Warship.Name)
	unknownWarship := NewUnknownWarship()
	assert.Equal(t, unknownWarship.Name, battle.Teams[0].Players[1].Warship.Name)
}
