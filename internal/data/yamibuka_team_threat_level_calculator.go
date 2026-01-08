package data

import "math"

// CalculateTeamThreatLevel calculates the threat level statistics for a team.
func CalculateTeamThreatLevel(players Players, statsPattern StatsPattern) TeamThreatLevel {
	if len(players) == 0 {
		return TeamThreatLevel{}
	}

	scores := make([]float64, 0, len(players))
	for _, player := range players {
		if player.PlayerInfo.ID == 0 {
			continue
		}

		if player.PlayerInfo.IsHidden {
			continue
		}

		var score float64
		switch statsPattern {
		case StatsPatternPvPSolo:
			score = player.PvPSolo.OverallStats.ThreatLevel.Modified
		case StatsPatternPvPAll:
			score = player.PvPAll.OverallStats.ThreatLevel.Modified
		case StatsPatternRankSolo:
			score = player.RankSolo.OverallStats.ThreatLevel.Modified
		default:
			continue
		}

		scores = append(scores, score)
	}

	if len(scores) == 0 {
		return TeamThreatLevel{}
	}

	maxScore := maxThreatLevelScore(scores)
	mean := geometricMeanThreatLevel(scores)

	return TeamThreatLevel{
		Average:            mean,
		DissociationDegree: (maxScore/mean - 1) * 100,
		Accuracy:           float64(len(scores)) / float64(len(players)) * 100,
	}
}

func maxThreatLevelScore(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	maxValue := values[0]
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

func geometricMeanThreatLevel(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	product := 1.0
	for _, value := range values {
		product *= value
	}
	return math.Pow(product, 1/float64(len(values)))
}
