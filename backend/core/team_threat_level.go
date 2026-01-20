package core

type TeamThreatLevel struct {
	Average            float64 `json:"average"`
	DissociationDegree float64 `json:"dissociationDegree"`
	Accuracy           float64 `json:"accuracy"`
}

func NewTeamThreatLevel(players Players, statsPattern StatsPattern) TeamThreatLevel {
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

	maxScore := max(scores)
	mean := geometricMean(scores)

	return TeamThreatLevel{
		Average:            mean,
		DissociationDegree: (safeDivide(maxScore, mean) - 1) * 100,
		Accuracy:           safeDivide(len(scores), len(players)) * 100,
	}
}
