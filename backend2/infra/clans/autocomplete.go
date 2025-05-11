package clans

type Autocomplete struct {
	SearchAutocompleteResult []struct {
		HexColor string `json:"hex_color"`
		Tag      string `json:"tag"`
		ID       int    `json:"id"`
	} `json:"search_autocomplete_result"`
}

func (a Autocomplete) HexColor(clanTag string) string {
	for _, v := range a.SearchAutocompleteResult {
		if clanTag == v.Tag {
			return v.HexColor
		}
	}

	return ""
}
