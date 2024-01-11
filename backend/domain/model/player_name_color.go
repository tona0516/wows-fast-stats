package model

type PlayerNameColor string

//nolint:gochecknoglobals
var PlayerNameColors = []string{
	PlayerNameColorShip,
	PlayerNameColorOverall,
	PlayerNameColorNone,
}

const (
	PlayerNameColorShip    = "ship"
	PlayerNameColorOverall = "overall"
	PlayerNameColorNone    = "none"
)
