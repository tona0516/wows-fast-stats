package yamibuka

import (
	"math"
	"wfs/backend/data"
)

type TeamThreatLevel struct {
}

func CalculateTeamThreatLevel(players data.Players, statsPattern data.StatsPattern) data.TeamThreatLevel {
	if len(players) == 0 {
		return data.TeamThreatLevel{}
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
		case data.StatsPatternPvPSolo:
			score = player.PvPSolo.OverallStats.ThreatLevel.Modified
		case data.StatsPatternPvPAll:
			score = player.PvPAll.OverallStats.ThreatLevel.Modified
		case data.StatsPatternRankSolo:
			score = player.RankSolo.OverallStats.ThreatLevel.Modified
		default:
			continue
		}

		scores = append(scores, score)
	}

	if len(scores) == 0 {
		return data.TeamThreatLevel{}
	}

	maxScore := max(scores)
	mean := geometricMean(scores)

	return data.TeamThreatLevel{
		Average:            mean,
		DissociationDegree: (maxScore/mean - 1) * 100,
		Accuracy:           float64(len(scores)) / float64(len(players)) * 100,
	}
}

func max(values []float64) float64 {
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

func geometricMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	product := 1.0
	for _, value := range values {
		product *= value
	}
	return math.Pow(product, 1/float64(len(values)))
}
