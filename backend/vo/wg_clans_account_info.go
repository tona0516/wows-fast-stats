package vo

type WGClansAccountInfo struct {
	Status string `json:"status"`
	Data   map[int]struct {
		ClanID int `json:"clan_id"`
	} `json:"data"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field"`
		Value   string `json:"value"`
	} `json:"error"`
}

func (w *WGClansAccountInfo) ClanIDs() []int {
	clanIDs := make([]int, 0)
	for i := range w.Data {
		clanID := w.Data[i].ClanID
		if clanID != 0 {
			clanIDs = append(clanIDs, w.Data[i].ClanID)
		}
	}
	return clanIDs
}
