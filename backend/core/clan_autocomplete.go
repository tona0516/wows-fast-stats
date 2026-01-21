package core

type ClanAutocomplete struct {
	SearchAutocompleteResult []struct {
		HexColor string `json:"hex_color"`
		ID       ClanID `json:"id"`
	} `json:"search_autocomplete_result"`
}

func (u ClanAutocomplete) ColorCode(clanID ClanID) string {
	for _, v := range u.SearchAutocompleteResult {
		if clanID == v.ID {
			return v.HexColor
		}
	}

	return ""
}
