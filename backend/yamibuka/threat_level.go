package yamibuka

import (
	"math"
	"wfs/backend/data"

	"github.com/shopspring/decimal"
)

type specialAAShipMap map[int]struct {
	avg  float64
	coef float64
}

type specialShipScoreMap map[int]float64

func specialAAShips() specialAAShipMap {
	return specialAAShipMap{
		4179539920: {avg: 4.2, coef: 0.1},   // Minotaur
		4273911792: {avg: 4.4, coef: 0.1},   // Des Moines
		4180588496: {avg: 3.3, coef: 0.075}, // Neptune
		4277057520: {avg: 3.6, coef: 0.05},  // Baltimore
		3762206160: {avg: 2.8, coef: 0.05},  // Kutuzov
		4280203248: {avg: 2.3, coef: 0.05},  // New Orleans
		3553540080: {avg: 5.9, coef: 0.1},   // Flint
		4288591856: {avg: 3.3, coef: 0.1},   // Atlanta
		3763255248: {avg: 1.5, coef: 0.025}, // Belfast
		4282300400: {avg: 1.8, coef: 0.05},  // Pensacola
		4287543280: {avg: 3.0, coef: 0.05},  // Cleveland
		4272830448: {avg: 1.1, coef: 0.025}, // Fletcher
		4264441840: {avg: 0.7, coef: 0.025}, // Sims
		4074649040: {avg: 1.9, coef: 0.025}, // Grozovoi
		4181604048: {avg: 0.7, coef: 0.025}, // Akizuki
	}
}

func specialShipScores() specialShipScoreMap {
	return specialShipScoreMap{
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
}

func CalculateThreatLevel(f ThreatLevelFactor) data.ThreatLevel {
	// 戦闘情報の取得
	isCVMatch, topTier, bottomTier := matchInfo(f.tempArenaInfo, f.warships)

	// プレイヤー総合補正指数
	playerOverallScore := playerOverallScore(
		f.overallBattles,
		f.overallDamage,
		f.overallKill,
		f.overallKdRate,
		f.overallWinRate,
	)

	// 艦成績補正
	playerShipScore := playerShipScore(
		f.warships,
		f.shipID,
		f.shipBattles,
		f.shipDamage,
		f.shipSurvivedRate,
		f.shipPlanesKilled,
		f.shipWinRate,
	)

	// 最後に数値の幅を作る係数を設定
	playerTotalSkillScore := (playerOverallScore + playerShipScore) * 0.5

	// 特に補正が必要だと思う艦級を含めた脅威度補正
	shipClassScore := shipClassScore(f.warships, f.shipID)

	// AA特化艦補正
	shipAAIndex := antiAirCoefficient(f.shipID, f.shipPlanesKilled)

	// 脅威レベルの算出
	raw := (playerTotalSkillScore + 1) * shipClassScore * 10000

	// マッチのおける脅威レベルの補正
	modified := correctBasedOnMatch(raw, f.warships, f.shipID, shipAAIndex, isCVMatch, topTier, bottomTier)

	return data.ThreatLevel{
		Raw:      raw,
		Modified: modified,
	}
}

//nolint:nonamedreturns
func matchInfo(
	tempArenaInfo data.TempArenaInfo,
	warships data.Warships,
) (isCVMatch bool, topTier uint, bottomTier uint) {
	isCVMatch = false
	topTier = 1
	bottomTier = 1

	for _, vehicle := range tempArenaInfo.Vehicles {
		warship, ok := warships[vehicle.ShipID]
		if !ok {
			continue
		}

		if warship.Type == data.ShipTypeCV {
			isCVMatch = true
		}

		tier := warship.Tier

		// Note: 超艦艇や1個差のみのマッチが考慮されてないが、闇深の実装に合わせる
		// Note: すべてが同じTierの場合、topもbottomも1になる
		if topTier < tier {
			topTier = tier
			if topTier > 2 {
				bottomTier = topTier - 2
			} else {
				bottomTier = 1
			}
		}
		if bottomTier > tier {
			bottomTier = tier
			if bottomTier > 8 {
				topTier = 10
			} else {
				topTier = bottomTier + 2
			}
		}
	}

	return isCVMatch, topTier, bottomTier
}

func baseDamageScore(overallAvgDamage float64) float64 {
	return floorU4((limitedValue(floor(overallAvgDamage)/40000, 1.5, 0.5) - 1) / 2)
}

func killScore(
	overallBattles uint,
	overallAvgKill float64,
	overallKdRate float64,
) (float64, float64) {
	if overallBattles == 0 {
		return floorU4(-0.3 / 5), floorU4(-0.5 / 5)
	}

	killScore := limitedValue(overallAvgKill, 1.5, 0.5) - 0.8
	kdRateScore := limitedValue(floor(round(overallKdRate, 2)), 3, 0.7) - 1.2

	// KDRは高いのにKPRは低い（≒芋）は逆補正に
	if killScore < 0 && kdRateScore > 0 {
		kdRateScore *= -1
	}

	return floorU4(killScore / 5), floorU4(kdRateScore / 5)
}

func winRateScore(overallWinRate float64) float64 {
	return floorU4((limitedValue(overallWinRate, 60, 30) - 50) / 100 * 3)
}

func playerOverallScore(
	overallBattles uint,
	overallAvgDamage float64,
	overallAvgKill float64,
	overallKdRate float64,
	overallWinRate float64,
) float64 {
	// プレイヤー総合補正指数を一旦算出
	baseDamageScore := baseDamageScore(overallAvgDamage)
	killScore, kdRateScore := killScore(overallBattles, overallAvgKill, overallKdRate)
	winRateScore := winRateScore(overallWinRate)

	playerGeneralScore := floorU4(
		baseDamageScore + killScore + kdRateScore + winRateScore)

	// プレイ回数関連補正処理
	battlesCountScore := limitedValue(float64(overallBattles)/1000, 3, 0.5)

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

	playerGeneralScore = floorU4((battlesCountScore-1.5)*0.05) + playerGeneralScore

	return playerGeneralScore
}

//nolint:cyclop
func playerShipScore(
	warships data.Warships,
	shipID int,
	shipBattles uint,
	shipAvgDamage float64,
	shipSurvivedRate float64,
	shipAvgPlanesKilled float64,
	shipWinRate float64,
) float64 {
	// 艦の戦績が無い場合は総合戦績から補正なしで
	if shipBattles <= 3 {
		return 0.75
	}

	warship, ok := warships[shipID]
	if !ok {
		return 0.75
	}
	shipType := warship.Type
	shipTier := float64(warship.Tier)

	type shipClassStd struct {
		damage       float64 // 目標与ダメージ基準値
		aaScore      float64 // 制空目標基準値
		survivedRate float64 // 生存率基準目標値
		influence    float64 // 影響度係数
	}

	// 艦の使用回数が3回より多い場合のみ参考に
	// 艦種ごとの基準を設定
	var std shipClassStd
	//nolint:exhaustive
	switch shipType {
	case data.ShipTypeDD:
		std = shipClassStd{
			damage:       shipTier * 4000,
			aaScore:      0,
			survivedRate: 0.4,
			influence:    1.3,
		}
	case data.ShipTypeCL:
		std = shipClassStd{
			damage:       shipTier * 6000,
			aaScore:      0,
			survivedRate: 0.5,
			influence:    1,
		}
	case data.ShipTypeBB:
		std = shipClassStd{
			damage:       shipTier * 7200,
			aaScore:      0,
			survivedRate: 0.45,
			influence:    1.1,
		}
	case data.ShipTypeCV:
		std = shipClassStd{
			damage:       shipTier * 8000,
			aaScore:      (shipTier - 2) * 3.5,
			survivedRate: 0.7,
			influence:    1.7,
		}
	default:
		std = shipClassStd{
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
	shipDamageScore := floorU4(limitedValue(floor(shipAvgDamage)/std.damage, 1.5, 0.1) - 1)

	// 生存率補正指数
	// 艦種ごとの影響度を考慮して補正
	surviveRate := shipSurvivedRate / 100 // %で与えられるため0~100に変換する
	limitedSurviveRate := limitedValue(surviveRate, std.survivedRate+0.1, std.survivedRate-0.1)
	shipServiveScore := floorU4(limitedSurviveRate * std.influence / 1.5)

	// 空母の場合、制空の補正を考慮
	var shipAAScore float64
	if shipType == data.ShipTypeCV {
		shipAAScore = shipAvgPlanesKilled / std.aaScore
		if shipAAScore > 1 {
			shipAAScore = floorU4((shipAAScore - 1) / 2)
		} else {
			shipAAScore = floorU4(shipAAScore - 1)
		}
		shipAAScore = limitedValue(shipAAScore, 0.2, -0.3)
	}

	// 艦勝率補正指数
	winRate := shipWinRate / 100 // %で与えられるため0~100に変換する
	shipWinRateScore := limitedValue(winRate, 0.6, 0.4) - 0.5
	// 空母の場合、勝率の影響を3割増加
	if shipType == data.ShipTypeCV {
		shipWinRateScore *= 1.3
	}
	shipWinRateScore = floorU4(shipWinRateScore * std.influence * 1.5)

	return floorU4(shipDamageScore + shipServiveScore + shipAAScore + shipWinRateScore)
}

func antiAirCoefficient(
	shipID int,
	shipAvgPlanesKilled float64,
) float64 {
	specialAAShips := specialAAShips()
	specialAAShip, ok := specialAAShips[shipID]
	if !ok {
		return 1.0
	}

	planesKilledRatio := shipAvgPlanesKilled / specialAAShip.avg * 1.5
	if planesKilledRatio > 1 {
		planesKilledRatio = floorU4((planesKilledRatio-1)/4) + 1
	}
	limitedPlanesKilledRatio := limitedValue(planesKilledRatio, 1.2, 1)

	return limitedPlanesKilledRatio + specialAAShip.coef
}

func shipClassScore(
	warships data.Warships,
	shipID int,
) float64 {
	result := 1.0

	// 存在しない艦の場合はデフォルト値
	warship, ok := warships[shipID]
	if !ok {
		return result
	}
	shipTier := warship.Tier

	// 特別補正艦はベースをその値にする
	specialShipScores := specialShipScores()
	score, ok := specialShipScores[shipID]
	if ok {
		result = score
	}

	// Tier10艦は性能ジャンプが大きいので追加補正
	if shipTier == 10 {
		return result + 0.15
	}

	return result
}

func correctBasedOnMatch(
	raw float64,
	warships data.Warships,
	shipID int,
	shipAAIndex float64,
	isCVMatch bool,
	topTier uint,
	bottomTier uint,
) float64 {
	result := raw
	if isCVMatch {
		result = raw * shipAAIndex
	}

	warship, ok := warships[shipID]
	if !ok {
		return result
	}
	shipTier := warship.Tier

	if shipTier == topTier {
		result *= 1.1
	}
	if shipTier == bottomTier {
		result *= 0.9
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

func floorU4(value float64) float64 {
	pow := math.Pow(10, 4)
	return floor(value*pow) / pow
}

func floor(value float64) float64 {
	result, _ := decimal.NewFromFloat(value).Floor().Float64()
	return result
}

func round(value float64, digits uint16) float64 {
	result, _ := decimal.NewFromFloat(value).Round(int32(digits)).Float64()
	return result
}
