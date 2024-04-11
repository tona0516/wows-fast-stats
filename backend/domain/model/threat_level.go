package model

import (
	"math"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
	"golang.org/x/exp/slices"
)

//nolint:gochecknoglobals
var specialAAShips = map[int]struct {
	average     float64
	coefficient float64
}{
	4179539920: {average: 4.2, coefficient: 0.1},   // Minotaur
	4273911792: {average: 4.4, coefficient: 0.1},   // Des Moines
	4180588496: {average: 3.3, coefficient: 0.075}, // Neptune
	4277057520: {average: 3.6, coefficient: 0.05},  // Baltimore
	3762206160: {average: 2.8, coefficient: 0.05},  // Kutuzov
	4280203248: {average: 2.3, coefficient: 0.05},  // New Orleans
	3553540080: {average: 5.9, coefficient: 0.1},   // Flint
	4288591856: {average: 3.3, coefficient: 0.1},   // Atlanta
	3763255248: {average: 1.5, coefficient: 0.025}, // Belfast
	4282300400: {average: 1.8, coefficient: 0.05},  // Pensacola
	4287543280: {average: 3.0, coefficient: 0.05},  // Cleveland
	4272830448: {average: 1.1, coefficient: 0.025}, // Fletcher
	4264441840: {average: 0.7, coefficient: 0.025}, // Sims
	4074649040: {average: 1.9, coefficient: 0.025}, // Grozovoi
	4181604048: {average: 0.7, coefficient: 0.025}, // Akizuki
}

//nolint:gochecknoglobals
var specialShipScores = map[int]float64{
	3553540080: 1.25,  // Flint
	3551410160: 1.3,   // Black
	3763255248: 1.2,   // Belfast
	3763320816: 1.25,  // Saipan
	4293866960: 1.15,  // Nikolai
	4267587280: 1.175, // KamikazeR
	4293801424: 1.2,   // Gremyashchy
	4255037136: 1.15,  // Atago
	4180555568: 1.075, // Z-46
	3764336624: 1.075, // Arizona
	3761190896: 1.125, // Missouri
	4272830448: 1.1,   // Fletcher
	4264441840: 1.075, // Sims
	4182718256: 1.1,   // Gneisenau
	3763287856: 1.075, // Scharnhorst
	4179572528: 1.1,   // Großer Kurfürst
	4281219056: 1.05,  // Gearing
	4279219920: 1.05,  // Taiho
	4277122768: 1.05,  // Hakuryu
	4179506992: 1.1,   // Z-52
	4273911792: 1.05,  // Des Moines
	3762206160: 1.15,  // Kutuzov
	3552491216: 1.1,   // ARP Takao

	// マイナス補正
	4281317360: 0.8,   // Essex
	4282365936: 0.85,  // Lexington
	4284463088: 0.8,   // Ranger
	4282300400: 0.9,   // Pensacola
	4277057520: 0.925, // Baltimore
	4288657392: 0.85,  // Independence
	4183701200: 0.875, // Fubuki
	4184749776: 0.875, // Mutsuki
	4076746448: 0.9,   // Kagero
	4282267344: 0.9,   // Shimakaze
	4280203248: 0.925, // New Orleans
	3553539792: 0.9,   // ARP Ashigara
	4286494416: 0.95,  // Myoko
	3522082512: 0.9,   // ARP Nachi
	3543054032: 0.95,  // Southern Dragon
	4182685488: 0.9,   // Yorck
	4288591856: 0.9,   // Atlanta
	4183734064: 0.9,   // Nürnberg
	3762206512: 0.9,   // Prinz Eugen
	4272895696: 0.9,   // Izumo
	4288559088: 0.925, // Mahan
	4182652624: 0.85,  // Akatsuki
	4180555216: 0.9,   // Kiev
	4076746192: 0.9,   // Ognevoi
	4288558800: 0.9,   // Hatsuharu
	3865982416: 0.85,  // Tashkent
	4074649040: 0.85,  // Grozovoi
	4184782672: 0.9,   // Émile Bertin
	4183734096: 0.85,  // La Galissonnière
	4182685520: 0.925, // Algérie
	4181636944: 0.9,   // Charles Martel
	4179539792: 0.875, // Henri IV
	3763320528: 0.9,   // Kaga
}

type shipClassStdValues struct {
	damage       float64 // 目標与ダメージ基準値
	aaScore      float64 // 制空目標基準値
	survivedRate float64 // 生存率基準目標値
	influence    float64 // 影響度係数
}

type ThreatLevel struct {
	useShipID     int
	warships      Warships
	tempArenaInfo TempArenaInfo
	statistics    ThreatLevelStatistics
	playerName    string
}

type ThreatLevelStatistics struct {
	PlayerAvgDamage  float64
	PlayerAvgKill    float64
	PlayerKdRate     float64
	PlayerWinRate    float64
	PlayerBattles    uint
	ShipAvgDamage    float64
	ShipWinRate      float64
	ShipSurvivedRate float64
	ShipPlanesKilled float64
	ShipBattles      uint
}

type threatLevelBattleDetail struct {
	IsCVMatch  bool
	BottomTier uint
	TopTier    uint
}

func NewThreatLevel(
	useShipID int,
	warships Warships,
	tempArenaInfo TempArenaInfo,
	statistics ThreatLevelStatistics,
	playerName string,
) *ThreatLevel {
	return &ThreatLevel{
		useShipID:     useShipID,
		warships:      warships,
		tempArenaInfo: tempArenaInfo,
		statistics:    statistics,
		playerName:    playerName,
	}
}

func (t *ThreatLevel) Calculate() float64 {
	// 戦闘情報の取得
	battleDetail, err := t.battleDetail()
	if err != nil {
		return -1
	}

	// プレイヤー総合補正指数
	playerOverallScore := t.playerOverallScore()

	// 艦成績補正
	playerShipScore := t.playerShipScore()

	// 最後に数値の幅を作る係数を設定
	playerTotalSkillScore := (playerOverallScore + playerShipScore) * 0.5

	// 特に補正が必要だと思う艦級を含めた脅威度補正
	shipClassScore := t.shipClassScore()

	// AA特化艦補正
	shipAAIndex := t.antiAirCoefficient()

	// 脅威レベルの算出
	threatLevel := math.Round((playerTotalSkillScore + 1) * shipClassScore * 10000)

	// マッチのおける脅威レベルの補正
	threatLevelInMatch := t.correctBasedOnMatch(threatLevel, shipAAIndex, battleDetail)

	return threatLevelInMatch
}

func (t *ThreatLevel) battleDetail() (threatLevelBattleDetail, error) {
	isCVMatch := false
	matchTiers := make([]uint, 0)
	for _, vehicle := range t.tempArenaInfo.Vehicles {
		warship, ok := t.warships[vehicle.ShipID]
		if !ok {
			continue
		}

		if warship.Type == ShipTypeCV {
			isCVMatch = true
		}

		matchTiers = append(matchTiers, warship.Tier)
	}

	if len(matchTiers) == 0 {
		return threatLevelBattleDetail{}, failure.New(apperr.InvalidTempArenaInfo)
	}

	return threatLevelBattleDetail{
		IsCVMatch:  isCVMatch,
		BottomTier: slices.Min(matchTiers),
		TopTier:    slices.Max(matchTiers),
	}, nil
}

func (t *ThreatLevel) baseDamageScore() float64 {
	avgDamage := toInteger(t.statistics.PlayerAvgDamage)
	return round((limitedValue(avgDamage/40000, 1.5, 0.5) - 1) / 2)
}

//nolint:nonamedreturns
func (t *ThreatLevel) killScore() (killScore float64, kdRateScore float64) {
	if t.statistics.PlayerBattles == 0 {
		return round(-0.3 / 5), round(-0.5 / 5)
	}

	killScore = limitedValue(t.statistics.PlayerAvgKill, 1.5, 0.5) - 0.8
	kdRateScore = limitedValue(toInteger(t.statistics.PlayerKdRate), 3, 0.7) - 1.2

	// KDRは高いのにKPRは低い（≒芋）は逆補正に
	if killScore < 0 && kdRateScore > 0 {
		kdRateScore *= -1
	}

	return round(killScore / 5), round(kdRateScore / 5)
}

func (t *ThreatLevel) winRateScore() float64 {
	return round((limitedValue(t.statistics.PlayerWinRate, 60, 30) - 50) / 100 * 3)
}

func (t *ThreatLevel) playerOverallScore() float64 {
	// プレイヤー総合補正指数を一旦算出
	killScore, kdRateScore := t.killScore()

	playerGeneralScore := round(t.baseDamageScore() + killScore + kdRateScore + t.winRateScore())

	// プレイ回数関連補正処理
	battlesCountScore := limitedValue(float64(t.statistics.PlayerBattles)/1000, 3, 0.5)

	switch {
	case battlesCountScore > 3:
		// 戦闘回数は多いのに低戦績(疑いの余地がない無能)
		if playerGeneralScore < 0 {
			battlesCountScore = 1 + playerGeneralScore*0.5
		}
	case battlesCountScore > 2:
		// 戦闘回数はそこそこあって、でも戦績低い
		if playerGeneralScore < 0 {
			battlesCountScore = 1 + playerGeneralScore*0.75
		}
	default:
		// 顕著に成績が高めの場合&リロールor戦績リセット勢と思われるパターンへの対応
		switch {
		case playerGeneralScore > 0.7:
			battlesCountScore = 3
		case playerGeneralScore > 0.5:
			battlesCountScore = 2
		case playerGeneralScore > 0.3:
			battlesCountScore = 1.5
		case playerGeneralScore > 0.1:
			battlesCountScore = 1.2
		}
	}

	playerGeneralScore = round((battlesCountScore-1.5)*0.05) + playerGeneralScore

	return playerGeneralScore
}

//nolint:cyclop
func (t *ThreatLevel) playerShipScore() float64 {
	// 艦の戦績が無い場合は総合戦績から補正なしで
	if t.statistics.ShipBattles <= 3 {
		return 0.75
	}

	warship, ok := t.warships[t.useShipID]
	if !ok {
		return 0.75
	}
	shipType := warship.Type
	shipTier := float64(warship.Tier)

	// 艦の使用回数が3回より多い場合のみ参考に
	// 艦種ごとの基準を設定
	var std shipClassStdValues
	//nolint:exhaustive
	switch shipType {
	case ShipTypeDD:
		std = shipClassStdValues{
			damage:       shipTier * 4000,
			aaScore:      0,
			survivedRate: 0.4,
			influence:    1.3,
		}
	case ShipTypeCL:
		std = shipClassStdValues{
			damage:       shipTier * 6000,
			aaScore:      0,
			survivedRate: 0.5,
			influence:    1,
		}
	case ShipTypeBB:
		std = shipClassStdValues{
			damage:       shipTier * 7200,
			aaScore:      0,
			survivedRate: 0.45,
			influence:    1.1,
		}
	case ShipTypeCV:
		std = shipClassStdValues{
			damage:       shipTier * 8000,
			aaScore:      (shipTier - 2) * 3.5,
			survivedRate: 0.7,
			influence:    1.7,
		}
	default:
		std = shipClassStdValues{
			damage:       0,
			aaScore:      0,
			survivedRate: 0,
			influence:    1,
		}
	}

	if std.damage < 25000 {
		std.damage = 25000
	}

	// ダメージ補正指数
	shipDamageScore := round(limitedValue(toInteger(t.statistics.ShipAvgDamage)/std.damage, 1.5, 0.1) - 1)

	// 生存率補正指数
	// 艦種ごとの影響度を考慮して補正
    surviveRate := t.statistics.ShipSurvivedRate / 100 // %で与えられるため0~100に変換する
	limitedSurviveRate := limitedValue(surviveRate, std.survivedRate+0.1, std.survivedRate-0.1)
	shipServiveScore := round(limitedSurviveRate * std.influence / 1.5)

	// 空母の場合、制空の補正を考慮
	var shipAAScore float64
	if shipType == ShipTypeCV {
		shipAAScore = t.statistics.ShipPlanesKilled / std.aaScore
		if shipAAScore > 1 {
			shipAAScore = round((shipAAScore - 1) / 2)
		} else {
			//nolint:gocritic
			shipAAScore = round(shipAAScore - 1)
		}
		shipAAScore = limitedValue(shipAAScore, 0.2, -0.3)
	}

	// 艦勝率補正指数
    winRate := t.statistics.ShipWinRate / 100 // %で与えられるため0~100に変換する
	shipWinRateScore := limitedValue(winRate, 0.6, 0.4) - 0.5
	// 空母の場合、勝率の影響を3割増加
	if shipType == ShipTypeCV {
		shipWinRateScore = shipWinRateScore * 1.3
	}
	shipWinRateScore = round(shipWinRateScore * std.influence * 1.5)

	return round(shipDamageScore + shipServiveScore + shipAAScore + shipWinRateScore)
}

func (t *ThreatLevel) antiAirCoefficient() float64 {
	specialAAShip, ok := specialAAShips[t.useShipID]
	if !ok {
		return 1.0
	}

	planesKilledRatio := t.statistics.ShipPlanesKilled / specialAAShip.average * 1.5
	if planesKilledRatio > 1 {
		planesKilledRatio = round((planesKilledRatio-1)/4) + 1
	}
	limitedPlanesKilledRatio := limitedValue(planesKilledRatio, 1.2, 1)

	return limitedPlanesKilledRatio + specialAAShip.coefficient
}

func (t *ThreatLevel) shipClassScore() float64 {
	result := 1.0

	// 存在しない艦の場合はデフォルト値
	warship, ok := t.warships[t.useShipID]
	if !ok {
		return result
	}
	shipTier := warship.Tier

	// 特別補正艦はベースをその値にする
	score, ok := specialShipScores[t.useShipID]
	if ok {
		result = score
	}

	// Tier10艦は性能ジャンプが大きいので追加補正
	if shipTier == 10 {
		return result + 0.15
	}

	return result
}

func (t *ThreatLevel) correctBasedOnMatch(
	threatLevel float64,
	shipAAIndex float64,
	battleDetail threatLevelBattleDetail,
) float64 {
	result := threatLevel
	if battleDetail.IsCVMatch {
		result = math.Round(threatLevel * shipAAIndex)
	}

	warship, ok := t.warships[t.useShipID]
	if !ok {
		return result
	}
	shipTier := warship.Tier

	if shipTier == battleDetail.TopTier {
		result = math.Round(result * 1.1)
	}
	if shipTier == battleDetail.BottomTier {
		result = math.Round(result * 0.9)
	}

	return result
}

func limitedValue(value float64, max float64, min float64) float64 {
	if value > max {
		return max
	}
	if value < min {
		return min
	}
	return value
}

func round(value float64) float64 {
	pow := math.Pow(10, 4)
	return toInteger(value*pow) / pow
}

func toInteger(value float64) float64 {
	return float64(int(value))
}
