package data

type UWGClansAutocomplete struct {
	SearchAutocompleteResult []struct {
		HexColor string `json:"hex_color"`
		Tag      string `json:"tag"`
		ID       int    `json:"id"`
	} `json:"search_autocomplete_result"`
}

func (u UWGClansAutocomplete) HexColor(clanTag string) string {
	for _, v := range u.SearchAutocompleteResult {
		if clanTag == v.Tag {
			return v.HexColor
		}
	}

	return ""
}
