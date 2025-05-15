package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUWGClansAutocomplete_HexColor(t *testing.T) {
	t.Parallel()

	instance := UWGClansAutocomplete{
		SearchAutocompleteResult: []struct {
			HexColor string `json:"hex_color"`
			Tag      string `json:"tag"`
			ID       int    `json:"id"`
		}{
			{HexColor: "#000000", Tag: "TEST", ID: 1},
			{HexColor: "#000001", Tag: "TEST2", ID: 2},
		},
	}

	assert.Equal(t, "#000000", instance.HexColor("TEST"))
	assert.Empty(t, instance.HexColor("INVALID"))
}
