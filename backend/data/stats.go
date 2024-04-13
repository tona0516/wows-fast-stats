package data

import (
	"math"
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

type matchDetail struct {
	isCVMatch  bool
	topTier    uint
	bottomTier uint
}

type Stats struct {
	useShipID        int
	accountInfo      WGAccountInfoData
	useShipStats     WGShipsStatsData
	allShipsStats    []WGShipsStatsData
	allExpectedStats ExpectedStats
	warships         Warships
	tempArenaInfo    TempArenaInfo
}

func NewStats(
	useShipID int,
	accountInfo WGAccountInfoData,
	allShipsStats []WGShipsStatsData,
	expectedStats ExpectedStats,
	warships Warships,
	tempArenaInfo TempArenaInfo,
) *Stats {
	var useShipStats WGShipsStatsData
	for _, v := range allShipsStats {
		if v.ShipID == useShipID {
			useShipStats = v
			break
		}
	}

	return &Stats{
		useShipID:        useShipID,
		accountInfo:      accountInfo,
		useShipStats:     useShipStats,
		allShipsStats:    allShipsStats,
		allExpectedStats: expectedStats,
		warships:         warships,
		tempArenaInfo:    tempArenaInfo,
	}
}

func (s *Stats) PR(category StatsCategory, pattern StatsPattern) float64 {
	switch category {
	case StatsCategoryShip:
		values, _ := s.statsValues(pattern)
		battles := values.Battles

		return s.pr(
			PRFactor{
				damage: avgDamage(values.DamageDealt, battles),
				frags:  avgKill(values.Frags, battles),
				wins:   winRate(values.Wins, battles),
			},
			PRFactor{
				damage: s.allExpectedStats[s.useShipID].AverageDamageDealt,
				frags:  s.allExpectedStats[s.useShipID].AverageFrags,
				wins:   s.allExpectedStats[s.useShipID].WinRate,
			},
			battles,
		)

	case StatsCategoryOverall:
		var (
			actual     PRFactor
			expected   PRFactor
			allBattles uint
		)

		for _, ship := range s.allShipsStats {
			values := s.statsValuesForm(ship, pattern)
			battles := values.Battles

			es, ok := s.allExpectedStats[ship.ShipID]
			if !ok {
				continue
			}

			actual.damage += float64(values.DamageDealt)
			actual.frags += float64(values.Frags)
			actual.wins += float64(values.Wins)

			expected.damage += es.AverageDamageDealt * float64(battles)
			expected.frags += es.AverageFrags * float64(battles)
			expected.wins += es.WinRate / 100 * float64(battles)

			allBattles += battles
		}

		return s.pr(actual, expected, allBattles)
	}

	return -1
}

func (s *Stats) Battles(category StatsCategory, pattern StatsPattern) uint {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return ship.Battles
	case StatsCategoryOverall:
		return player.Battles
	}

	return 0
}

func (s *Stats) AvgDamage(category StatsCategory, pattern StatsPattern) float64 {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return avgDamage(ship.DamageDealt, ship.Battles)
	case StatsCategoryOverall:
		return avgDamage(player.DamageDealt, player.Battles)
	}

	return 0
}

func (s *Stats) MaxDamage(category StatsCategory, pattern StatsPattern) MaxDamage {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return MaxDamage{
			Value: ship.MaxDamageDealt,
		}
	case StatsCategoryOverall:
		shipID := player.MaxDamageDealtShipID
		warship := s.warships[shipID]
		return MaxDamage{
			ShipID:   shipID,
			ShipName: warship.Name,
			ShipTier: warship.Tier,
			Value:    player.MaxDamageDealt,
		}
	}

	return MaxDamage{}
}

func (s *Stats) KdRate(category StatsCategory, pattern StatsPattern) float64 {
	var (
		survivedBattles uint
		frags           uint
		battles         uint
	)

	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		survivedBattles = ship.SurvivedBattles
		frags = ship.Frags
		battles = ship.Battles
	case StatsCategoryOverall:
		survivedBattles = player.SurvivedBattles
		frags = player.Frags
		battles = player.Battles
	}

	death := battles - survivedBattles
	if death < 1 {
		death = 1
	}

	return float64(frags) / float64(death)
}

func (s *Stats) AvgKill(category StatsCategory, pattern StatsPattern) float64 {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return avgKill(ship.Frags, ship.Battles)
	case StatsCategoryOverall:
		return avgKill(player.Frags, player.Battles)
	}

	return 0
}

func (s *Stats) AvgExp(category StatsCategory, pattern StatsPattern) float64 {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return div(ship.Xp, ship.Battles)
	case StatsCategoryOverall:
		return div(player.Xp, player.Battles)
	}

	return 0
}

func (s *Stats) WinRate(category StatsCategory, pattern StatsPattern) float64 {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return winRate(ship.Wins, ship.Battles)
	case StatsCategoryOverall:
		return winRate(player.Wins, player.Battles)
	}

	return 0
}

func (s *Stats) SurvivedRate(category StatsCategory, pattern StatsPattern) float64 {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return percentage(ship.SurvivedBattles, ship.Battles)
	case StatsCategoryOverall:
		return percentage(player.SurvivedBattles, player.Battles)
	}

	return 0
}

func (s *Stats) WinSurvivedRate(category StatsCategory, pattern StatsPattern) float64 {
	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		return percentage(ship.SurvivedWins, ship.Wins)
	case StatsCategoryOverall:
		return percentage(player.SurvivedWins, player.Wins)
	}

	return 0
}

func (s *Stats) LoseSurvivedRate(category StatsCategory, pattern StatsPattern) float64 {
	var (
		battles         uint
		wins            uint
		survivedBattles uint
		survivedWins    uint
	)

	ship, player := s.statsValues(pattern)
	switch category {
	case StatsCategoryShip:
		battles = ship.Battles
		wins = ship.Wins
		survivedBattles = ship.SurvivedBattles
		survivedWins = ship.SurvivedWins
	case StatsCategoryOverall:
		battles = player.Battles
		wins = player.Wins
		survivedBattles = player.SurvivedBattles
		survivedWins = player.SurvivedWins
	}

	loses := battles - wins
	return percentage(survivedBattles-survivedWins, loses)
}

func (s *Stats) MainBatteryHitRate(pattern StatsPattern) float64 {
	ship, _ := s.statsValues(pattern)
	return percentage(ship.MainBattery.Hits, ship.MainBattery.Shots)
}

func (s *Stats) TorpedoesHitRate(pattern StatsPattern) float64 {
	ship, _ := s.statsValues(pattern)
	return percentage(ship.Torpedoes.Hits, ship.Torpedoes.Shots)
}

func (s *Stats) PlanesKilled(pattern StatsPattern) float64 {
	ship, _ := s.statsValues(pattern)
	return div(ship.PlanesKilled, ship.Battles)
}

func (s *Stats) AvgTier(
	pattern StatsPattern,
) float64 {
	var (
		sum        uint
		allBattles uint
	)

	for _, stats := range s.allShipsStats {
		warship, ok := s.warships[stats.ShipID]
		if !ok {
			continue
		}

		values := s.statsValuesForm(stats, pattern)
		sum += values.Battles * warship.Tier
		allBattles += values.Battles
	}

	return div(sum, allBattles)
}

func (s *Stats) UsingTierRate(
	pattern StatsPattern,
) TierGroup {
	tierGroupMap := make(map[string]uint)

	for _, ship := range s.allShipsStats {
		warship, ok := s.warships[ship.ShipID]
		if !ok {
			continue
		}

		values := s.statsValuesForm(ship, pattern)
		tier := warship.Tier
		battles := values.Battles
		switch {
		case tier >= 1 && tier <= 4:
			tierGroupMap["low"] += battles
		case tier >= 5 && tier <= 7:
			tierGroupMap["middle"] += battles
		case tier >= 8:
			tierGroupMap["high"] += battles
		}
	}

	var allBattles uint
	for _, v := range tierGroupMap {
		allBattles += v
	}

	return TierGroup{
		Low:    percentage(tierGroupMap["low"], allBattles),
		Middle: percentage(tierGroupMap["middle"], allBattles),
		High:   percentage(tierGroupMap["high"], allBattles),
	}
}

func (s *Stats) UsingShipTypeRate(
	pattern StatsPattern,
) ShipTypeGroup {
	shipTypeMap := make(map[ShipType]uint)

	for _, ship := range s.allShipsStats {
		warship, ok := s.warships[ship.ShipID]
		if !ok {
			continue
		}

		values := s.statsValuesForm(ship, pattern)
		shipTypeMap[warship.Type] += values.Battles
	}

	var allBattles uint
	for _, v := range shipTypeMap {
		allBattles += v
	}

	return ShipTypeGroup{
		SS: percentage(shipTypeMap[ShipTypeSS], allBattles),
		DD: percentage(shipTypeMap[ShipTypeDD], allBattles),
		CL: percentage(shipTypeMap[ShipTypeCL], allBattles),
		BB: percentage(shipTypeMap[ShipTypeBB], allBattles),
		CV: percentage(shipTypeMap[ShipTypeCV], allBattles),
	}
}

func (s *Stats) PlatoonRate(
	category StatsCategory,
) float64 {
	switch category {
	case StatsCategoryShip:
		stats := s.useShipStats
		return platoonRate(
			stats.Pvp.Battles,
			stats.PvpSolo.Battles,
			stats.PvpDiv2.Battles,
			stats.PvpDiv3.Battles,
		)
	case StatsCategoryOverall:
		stats := s.accountInfo.Statistics
		return platoonRate(
			stats.Pvp.Battles,
			stats.PvpSolo.Battles,
			stats.PvpDiv2.Battles,
			stats.PvpDiv3.Battles,
		)
	}

	return 0
}

func (s *Stats) ThreatLevel() ThreatLevel {
	// 戦闘情報の取得
	battleDetail := s.battleDetail()

	// プレイヤー総合補正指数
	playerOverallScore := s.playerOverallScore()

	// 艦成績補正
	playerShipScore := s.playerShipScore()

	// 最後に数値の幅を作る係数を設定
	playerTotalSkillScore := (playerOverallScore + playerShipScore) * 0.5

	// 特に補正が必要だと思う艦級を含めた脅威度補正
	shipClassScore := s.shipClassScore()

	// AA特化艦補正
	shipAAIndex := s.antiAirCoefficient()

	// 脅威レベルの算出
	threatLevel := math.Round((playerTotalSkillScore + 1) * shipClassScore * 10000)

	// マッチのおける脅威レベルの補正
	threatLevelInMatch := s.correctBasedOnMatch(threatLevel, shipAAIndex, battleDetail)

	return ThreatLevel{
		Raw:      threatLevel,
		Modified: threatLevelInMatch,
	}
}

func (s *Stats) statsValues(pattern StatsPattern) (WGShipStatsValues, WGPlayerStatsValues) {
	switch pattern {
	case StatsPatternPvPAll:
		return s.useShipStats.Pvp, s.accountInfo.Statistics.Pvp
	case StatsPatternPvPSolo:
		return s.useShipStats.PvpSolo, s.accountInfo.Statistics.PvpSolo
	}

	return WGShipStatsValues{}, WGPlayerStatsValues{}
}

func (s *Stats) statsValuesForm(statsData WGShipsStatsData, pattern StatsPattern) WGShipStatsValues {
	switch pattern {
	case StatsPatternPvPAll:
		return statsData.Pvp
	case StatsPatternPvPSolo:
		return statsData.PvpSolo
	}

	return WGShipStatsValues{}
}

func (s *Stats) pr(
	actual PRFactor,
	expected PRFactor,
	battles uint,
) float64 {
	if battles < 1 {
		return -1
	}

	ratio := PRFactor{
		damage: actual.damage / expected.damage,
		frags:  actual.frags / expected.frags,
		wins:   actual.wins / expected.wins,
	}

	if !ratio.Valid() {
		return -1
	}

	norm := PRFactor{
		damage: math.Max(0, (ratio.damage-0.4)/(1-0.4)),
		frags:  math.Max(0, (ratio.frags-0.1)/(1-0.1)),
		wins:   math.Max(0, (ratio.wins-0.7)/(1-0.7)),
	}

	return 700*norm.damage + 300*norm.frags + 150*norm.wins
}

func (s *Stats) battleDetail() matchDetail {
	result := matchDetail{
		isCVMatch:  false,
		topTier:    1,
		bottomTier: 1,
	}

	for _, vehicle := range s.tempArenaInfo.Vehicles {
		warship, ok := s.warships[vehicle.ShipID]
		if !ok {
			continue
		}

		if warship.Type == ShipTypeCV {
			result.isCVMatch = true
		}

		tier := warship.Tier

		// Note: 超艦艇や1個差のみのマッチが考慮されてないが、闇深の実装に合わせる
		if result.topTier < tier {
			result.topTier = tier
			if result.topTier > 2 {
				result.bottomTier = result.topTier - 2
			} else {
				result.bottomTier = 1
			}
		}
		if result.bottomTier > tier {
			result.bottomTier = tier
			if result.bottomTier > 8 {
				result.topTier = 10
			} else {
				result.topTier = result.bottomTier + 2
			}
		}
	}

	return result
}

func (s *Stats) baseDamageScore() float64 {
	avgDamage := toInteger(s.AvgDamage(StatsCategoryOverall, StatsPatternPvPAll))
	return round((limitedValue(avgDamage/40000, 1.5, 0.5) - 1) / 2)
}

func (s *Stats) killScore() struct {
	killScore   float64
	kdRateScore float64
} {
	battles := s.Battles(StatsCategoryOverall, StatsPatternPvPAll)
	if battles == 0 {
		return struct {
			killScore   float64
			kdRateScore float64
		}{
			killScore:   round(-0.3 / 5),
			kdRateScore: round(-0.5 / 5),
		}
	}

	avgKill := s.AvgKill(StatsCategoryOverall, StatsPatternPvPAll)
	kdRate := s.KdRate(StatsCategoryOverall, StatsPatternPvPAll)

	killScore := limitedValue(avgKill, 1.5, 0.5) - 0.8
	kdRateScore := limitedValue(kdRate, 3, 0.7) - 1.2

	// KDRは高いのにKPRは低い（≒芋）は逆補正に
	if killScore < 0 && kdRateScore > 0 {
		kdRateScore *= -1
	}

	return struct {
		killScore   float64
		kdRateScore float64
	}{
		killScore:   round(killScore / 5),
		kdRateScore: round(kdRateScore / 5),
	}
}

func (s *Stats) winRateScore() float64 {
	winRate := s.WinRate(StatsCategoryOverall, StatsPatternPvPAll)
	return round((limitedValue(winRate, 60, 30) - 50) / 100 * 3)
}

func (s *Stats) playerOverallScore() float64 {
	// プレイヤー総合補正指数を一旦算出
	ks := s.killScore()

	playerGeneralScore := round(s.baseDamageScore() + ks.killScore + ks.kdRateScore + s.winRateScore())

	// プレイ回数関連補正処理
	battles := s.Battles(StatsCategoryOverall, StatsPatternPvPAll)
	battlesCountScore := limitedValue(float64(battles)/1000, 3, 0.5)

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
func (s *Stats) playerShipScore() float64 {
	// 艦の戦績が無い場合は総合戦績から補正なしで
	battles := s.Battles(StatsCategoryShip, StatsPatternPvPAll)
	if battles <= 3 {
		return 0.75
	}

	warship, ok := s.warships[s.useShipID]
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
	case ShipTypeDD:
		std = shipClassStd{
			damage:       shipTier * 4000,
			aaScore:      0,
			survivedRate: 0.4,
			influence:    1.3,
		}
	case ShipTypeCL:
		std = shipClassStd{
			damage:       shipTier * 6000,
			aaScore:      0,
			survivedRate: 0.5,
			influence:    1,
		}
	case ShipTypeBB:
		std = shipClassStd{
			damage:       shipTier * 7200,
			aaScore:      0,
			survivedRate: 0.45,
			influence:    1.1,
		}
	case ShipTypeCV:
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
	avgDamage := s.AvgDamage(StatsCategoryShip, StatsPatternPvPAll)
	shipDamageScore := round(limitedValue(toInteger(avgDamage)/std.damage, 1.5, 0.1) - 1)

	// 生存率補正指数
	// 艦種ごとの影響度を考慮して補正
	surviveRate := s.SurvivedRate(StatsCategoryShip, StatsPatternPvPAll) / 100 // %で与えられるため0~100に変換する
	limitedSurviveRate := limitedValue(surviveRate, std.survivedRate+0.1, std.survivedRate-0.1)
	shipServiveScore := round(limitedSurviveRate * std.influence / 1.5)

	// 空母の場合、制空の補正を考慮
	var shipAAScore float64
	if shipType == ShipTypeCV {
		shipAAScore = s.PlanesKilled(StatsPatternPvPAll) / std.aaScore
		if shipAAScore > 1 {
			shipAAScore = round((shipAAScore - 1) / 2)
		} else {
			shipAAScore = round(shipAAScore - 1)
		}
		shipAAScore = limitedValue(shipAAScore, 0.2, -0.3)
	}

	// 艦勝率補正指数
	winRate := s.WinRate(StatsCategoryShip, StatsPatternPvPAll) / 100 // %で与えられるため0~100に変換する
	shipWinRateScore := limitedValue(winRate, 0.6, 0.4) - 0.5
	// 空母の場合、勝率の影響を3割増加
	if shipType == ShipTypeCV {
		shipWinRateScore *= 1.3
	}
	shipWinRateScore = round(shipWinRateScore * std.influence * 1.5)

	return round(shipDamageScore + shipServiveScore + shipAAScore + shipWinRateScore)
}

func (s *Stats) antiAirCoefficient() float64 {
	specialAAShip, ok := specialAAShips[s.useShipID]
	if !ok {
		return 1.0
	}

	planesKilledRatio := s.PlanesKilled(StatsPatternPvPAll) / specialAAShip.average * 1.5
	if planesKilledRatio > 1 {
		planesKilledRatio = round((planesKilledRatio-1)/4) + 1
	}
	limitedPlanesKilledRatio := limitedValue(planesKilledRatio, 1.2, 1)

	return limitedPlanesKilledRatio + specialAAShip.coefficient
}

func (s *Stats) shipClassScore() float64 {
	result := 1.0

	// 存在しない艦の場合はデフォルト値
	warship, ok := s.warships[s.useShipID]
	if !ok {
		return result
	}
	shipTier := warship.Tier

	// 特別補正艦はベースをその値にする
	score, ok := specialShipScores[s.useShipID]
	if ok {
		result = score
	}

	// Tier10艦は性能ジャンプが大きいので追加補正
	if shipTier == 10 {
		return result + 0.15
	}

	return result
}

func (s *Stats) correctBasedOnMatch(
	threatLevel float64,
	shipAAIndex float64,
	matchDetail matchDetail,
) float64 {
	result := threatLevel
	if matchDetail.isCVMatch {
		result = math.Round(threatLevel * shipAAIndex)
	}

	warship, ok := s.warships[s.useShipID]
	if !ok {
		return result
	}
	shipTier := warship.Tier

	if shipTier == matchDetail.topTier {
		result = math.Round(result * 1.1)
	}
	if shipTier == matchDetail.bottomTier {
		result = math.Round(result * 0.9)
	}

	return result
}

func avgDamage(damageDealt uint, battles uint) float64 {
	return div(damageDealt, battles)
}

func avgKill(frags uint, battles uint) float64 {
	return div(frags, battles)
}

func winRate(wins uint, battles uint) float64 {
	return percentage(wins, battles)
}

func platoonRate(allBattles uint, soloBattles uint, div2Battles uint, div3Battles uint) float64 {
	soloRate := div(soloBattles, allBattles) * 1
	div2Rate := div(div2Battles, allBattles) * 2
	div3Rate := div(div3Battles, allBattles) * 3
	return soloRate + div2Rate + div3Rate
}

func div(a uint, b uint) float64 {
	if b <= 0 {
		return 0
	}

	return float64(a) / float64(b)
}

func percentage(a uint, b uint) float64 {
	return div(a, b) * 100
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
