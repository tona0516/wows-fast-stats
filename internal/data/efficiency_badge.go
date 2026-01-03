package data

type EfficiencyBadge string

const (
	EfficiencyBadgeExpert = "E"
	EfficiencyBadgeFirst  = "1"
	EfficiencyBadgeSecond = "2"
	EfficiencyBadgeThird  = "3"
	EfficiencyBadgeNone   = ""
)

func EfficiencyBadges() []string {
	return []string{
		EfficiencyBadgeExpert,
		EfficiencyBadgeFirst,
		EfficiencyBadgeSecond,
		EfficiencyBadgeThird,
		EfficiencyBadgeNone,
	}
}
