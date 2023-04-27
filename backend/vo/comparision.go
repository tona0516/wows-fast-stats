package vo

type ShipComp struct {
    PR Between `json:"pr"`
    Damage Between `json:"damage"`
    WinRate Between `json:"win_rate"`
    KdRate Between `json:"kd_rate"`
}

type OverallComp struct {
    Damage Between `json:"damage"`
    WinRate Between `json:"win_rate"`
    KdRate Between `json:"kd_rate"`
}

type Comparision struct {
    Ship ShipComp `json:"ship"`
    Overall OverallComp `json:"overall"`
}
