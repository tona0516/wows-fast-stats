package data

type PlayerNameColor string

const (
	PlayerNameColorShip    = "ship"
	PlayerNameColorOverall = "overall"
	PlayerNameColorNone    = "none"
)

func PlayerNameColors() []string {
	return []string{
		PlayerNameColorShip,
		PlayerNameColorOverall,
		PlayerNameColorNone,
	}
}
