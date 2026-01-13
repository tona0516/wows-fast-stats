package data

type ClanAutocomplete struct {
	SearchAutocompleteResult []struct {
		HexColor string `json:"hex_color"`
		Tag      string `json:"tag"`
		ID       ClanID `json:"id"`
	} `json:"search_autocomplete_result"`
}

func (u ClanAutocomplete) HexColor(clanID ClanID) string {
	for _, v := range u.SearchAutocompleteResult {
		if clanID == v.ID {
			return v.HexColor
		}
	}

	return ""
}
