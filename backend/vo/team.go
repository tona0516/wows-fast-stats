package vo

type Team struct {
    Players Players `json:"players"`
    Name string `json:"name"`
    WinRateByShip float64 `json:"win_rate_by_ship"`
    WinRateByPlayer float64 `json:"win_rate_by_player"`
}
