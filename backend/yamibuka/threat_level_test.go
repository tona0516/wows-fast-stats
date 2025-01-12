package yamibuka

import (
	"testing"
	"wfs/backend/data"
	"wfs/backend/domain/model"

	"github.com/stretchr/testify/assert"
)

const (
	allowableDelta = 1.0
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

	actual := CalculateThreatLevel(NewThreatLevelFactor(
		0,
		model.TempArenaInfo{
			Vehicles: []model.Vehicle{
				{ShipID: shipIDNagato},
				{ShipID: shipIDYorktown},
				{ShipID: shipIDKitakaze},
			},
		},
		model.Warships{
			shipIDNagato: {
				Name: "長門",
				Tier: 7,
				Type: model.ShipTypeBB,
			},
			shipIDYorktown: {
				Name: "Yorktown",
				Tier: 8,
				Type: model.ShipTypeCV,
			},
			shipIDKitakaze: {
				Name: "北風",
				Tier: 9,
				Type: model.ShipTypeDD,
			},
		},
		shipIDYorktown,
		17,
		71540,
		52.94117647058824,
		76.47058823529411,
		8.176470588235293,
		18940,
		67099,
		61.58,
		1.0761351636747625,
		2.34,
	))
	expected := data.ThreatLevel{
		Raw:      18111,
		Modified: 18111,
	}

	assert.InDelta(t, expected.Raw, actual.Raw, allowableDelta)
	assert.InDelta(t, expected.Modified, actual.Modified, allowableDelta)
}

func TestThreatLevel_CalculateThreatLevel_BB_CVあり_Tierトップ(t *testing.T) {
	t.Parallel()

	actual := CalculateThreatLevel(NewThreatLevelFactor(
		0,
		model.TempArenaInfo{
			Vehicles: []model.Vehicle{
				{ShipID: shipIDMutsuki},
				{ShipID: shipIDRanger},
				{ShipID: shipIDSinop},
			},
		},
		model.Warships{
			shipIDMutsuki: {
				Name: "睦月",
				Tier: 5,
				Type: model.ShipTypeDD,
			},
			shipIDRanger: {
				Name: "Ranger",
				Tier: 6,
				Type: model.ShipTypeCV,
			},
			shipIDSinop: {
				Name: "Sinop",
				Tier: 7,
				Type: model.ShipTypeBB,
			},
		},
		shipIDSinop,
		273,
		85237,
		67.76556776556777,
		55.31135531135531,
		5.47985347985348,
		18940,
		67099,
		61.58,
		1.0761351636747625,
		2.34,
	))
	expected := data.ThreatLevel{
		Raw:      19543,
		Modified: 21497,
	}

	assert.InDelta(t, expected.Raw, actual.Raw, allowableDelta)
	assert.InDelta(t, expected.Modified, actual.Modified, allowableDelta)
}

func TestThreatLevel_CalculateThreatLevel_CL_CVなし_Tierミドル(t *testing.T) {
	t.Parallel()

	actual := CalculateThreatLevel(NewThreatLevelFactor(
		0,
		model.TempArenaInfo{
			Vehicles: []model.Vehicle{
				{ShipID: shipIDYoshino},
			},
		},
		model.Warships{
			shipIDYoshino: {
				Name: "吉野",
				Tier: 10,
				Type: model.ShipTypeCL,
			},
		},
		shipIDYoshino,
		54,
		117010,
		55.55555555555556,
		61.11111111111112,
		5.888888888888889,
		18940,
		67099,
		61.58,
		1.0761351636747625,
		2.34,
	))
	expected := data.ThreatLevel{
		Raw:      21985,
		Modified: 24184,
	}

	assert.InDelta(t, expected.Raw, actual.Raw, allowableDelta)
	assert.InDelta(t, expected.Modified, actual.Modified, allowableDelta)
}

func TestThreatLevel_CalculateThreatLevel_DD_CVあり_Tierボトム_特殊補正艦(t *testing.T) {
	t.Parallel()

	shipIDSims := 4264441840
	shipIDYorktown := 4265588720
	shipIDAlaska := 3760109552

	actual := CalculateThreatLevel(NewThreatLevelFactor(
		0,
		model.TempArenaInfo{
			Vehicles: []model.Vehicle{
				{ShipID: shipIDSims},
				{ShipID: shipIDYorktown},
				{ShipID: shipIDAlaska},
			},
		},
		model.Warships{
			shipIDSims: {
				Name: "Sims",
				Tier: 7,
				Type: model.ShipTypeDD,
			},
			shipIDYorktown: {
				Name: "Yorktown",
				Tier: 8,
				Type: model.ShipTypeCV,
			},
			shipIDAlaska: {
				Name: "Alaska",
				Tier: 9,
				Type: model.ShipTypeCL,
			},
		},
		shipIDSims,
		82,
		36399,
		71.95,
		50,
		4.96,
		18940,
		67099,
		61.58,
		1.0761351636747625,
		2.34,
	))
	expected := data.ThreatLevel{
		Raw:      20255,
		Modified: 22331,
	}

	assert.InDelta(t, expected.Raw, actual.Raw, allowableDelta)
	assert.InDelta(t, expected.Modified, actual.Modified, allowableDelta)
}
