package vo

type TeamAverage struct {
    PersonalRating float64 `json:"personal_rating"`
    DamageByShip float64 `json:"damage_by_ship"`
    WinRateByShip float64 `json:"win_rate_by_ship"`
    KdRateByShip float64 `json:"kd_rate_by_ship"`
    DamageByPlayer float64 `json:"damage_by_player"`
    WinRateByPlayer float64 `json:"win_rate_by_player"`
    KdRateByPlayer float64 `json:"kd_rate_by_player"`
}
