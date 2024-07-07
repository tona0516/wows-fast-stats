package vo

import (
	"errors"
	"strings"
)

type PlayerName struct {
	ValueObject[string]
}

func NewPlayerName(value string) (PlayerName, error) {
	if len(strings.TrimSpace(value)) == 0 {
		return PlayerName{}, errors.New("empty_player_name")
	}

	return PlayerName{ValueObject[string]{value}}, nil
}

func (n *PlayerName) IsBot() bool {
	v := n.value
	return (strings.HasPrefix(v, ":") && strings.HasSuffix(v, ":")) || strings.HasPrefix(v, "IDS_OP")
}
