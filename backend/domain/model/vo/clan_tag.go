package vo

import (
	"errors"
	"strings"
)

type ClanTag struct {
	ValueObject[string]
}

func NewClanTag(value string) (ClanTag, error) {
	if len(strings.TrimSpace(value)) == 0 {
		return ClanTag{}, errors.New("empty_clan_tag")
	}

	return ClanTag{ValueObject[string]{value}}, nil
}
