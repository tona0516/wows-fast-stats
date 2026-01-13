package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClanAutocomplete_HexColor(t *testing.T) {
	t.Parallel()

	instance := ClanAutocomplete{
		SearchAutocompleteResult: []struct {
			HexColor string `json:"hex_color"`
			Tag      string `json:"tag"`
			ID       ClanID `json:"id"`
		}{
			{HexColor: "#000000", Tag: "TEST", ID: 1},
			{HexColor: "#000001", Tag: "TEST2", ID: 2},
		},
	}

	assert.Equal(t, "#000000", instance.HexColor(1))
	assert.Empty(t, instance.HexColor(999))
}
