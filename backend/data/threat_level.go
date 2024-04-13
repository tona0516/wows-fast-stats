package data

type ThreatLevel struct {
	Raw      float64
	Modified float64
}

// type shipClassStdValues struct {
// 	damage       float64 // 目標与ダメージ基準値
// 	aaScore      float64 // 制空目標基準値
// 	survivedRate float64 // 生存率基準目標値
// 	influence    float64 // 影響度係数
// }

// type ThreatLevelCalc struct {
// 	useShipID     int
// 	warships      Warships
// 	tempArenaInfo TempArenaInfo
// 	statistics    ThreatLevelStatistics
// 	playerName    string
// }

// type ThreatLevelStatistics struct {
// 	PlayerAvgDamage  float64
// 	PlayerAvgKill    float64
// 	PlayerKdRate     float64
// 	PlayerWinRate    float64
// 	PlayerBattles    uint
// 	ShipAvgDamage    float64
// 	ShipWinRate      float64
// 	ShipSurvivedRate float64
// 	ShipPlanesKilled float64
// 	ShipBattles      uint
// }

// type threatLevelBattleDetail struct {
// 	IsCVMatch  bool
// 	BottomTier uint
// 	TopTier    uint
// }

// func NewThreatLevel(
// 	useShipID int,
// 	warships Warships,
// 	tempArenaInfo TempArenaInfo,
// 	statistics ThreatLevelStatistics,
// 	playerName string,
// ) *ThreatLevelCalc {
// 	return &ThreatLevelCalc{
// 		useShipID:     useShipID,
// 		warships:      warships,
// 		tempArenaInfo: tempArenaInfo,
// 		statistics:    statistics,
// 		playerName:    playerName,
// 	}
// }

// func (t *ThreatLevelCalc) Calculate() (float64, float64) {
// 	// 戦闘情報の取得
// 	battleDetail := t.battleDetail()

// 	// プレイヤー総合補正指数
// 	playerOverallScore := t.playerOverallScore()

// 	// 艦成績補正
// 	playerShipScore := t.playerShipScore()

// 	// 最後に数値の幅を作る係数を設定
// 	playerTotalSkillScore := (playerOverallScore + playerShipScore) * 0.5

// 	// 特に補正が必要だと思う艦級を含めた脅威度補正
// 	shipClassScore := t.shipClassScore()

// 	// AA特化艦補正
// 	shipAAIndex := t.antiAirCoefficient()

// 	// 脅威レベルの算出
// 	threatLevel := math.Round((playerTotalSkillScore + 1) * shipClassScore * 10000)

// 	// マッチのおける脅威レベルの補正
// 	threatLevelInMatch := t.correctBasedOnMatch(threatLevel, shipAAIndex, battleDetail)

// 	return threatLevel, threatLevelInMatch
// }
