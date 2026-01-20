package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	shipIDMutsuki  = 4184749776
	shipIDRanger   = 4183799792
	shipIDSims     = 4264441840
	shipIDNagato   = 4284430032
	shipIDSinop    = 4182717904
	shipIDYorktown = 4265588720
	shipIDKitakaze = 4065212112
	shipIDAlaska   = 3760109552
	shipIDYoshino  = 3749623504
)

func TestThreatLevel_CalculateThreatLevel_CV_CVあり_Tierミドル(t *testing.T) {
	t.Parallel()

	instance := NewThreatLevelCalculator()
	input := ThreatLevelInput{
		Vehicles: []Vehicle{
			{ShipID: shipIDNagato},
			{ShipID: shipIDYorktown},
			{ShipID: shipIDKitakaze},
		},
		Warships: Warships{
			shipIDNagato: {
				Name: "長門",
				Tier: 7,
				Type: ShipTypeBB,
			},
			shipIDYorktown: {
				Name: "Yorktown",
				Tier: 8,
				Type: ShipTypeCV,
			},
			shipIDKitakaze: {
				Name: "北風",
				Tier: 9,
				Type: ShipTypeDD,
			},
		},
		ShipID:           shipIDYorktown,
		ShipBattles:      17,
		ShipDamage:       71540,
		ShipWinRate:      52.94117647058824,
		ShipSurvivedRate: 76.47058823529411,
		ShipPlanesKilled: 8.176470588235293,
		OverallBattles:   18940,
		OverallDamage:    67099,
		OverallWinRate:   61.58,
		OverallKill:      1.0761351636747625,
		OverallKdRate:    2.34,
	}
	actual := instance.Calculate(input)
	expected := ThreatLevel{
		Raw:      18111,
		Modified: 18111,
	}

	assert.Equal(t, ThreatLevelRankI, actual.Rank)
	assert.InDelta(t, expected.Raw, actual.Raw, 1.0)
	assert.InDelta(t, expected.Modified, actual.Modified, 1.0)
}

func TestThreatLevel_CalculateThreatLevel_BB_CVあり_Tierトップ(t *testing.T) {
	t.Parallel()

	instance := NewThreatLevelCalculator()
	input := ThreatLevelInput{
		Vehicles: []Vehicle{
			{ShipID: shipIDMutsuki},
			{ShipID: shipIDRanger},
			{ShipID: shipIDSinop},
		},
		Warships: Warships{
			shipIDMutsuki: {
				Name: "睦月",
				Tier: 5,
				Type: ShipTypeDD,
			},
			shipIDRanger: {
				Name: "Ranger",
				Tier: 6,
				Type: ShipTypeCV,
			},
			shipIDSinop: {
				Name: "Sinop",
				Tier: 7,
				Type: ShipTypeBB,
			},
		},
		ShipID:           shipIDSinop,
		ShipBattles:      273,
		ShipDamage:       85237,
		ShipWinRate:      67.76556776556777,
		ShipSurvivedRate: 55.31135531135531,
		ShipPlanesKilled: 5.47985347985348,
		OverallBattles:   18940,
		OverallDamage:    67099,
		OverallWinRate:   61.58,
		OverallKill:      1.0761351636747625,
		OverallKdRate:    2.34,
	}
	actual := instance.Calculate(input)
	expected := ThreatLevel{
		Raw:      19543,
		Modified: 21497,
	}

	assert.Equal(t, ThreatLevelRankI, actual.Rank)
	assert.InDelta(t, expected.Raw, actual.Raw, 1.0)
	assert.InDelta(t, expected.Modified, actual.Modified, 1.0)
}

func TestThreatLevel_CalculateThreatLevel_CL_CVなし_Tierミドル(t *testing.T) {
	t.Parallel()

	instance := NewThreatLevelCalculator()
	input := ThreatLevelInput{
		Vehicles: []Vehicle{
			{ShipID: shipIDYoshino},
		},
		Warships: Warships{
			shipIDYoshino: {
				Name: "吉野",
				Tier: 10,
				Type: ShipTypeCL,
			},
		},
		ShipID:           shipIDYoshino,
		ShipBattles:      54,
		ShipDamage:       117010,
		ShipWinRate:      55.55555555555556,
		ShipSurvivedRate: 61.11111111111112,
		ShipPlanesKilled: 5.888888888888889,
		OverallBattles:   18940,
		OverallDamage:    67099,
		OverallWinRate:   61.58,
		OverallKill:      1.0761351636747625,
		OverallKdRate:    2.34,
	}
	actual := instance.Calculate(input)
	expected := ThreatLevel{
		Raw:      21985,
		Modified: 24184,
	}

	assert.Equal(t, ThreatLevelRankV, actual.Rank)
	assert.InDelta(t, expected.Raw, actual.Raw, 1.0)
	assert.InDelta(t, expected.Modified, actual.Modified, 1.0)
}

func TestThreatLevel_CalculateThreatLevel_DD_CVあり_Tierボトム_特殊補正艦(t *testing.T) {
	t.Parallel()

	instance := NewThreatLevelCalculator()
	input := ThreatLevelInput{
		Vehicles: []Vehicle{
			{ShipID: shipIDSims},
			{ShipID: shipIDYorktown},
			{ShipID: shipIDAlaska},
		},
		Warships: Warships{
			shipIDSims: {
				Name: "Sims",
				Tier: 7,
				Type: ShipTypeDD,
			},
			shipIDYorktown: {
				Name: "Yorktown",
				Tier: 8,
				Type: ShipTypeCV,
			},
			shipIDAlaska: {
				Name: "Alaska",
				Tier: 9,
				Type: ShipTypeCL,
			},
		},
		ShipID:           shipIDSims,
		ShipBattles:      82,
		ShipDamage:       36399,
		ShipWinRate:      71.95,
		ShipSurvivedRate: 50,
		ShipPlanesKilled: 4.96,
		OverallBattles:   18940,
		OverallDamage:    67099,
		OverallWinRate:   61.58,
		OverallKill:      1.0761351636747625,
		OverallKdRate:    2.34,
	}
	actual := instance.Calculate(input)
	expected := ThreatLevel{
		Raw:      20255,
		Modified: 22331,
	}

	assert.InDelta(t, expected.Raw, actual.Raw, 1.0)
	assert.InDelta(t, expected.Modified, actual.Modified, 1.0)
}
